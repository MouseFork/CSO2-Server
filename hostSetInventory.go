package main

import (
	"log"
	"net"
)

type inHostSetInventoryPacket struct {
	userID uint32
}

func onHostSetUserInventory(p packet, client net.Conn) {
	//检查数据包
	var pkt inHostSetInventoryPacket
	if !praseSetUserInventoryPacket(p, &pkt) {
		log.Println("Error : Cannot prase a host setUserInventory packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : A user request to setUserInventory but not in server!")
		return
	}
	//找到玩家的房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User:", string(uPtr.username), "try to setUserInventory but in a null room !")
		return
	}
	//是不是房主
	if rm.hostUserID != uPtr.userid {
		log.Println("Error : User:", string(uPtr.username), "try to setUserInventory but isn't host !")
		return
	}
	//发送用户的装备给房主

}

func praseSetUserInventoryPacket(p packet, dest *inHostSetInventoryPacket) bool {
	if dest == nil ||
		p.length < 6 {
		return false
	}
	offset := 6
	(*dest).userID = ReadUint32(p.data, &offset)
	return true
}

func BuildSetUserInventory(u user) {
	//buf := make([]byte,8+)
}
