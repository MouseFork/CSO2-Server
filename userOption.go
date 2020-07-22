package main

import (
	"log"
	"net"
)

type InOptionPacket struct {
	OptionPacketType uint8
}

const (
	SaveBuyMenu = 1
)

func onOption(p packet, client net.Conn) {
	var pkt InOptionPacket
	if praseOptionPacket(p, &pkt) {
		switch pkt.OptionPacketType {
		case SaveBuyMenu:
			onSaveBuyMenu(p, client)
		default:
			log.Println("Unknown option packet", pkt.OptionPacketType, "from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal option packet from", client.RemoteAddr().String())
	}
}

func praseOptionPacket(p packet, dest *InOptionPacket) bool {
	if p.datalen-HeaderLen < 2 {
		return false
	}
	(*dest).OptionPacketType = p.data[5]
	return true
}
