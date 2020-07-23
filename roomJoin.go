package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inJoinRoomPacket struct {
	roomId        uint16
	lenOfPassWord uint8
	passWord      []byte
}

func onJoinRoom(seq *uint8, p packet, client net.Conn) {
	//检索数据包
	var pkt inJoinRoomPacket
	if !praseJoinRoomPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a error JoinRoom packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to join room but not in server !")
		return
	}
	//检索玩家房间
	if uPtr.currentRoomId != 0 {
		log.Println("Error : User", string(uPtr.username), "try to join room but in another room !")
		return
	}
	//找到对应房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uint16(pkt.roomId))
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User", string(uPtr.username), "try to join a null room !")
		return
	}
	//检索密码
	if rm.setting.lenOfPassWd > 0 {
		if !CompareBytes(pkt.passWord, rm.setting.PassWd) {
			onSendMessage(uPtr.currentSequence, uPtr.currentConnection, DialogBox,
				GAME_ROOM_JOIN_FAILED_BAD_PASSWORD)
			log.Println("User", string(uPtr.username), "try to join a room with error password!")
			return
		}
	}
	//检索房间状态
	if rm.getFreeSlots() <= 0 {
		log.Println("User", string(uPtr.username), "try to join a full room !")
		return
	}
	//玩家加进房间
	if !rm.joinUser(uPtr) {
		return
	}
	//发送数据
	p.id = TypeRoom
	rst := append(BuildHeader(uPtr.currentSequence, p), OUTCreateAndJoin)
	rst = BytesCombine(rst, buildCreateAndJoin(*rm))
	sendPacket(rst, client)
	log.Println("User", string(uPtr.username), "joined room", string(rm.setting.roomName), "id", rm.id)
	rst = BytesCombine(BuildHeader(uPtr.currentSequence, p), buildRoomSetting(*rm))
	sendPacket(rst, client)
	log.Println("Sent a room setting packet to", string(uPtr.username))
	//发送玩家状态
	for _, v := range rm.users {
		rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildUserReadyStatus(v))
		sendPacket(rst, uPtr.currentConnection)
		if v.userid != uPtr.userid {
			//发送给其他玩家该玩家信息
			rst = BytesCombine(BuildHeader(v.currentSequence, p), BuildPlayerJoin(*uPtr))
			sendPacket(rst, v.currentConnection)
			rst = BytesCombine(BuildHeader(v.currentSequence, p), BuildUserReadyStatus(*uPtr))
			sendPacket(rst, v.currentConnection)
		}
	}
	log.Println("Sync user status to all player in room", string(rm.setting.roomName), "id", rm.id)
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
