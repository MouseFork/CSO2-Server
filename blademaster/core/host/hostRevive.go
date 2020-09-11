package host

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
)

func OnHostRevivedPacket(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InRevivedPacket
	if !p.PraseInRevivedPacket(&pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error HostKill packet !")
		return
	}
	if pkt.UserID <= 0 {
		//log.Println("Bot revived at (", pkt.x, ",", pkt.y, ",", pkt.z, ")")
	} else {
		log.Println("UserID", pkt.UserID, "revived at (", pkt.X, ",", pkt.Y, ",", pkt.Z, ")")
	}
}
