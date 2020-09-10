package typestruct

import "sync"

type (
	//分区服务器，管理所拥有的频道
	ChannelServer struct {
		ServerIndex  uint8
		ServerStatus uint8
		ServerType   uint8
		ServerName   []byte
		ChannelCount uint8
		Channels     []*ChannelInfo
	}

	//主服务器，管理各个分区
	ServerManager struct {
		ServerNum uint8
		Servers   []*ChannelServer
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
		Rooms       map[uint16]*Room
		RoomNums    map[uint8]uint16

		ChannelMutex *sync.Mutex
	}
)

const (
	MAXCHANNELNUM       = 16
	MAXSERVERNUM        = 15
	MAXCHANNELROOMNUM   = 0xFF
	MAXROOMNUM          = 0xFFFF
	DefalutServerName   = "CSO2-Server[1/1]"
	DefalutChannelName1 = "CSO2-Channel[1/2]"
	DefalutChannelName2 = "CSO2-Channel[2/2]"
)
