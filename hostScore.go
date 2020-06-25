package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inGameScorePacket struct {
	winnerTeam uint8
	TrScore    uint8
	CtScore    uint8
	PacketType uint8 //maybe
	hostID     uint32
	unk00      uint32
}

func onHostGameScorePacket(p packet, client net.Conn) {
	//检索数据包
	var pkt inGameScorePacket
	if !praseInGameScorePacket(p, &pkt) {
		log.Println("Error : User from", client.RemoteAddr().String(), "sent a error GameScore packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : User from", client.RemoteAddr().String(), "sent GameScore but not in server !")
		return
	}
	//找到对应房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.currentRoomId)
	if rm == nil ||
		rm.id <= 0 {
		return
	}
	//修改房间结果
	//if pkt.PacketType != 0 {
	rm.setRoomScore(pkt.CtScore, pkt.TrScore)
	rm.setRoomWinner(pkt.winnerTeam)
	//}
}

func praseInGameScorePacket(p packet, dest *inGameScorePacket) bool {
	if p.datalen < 10 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).winnerTeam = ReadUint8(p.data, &offset)
	(*dest).TrScore = ReadUint8(p.data, &offset)
	(*dest).CtScore = ReadUint8(p.data, &offset)
	(*dest).PacketType = ReadUint8(p.data, &offset)
	if (*dest).PacketType != 0 {
		(*dest).hostID = ReadUint32(p.data, &offset)
		(*dest).unk00 = ReadUint32(p.data, &offset)
	}
	return true
}
