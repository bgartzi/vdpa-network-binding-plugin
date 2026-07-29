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

package virtio

const (
	VIRTIO_NET_F_CSUM                uint8 = 0  // Host handles pkts w/ partial csum
	VIRTIO_NET_F_GUEST_CSUM          uint8 = 1  // Guest handles pkts w/ partial csum
	VIRTIO_NET_F_CTRL_GUEST_OFFLOADS uint8 = 2  // Dynamic offload configuration.
	VIRTIO_NET_F_MTU                 uint8 = 3  // Initial MTU advice
	VIRTIO_NET_F_MAC                 uint8 = 5  // Host has given MAC address.
	VIRTIO_NET_F_GUEST_TSO4          uint8 = 7  // Guest can handle TSOv4 in.
	VIRTIO_NET_F_GUEST_TSO6          uint8 = 8  // Guest can handle TSOv6 in.
	VIRTIO_NET_F_GUEST_ECN           uint8 = 9  // Guest can handle TSO[6] w/ ECN in.
	VIRTIO_NET_F_GUEST_UFO           uint8 = 10 // Guest can handle UFO in.
	VIRTIO_NET_F_HOST_TSO4           uint8 = 11 // Host can handle TSOv4 in.
	VIRTIO_NET_F_HOST_TSO6           uint8 = 12 // Host can handle TSOv6 in.
	VIRTIO_NET_F_HOST_ECN            uint8 = 13 // Host can handle TSO[6] w/ ECN in.
	VIRTIO_NET_F_HOST_UFO            uint8 = 14 // Host can handle UFO in.
	VIRTIO_NET_F_MRG_RXBUF           uint8 = 15 // Host can merge receive buffers.
	VIRTIO_NET_F_STATUS              uint8 = 16 // virtio_net_config.status available
	VIRTIO_NET_F_CTRL_VQ             uint8 = 17 // Control channel available
	VIRTIO_NET_F_CTRL_RX             uint8 = 18 // Control channel RX mode support
	VIRTIO_NET_F_CTRL_VLAN           uint8 = 19 // Control channel VLAN filtering
	VIRTIO_NET_F_CTRL_RX_EXTRA       uint8 = 20 // Extra RX mode control support
	VIRTIO_NET_F_GUEST_ANNOUNCE      uint8 = 21 // Guest can announce device on the network
	VIRTIO_NET_F_MQ                  uint8 = 22 // Device supports Receive Flow Steering
	VIRTIO_NET_F_CTRL_MAC_ADDR       uint8 = 23 // Set MAC address
	VIRTIO_NET_F_DEVICE_STATS        uint8 = 50 // Device can provide device-level statistics.
	VIRTIO_NET_F_VQ_NOTF_COAL        uint8 = 52 // Device supports virtqueue notification coalescing
	VIRTIO_NET_F_NOTF_COAL           uint8 = 53 // Device supports notifications coalescing
	VIRTIO_NET_F_GUEST_USO4          uint8 = 54 // Guest can handle USOv4 in.
	VIRTIO_NET_F_GUEST_USO6          uint8 = 55 // Guest can handle USOv6 in.
	VIRTIO_NET_F_HOST_USO            uint8 = 56 // Host can handle USO in.
	VIRTIO_NET_F_HASH_REPORT         uint8 = 57 // Supports hash report
	VIRTIO_NET_F_GUEST_HDRLEN        uint8 = 59 // Guest provides the exact hdr_len value.
	VIRTIO_NET_F_RSS                 uint8 = 60 // Supports RSS RX steering
	VIRTIO_NET_F_RSC_EXT             uint8 = 61 // extended coalescing info
	VIRTIO_NET_F_STANDBY             uint8 = 62 // Act as standby for another device with the same MAC.
	VIRTIO_NET_F_SPEED_DUPLEX        uint8 = 63 // Device set linkspeed and duplex
)

func Contains(features uint64, feature uint8) bool {
	return features&(1<<feature) > 0
}
