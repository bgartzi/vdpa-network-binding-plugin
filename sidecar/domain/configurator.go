/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */

package domain

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"kubevirt.io/client-go/log"

	vmschema "kubevirt.io/api/core/v1"
	"kubevirt.io/kubevirt/pkg/apimachinery/wait"
	"kubevirt.io/kubevirt/pkg/network/downwardapi"
	netnamescheme "kubevirt.io/kubevirt/pkg/network/namescheme"
	"kubevirt.io/kubevirt/pkg/network/vmispec"

	libvirtxml "libvirt.org/go/libvirtxml"

	"kubevirt.io/vdpa-network-binding-plugin/sidecar/symlink"
)

type VdpaIfaceConfig struct {
	vmiSpecIface *vmschema.Interface
	symlinkName  string
	*downwardapi.Interface
}

type VdpaNetworkConfigurator struct {
	vdpaConfigs   []*VdpaIfaceConfig
	containerName string
}

const (
	// VdpaPluginName vdpa binding plugin name should be registered to Kubevirt through Kubevirt CR
	VdpaPluginName = "vdpa"
)

func readFileUntilNotEmpty(networkPCIMapPath string) ([]byte, error) {
	var networkPCIMapBytes []byte

	err := wait.PollImmediately(100*time.Millisecond, 10*time.Second, func(_ context.Context) (bool, error) {
		var err error
		networkPCIMapBytes, err = os.ReadFile(networkPCIMapPath)
		return len(networkPCIMapBytes) > 0, err
	})

	return networkPCIMapBytes, err
}

func GetDownwardAPINetworkInfo(filePath string) (*downwardapi.NetworkInfo, error) {
	networkPCIMapBytes, err := readFileUntilNotEmpty(filePath)
	if err != nil {
		return nil, err
	}

	result := &downwardapi.NetworkInfo{}
	err = json.Unmarshal(networkPCIMapBytes, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func lookupNetworkInfoByName(name string, netInfo *downwardapi.NetworkInfo) (*downwardapi.Interface, error) {
	for _, ifaceInfo := range netInfo.Interfaces {
		if name == ifaceInfo.Network {
			return &ifaceInfo, nil
		}
	}
	return nil, fmt.Errorf("failed to find networkinfo for interface %s", name)
}

func NewVdpaNetworkConfigurator(
	ifaces []vmschema.Interface,
	networks []vmschema.Network,
	netInfo *downwardapi.NetworkInfo,
	containerName string,
) (*VdpaNetworkConfigurator, error) {
	netNameSchema := netnamescheme.CreateHashedNetworkNameScheme(networks)
	var configs []*VdpaIfaceConfig
	for _, net := range networks {
		if net.Multus != nil {
			iface := vmispec.LookupInterfaceByName(ifaces, net.Name)
			if iface == nil {
				return nil, fmt.Errorf("no interface named %s found", net.Name)
			}

			if iface.Binding == nil || iface.Binding.Name != VdpaPluginName {
				log.Log.Infof("interface %q is not set with Vdpa network binding plugin", net.Name)
				continue
			}

			ifaceInfo, err := lookupNetworkInfoByName(net.Name, netInfo)
			if err != nil {
				return nil, err
			}

			podNetName := netNameSchema[net.Name]

			configs = append(configs,
				&VdpaIfaceConfig{
					vmiSpecIface: iface,
					symlinkName:  podNetName,
					Interface:    ifaceInfo,
				},
			)
		}
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("no vdpa interface found")
	}

	return &VdpaNetworkConfigurator{
			vdpaConfigs:   configs,
			containerName: containerName,
		},
		nil
}

func (p VdpaNetworkConfigurator) Mutate(domainSpec *libvirtxml.Domain) (*libvirtxml.Domain, error) {
	generatedIfaces, err := p.generateInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to generate domain interface spec: %v", err)
	}

	if domainSpec.Devices == nil {
		domainSpec.Devices = &libvirtxml.DomainDeviceList{}
	}

	for i, config := range p.vdpaConfigs {
		if iface := lookupIfaceByAliasName(domainSpec.Devices.Interfaces, config.vmiSpecIface.Name); iface != nil {
			*iface = *generatedIfaces[i]
		} else {
			domainSpec.Devices.Interfaces = append(domainSpec.Devices.Interfaces, *generatedIfaces[i])
		}

		ifaceInfo, _ := xml.Marshal(generatedIfaces[i])
		log.Log.Infof("vdpa interface %s is added to domain spec successfully: %s",
			config.vmiSpecIface.Name, string(ifaceInfo))
	}

	return domainSpec, nil
}

func lookupIfaceByAliasName(ifaces []libvirtxml.DomainInterface, name string) *libvirtxml.DomainInterface {
	for i, iface := range ifaces {
		if iface.Alias != nil && iface.Alias.Name == name {
			return &ifaces[i]
		}
	}

	return nil
}

func (p VdpaNetworkConfigurator) generateInterfaces() ([]*libvirtxml.DomainInterface, error) {
	var domainInterfaces []*libvirtxml.DomainInterface

	for _, cfg := range p.vdpaConfigs {
		var address *libvirtxml.DomainAddress
		var err error
		if cfg.vmiSpecIface.PciAddress != "" {
			address, err = parsePCIAddress(cfg.vmiSpecIface.PciAddress)
			if err != nil {
				return nil, err
			}
		}

		var mac *libvirtxml.DomainInterfaceMAC
		if cfg.vmiSpecIface.MacAddress != "" {
			mac = &libvirtxml.DomainInterfaceMAC{Address: cfg.vmiSpecIface.MacAddress}
		} else if cfg.Mac != "" {
			mac = &libvirtxml.DomainInterfaceMAC{Address: cfg.Mac}
		}

		var acpi *libvirtxml.DomainDeviceACPI
		if cfg.vmiSpecIface.ACPIIndex > 0 {
			acpi = &libvirtxml.DomainDeviceACPI{Index: uint(cfg.vmiSpecIface.ACPIIndex)}
		}

		driver := NewNetInterfaceDriver(
			uint(cfg.DeviceInfo.Vdpa.MaxVQP),
			cfg.DeviceInfo.Vdpa.VirtioFeatures,
		)

		vdpaPath := symlink.SharedComputeSymlinkPath(p.containerName, cfg.symlinkName)

		domainInterfaces = append(domainInterfaces, &libvirtxml.DomainInterface{
			Alias:   &libvirtxml.DomainAlias{Name: cfg.vmiSpecIface.Name},
			Model:   &libvirtxml.DomainInterfaceModel{Type: "virtio"},
			Address: address,
			MAC:     mac,
			ACPI:    acpi,
			Source: &libvirtxml.DomainInterfaceSource{
				VDPA: &libvirtxml.DomainInterfaceSourceVDPA{Device: vdpaPath},
			},
			Driver: driver,
		})
	}

	return domainInterfaces, nil
}

func (p *VdpaNetworkConfigurator) VdpaPathsToSymlinkNames() map[string]string {
	pathsToSymlinks := make(map[string]string, len(p.vdpaConfigs))
	for _, cfg := range p.vdpaConfigs {
		pathsToSymlinks[cfg.DeviceInfo.Vdpa.Path] = cfg.symlinkName
	}
	return pathsToSymlinks
}
