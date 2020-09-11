package playerinfo

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
)

const (
	SetSignature = 5
	SetTitle     = 6
	SetAvatar    = 7
)

func OnPlayerInfo(p *PacketData, client net.Conn) {
	var pkt InPlayerInfoPacket
	if p.PrasePlayerInfoPacket(&pkt) {
		switch pkt.InfoType {
		case SetSignature:
			OnSetSignature(p, client)
		case SetTitle:
			OnSetTitle(p, client)
		case SetAvatar:
			OnSetAvatar(p, client)
		default:
			log.Println("Unknown PlayerInfo packet", pkt.InfoType, "from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal PlayerInfo packet from", client.RemoteAddr().String())
	}
}
