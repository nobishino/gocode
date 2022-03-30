package arp

import (
	"syscall"

	"github.com/nobishino/gocode/network/ethernet"
)

// https://atmarkit.itmedia.co.jp/ait/articles/0305/09/news003_2.html

type HardWareType [2]byte
type Protocol [2]byte
type HLEN byte
type PLEN byte
type Opcode [2]byte
type IP [4]byte

// type SourceMacAddr ethernet.MacAddr

type ARP struct {
	HardwareType  HardWareType
	ProtocolType  Protocol
	HardwareSize  HLEN
	ProtocolSize  PLEN
	Opcode        Opcode
	SenderMacAddr ethernet.MacAddr
	SenderIPAddr  IP
	TargetMacAddr ethernet.MacAddr
	TargetIPAddr  IP
}

func NewRequest(senderMacAddr ethernet.MacAddr, senderIP, targetIP IP) ARP {
	return ARP{
		HardwareType:  [2]byte{0x00, 0x01},
		ProtocolType:  [2]byte{0x08, 0x00},
		HardwareSize:  0x06,
		ProtocolSize:  0x04,
		Opcode:        [2]byte{0x00, 0x01},
		SenderMacAddr: senderMacAddr,
		SenderIPAddr:  senderIP,
		// TargetMacAddr: targetMacAddr,
		TargetIPAddr: targetIP,
	}
}

func (*ARP) Send(ifindex int, packet []byte) ARP {

	addr := syscall.SockaddrLinklayer{}

}
