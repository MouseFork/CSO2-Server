package main

import (
	"log"
	"net"
)

type inPlayerInfoPacket struct {
	infoType uint8
}

const (
	SetSignature = 5
	SetTitle     = 6
	SetAvatar    = 7
)

func onPlayerInfo(p packet, client net.Conn) {
	var pkt inPlayerInfoPacket
	if prasePlayerInfoPacket(p, &pkt) {
		switch pkt.infoType {
		case SetSignature:
			onSetSignature(p, client)
		case SetTitle:
			onSetTitle(p, client)
		case SetAvatar:
			onSetAvatar(p, client)
		default:
			log.Println("Unknown PlayerInfo packet", pkt.infoType, "from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal PlayerInfo packet from", client.RemoteAddr().String())
	}
}

func prasePlayerInfoPacket(p packet, dest *inPlayerInfoPacket) bool {
	if p.datalen-HeaderLen < 2 {
		return false
	}
	(*dest).infoType = p.data[5]
	return true
}
