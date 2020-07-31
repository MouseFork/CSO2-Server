package handler

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/server/packet"
	. "github.com/KouKouChan/CSO2-Server/server/room"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

type inJoinRoomPacket struct {
	roomId        uint16
	lenOfPassWord uint8
	passWord      []byte
}

func onJoinRoom(seq *uint8, p Packet, client net.Conn) {
	//检索数据包
	var pkt inJoinRoomPacket
	if !praseJoinRoomPacket(p, &pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error JoinRoom packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to join room but not in server !")
		return
	}
	//检索玩家房间
	if uPtr.CurrentRoomId != 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to join room but in another room !")
		return
	}
	//找到对应房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uint16(pkt.roomId))
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to join a null room !")
		return
	}
	//检索密码
	if rm.Setting.LenOfPassWd > 0 {
		if !CompareBytes(pkt.passWord, rm.Setting.PassWd) {
			onSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, DialogBox,
				GAME_ROOM_JOIN_FAILED_BAD_PASSWORD)
			DebugInfo(2, "User", string(uPtr.Username), "try to join a room with error password!")
			return
		}
	}
	//检索房间状态
	if rm.GetFreeSlots() <= 0 {
		DebugInfo(2, "User", string(uPtr.Username), "try to join a full room !")
		return
	}
	//玩家加进房间
	if !rm.JoinUser(uPtr) {
		return
	}
	//发送数据
	p.Id = TypeRoom
	rst := append(BuildHeader(uPtr.CurrentSequence, p), OUTCreateAndJoin)
	rst = BytesCombine(rst, buildCreateAndJoin(*rm))
	SendPacket(rst, client)
	DebugInfo(2, "User", string(uPtr.Username), "joined room", string(rm.Setting.RoomName), "id", rm.Id)
	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, p), buildRoomSetting(*rm))
	SendPacket(rst, client)
	DebugInfo(2, "Sent a room setting packet to", string(uPtr.username))
	//发送玩家状态
	for _, v := range rm.UserIDs {
		u := GetUserFromID(v)
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, p), BuildUserReadyStatus(u))
		sendPacket(rst, uPtr.CurrentConnection)
		if u.Userid != uPtr.Userid {
			//发送给其他玩家该玩家信息
			rst = BytesCombine(BuildHeader(u.CurrentSequence, p), BuildPlayerJoin(*uPtr))
			sendPacket(rst, u.currentConnection)
			rst = BytesCombine(BuildHeader(u.CurrentSequence, p), BuildUserReadyStatus(*uPtr))
			sendPacket(rst, u.CurrentConnection)
		}
	}
	DebugInfo(2, "Sync user status to all player in room", string(rm.setting.roomName), "id", rm.id)
}

func praseJoinRoomPacket(p packet, dest *inJoinRoomPacket) bool {
	if p.datalen < 7 {
		return false
	}
	offset := 6
	(*dest).roomId = ReadUint16(p.data, &offset)
	(*dest).lenOfPassWord = ReadUint8(p.data, &offset)
	(*dest).passWord = ReadString(p.data, &offset, int((*dest).lenOfPassWord))
	return true
}

func BuildPlayerJoin(u user) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint8(&buf, OUTPlayerJoin, &offset)
	WriteUint32(&buf, u.userid, &offset)
	buf = BytesCombine(buf, u.buildUserNetInfo(),
		BuildUserInfo(newUserInfo(u), 0, false))
	return buf
}
