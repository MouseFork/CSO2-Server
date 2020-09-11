package host

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
)

func OnHostAssistPacket(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InAssistPacket
	if !p.PraseInAssistPacket(&pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error HostKill packet !")
		return
	}
	//log.Println(p.data)
	//找到对应用户
	uPtr := GetUserFromID(pkt.AssisterID)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		//log.Println("Error : Client from", client.RemoteAddr().String(), "sent HostKill but not in server or is bot !")
		return
	}
	//修改玩家当前数据
	uPtr.CountAssistNum()
	//log.Println("User", string(uPtr.username), "assisted")
}
