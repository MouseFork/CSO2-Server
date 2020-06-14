package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	UdpPacketSignature = 0x57
)

type InUDPmsg struct {
	signature          uint8
	userId             uint32
	portId             uint16
	ipAddress          uint32
	port               uint16
	packetData         []byte
	curOffset          uint32 //可能32位
	parsedSuccessfully bool
}

func startHolePunchServer(server *(net.UDPConn)) {
	defer server.Close()
	fmt.Println("Server holepunch is running at", "[AnyAdapter]:"+strconv.Itoa(HOLEPUNCHPORT))
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := server.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("Server holepunch read error from", server.RemoteAddr().String())
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		// client, err := (*server).Accept()
		// if err != nil {
		// 	log.Fatal("Server Accept data error !\n")
		// 	continue
		// }
		// log.Println("Server accept a new connection request at", client.RemoteAddr().String())
		// go RecvHolePunchMessage(client)
	}
}

//RecvHolePunchMessage 处理收到的包
func RecvHolePunchMessage(holeclient net.Conn) {

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
