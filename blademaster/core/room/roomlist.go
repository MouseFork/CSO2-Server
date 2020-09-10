package room

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnRoomList(p *PacketData, client net.Conn) {
	var pkt InRoomListRequestPacket
	if p.PraseChannelRequest(&pkt) {
		uPtr := GetUserFromConnection(client)
		if uPtr.Userid <= 0 {
			DebugInfo(2, "Error : A unknow Client from", client.RemoteAddr().String(), "request a RoomList !")
			return
		}

		//发送频道请求返回包
		chlsrv := GetChannelServerWithID(pkt.ChannelServerIndex)
		if chlsrv == nil {
			DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "request a unknown channelServer !")
			return
		}
		rst := BuildLobbyReply(uPtr.CurrentSequence, *p)
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(2, "Sent a lobbyReply packet to", client.RemoteAddr().String())

		//发送频道请求所得房间列表
		chl := GetChannelWithID(pkt.ChannelIndex, *chlsrv)
		if chl == nil {
			log.Println("Error : Client from", client.RemoteAddr().String(), "request a unknown channel !")
			return
		}
		rst = BuildRoomList(uPtr.CurrentSequence, *p, *chl)
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(2, "Sent a roomList packet to", client.RemoteAddr().String())

		//设置用户所在频道
		uPtr.SetUserChannelServer(chlsrv.ServerIndex)
		uPtr.SetUserChannel(chl.ChannelID)
	} else {
		log.Println("Recived a damaged packet from", client.RemoteAddr().String())
	}
}

func BuildLobbyReply(seq *uint8, p PacketData) []byte {
	rst := BuildHeader(seq, PacketTypeLobby)
	lob := OutLobbyJoinRoom{
		0, 2, 4,
	}
	rst = append(rst,
		JoinRoom,
		lob.Unk00,
		lob.Unk01,
		lob.Unk02)
	WriteLen(&rst)
	return rst
}

//暂定
func BuildRoomList(seq *uint8, p PacketData, chl ChannelInfo) []byte {
	rst := BuildHeader(seq, PacketTypeRoomList)
	rst = append(rst,
		SendFullRoomList,
	)
	buf := make([]byte, 2)
	tempoffset := 0
	WriteUint16(&buf, chl.RoomNum, &tempoffset)
	var i uint16
	for i = 0; i < chl.RoomNum; i++ {
		// name, err := GbkToUtf8(chl.rooms[i].setting.roomName)
		// if err != nil {
		// 	continue
		// }
		roombuf := make([]byte, 512)
		offset := 0
		WriteUint16(&roombuf, chl.Rooms[i].Id, &offset)
		WriteUint64(&roombuf, chl.Rooms[i].Flags, &offset)
		WriteString(&roombuf, chl.Rooms[i].Setting.RoomName, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].RoomNumber, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].PasswordProtected, &offset)
		WriteUint16(&roombuf, 0, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.GameModeID, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.MapID, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].NumPlayers, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.MaxPlayers, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk08, &offset)
		WriteUint32(&roombuf, chl.Rooms[i].HostUserID, &offset)
		hostname, _ := GbkToUtf8(chl.Rooms[i].HostUserName)
		WriteString(&roombuf, hostname, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk11, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk12, &offset)
		WriteUint32(&roombuf, chl.Rooms[i].Unk13, &offset)
		WriteUint16(&roombuf, chl.Rooms[i].Unk14, &offset)
		WriteUint16(&roombuf, chl.Rooms[i].Unk15, &offset)
		WriteUint32(&roombuf, chl.Rooms[i].Unk16, &offset)
		WriteUint16(&roombuf, chl.Rooms[i].Unk17, &offset)
		WriteUint16(&roombuf, chl.Rooms[i].Unk18, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk19, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk20, &offset)
		if chl.Rooms[i].Unk20 == 1 {
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
		}
		WriteUint8(&roombuf, chl.Rooms[i].Unk21, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.Status, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.AreBotsEnabled, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk24, &offset)
		WriteUint16(&roombuf, chl.Rooms[i].Setting.StartMoney, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk26, &offset)
		WriteUint8(&roombuf, 0, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk28, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk29, &offset)
		WriteUint64(&roombuf, chl.Rooms[i].Unk30, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.WinLimit, &offset)
		WriteUint16(&roombuf, chl.Rooms[i].Setting.KillLimit, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.ForceCamera, &offset)
		// WriteUint8(&roombuf, chl.rooms[i].botEnabled, &offset)
		// if chl.rooms[i].botEnabled == 1 {
		// 	WriteUint8(&roombuf, chl.rooms[i].botDifficulty, &offset)
		// 	WriteUint8(&roombuf, chl.rooms[i].numCtBots, &offset)
		// 	WriteUint8(&roombuf, chl.rooms[i].numTrBots, &offset)
		// }
		WriteUint8(&roombuf, chl.Rooms[i].Unk31, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Unk35, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.NextMapEnabled, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.ChangeTeams, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].AreFlashesDisabled, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].CanSpec, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].IsVipRoom, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].VipRoomLevel, &offset)
		WriteUint8(&roombuf, chl.Rooms[i].Setting.Difficulty, &offset)
		buf = BytesCombine(buf, roombuf[:offset])
	}
	return BytesCombine(rst, buf)
}
