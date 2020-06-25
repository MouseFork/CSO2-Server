package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inQuickList struct {
	//unk00     uint8
	gameModID   uint8
	IsEnableBot uint8
}

func onQuickList(p packet, client net.Conn) {
	//检索数据包
	var pkt inQuickList
	if !praseInQuickListPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error QuickList packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to request QuickList but not in server !")
		return
	}
	//发送房间数据,暂时发送空数据
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildQuickList(pkt))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("Sent a null quickList to User", string(uPtr.username))
}

func praseInQuickListPacket(p packet, dest *inQuickList) bool {
	if p.datalen < 8 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).gameModID = ReadUint8(p.data, &offset)
	(*dest).IsEnableBot = ReadUint8(p.data, &offset)
	return true
}

func BuildQuickList(pkt inQuickList) []byte {
	buf := make([]byte, 2)
	offset := 0
	WriteUint8(&buf, QuickList, &offset)
	WriteUint8(&buf, 0, &offset) //num of room

	return buf[:offset]
}
