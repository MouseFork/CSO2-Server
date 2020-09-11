package option

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
)

const (
	SaveBuyMenu = 1
)

func OnOption(p *PacketData, client net.Conn) {
	var pkt InOptionPacket
	if p.PraseOptionPacket(&pkt) {
		switch pkt.OptionPacketType {
		case SaveBuyMenu:
			OnSaveBuyMenu(p, client)
		default:
			log.Println("Unknown option packet", pkt.OptionPacketType, "from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal option packet from", client.RemoteAddr().String())
	}
}
