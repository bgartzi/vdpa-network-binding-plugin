/*
 * This file is part of the KubeVirt project *
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

package virtio_test

import (
	"kubevirt.io/vdpa-network-binding-plugin/sidecar/virtio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("virtio features", func() {
	DescribeTable("contain a specific feature",
		func(features uint64, feature uint8, expected bool) {
			Expect(virtio.Contains(features, feature)).To(Equal(expected))
		},
		Entry("All zeroes", uint64(0), virtio.VIRTIO_NET_F_CTRL_GUEST_OFFLOADS, false),
		Entry("Just one", uint64(1), virtio.VIRTIO_NET_F_CSUM, true),
		Entry("Multiple", uint64(11529215046068469760), virtio.VIRTIO_NET_F_SPEED_DUPLEX, true),
	)
})
