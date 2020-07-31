package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inDeathPacket struct {
	deadID     uint32
	unk00      uint32 //一直是0
	unk01      uint8  //貌似是死亡方式？
	deathNum   uint16 //死亡数,生化模式3倍
	playerTeam uint8  //待定
}

func onHostDeathPacket(p packet, client net.Conn) {
	//检索数据包
	var pkt inDeathPacket
	if !praseInDeathPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error HostDeath packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromID(pkt.deadID)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		//log.Println("Error : Client from", client.RemoteAddr().String(), "sent HostDeath but not in server or is bot !")
		return
	}
	//修改玩家当前数据
	uPtr.CountDeadNum(pkt.deathNum)
}

func praseInDeathPacket(p packet, dest *inDeathPacket) bool {
	if p.datalen < 18 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).deadID = ReadUint32(p.data, &offset)
	(*dest).unk00 = ReadUint32(p.data, &offset)
	(*dest).unk01 = ReadUint8(p.data, &offset)
	(*dest).deathNum = ReadUint16(p.data, &offset)
	(*dest).playerTeam = ReadUint8(p.data, &offset)
	return true
}
