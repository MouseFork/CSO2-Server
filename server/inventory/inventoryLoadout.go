package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type userLoadout struct {
	items []uint32
}

type inFavoriteSetLoadout struct {
	loadout    uint8
	weaponSlot uint8
	itemId     uint32
}

func onFavoriteSetLoadout(seq *uint8, p packet, client net.Conn) {
	//检索数据包
	var pkt inFavoriteSetLoadout
	if !praseFavoriteSetLoadoutPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error SetLoadout packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to SetLoadout but not in server !")
		return
	}
	//设置武器
	if pkt.loadout > 2 ||
		pkt.weaponSlot > 6 {
		log.Println("Error : User", string(uPtr.username), "try to SetLoadout with invalid data !")
		return
	}
	(*uPtr).inventory.loadouts[pkt.loadout].items[pkt.weaponSlot] = pkt.itemId
	log.Println("Setting User", string(uPtr.username), "new weapon", pkt.itemId, "to slot", pkt.weaponSlot, "in loadout", pkt.loadout)
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
	(*u).inventory.loadouts[pkt.loadout].items[pkt.weaponSlot] = pkt.itemId
}
func praseFavoriteSetLoadoutPacket(p packet, dest *inFavoriteSetLoadout) bool {
	if p.datalen < 12 {
		return false
	}
	offset := 6
	(*dest).loadout = ReadUint8(p.data, &offset)
	(*dest).weaponSlot = ReadUint8(p.data, &offset)
	(*dest).itemId = ReadUint32(p.data, &offset)
	return true
}
func BuildLoadout(inventory userInventory) []byte {
	buf := make([]byte, 5+len(inventory.loadouts)*96)
	offset := 0
	WriteUint8(&buf, FavoriteSetLoadout, &offset)
	WriteUint8(&buf, uint8(len(inventory.loadouts))*16, &offset)
	for i, v := range inventory.loadouts {
		for j, k := range v.items {
			WriteUint8(&buf, uint8(i), &offset)
			WriteUint8(&buf, uint8(j), &offset)
			WriteUint32(&buf, k, &offset)
		}
	}
	return buf[:offset]
}

func createNewLoadout() []userLoadout {
	return []userLoadout{
		{[]uint32{5336, 5356, 5330, 4, 23, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[]uint32{5285, 5294, 5231, 4, 23, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[]uint32{5206, 5356, 5365, 4, 23, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}
}
