package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inHostSetInventoryPacket struct {
	userID uint32
}

//onHostSetUserInventory 用户发来请求将自己的装备信息发给指定user
func onHostSetUserInventory(p packet, client net.Conn) {
	//检查数据包
	var pkt inHostSetInventoryPacket
	if !praseSetUserInventoryPacket(p, &pkt) {
		log.Println("Error : Cannot prase a send UserInventory packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : A user request to send UserInventory but not in server!")
		return
	}
	dest := getUserFromID(pkt.userID)
	if dest == nil ||
		dest.userid <= 0 {
		log.Println("Error : A user request to send UserInventory but dest user is null!")
		return
	}
	//找到玩家的房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User", string(uPtr.username), "try to send UserInventory but in a null room !")
		return
	}
	destRm := getRoomFromID(dest.getUserChannelServerID(),
		dest.getUserChannelID(),
		dest.getUserRoomID())
	if destRm == nil ||
		destRm.id <= 0 {
		log.Println("Error : User", string(dest.username), "try to send UserInventory but in a null room !")
		return
	}
	if rm.id != destRm.id {
		log.Println("Error : User", string(dest.username), "try to send UserInventory to", string(dest.username), "but not a same room !")
		return
	}
	//是不是房主
	if rm.hostUserID != uPtr.userid {
		log.Println("Error : User", string(uPtr.username), "try to send UserInventory but isn't host !")
		return
	}
	//发送用户的装备给目标user
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildSetUserInventory(*dest, dest.userid))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("Send User", string(dest.username), "Inventory to host", string(uPtr.username))
	rst = BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildSetUserLoadout(*dest))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("Send User", string(dest.username), "Loadout to host", string(uPtr.username))
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

//BuildSetUserInventory 建立要发给主机的玩家装备信息，按理来说应该是所有玩家的装备，待定，L-Leite是发的主机的装备加普通用户ID
func BuildSetUserInventory(u user, destid uint32) []byte {
	buf := make([]byte, 10+6*u.inventory.numOfItem)
	offset := 0
	WriteUint8(&buf, SetInventory, &offset)
	WriteUint32(&buf, destid, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint16(&buf, u.inventory.numOfItem, &offset)
	for _, v := range u.inventory.items {
		WriteUint32(&buf, v.id, &offset)
		WriteUint16(&buf, v.count, &offset)
	}
	return buf[:offset]
}
