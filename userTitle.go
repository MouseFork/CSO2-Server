package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inSetTitlePacket struct {
	TitleId uint16
}

func onSetTitle(p packet, client net.Conn) {
	var pkt inSetTitlePacket
	if !praseSetTitlePacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal SetTitle packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to SetTitle but not in server !")
		return
	}
	//修改数据
	uPtr.SetTitle(pkt.TitleId)
	//发送数据包
	p.id = TypeUserInfo
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildSetTitlePacket(uPtr.userid, pkt.TitleId))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("User", string(uPtr.username), "Set Title to", pkt.TitleId)
	//如果是在房间内
}

func praseSetTitlePacket(p packet, dest *inSetTitlePacket) bool {
	if p.datalen-HeaderLen < 4 {
		return false
	}
	offset := 6
	(*dest).TitleId = ReadUint16(p.data, &offset)
	return true
}

func BuildSetTitlePacket(id uint32, Title uint16) []byte {
	buf := make([]byte, 10)
	offset := 0
	WriteUint32(&buf, id, &offset)
	WriteUint32(&buf, 0x8000, &offset)
	WriteUint16(&buf, Title, &offset)
	return buf[:offset]
}
