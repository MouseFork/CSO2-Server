package channel

import (
	"sync"

	. "github.com/KouKouChan/CSO2-Server/model/channel"
	. "github.com/KouKouChan/CSO2-Server/model/room"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	DefalutChannelName1 = "CSO2-Channel[1/2]"
	DefalutChannelName2 = "CSO2-Channel[2/2]"
)

func newChannelInfo(name []byte, chlsrv ChannelServer) ChannelInfo {
	var mutex sync.Mutex
	return ChannelInfo{
		getNewChannelID(chlsrv),
		name,
		4,
		0x1F4,
		1,
		0,
		1,
		1,
		0,
		[]RoomInfo{},
		&mutex,
	}
}

//getNewChannelID() 暂定
func getNewChannelID(chlsrv ChannelServer) uint8 {
	if chlsrv.ChannelCount > MAXCHANNELNUM {
		DebugInfo(2, "Channel is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXCHANNELNUM + 2]uint8
	//哈希思想
	for i := 0; i < int(chlsrv.ChannelCount); i++ {
		intbuf[chlsrv.Channels[i].ChannelID] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXCHANNELNUM+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint8(i)
		}
	}
	return 0
}

//通过ID获取频道
func GetChannelWithID(id uint8, chlsrv ChannelServer) *ChannelInfo {
	count := chlsrv.ChannelCount
	for i := 0; i < int(count); i++ {
		if chlsrv.Channels[i].ChannelID == id {
			return &(chlsrv.Channels[i])
		}
	}
	return nil
}

//添加房间,一次只能一个协程修改该频道
func addChannelRoom(room RoomInfo, chlid uint8, chlsrvid uint8) bool {
	chlsrv := GetChannelServerWithID(chlsrvid)
	if chlsrv.ServerIndex <= 0 {
		DebugInfo(2, "Add room to a null channelServer!")
		return false
	}
	chl := GetChannelWithID(chlid, *chlsrv)
	if chl.ChannelID <= 0 {
		DebugInfo(2, "Add room to a null channel!")
		return false
	}
	if chl.RoomNum > MAXROOMNUMS {
		DebugInfo(2, "Room is too much ! Unable to add more !")
		return false
	}
	if room.Id <= 0 {
		DebugInfo(2, "ID of room is illegal !")
		return false
	}
	for _, v := range chl.Rooms {
		if v.Id == room.Id {
			DebugInfo(2, "Room is already existed in Channel!")
			return false
		}
	}
	chl.ChannelMutex.Lock()
	defer chl.ChannelMutex.Unlock()
	chl.RoomNum++
	chl.Rooms = append(chl.Rooms, room)
	return true
}

//删除频道房间
func delChannelRoom(roomid uint16, chlid uint8, chlsrvid uint8) bool {
	chlsrv := GetChannelServerWithID(chlsrvid)
	if chlsrv.ServerIndex <= 0 {
		DebugInfo(2, "Remove room to a null channelServer!")
		return false
	}
	chl := GetChannelWithID(chlid, *chlsrv)
	if chl.ChannelID <= 0 {
		DebugInfo(2, "Remove room to a null channel!")
		return false
	}
	if chl.RoomNum <= 0 {
		DebugInfo(2, "There is no room in this channel , unable to remove!")
		return false
	}
	if roomid <= 0 {
		DebugInfo(2, "ID of room is illegal !")
		return false
	}
	for k, v := range chl.Rooms {
		if v.Id == roomid {
			chl.ChannelMutex.Lock()
			defer chl.ChannelMutex.Unlock()
			chl.RoomNum--
			chl.Rooms = append(chl.Rooms[:k], chl.Rooms[k+1:]...)
			DebugInfo(1, "Room", string(v.Setting.RoomName), "id", roomid, "had been deleted !")
			return true
		}
	}
	return false
}
