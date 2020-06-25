package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inRevivedPacket struct {
	userID uint32
	x      uint32 //待定，但是极像坐标
	y      uint32
	z      uint32
	unk00  uint8
}

func onHostRevivedPacket(p packet, client net.Conn) {
	//检索数据包
	var pkt inRevivedPacket
	if !praseInRevivedPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error HostKill packet !")
		return
	}
	if pkt.userID <= 0 {
		//log.Println("Bot revived at (", pkt.x, ",", pkt.y, ",", pkt.z, ")")
	} else {
		log.Println("UserID", pkt.userID, "revived at (", pkt.x, ",", pkt.y, ",", pkt.z, ")")
	}
}

func praseInRevivedPacket(p packet, dest *inRevivedPacket) bool {
	if p.datalen < 23 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).userID = ReadUint32(p.data, &offset)
	(*dest).x = ReadUint32(p.data, &offset)
	(*dest).y = ReadUint32(p.data, &offset)
	(*dest).z = ReadUint32(p.data, &offset)
	(*dest).unk00 = ReadUint8(p.data, &offset)
	return true
}
