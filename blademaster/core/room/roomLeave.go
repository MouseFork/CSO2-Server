package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnLeaveRoom(p *PacketData, client net.Conn) {
	//找到玩家
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to leave room but not in server !")
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to leave a null room !")
		return
	}
	//检查玩家游戏状态，准备情况下并且开始倒计时了，那么就不允许离开房间
	if uPtr.IsUserReady() &&
		rm.IsGlobalCountdownInProgress() {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to leave room but is started !")
		return
	}
	//房间移除玩家
	rm.RoomRemoveUser(uPtr.Userid)
	//检查房间是否为空
	if rm.NumPlayers <= 0 {
		DelChannelRoom(rm.Id,
			uPtr.GetUserChannelID(),
			uPtr.GetUserChannelServerID())

	} else {
		//向其他玩家发送离开信息
		SentUserLeaveMes(uPtr, rm)
	}
	//设置玩家状态
	p.Data = make([]byte, 3)
	p.Length = 3
	p.Data[1] = uPtr.GetUserChannelServerID()
	p.Data[2] = uPtr.GetUserChannelID()
	p.CurOffset = 1
	uPtr.QuitRoom()
	//房间状态
	rm.CheckIngameStatus()
	//发送房间列表给玩家
	OnRoomList(p, client)
	DebugInfo(2, "User", string(uPtr.Username), "left room", string(rm.Setting.RoomName), "id", rm.Id)
}
func SentUserLeaveMes(uPtr *User, rm *Room) {
	//如果玩家是房主
	for _, v := range rm.Users {
		rm.SetRoomHost(v)
		break
	}
	if rm.HostUserID == uPtr.Userid {
		for _, v := range rm.Users {
			rst1 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTPlayerLeave)
			rst1 = BytesCombine(rst1, BuildUserLeave(uPtr.Userid))
			rst2 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTSetHost)
			rst2 = BytesCombine(rst2, BuildSetHost(rm.HostUserID))
			SendPacket(rst1, v.CurrentConnection)
			SendPacket(rst2, v.CurrentConnection)
		}
		DebugInfo(2, "Sent a set roomHost packet to other users")
	} else {
		for _, v := range rm.Users {
			rst1 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTPlayerLeave)
			rst1 = BytesCombine(rst1, BuildUserLeave(uPtr.Userid))
			SendPacket(rst1, v.CurrentConnection)
		}
		DebugInfo(2, "Sent a leave room packet to other users")
	}
}
func BuildUserLeave(id uint32) []byte {
	buf := make([]byte, 4)
	offset := 0
	WriteUint32(&buf, id, &offset)
	return buf
}
func BuildSetHost(id uint32) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint32(&buf, id, &offset)
	buf[4] = 0
	return buf
}
