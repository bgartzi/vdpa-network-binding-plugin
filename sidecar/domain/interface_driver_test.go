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

package domain_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"libvirt.org/go/libvirtxml"

	"kubevirt.io/vdpa-network-binding-plugin/sidecar/domain"
	"kubevirt.io/vdpa-network-binding-plugin/sidecar/virtio"
)

var _ = Describe("driver configurator", func() {
	DescribeTable("NewNetInterfaceDriver",
		func(maxVQP uint, features uint64, expected *libvirtxml.DomainInterfaceDriver) {
			driver := domain.NewNetInterfaceDriver(maxVQP, features)
			Expect(driver).To(Equal(expected))
		},
		Entry(
			"should return nil if empty parameters are given",
			uint(0),
			uint64(0),
			nil,
		),
		Entry(
			"should configure just the virtqueues if features not provided",
			uint(10),
			uint64(0),
			&libvirtxml.DomainInterfaceDriver{Name: "vhost", Queues: 10},
		),
		Entry(
			"should return nil if irrelevant virtio features are given",
			uint(0),
			uint64(13027704872),
			nil,
		),
		Entry(
			"should configure just the core driver features if guest or host features are not present",
			uint(0),
			uint64(1<<virtio.VIRTIO_NET_F_RSS),
			&libvirtxml.DomainInterfaceDriver{Name: "vhost", RSS: "on"},
		),
		Entry(
			"should configure just the host features if guest features are not present",
			uint(0),
			uint64(1<<virtio.VIRTIO_NET_F_CSUM),
			&libvirtxml.DomainInterfaceDriver{
				Name: "vhost",
				Host: &libvirtxml.DomainInterfaceDriverHost{
					CSum: "on",
				},
			},
		),
		Entry(
			"should configure just the guest features if host features are not present",
			uint(0),
			uint64(1<<virtio.VIRTIO_NET_F_GUEST_CSUM),
			&libvirtxml.DomainInterfaceDriver{
				Name: "vhost",
				Guest: &libvirtxml.DomainInterfaceDriverGuest{
					CSum: "on",
				},
			},
		),
		Entry(
			"should work properly if everything is affected too",
			uint(16),
			uint64((1<<virtio.VIRTIO_NET_F_HOST_TSO4)|(1<<virtio.VIRTIO_NET_F_GUEST_ECN)|(1<<virtio.VIRTIO_NET_F_HASH_REPORT)),
			&libvirtxml.DomainInterfaceDriver{
				Name:          "vhost",
				Queues:        16,
				RSSHashReport: "on",
				Host: &libvirtxml.DomainInterfaceDriverHost{
					TSO4: "on",
				},
				Guest: &libvirtxml.DomainInterfaceDriverGuest{
					ECN: "on",
				},
			},
		),
	)
})
