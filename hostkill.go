package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inKillPacket struct {
	unk00      uint8 //一直是0
	killerID   uint32
	unk01      uint32 //一直是0
	unk02      uint8
	killType   uint8  //貌似是击杀方式？
	killNum    uint16 //杀敌数,生化模式3倍
	playerTeam uint8  //待定
}

func onHostKillPacket(p packet, client net.Conn) {
	//检索数据包
	var pkt inKillPacket
	if !praseInKillPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error HostKill packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromID(pkt.killerID)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		//log.Println("Error : Client from", client.RemoteAddr().String(), "sent HostKill but not in server or is bot !")
		return
	}
	//修改玩家当前数据
	uPtr.CountKillNum(pkt.killNum)
	//修改房间数据
	uPtr = getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		return
	}
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		return
	}
	if pkt.playerTeam == CounterTerrorist {
		rm.CountRoomCtKill()
	} else {
		rm.CountRoomTrKill()
	}

}

func praseInKillPacket(p packet, dest *inKillPacket) bool {
	if p.datalen < 20 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).unk00 = ReadUint8(p.data, &offset)
	(*dest).killerID = ReadUint32(p.data, &offset)
	(*dest).unk01 = ReadUint32(p.data, &offset)
	(*dest).unk02 = ReadUint8(p.data, &offset)
	(*dest).killType = ReadUint8(p.data, &offset)
	(*dest).killNum = ReadUint16(p.data, &offset)
	(*dest).playerTeam = ReadUint8(p.data, &offset)
	return true
}
