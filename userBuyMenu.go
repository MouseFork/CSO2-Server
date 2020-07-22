package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type InOptionBuyMenu struct {
	menuLength uint16
	unk00      uint8
	buymenu    userBuyMenu
}

const (
	SUBMENU_ITEM_NUM = 9
)

func onSaveBuyMenu(p packet, client net.Conn) {
	var pkt InOptionBuyMenu
	if !praseSaveBuyMenu(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal save buymenu packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to save buymenu but not in server !")
		return
	}
	//修改数据
	(*uPtr).inventory.buyMenu = pkt.buymenu
	log.Println("User", string(uPtr.username), "updated buymenu ...")
}

func praseSaveBuyMenu(p packet, dest *InOptionBuyMenu) bool {
	if p.datalen < 7 {
		return false
	}
	offset := 6
	(*dest).menuLength = ReadUint16(p.data, &offset)
	(*dest).unk00 = ReadUint8(p.data, &offset)
	(*dest).buymenu.pistols = ReadSubMenu(p.data, &offset)
	(*dest).buymenu.shotguns = ReadSubMenu(p.data, &offset)
	(*dest).buymenu.smgs = ReadSubMenu(p.data, &offset)
	(*dest).buymenu.rifles = ReadSubMenu(p.data, &offset)
	(*dest).buymenu.snipers = ReadSubMenu(p.data, &offset)
	(*dest).buymenu.machineguns = ReadSubMenu(p.data, &offset)
	(*dest).buymenu.melees = ReadSubMenu(p.data, &offset)
	(*dest).buymenu.equipment = ReadSubMenu(p.data, &offset)
	return true
}

func ReadSubMenu(b []byte, offset *int) []uint32 {
	len := ReadUint8(b, offset)
	if len != SUBMENU_ITEM_NUM {
		log.Println("Length of submenu is illegal !")
	}
	var submenu []uint32
	for i := 0; i < SUBMENU_ITEM_NUM; i++ {
		ReadUint8(b, offset)
		submenu = append(submenu, ReadUint32(b, offset))
	}
	return submenu
}
