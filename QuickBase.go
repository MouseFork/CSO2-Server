package main

import (
	"log"
	"net"
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
	if praseQuickPacket(p, &pkt) {
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
func praseQuickPacket(p packet, dest *inQuickPacket) bool {
	if p.datalen-HeaderLen < 2 {
		return false
	}
	(*dest).inQuickType = p.data[5]
	return true
}
