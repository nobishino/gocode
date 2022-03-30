package ethernet

type MacAddr [6]byte

type Frame struct {
	Destination MacAddr
	Source      MacAddr
	Type        [2]byte
}

type FrameType [2]byte

var (
	IPv4 = FrameType{0x08, 0x00}
	ARP  = FrameType{0x08, 0x06}
	// ICMP = FrameType{}
)
