package option

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
)

func OnSaveBuyMenu(p *PacketData, client net.Conn) {
	var pkt InOptionBuyMenu
	if !p.PraseSaveBuyMenu(&pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal save buymenu packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to save buymenu but not in server !")
		return
	}
	//修改数据
	uPtr.SetBuyMenu(pkt.Buymenu)
	log.Println("User", string(uPtr.Username), "updated buymenu ...")
}
