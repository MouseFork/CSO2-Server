package handler

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

type inChangeTeamPacket struct {
	newTeam uint8
}

func onChangeTeam(seq *uint8, p packet, client net.Conn) {
	var pkt inChangeTeamPacket
	if !praseChangeTeamPacket(p, &pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a illegal change team packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to change team but not in server !")
		return
	}
	//检查房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to change team in a null room !")
		uPtr.quitRoom()
		return
	}
	u := rm.roomGetUser(uPtr.userid)
	if u == nil ||
		u.userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to change team but not in server !")
		return
	}
	//检查用户所在房间
	if rm.id != uPtr.currentRoomId {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to change team but in another room !")
		return
	}
	//检查用户状态
	if uPtr.isUserReady() {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to change team but is ready !")
		return
	}
	//检查是否是房主
	if uPtr.userid != rm.hostUserID &&
		rm.setting.areBotsEnabled != 0 {
		DebugInfo(2, "Error : User", string(uPtr.username), "try to change team in bot mode and isn't host !")
		return
	}
	//更换队伍
	(*uPtr).currentTeam = pkt.newTeam
	(*u).currentTeam = pkt.newTeam
	//发送数据包
	setteam := BuildChangTeam(u.userid, pkt.newTeam)
	for i, v := range rm.users {
		rst := BytesCombine(BuildHeader(v.currentSequence, p), setteam)
		sendPacket(rst, v.currentConnection)
		if rm.setting.areBotsEnabled != 0 {
			tempu := getUserFromID(v.userid)
			if tempu == nil ||
				tempu.userid <= 0 {
				continue
			}
			(*tempu).currentTeam = pkt.newTeam
			(*rm).users[i].currentTeam = pkt.newTeam
			setteam := BuildChangTeam(v.userid, pkt.newTeam)
			for _, k := range rm.users {
				rst = BytesCombine(BuildHeader(k.currentSequence, p), setteam)
				sendPacket(rst, k.currentConnection)
			}
		}
	}
	DebugInfo(2, "User", string(uPtr.username), "changed team to", pkt.newTeam, "in room", string(rm.setting.roomName), "id", rm.id)
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

func praseChangeTeamPacket(p packet, dest *inChangeTeamPacket) bool {
	if p.datalen < 7 {
		return false
	}
	offset := 6
	(*dest).newTeam = ReadUint8(p.data, &offset)
	return true
}
