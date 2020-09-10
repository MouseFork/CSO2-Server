package servermanager

import (
	"sync"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

type RM struct {
	Rooms map[uint16]*Room
	Lock  *sync.Mutex
}

var (
	rmlock       sync.Mutex
	RoomsManager = RM{
		map[uint16]*Room{},
		&rmlock,
	}
)

//getNewRoomNumber() 获取房间在某个频道下的标号
func GetNewRoomID(chl ChannelInfo) uint16 {
	if chl.RoomNum > MAXROOMNUM {
		DebugInfo(2, "Error : Room is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	RoomsManager.Lock.Lock()
	defer RoomsManager.Lock.Unlock()
	//第一个不存在的ID
	for i := 1; i < MAXROOMNUM; i++ {
		if _, ok := RoomsManager.Rooms[uint16(i)]; !ok {
			return uint16(i)
		}
	}
	return 0
}

//添加房间,一次只能一个协程修改该频道
func AddRoomToManager(room *Room) bool {
	if room.Id <= 0 {
		DebugInfo(2, "RoomManager Error: ID of room is illegal !")
		return false
	}
	RoomsManager.Lock.Lock()
	defer RoomsManager.Lock.Unlock()
	if _, ok := RoomsManager.Rooms[room.Id]; ok {
		DebugInfo(2, "RoomManager Error: Room is already existed in Channel!")
		return false
	}
	RoomsManager.Rooms[room.Id] = room
	return true
}

//删除频道房间
func DelRoomFromManager(roomid uint16) bool {
	if roomid <= 0 {
		DebugInfo(2, "RoomManager Error: ID of room is illegal !")
		return false
	}
	if _, ok := RoomsManager.Rooms[roomid]; ok {
		RoomsManager.Lock.Lock()
		defer RoomsManager.Lock.Unlock()
		delete(RoomsManager.Rooms, roomid)
		return true
	}
	return false
}
