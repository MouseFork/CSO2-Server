package servermanager

import (
	"log"
	"net"
	"sync"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	//MainServer 主服务器
	MainServer = ServerManager{
		0,
		[]*ChannelServer{},
	}
	mainServerMutex sync.Mutex
)

//初始化主服务器，仅运行一次，不考虑互斥问题
func NewMainServer() ServerManager {
	srvmgr := ServerManager{
		0, //初始化，默认数量为0
		[]*ChannelServer{},
	}
	chlsrv := NewChannelServer([]byte(DefalutServerName))
	chl1 := NewChannelInfo([]byte(DefalutChannelName1), chlsrv)
	if !AddChannel(&chlsrv, &chl1) ||
		!AddChannelServer(&srvmgr, &chlsrv) {
		DebugInfo(2, "Error : Unable to initializing main server !")
		log.Fatalln("")
	}
	chl2 := NewChannelInfo([]byte(DefalutChannelName2), chlsrv)
	if !AddChannel(&chlsrv, &chl2) {
		DebugInfo(2, "Error : Unable to initializing main server !")
		log.Fatalln("")
	}
	return srvmgr
}

//主服务器里加频道服务器,一次只能一个协程访问
func AddChannelServer(dest *ServerManager, src *ChannelServer) bool {
	mainServerMutex.Lock()
	defer mainServerMutex.Unlock()
	if dest.ServerNum > MAXSERVERNUM {
		DebugInfo(2, "Error : ChannelServer is too much ! Unable to add more !")
		return false
	}
	if src.ServerIndex == 0 {
		DebugInfo(2, "Error : ID of ChannelServer is illegal !")
		return false
	}
	for _, v := range dest.Servers {
		if v.ServerIndex == src.ServerIndex {
			DebugInfo(2, "Error : ChannelServer is already existed in MainServer!")
			return false
		}
	}
	dest.ServerNum++
	dest.Servers = append(dest.Servers, src)
	return true
}

//新建频道服务器
func NewChannelServer(name []byte) ChannelServer {
	chlsrv := ChannelServer{
		GetNewChannelServerID(),
		1,
		3,
		name,
		0,
		[]*ChannelInfo{},
	}
	return chlsrv
}

//频道服务器里加频道,一次只能一个协程访问
func AddChannel(dest *ChannelServer, src *ChannelInfo) bool {
	mainServerMutex.Lock()
	defer mainServerMutex.Unlock()
	if dest.ChannelCount > MAXCHANNELNUM {
		DebugInfo(2, "Error : Channel is too much ! Unable to add more !")
		return false
	}
	if src.ChannelID <= 0 {
		DebugInfo(2, "Error : ID of Channel is illegal !")
		return false
	}
	for _, v := range dest.Channels {
		if v.ChannelID == src.ChannelID {
			DebugInfo(2, "Error : Channel is already existed in ChannelServer!")
			return false
		}
	}
	dest.ChannelCount++
	dest.Channels = append(dest.Channels, src)
	return true
}

//新的频道服务器ID,一次只能一个协程访问
func GetNewChannelServerID() uint8 {
	mainServerMutex.Lock()
	defer mainServerMutex.Unlock()
	if MainServer.ServerNum > MAXSERVERNUM {
		DebugInfo(2, "Error : ChannelServer is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXSERVERNUM + 1]uint8
	//哈希思想
	for i := 0; i < int(MainServer.ServerNum); i++ {
		intbuf[MainServer.Servers[i].ServerIndex] = 1
	}
	//找到空闲的ID
	for i := 1; i <= int(MAXSERVERNUM); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint8(i)
		}
	}
	return 0
}

//通过ID得到对应频道服务器
func GetChannelServerWithID(id uint8) *ChannelServer {
	if id <= 0 {
		return nil
	}
	count := MainServer.ServerNum
	for i := 0; i < int(count); i++ {
		if MainServer.Servers[i].ServerIndex == id {
			return MainServer.Servers[i]
		}
	}
	return nil
}

//处理请求服务器列表
func OnServerList(client net.Conn) {
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : A unknow Client from", client.RemoteAddr().String(), "request a ChannelList !")
		return
	}
	rst := BuildHeader(uPtr.CurrentSequence, PacketTypeServerList)
	rst = append(rst, MainServer.ServerNum)
	for i := 0; i < int(MainServer.ServerNum); i++ {
		rst = BytesCombine(rst, BuildChannelServerList(*MainServer.Servers[i]))
	}
	SendPacket(rst, uPtr.CurrentConnection)

	uPtr.SetUserChannelServer(0)
	uPtr.QuitChannel()
	DebugInfo(2, "Sent a server list packet to", uPtr.CurrentConnection.RemoteAddr().String())

}

//建立某个频道服务器数据包
func BuildChannelServerList(chlsrv ChannelServer) []byte {
	var list []byte
	l := len(chlsrv.ServerName)
	list = append(list, chlsrv.ServerIndex,
		chlsrv.ServerStatus,
		chlsrv.ServerType,
		uint8(l),
	)
	list = BytesCombine(list, chlsrv.ServerName)
	list = append(list, chlsrv.ChannelCount)
	list = BytesCombine(list, BuildChannelList(chlsrv.ChannelCount, chlsrv.Channels))
	return list
}

func BuildChannelList(num uint8, channels []*ChannelInfo) []byte {
	var list []byte
	for i := 0; i < int(num); i++ {
		temp := make([]byte, 9+len(channels[i].ChannelName))
		offset := 0
		WriteUint8(&temp, channels[i].ChannelID, &offset)
		WriteString(&temp, channels[i].ChannelName, &offset)
		WriteUint16(&temp, channels[i].Unk00, &offset)
		WriteUint16(&temp, channels[i].Unk01, &offset)
		WriteUint8(&temp, channels[i].Unk02, &offset)
		WriteUint8(&temp, channels[i].Unk03, &offset)
		WriteUint8(&temp, channels[i].Unk04, &offset)
		list = BytesCombine(list, temp[:offset])
	}
	return list
}
