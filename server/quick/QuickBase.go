package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

const (
	QuickList = 0
	QuickJoin = 1
)

type inQuickPacket struct {
	inQuickType uint8
}

func onQuick(seq *uint8, p packet, client net.Conn) {
	var pkt inQuickPacket
	if p.PraseQuickPacket(&pkt) {
		switch pkt.inQuickType {
		case QuickList:
			onQuickList(p, client)
		case QuickJoin:
		default:
			log.Println("Unknown Quick packet", pkt.inQuickType, "from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal Quick packet from", client.RemoteAddr().String())
	}
}
