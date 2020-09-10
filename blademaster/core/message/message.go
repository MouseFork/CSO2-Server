package message

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

func OnSendMessage(seq *uint8, client net.Conn, tp uint8, msg []byte) {
	rst := BuildHeader(seq, PacketTypeChat)
	rst = append(rst, tp)
	rst = BytesCombine(rst, BuildMessage(msg, tp))
	SendPacket(rst, client)
}

func BuildMessage(msg []byte, tp uint8) []byte {
	if tp == Congratulate {
		buf := make([]byte, 1)
		buf[0] = 0
		return BytesCombine(buf, BuildString(msg))
	}
	return BuildLongString(msg)
}
