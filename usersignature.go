package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inSetSignaturePacket struct {
	len       uint8
	signature []byte
}

func onSetSignature(p packet, client net.Conn) {
	var pkt inSetSignaturePacket
	if !praseSetSignaturePacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal SetSignature packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to SetSignature but not in server !")
		return
	}
	//修改数据
	uPtr.SetSignature(pkt.signature)
	//发送数据包
	p.id = TypeUserInfo
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildSetSignaturePacket(uPtr.userid, pkt.signature, pkt.len))
	sendPacket(rst, uPtr.currentConnection)
	log.Println("User", string(uPtr.username), "Set Signature to", string(pkt.signature))
	//如果是在房间内
}

func praseSetSignaturePacket(p packet, dest *inSetSignaturePacket) bool {
	if p.datalen-HeaderLen < 3 {
		return false
	}
	offset := 6
	(*dest).len = ReadUint8(p.data, &offset)
	(*dest).signature = ReadString(p.data, &offset, int(dest.len))
	return true
}

func BuildSetSignaturePacket(id uint32, Signature []byte, len uint8) []byte {
	buf := make([]byte, 10+len)
	offset := 0
	WriteUint32(&buf, id, &offset)
	WriteUint32(&buf, 0x40000, &offset)
	WriteString(&buf, Signature, &offset)
	return buf[:offset]
}
