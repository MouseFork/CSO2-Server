package main

import (
	"log"
	"net"
)

const (
	InProgress = 0
	Stop       = 1
)

type InRoomCountdownPacket struct {
	CountdownType uint8
	count         uint8
}

func onGameStartCountdown(p packet, client net.Conn) {
	//检索数据包
	var pkt InRoomCountdownPacket
	if !praseRoomCountdownPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error Countdown packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to start counting but not in server !")
		return
	}
	//检查用户是不是房主
	curroom := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if curroom == nil ||
		curroom.id <= 0 {
		log.Println("Error : User:", string(uPtr.username), "try to start counting in a null room !")
		return
	}
	if curroom.hostUserID != uPtr.userid {
		log.Println("Error : User:", string(uPtr.username), "try to start counting but is not the host !")
		return
	}
	//检查用户所在房间
	if curroom.id != uPtr.currentRoomId {
		log.Println("Error : User:", string(uPtr.username), "try to start counting but in another room !")
		return
	}
	//检查当前游戏模式
	if !curroom.canStartGame() {
		log.Println("Error : User:", string(uPtr.username), "try to start countdown but is illegal !")
		return
	}
	//检查是否应该继续倒计时
	should := pkt.shouldCountdown()
	if should {
		curroom.progressCountdown(pkt.count)
		log.Println("User:", string(uPtr.username), "countdown at", curroom.getCountdown(), "host is", pkt.count)
	} else {
		curroom.stopCountdown()
		log.Println("User:", string(uPtr.username), "cancled room countdown")
	}
	//所有玩家发送倒计时数据
	for _, v := range curroom.users {
		rst := BytesCombine(BuildHeader(v.currentSequence, p), BuildCountdown(pkt, should))
		sendPacket(rst, v.currentConnection)
	}
}

func praseRoomCountdownPacket(p packet, dest *InRoomCountdownPacket) bool {
	if p.datalen < 7 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).CountdownType = ReadUint8(p.data, &offset)
	if p.datalen < 8 {
		return false
	}
	if (*dest).CountdownType == InProgress {
		(*dest).count = ReadUint8(p.data, &offset)
	}
	return true
}

func (p InRoomCountdownPacket) shouldCountdown() bool {
	return p.CountdownType == InProgress
}

func BuildCountdown(p InRoomCountdownPacket, should bool) []byte {
	buf := make([]byte, 3)
	offset := 0
	WriteUint8(&buf, OUTCountdown, &offset)
	WriteUint8(&buf, p.CountdownType, &offset)
	if should {
		WriteUint8(&buf, p.count, &offset)
	}
	return buf[:offset]
}
