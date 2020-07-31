package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type inHostTeamChangingPacket struct {
	userId uint32
	unk00  uint8
	//unk01   uint8
	newTeam uint8
}

func onChangingTeam(p packet, client net.Conn) {
	//检索数据包
	var pkt inHostTeamChangingPacket
	if !praseInTeamChangingPacket(p, &pkt) {
		log.Println("Error : User from", client.RemoteAddr().String(), "sent a error TeamChanging packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : User from", client.RemoteAddr().String(), "sent TeamChanging but not in server !")
		return
	}
	destUser := getUserFromID(uint32(pkt.userId))
	if destUser == nil ||
		destUser.userid <= 0 {
		log.Println("Error : User from", client.RemoteAddr().String(), "sent TeamChanging but dester is not in server !")
		return
	}
	//找到房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.currentRoomId)
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User from", client.RemoteAddr().String(), "sent TeamChanging but is not host !")
		return
	}
	//更新数据
	//log.Println(p.data)
	(*destUser).currentTeam = pkt.newTeam
	result := BuildChangingTeam(destUser.userid, pkt.newTeam)
	for _, v := range rm.users {
		u := getUserFromID(v.userid)
		if u == nil ||
			u.userid <= 0 {
			continue
		}
		rst := BytesCombine(BuildHeader(v.currentSequence, p), result)
		sendPacket(rst, v.currentConnection)
		// if rm.setting.areBotsEnabled != 0 {
		// 	//(*u).currentTeam = pkt.newTeam
		// 	//log.Println("User", string(u.username), "changed team", u.currentTeam)
		// 	//(*rm).users[k].currentTeam = pkt.newTeam
		// 	//向每名玩家表示该玩家更换了队伍
		// 	rst1 := BuildChangingTeam(u.userid, pkt.newTeam)
		// 	for _, i := range rm.users {
		// 		rst = BytesCombine(BuildHeader(i.currentSequence, p), rst1)
		// 		sendPacket(rst, i.currentConnection)
		// 	}
		// }
	}
}

func praseInTeamChangingPacket(p packet, dest *inHostTeamChangingPacket) bool {
	if p.datalen < 12 ||
		dest == nil {
		return false
	}
	offset := 6
	(*dest).userId = ReadUint32(p.data, &offset)
	(*dest).newTeam = ReadUint8(p.data, &offset)
	(*dest).unk00 = ReadUint8(p.data, &offset)
	return true
}

func BuildChangingTeam(id uint32, team uint8) []byte {
	buf := make([]byte, 7)
	offset := 0
	WriteUint8(&buf, OUTsetUserTeam, &offset)
	WriteUint8(&buf, 1, &offset) //numberOfUser
	WriteUint32(&buf, id, &offset)
	WriteUint8(&buf, team, &offset)
	return buf[:offset]
}
