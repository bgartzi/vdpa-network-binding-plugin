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

package process

import (
	"log"

	"github.com/mitchellh/go-ps"

	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/kubevirt/pkg/virt-handler/isolation"
)

type VmiProcessFinder interface {
	// Returns the Pid of a process that was found for vmi
	// nil if it was not found,
	// or an error if something unexpected occurred
	Find(vmi *v1.VirtualMachineInstance) (*int, error)

	// Returns a string that the Finder can be identified with
	Info() string
}

type VirtLauncherProcessFinder struct {
	launcher      isolation.PodIsolationDetector
	processFilter ProcessFilter
}

func NewVirtLauncherProcessFinder(
	launcherDetector isolation.PodIsolationDetector,
	processFilter ProcessFilter,
) *VirtLauncherProcessFinder {
	return &VirtLauncherProcessFinder{
		launcherDetector,
		processFilter,
	}
}

func NewVirtqemudProcessFinder(
	launcherDetector isolation.PodIsolationDetector,
) *VirtLauncherProcessFinder {
	return NewVirtLauncherProcessFinder(
		launcherDetector, NewVirtqemudProcessFilter(),
	)
}

func (f *VirtLauncherProcessFinder) Info() string {
	return f.processFilter.Info()
}

func (f *VirtLauncherProcessFinder) Find(vmi *v1.VirtualMachineInstance) (*int, error) {
	// Find virt launcher process based on the isolation finder
	virtLauncher, err := f.launcher.Detect(vmi)
	if err != nil {
		log.Printf("ERROR: did not find virt launcher for vmi %s: %v", vmi.Name, err)
		return nil, err
	}
	// Then list all of its child processes
	launcherProcesses, err := listChildProcesses(virtLauncher.Pid())
	if err != nil {
		log.Printf("ERROR: could not find child processes for the virt launcher of vmi %s: %v", vmi.Name, err)
		return nil, err
	}

	log.Printf("found %d child processes under vmi %s virt-launcher with pid %d",
		len(launcherProcesses), vmi.Name, virtLauncher.Pid())

	for _, process := range launcherProcesses {
		if f.processFilter.Matches(process) {
			log.Printf("vmi/%s,process/%d,exec/%s: it's a match", vmi.Name, process.Pid, process.ExecName)
			return &process.Pid, nil
		}
		log.Printf("vmi/%s,process/%d,exec/%s: not a match", vmi.Name, process.Pid, process.ExecName)
	}

	log.Printf("INFO: did not find child processes matching the search criteria for vmi %s", vmi.Name)
	return nil, nil
}

// Returns a list of child process executable names
func listChildProcesses(pid int) ([]*ProcessInfo, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}

	var res []*ProcessInfo
	for _, p := range processes {
		if p.PPid() == pid {
			res = append(res, &ProcessInfo{p.Pid(), p.Executable()})
		}
	}

	return res, nil
}
