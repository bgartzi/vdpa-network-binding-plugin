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
 */

package memlock

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/unix"

	v1 "kubevirt.io/api/core/v1"

	"kubevirt.io/vdpa-network-binding-plugin/nodehook/process"
)

func ConfigureMemlockRLimits(
	vmi *v1.VirtualMachineInstance,
	processFinder process.VmiProcessFinder,
) error {
	targetPid, err := processFinder.Find(vmi)
	if err != nil {
		return err
	} else if targetPid == nil {
		return fmt.Errorf("did not find %s", processFinder.Info())
	}

	log.Printf("setting memlock limits to process %d to handle vmi %q", targetPid, vmi.Name)
	return setProcessMemlockLimits(*targetPid, unix.RLIM_INFINITY)
}

// Taken from kubevirt/pkg/hypervisor/common
func setProcessMemlockLimits(pid int, size uint64) error {
	// standard golang libraries don't provide API to set runtime limits
	// for other processes, so we have to directly call to kernel
	rlimit := unix.Rlimit{
		Cur: size,
		Max: size,
	}
	_, _, errno := unix.RawSyscall6(unix.SYS_PRLIMIT64,
		uintptr(pid),
		uintptr(unix.RLIMIT_MEMLOCK),
		uintptr(unsafe.Pointer(&rlimit)), // #nosec used in unix RawSyscall6
		0, 0, 0)
	if errno != 0 {
		log.Printf("ERROR: could not set prlimit of process %d to %d: errno %v", pid, size, errno)
		return fmt.Errorf("error setting prlimit: %v", errno)
	}

	log.Printf("prlimit of process %d set to %d", pid, size)

	return nil
}
