package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostDeathPacket(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InDeathPacket
	if !p.PraseInDeathPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error HostDeath packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromID(pkt.DeadID)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		//log.Println("Error : Client from", client.RemoteAddr().String(), "sent HostDeath but not in server or is bot !")
		return
	}
	//修改玩家当前数据
	uPtr.CountDeadNum(pkt.DeathNum)

}
