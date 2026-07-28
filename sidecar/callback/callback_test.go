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

package callback_test

import (
	"fmt"

	"kubevirt.io/vdpa-network-binding-plugin/sidecar/callback"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"libvirt.org/go/libvirtxml"
)

var _ = Describe("vdpa hook callback handler", func() {
	Context("on define domain", func() {
		It("should fail given empty byte slice stream", func() {
			_, err := callback.OnDefineDomain([]byte{}, noopMutatorStub(nil))
			Expect(err).To(HaveOccurred())
		})

		It("should fail given invalid domain XML", func() {
			_, err := callback.OnDefineDomain([]byte("invalid-domain-xml"), noopMutatorStub(nil))
			Expect(err).To(HaveOccurred())
		})

		It("should fail when domain spec mutator fails", func() {
			domain := &libvirtxml.Domain{Name: "test-failure"}
			domainXML, err := domain.Marshal()
			Expect(err).ToNot(HaveOccurred())

			expectedErr := fmt.Errorf("test error")
			domSpecMutator := noopMutatorStub(expectedErr)

			_, err = callback.OnDefineDomain([]byte(domainXML), domSpecMutator)
			Expect(err).To(Equal(expectedErr))
		})

		It("given no-op mutator, domain spec should not change", func() {
			domain := &libvirtxml.Domain{Name: "test-noop"}
			domainXML, err := domain.Marshal()
			Expect(err).ToNot(HaveOccurred())

			domSpecMutator := noopMutatorStub(nil)

			res, err := callback.OnDefineDomain([]byte(domainXML), domSpecMutator)
			Expect(err).To(BeNil())
			Expect(string(res)).To(Equal(domainXML))
		})

		It("domain spec should mutate successfully", func() {
			domain := &libvirtxml.Domain{Name: "test-mutate"}
			domainXML, err := domain.Marshal()
			Expect(err).ToNot(HaveOccurred())

			newInterfaceName := "new-interface"

			mutator := func(dom *libvirtxml.Domain) *libvirtxml.Domain {
				if dom.Devices == nil {
					dom.Devices = &libvirtxml.DomainDeviceList{}
				}
				dom.Devices.Interfaces = append(
					dom.Devices.Interfaces,
					libvirtxml.DomainInterface{
						Alias: &libvirtxml.DomainAlias{Name: newInterfaceName},
					},
				)
				return dom
			}
			domSpecMutator := mutatorStub{mutator: mutator}

			domain.Devices = &libvirtxml.DomainDeviceList{}
			domain.Devices.Interfaces = append(
				domain.Devices.Interfaces,
				libvirtxml.DomainInterface{
					Alias: &libvirtxml.DomainAlias{Name: newInterfaceName},
				},
			)
			mutatedDomainSpecXML, err := domain.Marshal()
			Expect(err).ToNot(HaveOccurred())

			res, err := callback.OnDefineDomain([]byte(domainXML), domSpecMutator)
			Expect(err).To(BeNil())
			Expect(string(res)).To(Equal(mutatedDomainSpecXML))
		})
	})
})

type mutatorStub struct {
	mutator    func(dom *libvirtxml.Domain) *libvirtxml.Domain
	failMutate error
}

func (s mutatorStub) Mutate(dom *libvirtxml.Domain) (*libvirtxml.Domain, error) {
	return s.mutator(dom), s.failMutate
}

func noopMutatorStub(failMutate error) *mutatorStub {
	return &mutatorStub{
		mutator: func(dom *libvirtxml.Domain) *libvirtxml.Domain {
			return dom
		},
		failMutate: failMutate,
	}
}
