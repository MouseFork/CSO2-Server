package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

//InNewRoomPacket 新建房间时传进来的数据包
type InNewRoomPacket struct {
	lenOfName  uint8
	roomName   []byte
	unk00      uint16
	unk01      uint8
	gameModeID uint8
	mapID      uint8
	winLimit   uint8
	killLimit  uint16
	unk02      uint8
	unk03      uint8
	unk04      uint8
	lenOfUnk05 uint8
	unk05      []byte
	unk06      uint8
	unk07      uint8
	unk08      uint8
	unk09      uint8
	unk10      uint8
	unk11      uint32
}

func onNewRoom(seq *uint8, p packet, client net.Conn) {
	//检索房间数据报
	var roompkt InNewRoomPacket
	if !praseNewRoomQuest(p, &roompkt) {
		log.Println("Error : Cannot prase a new room request !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : A user request a new room but not in server!")
		return
	}
	//检索玩家当前房间
	if uPtr.currentRoomId > 0 {
		log.Println("Error :", uPtr.username, "request a new room but already in a room!")
		uPtr.quitRoom()
		return
	}
	//创建房间
	rm := CreateRoom(roompkt, uPtr)
	if rm.id <= 0 {
		log.Println("Error :", uPtr.username, "cannot create a new room !")
		return
	}
	//修改用户相关信息
	rm.hostUserID = uPtr.userid
	rm.hostUserName = uPtr.username
	rm.users = append(rm.users, *uPtr)
	rm.numPlayers = 1
	u := rm.roomGetUser(uPtr.userid)
	if u == nil {
		log.Println("Error : Cannot add host ", uPtr.username, "to new room !")
		return
	}
	uPtr.setUserRoom(rm.id)
	u.setUserRoom(rm.id)
	uPtr.currentTeam = CounterTerrorist
	uPtr.setUserStatus(UserNotReady)
	u.currentTeam = CounterTerrorist
	u.setUserStatus(UserNotReady)
	//把房间加进服务器
	if !addChannelRoom(rm,
		uPtr.getUserChannelID(),
		uPtr.getUserChannelServerID()) {
		uPtr.quitRoom()
		return
	}
	//生成返回数据报
	p.id = TypeRoom
	rst := append(BuildHeader(seq, p), OUTCreateAndJoin)
	rst = BytesCombine(rst, buildCreateAndJoin(rm))
	sendPacket(rst, client)
	log.Println("Sent a new room packet to", string(u.username))
	//生成房间设置数据包
	rst = BytesCombine(BuildHeader(seq, p), buildRoomSetting(rm))
	sendPacket(rst, client)
	log.Println("Sent a room setting packet to", string(u.username))
	log.Println("User", string(uPtr.username), "created room", string(rm.setting.roomName), "id", rm.id)
}

func praseNewRoomQuest(p packet, roompkt *InNewRoomPacket) bool {
	if p.datalen < 21 {
		return false
	}
	offset := 6
	(*roompkt).lenOfName = ReadUint8(p.data, &offset)
	(*roompkt).roomName = ReadString(p.data, &offset, int((*roompkt).lenOfName))
	(*roompkt).unk00 = ReadUint16(p.data, &offset)
	(*roompkt).unk01 = ReadUint8(p.data, &offset)
	(*roompkt).gameModeID = ReadUint8(p.data, &offset)
	(*roompkt).mapID = ReadUint8(p.data, &offset)
	(*roompkt).winLimit = ReadUint8(p.data, &offset)
	(*roompkt).killLimit = ReadUint16(p.data, &offset)
	(*roompkt).unk02 = ReadUint8(p.data, &offset)
	(*roompkt).unk03 = ReadUint8(p.data, &offset)
	(*roompkt).unk04 = ReadUint8(p.data, &offset)
	(*roompkt).lenOfUnk05 = ReadUint8(p.data, &offset)
	(*roompkt).unk05 = ReadString(p.data, &offset, int((*roompkt).lenOfUnk05))
	(*roompkt).unk06 = ReadUint8(p.data, &offset)
	(*roompkt).unk07 = ReadUint8(p.data, &offset)
	(*roompkt).unk08 = ReadUint8(p.data, &offset)
	(*roompkt).unk09 = ReadUint8(p.data, &offset)
	(*roompkt).unk10 = ReadUint8(p.data, &offset)
	(*roompkt).unk11 = ReadUint32(p.data, &offset)
	return true
}

func buildCreateAndJoin(rm roomInfo) []byte {
	buf := make([]byte, 128+rm.setting.lenOfName)
	offset := 0
	WriteUint32(&buf, rm.hostUserID, &offset)
	WriteUint8(&buf, 2, &offset)
	WriteUint8(&buf, 2, &offset)
	WriteUint16(&buf, rm.id, &offset)
	WriteUint8(&buf, 5, &offset)
	// special class start?
	WriteUint64(&buf, 0xFFFFFFFFFFFFFFFF, &offset)
	WriteString(&buf, rm.setting.roomName, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint32(&buf, 0, &offset)
	WriteUint32(&buf, 0, &offset)
	//WriteString(&buf, []byte(""), &offset)
	WriteUint8(&buf, 0, &offset) //字符串长度为0
	WriteUint16(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, rm.setting.gameModeID, &offset)
	WriteUint8(&buf, rm.setting.mapID, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, rm.setting.winLimit, &offset)
	WriteUint16(&buf, rm.setting.killLimit, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, 0xA, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, rm.setting.status, &offset)
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
	if rm.setting.status == StatusIngame {
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
	WriteUint8(&buf, rm.numPlayers, &offset)
	buf = buf[:offset]
	for i := 0; i < int(rm.numPlayers); i++ {
		temp := make([]byte, 4)
		idx := 0
		WriteUint32(&temp, rm.users[i].userid, &idx)
		buf = BytesCombine(buf, temp,
			rm.users[i].buildUserNetInfo(),
			BuildUserInfo(newUserInfo(rm.users[i]), 0, false))
	}
	return buf
}

func CreateRoom(pkt InNewRoomPacket, u *user) roomInfo {
	var rm roomInfo
	srv := getChannelServerWithID(u.getUserChannelServerID())
	if srv.serverIndex <= 0 {
		rm.id = 0
		return rm
	}
	chl := getChannelWithID(u.getUserChannelID(), *srv)
	if chl.channelID <= 0 {
		rm.id = 0
		return rm
	}
	id := getNewRoomID(*chl)
	if id <= 0 {
		return rm
	}
	rm.id = id
	rm.roomNumber = uint8(rm.id)
	rm.flags = 0xFFFFFFFFFFFFFFFF
	rm.hostUserID = 0
	rm.users = []user{}
	rm.parentChannel = chl.channelID
	rm.countingDown = false
	rm.countdown = DefaultCountdownNum
	rm.numPlayers = 0
	rm.passwordProtected = 0
	rm.unk13 = 0xD73DA43D
	rm.unk14 = 0x9F31
	rm.unk15 = 0xB2B9
	rm.unk16 = 0xD73DA43D
	rm.unk17 = 0x9F31
	rm.unk18 = 0xB2B9
	rm.unk19 = 5
	rm.unk20 = 0
	rm.unk21 = 5
	rm.unk29 = 1
	rm.unk30 = 0x5AF6F7BF
	rm.unk31 = 4
	rm.unk35 = 0
	if u.isVIP() {
		rm.isVipRoom = 1
	} else {
		rm.isVipRoom = 0
	}
	rm.vipRoomLevel = u.vipLevel
	rm.setting.roomName = pkt.roomName
	rm.setting.gameModeID = pkt.gameModeID
	rm.setting.mapID = pkt.mapID
	rm.setting.winLimit = pkt.winLimit
	rm.setting.killLimit = pkt.killLimit
	rm.setting.startMoney = 16000
	rm.setting.forceCamera = 1
	rm.setting.nextMapEnabled = 0
	rm.setting.changeTeams = 0
	rm.setting.areBotsEnabled = 0 //false = 0,true = 1
	rm.setting.maxPlayers = 16    //enableBot = 8
	rm.setting.respawnTime = 3
	rm.setting.difficulty = 0
	rm.setting.teamBalanceType = 0
	rm.setting.weaponRestrictions = 0
	rm.setting.status = StatusWaiting
	rm.setting.hltvEnabled = 0
	rm.setting.mapCycleType = 1
	rm.setting.numCtBots = 0
	rm.setting.numTrBots = 0
	rm.setting.botDifficulty = 0
	rm.setting.isIngame = 0 //false = 0,true = 1
	return rm
}
