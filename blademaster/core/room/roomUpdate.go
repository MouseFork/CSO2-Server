package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

//onUpdateRoom 房主更新房间信息
func OnUpdateRoom(p *PacketData, client net.Conn) {
	//检索数据报
	var pkt InUpSettingReq
	if !p.PraseUpdateRoomPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a illegal packet !")
		return
	}
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
		DebugInfo(2, "Error : User:", string(uPtr.Username), "try to update a null room but in server !")
		return
	}
	if curroom.HostUserID != uPtr.Userid {
		DebugInfo(2, "Error : User:", string(uPtr.Username), "try to update a room but isn't host !")
		return
	}
	//检查用户所在房间
	if curroom.Id != uPtr.CurrentRoomId {
		DebugInfo(2, "Error : User:", string(uPtr.Username), "try to update a room but not in !")
		return
	}
	//检查当前是不是正在倒计时
	if curroom.IsGlobalCountdownInProgress() {
		DebugInfo(2, "Error : User:", string(uPtr.Username), "try to update a room but is counting !")
		return
	}
	//更新房间设置
	curroom.ToUpdateSetting(&pkt)
	//向房间所有玩家发送更新报文
	settingpkt := BuildRoomSetting(curroom)
	for _, v := range curroom.Users {
		rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), settingpkt)
		SendPacket(rst, v.CurrentConnection)
		//log.Println("["+strconv.Itoa(k+1)+"/"+strconv.Itoa(int((*curroom).numPlayers))+"] Updated room for", v.currentConnection.RemoteAddr().String(), "!")
	}
	curroom.Lastflags = curroom.Flags
	DebugInfo(2, "Host", string(uPtr.Username), "updated room", string(curroom.Setting.RoomName), "id", curroom.Id)
}
