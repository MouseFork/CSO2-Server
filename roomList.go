package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

func onRoomList(seq *uint8, p *packet, client net.Conn) {
	var pkt roomsRequestPacket
	if praseChannelRequest(*p, &pkt) {
		uPtr := getUserFromConnection(client)
		if uPtr.userid <= 0 {
			log.Println("Error : A unknow Client from", client.RemoteAddr().String(), "request a RoomList !")
			return
		}
		//发送频道请求返回包
		chlsrv := getChannelServerWithID(pkt.channelServerIndex)
		if chlsrv == nil {
			log.Println("Error : Client from", client.RemoteAddr().String(), "request a unknown channelServer !")
			return
		}
		rst := BuildLobbyReply(seq, *p)
		sendPacket(rst, client)
		log.Println("Sent a lobbyReply packet to", client.RemoteAddr().String())
		//发送频道请求所得房间列表
		chl := getChannelWithID(pkt.channelIndex, *chlsrv)
		if chl == nil {
			log.Println("Error : Client from", client.RemoteAddr().String(), "request a unknown channel !")
			return
		}
		rst = BuildRoomList(seq, *p, *chl)
		sendPacket(rst, client)
		log.Println("Sent a roomList packet to", client.RemoteAddr().String())
		//设置用户所在频道
		uPtr.setUserChannelServer(chlsrv.serverIndex)
		uPtr.setUserChannel(chl.channelID)
	} else {
		log.Println("Recived a damaged packet from", client.RemoteAddr().String())
	}
}

func praseChannelRequest(p packet, dest *roomsRequestPacket) bool {
	if p.datalen-5 < 2 {
		return false
	}
	(*dest).channelServerIndex = p.data[5]
	(*dest).channelIndex = p.data[6]
	return true
}

func BuildLobbyReply(seq *uint8, p packet) []byte {
	p.id = TypeLobby
	rst := BuildHeader(seq, p)
	lob := lobbyJoinRoom{
		0, 2, 4,
	}
	rst = append(rst,
		JoinRoom,
		lob.unk00,
		lob.unk01,
		lob.unk02)
	WriteLen(&rst)
	return rst
}

//暂定
func BuildRoomList(seq *uint8, p packet, chl channelInfo) []byte {
	p.id = TypeRoomList
	rst := BuildHeader(seq, p)
	rst = append(rst,
		SendFullRoomList,
	)
	buf := make([]byte, 2)
	tempoffset := 0
	WriteUint16(&buf, chl.roomNum, &tempoffset)
	for i := 0; i < int(chl.roomNum); i++ {
		roombuf := make([]byte, 512)
		offset := 0
		WriteUint16(&roombuf, chl.rooms[i].id, &offset)
		WriteUint64(&roombuf, chl.rooms[i].flags, &offset)
		WriteString(&roombuf, chl.rooms[i].setting.roomName, &offset)
		WriteUint8(&roombuf, chl.rooms[i].roomNumber, &offset)
		WriteUint8(&roombuf, chl.rooms[i].passwordProtected, &offset)
		WriteUint16(&roombuf, 0, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.gameModeID, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.mapID, &offset)
		WriteUint8(&roombuf, chl.rooms[i].numPlayers, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.maxPlayers, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk08, &offset)
		WriteUint32(&roombuf, chl.rooms[i].hostUserID, &offset)
		WriteString(&roombuf, chl.rooms[i].hostUserName, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk11, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk12, &offset)
		WriteUint32(&roombuf, chl.rooms[i].unk13, &offset)
		WriteUint16(&roombuf, chl.rooms[i].unk14, &offset)
		WriteUint16(&roombuf, chl.rooms[i].unk15, &offset)
		WriteUint32(&roombuf, chl.rooms[i].unk16, &offset)
		WriteUint16(&roombuf, chl.rooms[i].unk17, &offset)
		WriteUint16(&roombuf, chl.rooms[i].unk18, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk19, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk20, &offset)
		if chl.rooms[i].unk20 == 1 {
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
		}
		WriteUint8(&roombuf, chl.rooms[i].unk21, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.status, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.areBotsEnabled, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk24, &offset)
		WriteUint16(&roombuf, chl.rooms[i].setting.startMoney, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk26, &offset)
		WriteUint8(&roombuf, 0, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk28, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk29, &offset)
		WriteUint64(&roombuf, chl.rooms[i].unk30, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.winLimit, &offset)
		WriteUint16(&roombuf, chl.rooms[i].setting.killLimit, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.forceCamera, &offset)
		// WriteUint8(&roombuf, chl.rooms[i].botEnabled, &offset)
		// if chl.rooms[i].botEnabled == 1 {
		// 	WriteUint8(&roombuf, chl.rooms[i].botDifficulty, &offset)
		// 	WriteUint8(&roombuf, chl.rooms[i].numCtBots, &offset)
		// 	WriteUint8(&roombuf, chl.rooms[i].numTrBots, &offset)
		// }
		WriteUint8(&roombuf, chl.rooms[i].unk31, &offset)
		WriteUint8(&roombuf, chl.rooms[i].unk35, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.nextMapEnabled, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.changeTeams, &offset)
		WriteUint8(&roombuf, chl.rooms[i].areFlashesDisabled, &offset)
		WriteUint8(&roombuf, chl.rooms[i].canSpec, &offset)
		WriteUint8(&roombuf, chl.rooms[i].isVipRoom, &offset)
		WriteUint8(&roombuf, chl.rooms[i].vipRoomLevel, &offset)
		WriteUint8(&roombuf, chl.rooms[i].setting.difficulty, &offset)
		buf = BytesCombine(buf, roombuf[:offset])
	}
	return BytesCombine(rst, buf)
}
