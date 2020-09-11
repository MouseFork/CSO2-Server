package playerinfo

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
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
			DebugInfo(2, "Unknown PlayerInfo packet", pkt.InfoType, "from", client.RemoteAddr().String())
		}
	} else {
		DebugInfo(2, "Error : Recived a illegal PlayerInfo packet from", client.RemoteAddr().String())
	}
}
