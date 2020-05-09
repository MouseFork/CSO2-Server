package main

import (
	"log"
	"net"
)

var (
	DefalutServerName = "CSO2-Server"
)

//服务器，用于请求服务器
type server struct {
	serverIndex  uint8
	serverStatus uint8
	serverType   uint8
	serverName   []byte
	channelCount uint8
	channels     []channelInfo
}

func onServerList(seq *uint8, p *packet, client *(net.Conn)) {
	(*p).id = TypeServerList
	rst := BytesCombine(BuildHeader(seq, *p), BuildServerList())
	WriteLen(&rst)       //写入长度
	(*client).Write(rst) //发送UserStart消息
	log.Println("Sent a server list packet to", (*client).RemoteAddr().String())
}

func BuildServerList() []byte {
	var list []byte
	list = append(list, 0x01, GameServer.serverIndex,
		GameServer.serverStatus,
		GameServer.serverType,
		uint8(len(GameServer.serverName)),
	)
	for i := 0; i < len(GameServer.serverName); i++ {
		list = append(list, GameServer.serverName[i])
	}
	list = append(list, GameServer.channelCount)
	chl := BuildChannelList(GameServer.channelCount, GameServer.channels)
	for i := 0; i < len(chl); i++ {
		list = append(list, chl[i])
	}
	return list
}

func newServer(name []byte) server {
	return server{
		1,
		1,
		3,
		name,
		1,
		[]channelInfo{newChannelInfo([]byte(DefalutChannelName))},
	}
}

func getChannelServerWithID(id uint8) *server {
	// count := GameServer.channelCount
	// for i := 0; i < int(count); i++ {
	// 	if GameServer.channels[i].channelID == id {
	// 		return &(GameServer.channels[i])
	// 	}
	// }
	//log.Println("GameServerIndex:", strconv.Itoa(int(GameServer.serverIndex)))
	if GameServer.serverIndex == id {
		return &GameServer
	}
	return nil
}
