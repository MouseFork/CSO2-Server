package option

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnSaveBuyMenu(p *PacketData, client net.Conn) {
	var pkt InOptionBuyMenu
	if !p.PraseSaveBuyMenu(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a illegal save buymenu packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to save buymenu but not in server !")
		return
	}
	//修改数据
	uPtr.SetBuyMenu(pkt.Buymenu)
	DebugInfo(1, "User", string(uPtr.Username), "updated buymenu ...")
}
