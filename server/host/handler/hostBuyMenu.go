package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inHostSetBuyMenu struct {
	userid uint32
}

func onHostSetUserBuyMenu(p packet, client net.Conn) {
	//检查数据包
	var pkt inHostSetBuyMenu
	if !praseSetBuyMenuPacket(p, &pkt) {
		log.Println("Error : Cannot prase a send BuyMenu packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : A user request to send BuyMenu but not in server!")
		return
	}
	dest := getUserFromID(pkt.userid)
	if dest == nil ||
		dest.userid <= 0 {
		log.Println("Error : A user request to send BuyMenu but dest user is null!")
		return
	}
	//找到玩家的房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User", string(uPtr.username), "try to send BuyMenu but in a null room !")
		return
	}
	destRm := getRoomFromID(dest.getUserChannelServerID(),
		dest.getUserChannelID(),
		dest.getUserRoomID())
	if destRm == nil ||
		destRm.id <= 0 {
		log.Println("Error : User", string(dest.username), "try to send BuyMenu but in a null room !")
		return
	}
	if rm.id != destRm.id {
		log.Println("Error : User", string(dest.username), "try to send BuyMenu to", string(dest.username), "but not a same room !")
		return
	}
	//是不是房主
	if rm.hostUserID != uPtr.userid {
		log.Println("Error : User", string(uPtr.username), "try to send BuyMenu but isn't host !")
		return
	}
	//发送数据包
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildSetBuyMenu(dest.userid, dest.inventory))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("Send User", string(dest.username), "BuyMenu to host", string(uPtr.username))
}

func praseSetBuyMenuPacket(p packet, dest *inHostSetBuyMenu) bool {
	if dest == nil ||
		p.length < 6 {
		return false
	}
	offset := 6
	(*dest).userid = ReadUint32(p.data, &offset)
	return true
}
func BuildSetBuyMenu(id uint32, inventory userInventory) []byte {
	l := 6 * (len(inventory.buyMenu.pistols) +
		len(inventory.buyMenu.shotguns) +
		len(inventory.buyMenu.smgs) +
		len(inventory.buyMenu.rifles) +
		len(inventory.buyMenu.snipers) +
		len(inventory.buyMenu.machineguns) +
		len(inventory.buyMenu.melees) +
		len(inventory.buyMenu.equipment))
	buf := make([]byte, 8+l)
	offset := 0
	WriteUint8(&buf, SetBuyMenu, &offset)
	WriteUint32(&buf, id, &offset)
	WriteUint16(&buf, 369, &offset) //buyMenuByteLength
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.pistols)), &offset)
	for k, v := range inventory.buyMenu.pistols {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}
	WriteUint8(&buf, uint8(len(inventory.buyMenu.shotguns)), &offset)
	for k, v := range inventory.buyMenu.shotguns {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.smgs)), &offset)
	for k, v := range inventory.buyMenu.smgs {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.rifles)), &offset)
	for k, v := range inventory.buyMenu.rifles {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.snipers)), &offset)
	for k, v := range inventory.buyMenu.snipers {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.machineguns)), &offset)
	for k, v := range inventory.buyMenu.machineguns {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.melees)), &offset)
	for k, v := range inventory.buyMenu.melees {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.equipment)), &offset)
	for k, v := range inventory.buyMenu.equipment {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}
	return buf[:offset]
}
