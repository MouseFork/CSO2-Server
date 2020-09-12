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
	for _, v := range chl.Rooms {
		if v == nil {
			DebugInfo(1, "Waring! here is a null room in channelID", chl.ChannelID)
			continue
		}
		// name, err := GbkToUtf8(chl.rooms[i].setting.roomName)
		// if err != nil {
		// 	continue
		// }
		roombuf := make([]byte, 512)
		offset := 0
		WriteUint16(&roombuf, v.Id, &offset)
		WriteUint64(&roombuf, 0XFFFFFFFFFFFFFFFF, &offset)
		WriteString(&roombuf, v.Setting.RoomName, &offset)
		WriteUint8(&roombuf, v.RoomNumber, &offset)
		WriteUint8(&roombuf, v.PasswordProtected, &offset)
		WriteUint16(&roombuf, 0, &offset)
		WriteUint8(&roombuf, v.Setting.GameModeID, &offset)
		WriteUint8(&roombuf, v.Setting.MapID, &offset)
		WriteUint8(&roombuf, v.NumPlayers, &offset)
		WriteUint8(&roombuf, v.Setting.MaxPlayers, &offset)
		WriteUint8(&roombuf, v.Unk08, &offset)
		WriteUint32(&roombuf, v.HostUserID, &offset)
		hostname, _ := GbkToUtf8(v.HostUserName)
		WriteString(&roombuf, hostname, &offset)
		WriteUint8(&roombuf, v.Unk11, &offset)
		WriteUint8(&roombuf, v.Unk12, &offset)
		WriteUint32(&roombuf, v.Unk13, &offset)
		WriteUint16(&roombuf, v.Unk14, &offset)
		WriteUint16(&roombuf, v.Unk15, &offset)
		WriteUint32(&roombuf, v.Unk16, &offset)
		WriteUint16(&roombuf, v.Unk17, &offset)
		WriteUint16(&roombuf, v.Unk18, &offset)
		WriteUint8(&roombuf, v.Unk19, &offset)
		WriteUint8(&roombuf, v.Unk20, &offset)
		if v.Unk20 == 1 {
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
		}
		WriteUint8(&roombuf, v.Unk21, &offset)
		WriteUint8(&roombuf, v.Setting.Status, &offset)
		WriteUint8(&roombuf, v.Setting.AreBotsEnabled, &offset)
		WriteUint8(&roombuf, v.Unk24, &offset)
		WriteUint16(&roombuf, v.Setting.StartMoney, &offset)
		WriteUint8(&roombuf, v.Unk26, &offset)
		WriteUint8(&roombuf, 0, &offset)
		WriteUint8(&roombuf, v.Unk28, &offset)
		WriteUint8(&roombuf, v.Unk29, &offset)
		WriteUint64(&roombuf, v.Unk30, &offset)
		WriteUint8(&roombuf, v.Setting.WinLimit, &offset)
		WriteUint16(&roombuf, v.Setting.KillLimit, &offset)
		WriteUint8(&roombuf, v.Setting.ForceCamera, &offset)
		// WriteUint8(&roombuf, v.botEnabled, &offset)
		// if v.botEnabled == 1 {
		// 	WriteUint8(&roombuf, v.botDifficulty, &offset)
		// 	WriteUint8(&roombuf, v.numCtBots, &offset)
		// 	WriteUint8(&roombuf, v.numTrBots, &offset)
		// }
		WriteUint8(&roombuf, v.Unk31, &offset)
		WriteUint8(&roombuf, v.Unk35, &offset)
		WriteUint8(&roombuf, v.Setting.NextMapEnabled, &offset)
		WriteUint8(&roombuf, v.Setting.ChangeTeams, &offset)
		WriteUint8(&roombuf, v.AreFlashesDisabled, &offset)
		WriteUint8(&roombuf, v.CanSpec, &offset)
		WriteUint8(&roombuf, v.IsVipRoom, &offset)
		WriteUint8(&roombuf, v.VipRoomLevel, &offset)
		WriteUint8(&roombuf, v.Setting.Difficulty, &offset)
		buf = BytesCombine(buf, roombuf[:offset])
	}
	return BytesCombine(rst, buf)
}
