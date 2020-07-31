package room

import (
	"log"

	. "github.com/KouKouChan/CSO2-Server/model/channel"
	. "github.com/KouKouChan/CSO2-Server/model/room"
	. "github.com/KouKouChan/CSO2-Server/server/channel"
)

//getNewRoomNumber() 获取房间在某个频道下的标号
func GetNewRoomNumber(chl ChannelInfo) uint16 {
	if chl.RoomNum > MAXROOMNUMS {
		log.Println("Error : Room is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXROOMNUMS + 2]uint16
	//哈希思想
	for i := 0; i < int(chl.RoomNum); i++ {
		intbuf[chl.Rooms[i].Id] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXROOMNUMS+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint16(i)
		}
	}
	return 0
}

//getNewRoomNumber() 获取房间在某个频道下的标号
func GetNewRoomID(chl ChannelInfo) uint16 {
	if chl.RoomNum > MAXROOMNUMS {
		log.Println("Error : Room is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXROOMNUMS + 2]uint16
	//哈希思想
	for i := 0; i < int(chl.RoomNum); i++ {
		intbuf[chl.Rooms[i].Id] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXROOMNUMS+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint16(i)
		}
	}
	return 0
}

func GetRoomFromID(chlsrvID uint8, chlID uint8, roomID uint16) *RoomInfo {
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
	for k, v := range chl.Rooms {
		if v.Id == roomID {
			return &chl.Rooms[k]
		}
	}
	return nil
}
