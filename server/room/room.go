package room

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/model/channel"
	. "github.com/KouKouChan/CSO2-Server/model/room"
	. "github.com/KouKouChan/CSO2-Server/server/channel"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func onRoomRequest(seq *uint8, p packet, client net.Conn) {
	var pkt inRoomPaket
	if praseRoomPacket(p, &pkt) {
		switch pkt.InRoomType {
		case NewRoomRequest:
			//log.Println("Recived a new room request from", client.RemoteAddr().String())
			onNewRoom(seq, p, client)
		case JoinRoomRequest:
			//log.Println("Recived a join room request from", client.RemoteAddr().String())
			onJoinRoom(seq, p, client)
		case LeaveRoomRequest:
			//log.Println("Recived a leave room request from", client.RemoteAddr().String())
			onLeaveRoom(seq, p, client)
		case ToggleReadyRequest:
			//log.Println("Recived a ready request from", client.RemoteAddr().String())
			onToggleReady(seq, p, client)
		case GameStartRequest:
			//log.Println("Recived a start game request from", client.RemoteAddr().String())
			onGameStart(seq, p, client)
		case UpdateSettings:
			//log.Println("Recived a update room setting request from", client.RemoteAddr().String())
			onUpdateRoom(seq, p, client)
		case OnCloseResultWindow:
			//log.Println("Recived a close resultWindow request from", client.RemoteAddr().String())
			onCloseResultRequest(seq, p, client)
		case SetUserTeamRequest:
			//log.Println("Recived a set user team request from", client.RemoteAddr().String())
			onChangeTeam(seq, p, client)
		case GameStartCountdownRequest:
			//log.Println("Recived a begin start game request from", client.RemoteAddr().String())
			onGameStartCountdown(p, client)
		default:
			log.Println("Unknown room packet", pkt.InRoomType, "from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal room packet from", client.RemoteAddr().String())
	}
}

func praseRoomPacket(p packet, dest *inRoomPaket) bool {
	if p.datalen-HeaderLen < 2 {
		return false
	}
	(*dest).InRoomType = p.data[5]
	return true
}

//getNewRoomNumber() 获取房间在某个频道下的标号
func GetNewRoomNumber(chl ChannelInfo) uint16 {
	if chl.RoomNum > MAXROOMNUMS {
		DebugInfo(2, "Error : Room is too much ! Unable to create more !")
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
		DebugInfo(2, "Error : Room is too much ! Unable to create more !")
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
