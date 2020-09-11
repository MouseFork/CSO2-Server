package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostKillPacket(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InKillPacket
	if !p.PraseInKillPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error HostKill packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromID(pkt.KillerID)
	if uPtr != nil &&
		uPtr.Userid > 0 {
		//log.Println("Error : Client from", client.RemoteAddr().String(), "sent HostKill but not in server or is bot !")
		//return
		//修改玩家当前数据
		uPtr.CountKillNum(pkt.KillNum)
	}
	//修改房间数据
	uPtr = GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		return
	}
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id < 0 {
		return
	}
	if pkt.PlayerTeam == UserForceCounterTerrorist {
		rm.CountRoomCtKill()
	} else {
		rm.CountRoomTrKill()
	}

}
