package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnChangeTeam(p *PacketData, client net.Conn) {
	var pkt InChangeTeamPacket
	if !p.PraseChangeTeamPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a illegal change team packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to change team but not in server !")
		return
	}
	//检查房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to change team in a null room !")
		uPtr.QuitRoom()
		return
	}
	//检查用户所在房间
	if rm.Id != uPtr.CurrentRoomId {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to change team but in another room !")
		return
	}
	//检查用户状态
	if uPtr.IsUserReady() {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to change team but is ready !")
		return
	}

	// //检查是否是房主
	// if uPtr.Userid != rm.HostUserID &&
	// 	rm.Setting.AreBotsEnabled != 0 {
	// 	OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, DialogBox,
	// 		GAME_ROOM_CHANGETEAM_FAILED)
	// 	DebugInfo(2, "Error : User", string(uPtr.Username), "try to change team in bot mode and isn't host !")
	// 	return
	// }

	//更换队伍
	uPtr.SetUserTeam(pkt.NewTeam)
	//发送数据包
	setteam := BuildChangTeam(uPtr.Userid, uPtr.CurrentTeam)
	for _, v := range rm.Users {
		//发送玩家更换队伍消息
		rst := BytesCombine(BuildHeader(v.CurrentSequence, p.Id), setteam)
		SendPacket(rst, v.CurrentConnection)
		// if rm.Setting.AreBotsEnabled != 0 {
		// 	v.SetUserTeam(pkt.NewTeam)
		// 	setteam := BuildChangTeam(v.Userid, pkt.NewTeam)
		// 	for _, k := range rm.Users {
		// 		rst = BytesCombine(BuildHeader(k.currentSequence, p), setteam)
		// 		sendPacket(rst, k.currentConnection)
		// 	}
		// }
	}
	DebugInfo(2, "User", string(uPtr.Username), "changed team to", uPtr.CurrentTeam, "in room", string(rm.Setting.RoomName), "id", rm.Id)
}

func BuildChangTeam(id uint32, team uint8) []byte {
	buf := make([]byte, 7)
	offset := 0
	WriteUint8(&buf, OUTsetUserTeam, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint32(&buf, id, &offset)
	WriteUint8(&buf, team, &offset)
	return buf[:offset]
}
