package handler

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func onToggleReady(seq *uint8, p packet, client net.Conn) {
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to toggle ready status but not in server !")
		return
	}
	//检查用户是不是房主
	curroom := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if curroom == nil ||
		curroom.id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to toggle in a null room !")
		return
	}
	if curroom.hostUserID == uPtr.userid {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to toggle but is host !")
		return
	}
	//检查用户所在房间
	if curroom.id != uPtr.currentRoomId {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to toggle but in another room !")
		return
	}
	// if uPtr.currentIsIngame {
	// 	log.Println("Error : User", string(uPtr.username), "try to toggle but is ingame !")
	// 	return
	// }
	u := curroom.roomGetUser(uPtr.userid)
	if u == nil {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to toggle but in null in room !")
		return
	}
	//设置新的状态
	if uPtr.currentstatus == UserNotReady {
		uPtr.setUserStatus(UserReady)
		u.setUserStatus(UserReady)
		DebugInfo(2, "User", string(uPtr.username), "unreadied in room", string(curroom.setting.roomName), "id", curroom.id)
	} else {
		uPtr.setUserStatus(UserNotReady)
		u.setUserStatus(UserNotReady)
		DebugInfo(2, "User", string(uPtr.username), "readied in room", string(curroom.setting.roomName), "id", curroom.id)
	}
	//对房间所有玩家发送该玩家的状态
	for _, v := range curroom.users {
		rst := BytesCombine(BuildHeader(v.currentSequence, p), BuildUserReadyStatus(*uPtr))
		sendPacket(rst, v.currentConnection)
	}
}

func BuildUserReadyStatus(u user) []byte {
	buf := make([]byte, 6)
	offset := 0
	WriteUint8(&buf, OUTSetPlayerReady, &offset)
	WriteUint32(&buf, u.userid, &offset)
	WriteUint8(&buf, u.currentstatus, &offset)
	return buf[:offset]
}
