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
 *
 */

package main

import (
	"log"
	"time"

	"github.com/bgartzi/plugins/pkg/sdk/plugin"
	"kubevirt.io/kubevirt/pkg/virt-handler/isolation"

	"kubevirt.io/vdpa-network-binding-plugin/nodehook/callback"
	"kubevirt.io/vdpa-network-binding-plugin/nodehook/process"
	"kubevirt.io/vdpa-network-binding-plugin/nodehook/vmi"
)

const (
	VdpaPluginName string = "vdpa"
)

func main() {
	log.Printf("Starting nodehook")
	//TODO: parametrize the NBP name
	vmiFilter := vmi.NewNBPVMIFilter(VdpaPluginName)

	launcherDetector := isolation.NewSocketBasedIsolationDetector()
	virtqemudProcessFilter := process.NewVirtqemudProcessFinder(launcherDetector)

	preStartCallback := callback.NewMemlockNodeHookCallback(
		string(plugin.PreVMStart),
		vmiFilter,
		virtqemudProcessFilter,
	)
	preMigrationTargetCallback := callback.NewMemlockNodeHookCallback(
		string(plugin.PreMigrationTarget),
		vmiFilter,
		virtqemudProcessFilter,
	)

	preStartVMHandler := plugin.NodeHandler(preStartCallback).
		WithFailureStrategy(plugin.Fail).
		WithTimeout(10 * time.Second)
	preMigrationTargetHandler := plugin.NodeHandler(preMigrationTargetCallback).
		WithFailureStrategy(plugin.Fail).
		WithTimeout(10 * time.Second)

	log.Printf("Handlers set up")
	p := plugin.New("vdpa-network-binding-plugin-node-plugin-hook").
		// TODO: ISn't it weird that we need to declare the integration
		// point here, yet req contains which hook its get triggered
		// for?
		WithNodeHook(plugin.PreVMStart, preStartVMHandler).
		WithNodeHook(plugin.PreMigrationTarget, preMigrationTargetHandler)
	log.Printf("Serving")
	p.Execute()
}
