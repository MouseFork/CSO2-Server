package servermanager

import (
	"sync"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func NewChannelInfo(name []byte, chlsrv ChannelServer) ChannelInfo {
	var mutex sync.Mutex
	return ChannelInfo{
		GetNewChannelID(chlsrv),
		name,
		4,
		0x1F4,
		1,
		0,
		1,
		1,
		0,
		map[uint16]*Room{},
		map[uint8]uint16{},
		&mutex,
	}
}

//getNewChannelID() 暂定
func GetNewChannelID(chlsrv ChannelServer) uint8 {
	if chlsrv.ChannelCount > MAXCHANNELNUM {
		DebugInfo(2, "Channel is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXCHANNELNUM + 1]uint8
	//哈希思想
	for i := 0; i < int(chlsrv.ChannelCount); i++ {
		intbuf[chlsrv.Channels[i].ChannelID] = 1
	}
	//找到空闲的ID
	for i := 1; i <= int(MAXCHANNELNUM); i++ {
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
			return chlsrv.Channels[i]
		}
	}
	return nil
}

//添加房间,一次只能一个协程修改该频道
func AddChannelRoom(room *Room, chlid uint8, chlsrvid uint8) bool {
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
	if chl.RoomNum > MAXCHANNELROOMNUM {
		DebugInfo(2, "Room is too much ! Unable to add more !")
		return false
	}
	if room.Id <= 0 {
		DebugInfo(2, "ID of room is illegal !")
		return false
	}
	chl.ChannelMutex.Lock()
	defer chl.ChannelMutex.Unlock()
	if _, ok := chl.Rooms[room.Id]; ok {
		DebugInfo(2, "Room is already existed in Channel!")
		return false
	}
	if !AddRoomToManager(room) {
		return false
	}
	chl.RoomNum++
	chl.Rooms[room.Id] = room
	chl.RoomNums[room.RoomNumber] = room.Id
	return true
}

//删除频道房间
func DelChannelRoom(roomid uint16, chlid uint8, chlsrvid uint8) bool {
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
	if _, ok := chl.Rooms[roomid]; ok {
		chl.ChannelMutex.Lock()
		defer chl.ChannelMutex.Unlock()
		DebugInfo(1, "Room", string(chl.Rooms[roomid].Setting.RoomName), "id", roomid, "had been deleted !")
		if !DelRoomFromManager(roomid) {
			return false
		}
		chl.RoomNum--
		delete(chl.RoomNums, chl.Rooms[roomid].RoomNumber)
		delete(chl.Rooms, roomid)
		return true
	}
	return false
}

//getNewRoomNumber() 获取房间在某个频道下的标号
func GetNewRoomNumber(chl ChannelInfo) uint8 {
	if chl.RoomNum > MAXCHANNELROOMNUM {
		DebugInfo(2, "Error : Room is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	chl.ChannelMutex.Lock()
	defer chl.ChannelMutex.Unlock()
	for i := 1; i <= MAXCHANNELROOMNUM; i++ {
		if _, ok := chl.RoomNums[uint8(i)]; !ok {
			return uint8(i)
		}
	}
	return 0
}

func GetRoomFromID(chlsrvID uint8, chlID uint8, roomID uint16) *Room {
	if chlsrvID <= 0 ||
		chlID <= 0 ||
		roomID <= 0 {
		return nil
	}
	chlsrv := GetChannelServerWithID(chlsrvID)
	if chlsrv.ServerIndex <= 0 {
		return nil
	}
	chl := GetChannelWithID(chlID, *chlsrv)
	if chl.ChannelID <= 0 || chl.RoomNum <= 0 {
		return nil
	}
	chl.ChannelMutex.Lock()
	defer chl.ChannelMutex.Unlock()
	if v, ok := chl.Rooms[roomID]; ok {
		return v
	}
	return nil
}
