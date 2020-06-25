package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inAssistPacket struct {
	killerID     uint32
	unk00        uint8  //可能是辅助击杀人数？
	AssisterID   uint32 //貌似是击杀方式？
	unk01        uint16
	unk02        uint16
	unk03        uint16
	AssisterTeam uint8 //待定,也可能是杀手的队伍
}

func onHostAssistPacket(p packet, client net.Conn) {
	//检索数据包
	var pkt inAssistPacket
	if !praseInAssistPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error HostKill packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromID(pkt.AssisterID)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		//log.Println("Error : Client from", client.RemoteAddr().String(), "sent HostKill but not in server or is bot !")
		return
	}
	//修改玩家当前数据
	uPtr.CountAssistNum()
}

func praseInAssistPacket(p packet, dest *inAssistPacket) bool {
	if p.datalen < 22 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).killerID = ReadUint32(p.data, &offset)
	(*dest).unk00 = ReadUint8(p.data, &offset)
	(*dest).AssisterID = ReadUint32(p.data, &offset)
	(*dest).unk01 = ReadUint16(p.data, &offset)
	(*dest).unk02 = ReadUint16(p.data, &offset)
	(*dest).unk03 = ReadUint16(p.data, &offset)
	(*dest).AssisterTeam = ReadUint8(p.data, &offset)
	return true
}
