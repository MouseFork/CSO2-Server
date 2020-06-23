package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type InFavoriteSetCosmetics struct {
	slot   uint8
	itemId uint32
}

func onFavoriteSetCosmetics(seq *uint8, p packet, client net.Conn) {
	//检索数据包
	var pkt InFavoriteSetCosmetics
	if !praseFavoriteSetCosmeticsPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error SetCosmetics packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to SetCosmetics but not in server !")
		return
	}
	//设置武器
	setCosmetic(pkt.slot, pkt.itemId, uPtr)
	log.Println("Setting User", string(uPtr.username), "new Cosmetic", pkt.itemId, "to slot", pkt.slot)
	//找到对应房间玩家
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.currentRoomId)
	if rm == nil ||
		rm.id <= 0 {
		return
	}
	u := rm.roomGetUser(uPtr.userid)
	if u == nil {
		return
	}
	setCosmetic(pkt.slot, pkt.itemId, u)
}
func praseFavoriteSetCosmeticsPacket(p packet, dest *InFavoriteSetCosmetics) bool {
	if p.datalen < 11 {
		return false
	}
	offset := 6
	(*dest).slot = ReadUint8(p.data, &offset)
	(*dest).itemId = ReadUint32(p.data, &offset)
	return true
}
func setCosmetic(slot uint8, itemId uint32, u *user) {
	switch slot {
	case 0:
		(*u).inventory.CTModel = itemId
	case 1:
		(*u).inventory.TModel = itemId
	case 2:
		(*u).inventory.headItem = itemId
	case 3:
		(*u).inventory.gloveItem = itemId
	case 4:
		(*u).inventory.backItem = itemId
	case 5:
		(*u).inventory.stepsItem = itemId
	case 6:
		(*u).inventory.cardItem = itemId
	case 7:
		(*u).inventory.sprayItem = itemId
	default:
		log.Println("Error : User", string(u.username), "try to setCosmetic invalid slot", slot)
		return
	}
}
func BuildCosmetics(inventory userInventory) []byte {
	buf := make([]byte, 2)
	offset := 0
	curItem := uint8(0)
	WriteUint8(&buf, FavoriteSetCosmetics, &offset)
	WriteUint8(&buf, 10, &offset)
	temp := WriteItem(inventory.CTModel, &curItem)
	temp = BytesCombine(temp, WriteItem(inventory.TModel, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.headItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.gloveItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.backItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.stepsItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.cardItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.sprayItem, &curItem))
	temp = BytesCombine(temp, WriteItem(0, &curItem))
	temp = BytesCombine(temp, WriteItem(0, &curItem))
	buf = BytesCombine(buf[:offset], temp)
	return buf
}
