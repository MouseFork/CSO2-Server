package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostRevivedPacket(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InRevivedPacket
	if !p.PraseInRevivedPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error HostKill packet !")
		return
	}
	if pkt.UserID <= 0 {
		//log.Println("Bot revived at (", pkt.x, ",", pkt.y, ",", pkt.z, ")")
	} else {
		DebugInfo(2, "UserID", pkt.UserID, "revived at (", pkt.X, ",", pkt.Y, ",", pkt.Z, ")")
	}
}
