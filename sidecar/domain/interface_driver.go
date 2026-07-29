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
	"kubevirt.io/vdpa-network-binding-plugin/sidecar/virtio"

	"libvirt.org/go/libvirtxml"
)

const (
	LIBVIRT_DOMXL_TRISTATE_ENABLED  = "on"
	LIBVIRT_DOMXL_TRISTATE_DISABLED = "off"
)

func defaultInterfaceDriver() *libvirtxml.DomainInterfaceDriver {
	return &libvirtxml.DomainInterfaceDriver{
		Name: "vhost",
	}
}

func driverWithCoreFeatures(driver *libvirtxml.DomainInterfaceDriver, features uint64) *libvirtxml.DomainInterfaceDriver {
	if driver == nil {
		driver = defaultInterfaceDriver()
	}

	if virtio.Contains(features, virtio.VIRTIO_NET_F_HASH_REPORT) {
		driver.RSSHashReport = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_RSS) {
		driver.RSS = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	return driver
}

func guestDriverConfigWithFeatures(features uint64) *libvirtxml.DomainInterfaceDriverGuest {
	guest := libvirtxml.DomainInterfaceDriverGuest{}

	if virtio.Contains(features, virtio.VIRTIO_NET_F_GUEST_CSUM) {
		guest.CSum = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_GUEST_TSO4) {
		guest.TSO4 = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_GUEST_TSO6) {
		guest.TSO6 = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_GUEST_ECN) {
		guest.ECN = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_GUEST_UFO) {
		guest.UFO = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}

	if guest != (libvirtxml.DomainInterfaceDriverGuest{}) {
		return &guest
	}

	return nil
}

func hostDriverConfigWithFeatures(features uint64) *libvirtxml.DomainInterfaceDriverHost {
	host := libvirtxml.DomainInterfaceDriverHost{}

	if virtio.Contains(features, virtio.VIRTIO_NET_F_CSUM) {
		host.CSum = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_CTRL_GUEST_OFFLOADS) {
		host.GSO = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_HOST_TSO4) {
		host.TSO4 = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_HOST_TSO6) {
		host.TSO6 = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_HOST_ECN) {
		host.ECN = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}
	if virtio.Contains(features, virtio.VIRTIO_NET_F_HOST_UFO) {
		host.UFO = LIBVIRT_DOMXL_TRISTATE_ENABLED
	}

	if host != (libvirtxml.DomainInterfaceDriverHost{}) {
		return &host
	}
	return nil
}

func driverWithVirtioFeatures(driver *libvirtxml.DomainInterfaceDriver, features uint64) *libvirtxml.DomainInterfaceDriver {
	if driver == nil {
		driver = defaultInterfaceDriver()
	}

	driver = driverWithCoreFeatures(driver, features)

	guestFeatures := guestDriverConfigWithFeatures(features)
	if guestFeatures != nil {
		driver.Guest = guestFeatures
	}

	hostFeatures := hostDriverConfigWithFeatures(features)
	if hostFeatures != nil {
		driver.Host = hostFeatures
	}

	if *driver != *defaultInterfaceDriver() {
		return driver
	}
	return nil
}

func driverWithQueues(driver *libvirtxml.DomainInterfaceDriver, queues uint) *libvirtxml.DomainInterfaceDriver {
	if queues > 0 {
		if driver == nil {
			driver = defaultInterfaceDriver()
		}

		driver.Queues = queues
	}

	return driver
}

func NewNetInterfaceDriver(maxVirtQueuePairs uint, virtioFeatures uint64) *libvirtxml.DomainInterfaceDriver {
	driver := driverWithVirtioFeatures(nil, virtioFeatures)
	return driverWithQueues(driver, maxVirtQueuePairs)
}
