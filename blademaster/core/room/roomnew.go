package room

import (
	"net"
	"sync"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnNewRoom(p *PacketData, client net.Conn) {
	//检索房间数据报
	var roompkt InNewRoomPacket
	if !p.PraseNewRoomQuest(&roompkt) {
		DebugInfo(2, "Error : Cannot prase a new room request !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : A user request a new room but not in server!")
		return
	}
	//检索玩家当前房间
	if uPtr.CurrentRoomId > 0 {
		DebugInfo(2, "Error :", string(uPtr.Username), "request a new room but already in a room!")
		uPtr.QuitRoom()
		return
	}
	//创建房间
	rm := CreateRoom(roompkt, uPtr)
	if rm.Id <= 0 {
		DebugInfo(2, "Error :", string(uPtr.Username), "cannot create a new room !")
		return
	}
	//修改用户相关信息
	rm.HostUserID = uPtr.Userid
	rm.HostUserName = uPtr.Username
	rm.Users[rm.HostUserID] = uPtr
	rm.NumPlayers = 1

	// u := rm.RoomGetUser(uPtr.Userid)
	// if u == nil {
	// 	DebugInfo(2, "Error : Cannot add host ", string(uPtr.Username), "to new room !")
	// 	return
	// }

	uPtr.SetUserRoom(rm.Id)
	uPtr.CurrentTeam = UserForceCounterTerrorist
	uPtr.SetUserStatus(UserNotReady)

	//把房间加进服务器
	if !AddChannelRoom(&rm,
		uPtr.GetUserChannelID(),
		uPtr.GetUserChannelServerID()) {
		uPtr.QuitRoom()
		return
	}
	//生成返回数据报
	rst := append(BuildHeader(uPtr.CurrentSequence, PacketTypeRoom), OUTCreateAndJoin)
	rst = BytesCombine(rst, BuildCreateAndJoin(&rm))
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Sent a new room packet to", string(uPtr.Username))
	//生成房间设置数据包
	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, p.Id), BuildRoomSetting(&rm, 0XFFFFFFFFFFFFFFFF))
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Sent a room setting packet to", string(uPtr.Username))
}

func BuildCreateAndJoin(rm *Room) []byte {
	buf := make([]byte, 128+rm.Setting.LenOfName)
	offset := 0
	WriteUint32(&buf, rm.HostUserID, &offset)
	WriteUint8(&buf, 2, &offset)
	WriteUint8(&buf, 2, &offset)
	WriteUint16(&buf, rm.Id, &offset)
	WriteUint8(&buf, 5, &offset)
	// special class start?
	WriteUint64(&buf, 0xFFFFFFFFFFFFFFFF, &offset)
	WriteString(&buf, rm.Setting.RoomName, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint32(&buf, 0, &offset)
	WriteUint32(&buf, 0, &offset)
	//WriteString(&buf, []byte(""), &offset)
	WriteUint8(&buf, 0, &offset) //字符串长度为0
	WriteUint16(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, rm.Setting.GameModeID, &offset)
	WriteUint8(&buf, rm.Setting.MapID, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, rm.Setting.WinLimit, &offset)
	WriteUint16(&buf, rm.Setting.KillLimit, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, 0xA, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, rm.Setting.Status, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0x5A, &offset)
	WriteUint8(&buf, 0, &offset)
	// for i:=0 ;i< 0 ;i++ {
	// 	WriteUint8(&buf,rm.unk27[i],&offset)
	// }
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, 0, &offset)
	// if == 1, it can have 3 more bytes
	WriteUint8(&buf, 0, &offset)
	// if (this.botEnabled) { == 0
	// 	WriteUint8(&buf,this.botDifficulty,&offset)
	// 	WriteUint8(&buf,this.numCtBots,&offset)
	// 	WriteUint8(&buf,this.numTrBots,&offset)
	// }
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	if rm.Setting.Status == StatusIngame {
		WriteUint8(&buf, 1, &offset)
	} else {
		WriteUint8(&buf, 0, &offset)
	}
	WriteUint16(&buf, 0x3E80, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 3, &offset)
	// special class end?
	WriteUint8(&buf, rm.NumPlayers, &offset)
	buf = buf[:offset]
	for k, v := range rm.Users {
		temp := make([]byte, 4)
		idx := 0
		WriteUint32(&temp, k, &idx)
		buf = BytesCombine(buf, temp,
			v.BuildUserNetInfo(),
			BuildUserInfo(0xFFFFFFFF, NewUserInfo(v), 0, false))
	}
	return buf
}

func CreateRoom(pkt InNewRoomPacket, u *User) Room {
	var rm Room
	srv := GetChannelServerWithID(u.GetUserChannelServerID())
	if srv.ServerIndex <= 0 {
		rm.Id = 0
		return rm
	}
	chl := GetChannelWithID(u.GetUserChannelID(), *srv)
	if chl.ChannelID <= 0 {
		rm.Id = 0
		return rm
	}
	id := GetNewRoomID(*chl)
	if id <= 0 {
		return rm
	}
	rm.Id = id
	rm.RoomNumber = uint8(rm.Id)
	rm.HostUserID = 0
	rm.Users = map[uint32]*User{}
	rm.ParentChannel = chl.ChannelID
	rm.CountingDown = false
	rm.Countdown = DefaultCountdownNum
	rm.NumPlayers = 0
	rm.PasswordProtected = 0
	rm.Unk13 = 0xD73DA43D
	rm.Unk14 = 0x9F31
	rm.Unk15 = 0xB2B9
	rm.Unk16 = 0xD73DA43D
	rm.Unk17 = 0x9F31
	rm.Unk18 = 0xB2B9
	rm.Unk19 = 5
	rm.Unk20 = 0
	rm.Unk21 = 5
	rm.Unk29 = 1
	rm.Unk30 = 0x5AF6F7BF
	rm.Unk31 = 4
	rm.Unk35 = 0
	if u.IsVIP() {
		rm.IsVipRoom = 1
	} else {
		rm.IsVipRoom = 0
	}
	rm.VipRoomLevel = u.VipLevel
	rm.Setting.RoomName = pkt.RoomName
	rm.Setting.GameModeID = pkt.GameModeID
	rm.Setting.MapID = pkt.MapID
	rm.Setting.WinLimit = pkt.WinLimit
	rm.Setting.KillLimit = pkt.KillLimit
	rm.Setting.StartMoney = 16000
	rm.Setting.ForceCamera = 1
	rm.Setting.NextMapEnabled = 0
	rm.Setting.ChangeTeams = 0
	rm.Setting.AreBotsEnabled = 0 //false = 0,true = 1
	rm.Setting.MaxPlayers = 16    //enableBot = 8
	rm.Setting.RespawnTime = 3
	rm.Setting.Difficulty = 0
	rm.Setting.TeamBalanceType = 0
	rm.Setting.WeaponRestrictions = 0
	rm.Setting.Status = StatusWaiting
	rm.Setting.HltvEnabled = 0
	rm.Setting.MapCycleType = 1
	rm.Setting.NumCtBots = 0
	rm.Setting.NumTrBots = 0
	rm.Setting.BotDifficulty = 0
	rm.Setting.IsIngame = 0 //false = 0,true = 1
	var mutex sync.Mutex
	rm.RoomMutex = &mutex
	return rm
}
