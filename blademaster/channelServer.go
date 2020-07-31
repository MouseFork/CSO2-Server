package blademaster

import (
	"sync"
)

type (
	//分区服务器，管理所拥有的频道
	ChannelServer struct {
		ServerIndex  uint8
		ServerStatus uint8
		ServerType   uint8
		ServerName   []byte
		ChannelCount uint8
		Channels     []ChannelInfo
	}

	//主服务器，管理各个分区
	ServerManager struct {
		ServerNum uint8
		Servers   []ChannelServer
	}

	//频道信息，隶属于分区服务器,用于请求服务器和请求频道
	ChannelInfo struct {
		ChannelID   uint8
		ChannelName []byte
		Unk00       uint16
		Unk01       uint16
		Unk02       uint8
		Unk03       uint8
		Unk04       uint8
		NextRoomID  uint8
		RoomNum     uint16
		Rooms       []RoomInfo

		ChannelMutex *sync.Mutex
	}
)

const (
	MAXCHANNELNUM uint8 = 16
	MAXSERVERNUM  uint8 = 8
)

var (
	//MainServer 主服务器
	MainServer = ServerManager{
		0,
		[]ChannelServer{},
	}
)
