package udp

import (
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

func BuildUDPHolepunch(index uint16) []byte {
	buf := make([]byte, 3)
	offset := 0
	WriteUint8(&buf, UdpPacketSignature, &offset)
	WriteUint16(&buf, index, &offset)
	return buf[:offset]
}
func UDPBuild(seq *uint8, isHost uint8, userid uint32, ip uint32, port uint16) []byte {
	rst := BuildHeader(seq, PacketTypeUdp)
	buf := make([]byte, 12)
	offset := 0
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, isHost, &offset)
	WriteUint32(&buf, userid, &offset)
	WriteUint32BE(&buf, ip, &offset)
	WriteUint16(&buf, port, &offset)
	rst = BytesCombine(rst, buf)
	return rst
}
