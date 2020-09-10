package quick

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	QuickList = 0
	QuickJoin = 1
)

func OnQuick(p *PacketData, client net.Conn) {
	var pkt InQuickPacket
	if p.PraseQuickPacket(&pkt) {
		switch pkt.InQuickType {
		case QuickList:
			OnQuickList(p, client)
		case QuickJoin:
		default:
			DebugInfo(2, "Unknown Quick packet", pkt.InQuickType, "from", client.RemoteAddr().String())
		}
	} else {
		DebugInfo(1, "Error : Recived a illegal Quick packet from", client.RemoteAddr().String())
	}
}
