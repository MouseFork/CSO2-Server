package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inHostSetLoadoutPacket struct {
	userID uint32
}

func onHostSetUserLoadout(p packet, client net.Conn) {
	//检查数据包
	var pkt inHostSetLoadoutPacket
	if !praseSetUserLoadoutPacket(p, &pkt) {
		log.Println("Error : Cannot prase a send UserLoadout packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : A user request to send UserLoadout but not in server!")
		return
	}
	dest := getUserFromID(pkt.userID)
	if dest == nil ||
		dest.userid <= 0 {
		log.Println("Error : A user request to send UserLoadout but dest user is null!")
		return
	}
	//找到玩家的房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User", string(uPtr.username), "try to send UserLoadout but in a null room !")
		return
	}
	destRm := getRoomFromID(dest.getUserChannelServerID(),
		dest.getUserChannelID(),
		dest.getUserRoomID())
	if destRm == nil ||
		destRm.id <= 0 {
		log.Println("Error : User", string(dest.username), "try to send UserLoadout but in a null room !")
		return
	}
	if rm.id != destRm.id {
		log.Println("Error : User", string(dest.username), "try to send UserLoadout to", string(dest.username), "but not a same room !")
		return
	}
	//是不是房主
	if rm.hostUserID != uPtr.userid {
		log.Println("Error : User", string(uPtr.username), "try to send UserLoadout but isn't host !")
		return
	}
	//发送用户背包数据
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildSetUserLoadout(*dest))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("Send User", string(dest.username), "Loadout to host", string(uPtr.username))
}

func praseSetUserLoadoutPacket(p packet, dest *inHostSetLoadoutPacket) bool {
	if dest == nil ||
		p.length < 6 {
		return false
	}
	offset := 6
	(*dest).userID = ReadUint32(p.data, &offset)
	return true
}

func BuildSetUserLoadout(u user) []byte {
	buf := make([]byte, 6)
	offset := 0
	WriteUint8(&buf, SetLoadout, &offset)
	WriteUint32(&buf, u.userid, &offset)
	WriteUint8(&buf, 8, &offset) //类型数量
	//当前8个类型的装备
	curItem := uint8(0)
	temp := WriteItem(u.inventory.CTModel, &curItem)
	temp = BytesCombine(temp, WriteItem(u.inventory.TModel, &curItem))
	temp = BytesCombine(temp, WriteItem(u.inventory.headItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.inventory.gloveItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.inventory.backItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.inventory.stepsItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.inventory.cardItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.inventory.sprayItem, &curItem))
	buf = BytesCombine(buf[:offset], temp)
	buf = append(buf, uint8(len(u.inventory.loadouts)))
	for _, v := range u.inventory.loadouts {
		buf = append(buf, uint8(len(v.items)))
		curItem = 0
		for _, j := range v.items {
			buf = BytesCombine(buf, WriteItem(j, &curItem))
		}
	}
	buf = append(buf, 0)
	return buf
}
