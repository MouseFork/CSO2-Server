package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnToggleReady(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to toggle ready status but not in server !")
		return
	}
	//检查用户是不是房主
	curroom := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if curroom == nil ||
		curroom.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to toggle in a null room !")
		return
	}
	if curroom.HostUserID == uPtr.Userid {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to toggle but is host !")
		return
	}
	//检查用户所在房间
	if curroom.Id != uPtr.CurrentRoomId {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to toggle but in another room !")
		return
	}
	// if uPtr.currentIsIngame {
	// 	log.Println("Error : User", string(uPtr.username), "try to toggle but is ingame !")
	// 	return
	// }
	u := curroom.RoomGetUser(uPtr.Userid)
	if u == nil {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to toggle but in null in room !")
		return
	}
	//设置新的状态
	if uPtr.Currentstatus == UserReady {
		uPtr.SetUserStatus(UserNotReady)
		uPtr.SetUserIngame(false)
		DebugInfo(2, "User", string(uPtr.Username), "unreadied in room", string(curroom.Setting.RoomName), "id", curroom.Id)
	} else {
		uPtr.SetUserStatus(UserReady)
		DebugInfo(2, "User", string(uPtr.Username), "readied in room", string(curroom.Setting.RoomName), "id", curroom.Id)

	}
	//对房间所有玩家发送该玩家的状态
	ustatus := BuildUserReadyStatus(uPtr)
	for _, v := range curroom.Users {
		rst := BytesCombine(BuildHeader(v.CurrentSequence, p.Id), ustatus)
		SendPacket(rst, v.CurrentConnection)
	}
}

func BuildUserReadyStatus(u *User) []byte {
	buf := make([]byte, 6)
	offset := 0
	WriteUint8(&buf, OUTSetPlayerReady, &offset)
	WriteUint32(&buf, u.Userid, &offset)
	WriteUint8(&buf, u.Currentstatus, &offset)
	return buf[:offset]
}
