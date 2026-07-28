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
	"strconv"

	hwutil "kubevirt.io/kubevirt/pkg/util/hardware"

	libvirtxml "libvirt.org/go/libvirtxml"
)

func parsePCIField(field string) (uint, error) {
	val, err := strconv.ParseUint(field, 10, 0)
	return uint(val), err
}

func parsePCIAddress(address string) (*libvirtxml.DomainAddress, error) {
	res := libvirtxml.DomainAddress{}

	pciFields, err := hwutil.ParsePciAddress(address)
	if err != nil {
		return nil, err
	}

	domain, err := parsePCIField(pciFields[0])
	if err != nil {
		return nil, err
	}
	bus, err := parsePCIField(pciFields[1])
	if err != nil {
		return nil, err
	}
	slot, err := parsePCIField(pciFields[2])
	if err != nil {
		return nil, err
	}
	function, err := parsePCIField(pciFields[3])
	if err != nil {
		return nil, err
	}

	res.PCI = &libvirtxml.DomainAddressPCI{
		Domain:   &domain,
		Bus:      &bus,
		Slot:     &slot,
		Function: &function,
	}

	return &res, nil
}
