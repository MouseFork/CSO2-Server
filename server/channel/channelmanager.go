package channel

import (
	"log"

	. "github.com/KouKouChan/CSO2-Server/model/channel"
	. "github.com/KouKouChan/CSO2-Server/model/lock"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	DefalutServerName = "CSO2-Server[1/1]"
)

//初始化主服务器，仅运行一次，不考虑互斥问题
func NewMainServer() ServerManager {
	srvmgr := ServerManager{
		0, //初始化，默认数量为0
		[]ChannelServer{},
	}
	chlsrv := newChannelServer([]byte(DefalutServerName))
	chl1 := newChannelInfo([]byte(DefalutChannelName1), chlsrv)
	chl2 := newChannelInfo([]byte(DefalutChannelName1), chlsrv)
	if !addChannel(&chlsrv, &chl1) ||
		!addChannel(&chlsrv, &chl2) ||
		!addChannelServer(&srvmgr, &chlsrv) {
		DebugInfo(2, "Error : Unable to initializing main server !")
		log.Fatalln("")
	}
	return srvmgr
}

//主服务器里加频道服务器,一次只能一个协程访问
func addChannelServer(dest *ServerManager, src *ChannelServer) bool {
	GlobalMutex.Lock()
	MainServerMutex.Lock()
	defer GlobalMutex.Unlock()
	defer MainServerMutex.Unlock()
	if (*dest).ServerNum > MAXSERVERNUM {
		DebugInfo(2, "Error : ChannelServer is too much ! Unable to add more !")
		return false
	}
	if (*src).ServerIndex == 0 {
		DebugInfo(2, "Error : ID of ChannelServer is illegal !")
		return false
	}
	for _, v := range dest.Servers {
		if v.ServerIndex == src.ServerIndex {
			DebugInfo(2, "Error : ChannelServer is already existed in MainServer!")
			return false
		}
	}
	(*dest).ServerNum++
	(*dest).Servers = append((*dest).Servers, *src)
	return true
}

//新建频道服务器
func newChannelServer(name []byte) ChannelServer {
	chlsrv := ChannelServer{
		getNewChannelServerID(),
		1,
		3,
		name,
		0,
		[]ChannelInfo{},
	}
	return chlsrv
}

//频道服务器里加频道,一次只能一个协程访问
func addChannel(dest *ChannelServer, src *ChannelInfo) bool {
	GlobalMutex.Lock()
	MainServerMutex.Lock()
	defer GlobalMutex.Unlock()
	defer MainServerMutex.Unlock()
	if (*dest).ChannelCount > MAXCHANNELNUM {
		DebugInfo(2, "Error : Channel is too much ! Unable to add more !")
		return false
	}
	if (*src).ChannelID == 0 {
		DebugInfo(2, "Error : ID of Channel is illegal !")
		return false
	}
	for _, v := range dest.Channels {
		if v.ChannelID == src.ChannelID {
			DebugInfo(2, "Error : Channel is already existed in ChannelServer!")
			return false
		}
	}
	(*dest).ChannelCount++
	(*dest).Channels = append((*dest).Channels, *src)
	return true
}

//新的频道服务器ID,一次只能一个协程访问
func getNewChannelServerID() uint8 {
	GlobalMutex.Lock()
	MainServerMutex.Lock()
	defer GlobalMutex.Unlock()
	defer MainServerMutex.Unlock()
	if MainServer.ServerNum > MAXSERVERNUM {
		DebugInfo(2, "Error : ChannelServer is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXSERVERNUM + 2]uint8
	//哈希思想
	for i := 0; i < int(MainServer.ServerNum); i++ {
		intbuf[MainServer.Servers[i].ServerIndex] = 1
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
func GetChannelServerWithID(id uint8) *ChannelServer {
	count := MainServer.ServerNum
	for i := 0; i < int(count); i++ {
		if MainServer.Servers[i].ServerIndex == id {
			return &(MainServer.Servers[i])
		}
	}
	return nil
}
