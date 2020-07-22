package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

const (
	UdpPacketSignature = 0x57
	UDPTypeClient      = 0
	UDPTypeServer      = 256
	UDPTypeSourceTV    = 512
)

type InUDPmsg struct {
	signature uint8
	userId    uint32
	portId    uint16
	ipAddress uint32
	port      uint16

	packetData         []byte
	datalen            int
	curOffset          int //可能32位
	parsedSuccessfully bool
}

func startHolePunchServer(server *net.UDPConn) {
	defer server.Close()
	fmt.Println("Server UDPholepunch is running at", "[AnyAdapter]:"+strconv.Itoa(HOLEPUNCHPORT))
	for {
		data := make([]byte, 1024)
		n, ClientAddr, err := server.ReadFromUDP(data)
		if err != nil {
			log.Println("UDP read error from", ClientAddr.String())
			continue
		}
		//log.Println("UDP packet from", ClientAddr.String())
		go RecvHolePunchMessage(data[:n], n, ClientAddr, server)
	}
}

//RecvHolePunchMessage 处理收到的包
func RecvHolePunchMessage(data []byte, len int, client *net.UDPAddr, server *net.UDPConn) {
	//分析数据包
	var p InUDPmsg
	if !praseUDPpacket(data, len, &p) {
		log.Println("UDP had a illegal packet from", client.String())
		return
	}
	if p.isHeartbeat() {
		return
	}
	cliadr := client.IP.To4().String()
	externalIPAddress, err := IPToUint32(cliadr)
	if err != nil {
		log.Println("Error : Prasing externalIpAddress error !")
		return
	}
	//找到对应玩家
	uPtr := getUserFromID(p.userId)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		//log.Println("UDP had a packet from", client.String(), "but not logged in !")
		return
	}
	//更新netinfo
	index := uPtr.updateHolepunch(p.portId, p.port, uint16(client.Port))
	if index == 0xFFFF {
		log.Println("Error : User", uPtr.username, "update Holepunch failed !")
		return
	}
	(*uPtr).netInfo.ExternalIpAddress = externalIPAddress
	(*uPtr).netInfo.LocalIpAddress = p.ipAddress
	//发送返回数据
	rst := BuildUDPHolepunch(index)
	server.WriteToUDP(rst, client)
	// log.Println(uPtr.netInfo.ExternalIpAddress,
	// 	uPtr.netInfo.ExternalServerPort,
	// 	uPtr.netInfo.ExternalClientPort,
	// 	uPtr.netInfo.ExternalTvPort,
	// 	uPtr.netInfo.LocalIpAddress,
	// 	uPtr.netInfo.LocalServerPort,
	// 	uPtr.netInfo.LocalClientPort,
	// 	uPtr.netInfo.LocalTvPort)
}

func praseUDPpacket(data []byte, len int, dest *InUDPmsg) bool {
	(*dest).curOffset = 0
	(*dest).signature = ReadUint8(data, &dest.curOffset)
	if (*dest).signature != UdpPacketSignature {
		(*dest).parsedSuccessfully = false
		return false
	}
	(*dest).datalen = len
	(*dest).packetData = data
	if dest.isHeartbeat() {
	} else {
		(*dest).userId = ReadUint32(data, &dest.curOffset)
		(*dest).portId = ReadUint16(data, &dest.curOffset)
		(*dest).ipAddress = ReadUint32BE(data, &dest.curOffset)
		(*dest).port = ReadUint16(data, &dest.curOffset)
	}
	(*dest).parsedSuccessfully = true
	return true
}

func BuildUDPHolepunch(index uint16) []byte {
	buf := make([]byte, 3)
	offset := 0
	WriteUint8(&buf, UdpPacketSignature, &offset)
	WriteUint16(&buf, index, &offset)
	return buf[:offset]
}
func UDPBuild(seq *uint8, p packet, isHost uint8, userid uint32, ip uint32, port uint16) []byte {
	p.id = TypeUdp
	rst := BuildHeader(seq, p)
	buf := make([]byte, 12)
	offset := 0
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, isHost, &offset)
	WriteUint32(&buf, userid, &offset)
	WriteUint32BE(&buf, ip, &offset)
	WriteUint16(&buf, port, &offset)
	rst = BytesCombine(rst, buf)
	return rst
}

func (p InUDPmsg) isHeartbeat() bool {
	return p.datalen == 6
}
