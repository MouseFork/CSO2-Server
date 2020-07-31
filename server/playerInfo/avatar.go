package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inSetAvatarPacket struct {
	avatarId uint16
}

func onSetAvatar(p packet, client net.Conn) {
	var pkt inSetAvatarPacket
	if !praseSetAvatarPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal SetAvatar packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to SetAvatar but not in server !")
		return
	}
	//修改数据
	uPtr.SetAvatar(pkt.avatarId)
	//发送数据包
	p.id = TypeUserInfo
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildSetAvatarPacket(uPtr.userid, pkt.avatarId))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("User", string(uPtr.username), "Set Avatar to", pkt.avatarId)
	//如果是在房间内
}

func praseSetAvatarPacket(p packet, dest *inSetAvatarPacket) bool {
	if p.datalen-HeaderLen < 4 {
		return false
	}
	offset := 6
	(*dest).avatarId = ReadUint16(p.data, &offset)
	return true
}

func BuildSetAvatarPacket(id uint32, avatar uint16) []byte {
	buf := make([]byte, 10)
	offset := 0
	WriteUint32(&buf, id, &offset)
	WriteUint32(&buf, 0x800000, &offset)
	WriteUint16(&buf, avatar, &offset)
	return buf[:offset]
}
