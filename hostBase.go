package main

import (
	"log"
	"net"
)

const (
	GameStart         = 0 // when a host starts a new game
	HostJoin          = 1 // when someone joins some host's game
	HostStop          = 3
	LeaveResultWindow = 4

	TeamChanging = 11
	// logging packet types
	OnGameEnd = 21

	SetInventory = 101 // there are 2 or 3 other host packet types that send this
	ItemUsing    = 105
	SetLoadout   = 107
	SetBuyMenu   = 111
)

type inHostPacket struct {
	inHostType uint8
}

func onHost(seq *uint8, p packet, client net.Conn) {
	var pkt inHostPacket
	if praseHostPacket(p, &pkt) {
		getUserFromConnection(client)
		switch pkt.inHostType {
		case OnGameEnd:
			log.Println("Recived a game end packet from", client.RemoteAddr().String())

		case SetInventory:
			log.Println("Recived a setInventory packet from", client.RemoteAddr().String())
			onHostSetUserInventory(p, client)
		case SetLoadout:
			log.Println("Recived a setLoadout packet from", client.RemoteAddr().String())

		case SetBuyMenu:
			log.Println("Recived a setBuyMenu packet from", client.RemoteAddr().String())

		case TeamChanging:
			log.Println("Recived a change team packet from", client.RemoteAddr().String())

		case ItemUsing:
			log.Println("Recived a use item packet from", client.RemoteAddr().String())

		default:
			log.Println("Recived a unknown host packet from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal host packet from", client.RemoteAddr().String())
	}
}

func praseHostPacket(p packet, dest *inHostPacket) bool {
	if p.datalen-HeaderLen < 2 {
		return false
	}
	(*dest).inHostType = p.data[5]
	return true
}
