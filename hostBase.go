package main

import (
	"fmt"
	"log"
	"net"
)

const (
	DataCheck    = 1 //也许是
	UserKillOne  = 7
	UserDeath    = 8
	TeamChanging = 11
	UserRevived  = 20
	OnGameEnd    = 21

	SetInventory = 101 // 不一定是101，可能其他的数据包
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
		case DataCheck:

		case OnGameEnd:
			//log.Println("Recived a game end packet from", client.RemoteAddr().String())
			onHostGameEnd(p, client)
		case SetInventory:
			//log.Println("Recived a setInventory packet from", client.RemoteAddr().String())
			onHostSetUserInventory(p, client)
		case SetLoadout:
			//log.Println("Recived a setLoadout packet from", client.RemoteAddr().String())
			onHostSetUserLoadout(p, client)
		case SetBuyMenu:
			log.Println("Recived a setBuyMenu packet from", client.RemoteAddr().String())

		case TeamChanging:
			log.Println("Recived a change team packet from", client.RemoteAddr().String())

		case ItemUsing:
			log.Println("Recived a use item packet from", client.RemoteAddr().String())
		case UserKillOne:
			fmt.Println("Kill   packet", p.data[:p.datalen], "from", client.RemoteAddr().String())
		case UserDeath:
			fmt.Println("death  packet", p.data[:p.datalen], "from", client.RemoteAddr().String())
		case UserRevived:
			fmt.Println("Revive packet", p.data[:p.datalen], "from", client.RemoteAddr().String())
		default:
			log.Println("Unknown host packet", pkt.inHostType, "from", client.RemoteAddr().String())
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
