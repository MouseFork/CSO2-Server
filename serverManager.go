package main

import (
	"log"
	"net"
	"os"
)

const (
	DefalutServerName       = "CSO2-Server"
	MAXSERVERNUM      uint8 = 8
)

//服务器，用于请求服务器
type channelServer struct {
	serverIndex  uint8
	serverStatus uint8
	serverType   uint8
	serverName   []byte
	channelCount uint8
	channels     []channelInfo
}

type serverManager struct {
	serverNum uint8
	servers   []channelServer
}

//处理请求服务器列表
func onServerList(seq *uint8, p *packet, client *(net.Conn)) {
	(*p).id = TypeServerList
	rst := BuildHeader(seq, *p)
	rst = append(rst, MainServer.serverNum)
	for i := 0; i < int(MainServer.serverNum); i++ {
		rst = BytesCombine(rst, BuildChannelServerList(MainServer.servers[i]))
	}
	WriteLen(&rst)       //写入长度
	(*client).Write(rst) //发送UserStart消息
	log.Println("Sent a server list packet to", (*client).RemoteAddr().String())
}

//建立某个频道服务器数据包
func BuildChannelServerList(chlsrv channelServer) []byte {
	var list []byte
	l := len(chlsrv.serverName)
	list = append(list, chlsrv.serverIndex,
		chlsrv.serverStatus,
		chlsrv.serverType,
		uint8(l),
	)
	for i := 0; i < l; i++ {
		list = append(list, chlsrv.serverName[i])
	}
	list = append(list, chlsrv.channelCount)
	chl := BuildChannelList(chlsrv.channelCount, chlsrv.channels)
	for i := 0; i < len(chl); i++ {
		list = append(list, chl[i])
	}
	return list
}

//初始化主服务器，仅运行一次
func newMainServer() serverManager {
	srvmgr := serverManager{
		0, //初始化，默认数量为0
		[]channelServer{},
	}
	chl := newChannelServer([]byte(DefalutServerName))
	if !addChannelServer(&srvmgr, &chl) {
		log.Fatalln("Unable to initializing main server !")
		os.Exit(-1)
	}
	return srvmgr
}

//主服务器里加频道服务器
func addChannelServer(dest *serverManager, src *channelServer) bool {
	if (*dest).serverNum > MAXSERVERNUM {
		log.Fatalln("ChannelServer is too much ! Unable to add more !")
		return false
	}
	if (*src).serverIndex == 0 {
		log.Fatalln("ID of ChannelServer is illegal !")
		return false
	}
	(*dest).serverNum++
	(*dest).servers = append((*dest).servers, *src)
	return true
}

//新建频道服务器
func newChannelServer(name []byte) channelServer {
	chlsrv := channelServer{
		getNewChannelServerID(),
		1,
		3,
		name,
		0,
		[]channelInfo{},
	}
	chl := newChannelInfo([]byte(DefalutChannelName), chlsrv)
	if !addChannel(&chlsrv, &chl) {
		log.Fatalln("Unable to initializing main server !")
		return channelServer{
			0,
			1,
			3,
			name,
			0,
			[]channelInfo{},
		}
	}
	return chlsrv
}

//频道服务器里加频道
func addChannel(dest *channelServer, src *channelInfo) bool {
	if (*dest).channelCount > MAXCHANNELNUM {
		log.Fatalln("Channel is too much ! Unable to add more !")
		return false
	}
	if (*src).channelID == 0 {
		log.Fatalln("ID of Channel is illegal !")
		return false
	}
	(*dest).channelCount++
	(*dest).channels = append((*dest).channels, *src)
	return true
}

//新的频道服务器ID
func getNewChannelServerID() uint8 {
	if MainServer.serverNum > MAXSERVERNUM {
		log.Fatalln("ChannelServer is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXSERVERNUM + 2]uint8
	//哈希思想
	for i := 0; i < int(MainServer.serverNum); i++ {
		intbuf[MainServer.servers[i].serverIndex] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXSERVERNUM+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint8(i)
		}
	}
	return 0
}

//通过ID得到对应频道服务器
func getChannelServerWithID(id uint8) *channelServer {
	count := MainServer.serverNum
	for i := 0; i < int(count); i++ {
		if MainServer.servers[i].serverIndex == id {
			return &(MainServer.servers[i])
		}
	}
	return nil
}
