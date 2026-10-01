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

package callback

import (
	"context"
	"fmt"
	"log"

	"github.com/bgartzi/plugins/pkg/sdk/plugin"

	"kubevirt.io/vdpa-network-binding-plugin/nodehook/memlock"
	"kubevirt.io/vdpa-network-binding-plugin/nodehook/process"
	"kubevirt.io/vdpa-network-binding-plugin/nodehook/vmi"
)

type MemlockNodeHookCallback struct {
	hookPoint     string
	vmiFilter     vmi.VMIFilter
	processFilter process.VmiProcessFinder
}

func NewMemlockNodeHookCallback(
	hookPoint string,
	vmiFilter vmi.VMIFilter,
	processFilter process.VmiProcessFinder,
) *MemlockNodeHookCallback {
	return &MemlockNodeHookCallback{hookPoint, vmiFilter, processFilter}
}

func (c *MemlockNodeHookCallback) ExecuteNodeHook(ctx context.Context, req *plugin.NodeHookRequest) error {
	log.Printf("handling node hook, node=%s, point=%s, vmi=%s", req.NodeName, req.HookPoint, req.VMI.Name)

	if req.HookPoint != c.hookPoint {
		log.Printf("expected to handle hook point %q got, %q. Ignoring request", c.hookPoint, req.HookPoint)
		return fmt.Errorf("expected to handle hook point %q. Got %q", c.hookPoint, req.HookPoint)
	}

	if !c.vmiFilter.Matches(req.VMI) {
		log.Printf("vmi %q does not match vmi Filter. Ignoring vmi", req.VMI.Name)
		return nil
	}

	log.Printf("configuring memlock limits")
	return memlock.ConfigureMemlockRLimits(req.VMI, c.processFilter)
}
