package typestruct

import (
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

const (
	UdpPacketSignature = 87
	UDPTypeClient      = 0
	UDPTypeServer      = 256
	UDPTypeSourceTV    = 512
)

type InUDPmsg struct {
	Signature uint8
	UserId    uint32
	PortId    uint16
	IpAddress uint32
	Port      uint16

	PacketData         []byte
	Datalen            int
	CurOffset          int //可能32位
	ParsedSuccessfully bool
}

func (p InUDPmsg) IsHeartbeat() bool {
	return p.Datalen == 6
}

func (dest *InUDPmsg) PraseUDPpacket(data []byte, len int) bool {
	dest.CurOffset = 0
	dest.Signature = ReadUint8(data, &dest.CurOffset)
	if dest.Signature != UdpPacketSignature {
		dest.ParsedSuccessfully = false
		return false
	}
	dest.Datalen = len
	dest.PacketData = data
	if dest.IsHeartbeat() {
	} else {
		dest.UserId = ReadUint32(data, &dest.CurOffset)
		dest.PortId = ReadUint16(data, &dest.CurOffset)
		dest.IpAddress = ReadUint32BE(data, &dest.CurOffset)
		dest.Port = ReadUint16(data, &dest.CurOffset)
	}
	dest.ParsedSuccessfully = true
	return true
}
