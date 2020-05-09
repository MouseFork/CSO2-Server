package main

var (
	DefalutChannelName = "CSO2-Channel"
)

//频道信息，隶属于Server频道,用于请求服务器和请求频道
type channelInfo struct {
	channelID   uint8
	channelName []byte
	unk00       uint16
	unk01       uint16
	unk02       uint8
	unk03       uint8
	unk04       uint8
	nextRoomID  uint8
	roomNum     uint16
	rooms       []roomInfo
}

// //频道内容,包括房间,用于请求频道
// type channel struct {
// 	channelID uint8
// 	roomNum   uint16
// 	rooms     []room
// }

func BuildChannelList(num uint8, channels []channelInfo) []byte {
	var list []byte
	for i := 0; i < int(num); i++ {
		temp := make([]byte, 9+len(channels[i].channelName))
		offset := 0
		WriteUint8(&temp, channels[i].channelID, &offset)
		WriteString(&temp, channels[i].channelName, &offset)
		WriteUint16(&temp, channels[i].unk00, &offset)
		WriteUint16(&temp, channels[i].unk01, &offset)
		WriteUint8(&temp, channels[i].unk02, &offset)
		WriteUint8(&temp, channels[i].unk03, &offset)
		WriteUint8(&temp, channels[i].unk04, &offset)
		for j := 0; j < len(temp); j++ {
			list = append(list, temp[j])
		}
	}
	return list
}

func newChannelInfo(name []byte) channelInfo {
	return channelInfo{
		getNewChannelID(),
		name,
		4,
		0x1F4,
		1,
		0,
		1,
		1,
		0,
		[]roomInfo{},
	}
}

//getNewChannelID() 暂定
func getNewChannelID() uint8 {
	return 1
}

func getChannelWithID(id uint8) *channelInfo {
	count := GameServer.channelCount
	for i := 0; i < int(count); i++ {
		//log.Println("ChannelIndex:", strconv.Itoa(int(GameServer.channels[i].channelID)))
		if GameServer.channels[i].channelID == id {
			return &(GameServer.channels[i])
		}
	}
	return nil
}

func newChannelRoom(host uint32, id uint8) {
	chlptr := getChannelWithID(id)
	// room := roomInfo{
	// 	(*chlptr).nextRoomID,

	// }
	//(*chlptr).rooms = append((*chlptr).rooms, room)
	(*chlptr).roomNum++

}
