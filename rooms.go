package main

import (
	"log"
	"net"
	"strconv"
	"unsafe"
)

const (
	//频道以及房间
	SendFullRoomList = 0
	JoinRoom         = 1
	UpdateUserInfo   = 2

	//房间操作
	NewRoomRequest            = 0
	JoinRoomRequest           = 1
	LeaveRoomRequest          = 3
	ToggleReadyRequest        = 4
	GameStartRequest          = 5
	UpdateSettings            = 6
	OnCloseResultWindow       = 7
	SetUserTeamRequest        = 9
	GameStartCountdownRequest = 19

	//游戏模式
	original            = 1
	teamdeath           = 2
	zombie              = 3
	stealth             = 4
	gunteamdeath        = 5
	tutorial            = 6
	hide                = 7
	pig                 = 8
	animationtest_vcd   = 9
	gz_survivor         = 10
	devtest             = 11
	originalmr          = 12
	originalmrdraw      = 13
	casualbomb          = 14
	deathmatch          = 15
	scenario_test       = 16
	gz                  = 17
	gz_intro            = 18
	gz_tour             = 19
	gz_pve              = 20
	eventmod01          = 21
	duel                = 22
	gz_ZB               = 23
	heroes              = 24
	eventmod02          = 25
	zombiecraft         = 26
	campaign1           = 27
	campaign2           = 28
	campaign3           = 29
	campaign4           = 30
	campaign5           = 31
	campaign6           = 32
	campaign7           = 33
	campaign8           = 34
	campaign9           = 35
	z_scenario          = 36
	zombie_prop         = 37
	ghost               = 38
	tag                 = 39
	hide_match          = 40
	hide_ice            = 41
	diy                 = 42
	hide_Item           = 43
	zd_boss1            = 44
	zd_boss2            = 45
	zd_boss3            = 46
	practice            = 47
	zombie_commander    = 48
	casualoriginal      = 49
	hide2               = 50
	gunball             = 51
	zombie_zeta         = 53
	tdm_small           = 54
	de_small            = 55
	gunteamdeath_re     = 56
	endless_wave        = 57
	rankmatch_original  = 58
	rankmatch_teamdeath = 59
	play_ground         = 60
	madcity             = 61
	hide_origin         = 62
	teamdeath_mutation  = 63
	giant               = 64
	z_scenario_side     = 65
	hide_multi          = 66
	madcity_team        = 67
	rankmatch_stealth   = 68

	//阵营
	Unknown          = 0
	Terrorist        = 1
	CounterTerrorist = 2

	//房间status
	StatusWaiting = 1
	StatusIngame  = 2

	//队伍平衡
	Disabled   = 0
	Enabled    = 1
	WithBots   = 2
	ByKadRatio = 4

	//房间包表示
	OUTCreateAndJoin  = 0
	OUTPlayerJoin     = 1
	OUTPlayerLeave    = 2
	OUTSetPlayerReady = 3
	OUTUpdateSettings = 4
	OUTSetHost        = 5
	OUTSetGameResult  = 6
	OUTsetUserTeam    = 7
	OUTCountdown      = 14

	//最大房间数
	MAXROOMNUMS         = 1024
	DefaultCountdownNum = 7
)

//房间信息,用于请求频道
type roomInfo struct {
	id    uint16
	flags uint64
	//roomName          []byte
	roomNumber        uint8
	passwordProtected uint8
	//unk03   = roomid          uint16
	// gameModeID        uint8
	// mapID             uint8
	//maxPlayers uint8
	unk08        uint8
	hostUserID   uint32
	hostUserName []byte
	unk11        uint8
	unk12        uint8
	unk13        uint32
	unk14        uint16
	unk15        uint16
	unk16        uint32
	unk17        uint16
	unk18        uint16
	unk19        uint8
	unk20        uint8
	unk21        uint8
	// roomStatus   uint8
	// enableBots   uint8
	unk24 uint8
	// startMoney   uint16
	unk26 uint8
	unk27 []uint8
	unk28 uint8
	unk29 uint8
	unk30 uint64
	// winLimit          uint8
	// killLimit         uint16
	// forceCamera    uint8
	// botEnabled     uint8
	// botDifficulty  uint8
	// numCtBots      uint8
	// numTrBots      uint8
	unk31 uint8
	unk35 uint8
	// nextMapEnabled uint8
	// changeTeams    uint8
	areFlashesDisabled uint8
	canSpec            uint8
	isVipRoom          uint8
	vipRoomLevel       uint8
	// difficulty     uint8

	//设置
	setting       roomSettings
	countingDown  bool
	countdown     uint8
	numPlayers    uint8
	users         []user
	parentChannel uint8
}

//每个房间的设置数据
type roomSettings struct {
	lenOfName          uint8
	roomName           []byte
	unk00              uint8
	unk01              uint8
	unk02              uint32
	unk03              uint32
	lenOfunk09         uint8
	unk09              []byte
	unk10              uint16
	forceCamera        uint8
	gameModeID         uint8
	mapID              uint8
	unk13              uint8
	maxPlayers         uint8
	winLimit           uint8
	killLimit          uint16
	unk17              uint8
	unk18              uint8
	weaponRestrictions uint8
	status             uint8
	unk21              uint8
	mapCycleType       uint8
	unk23              uint8
	unk24              uint8
	unk25              uint8
	lenOfMultiMaps     uint8
	multiMaps          []byte
	teamBalanceType    uint8
	unk29              uint8
	unk30              uint8
	unk31              uint8
	unk32              uint8
	unk33              uint8
	areBotsEnabled     uint8
	botDifficulty      uint8
	numCtBots          uint8
	numTrBots          uint8
	unk35              uint8
	unk36              uint8
	unk37              uint8
	unk38              uint8
	unk39              uint8
	startMoney         uint16
	changeTeams        uint8
	unk43              uint8
	hltvEnabled        uint8
	unk45              uint8
	respawnTime        uint8
	nextMapEnabled     uint8
	difficulty         uint8
	isIngame           uint8
}

//InRoomPacket 新建房间时传进来的数据包
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

//房间请求
type inRoomPaket struct {
	InRoomType uint8
}

//房间所属频道，用于请求频道
type roomsRequestPacket struct {
	channelServerIndex uint8
	channelIndex       uint8
}

//未知，用于请求频道
type lobbyJoinRoom struct {
	unk00 uint8
	unk01 uint8
	unk02 uint8
}

func onRoomRequest(seq *uint8, p packet, client net.Conn) {
	var pkt inRoomPaket
	if praseRoomPacket(p, &pkt) {
		getUserFromConnection(client)
		switch pkt.InRoomType {
		case NewRoomRequest:
			log.Println("Recived a new room request from", client.RemoteAddr().String())
			onNewRoom(seq, p, client)
		case JoinRoomRequest:
			log.Println("Recived a join room request from", client.RemoteAddr().String())
		case LeaveRoomRequest:
			log.Println("Recived a leave room request from", client.RemoteAddr().String())
		case ToggleReadyRequest:
			log.Println("Recived a ready request from", client.RemoteAddr().String())
		case GameStartRequest:
			log.Println("Recived a start game request from", client.RemoteAddr().String())
		case UpdateSettings:
			log.Println("Recived a update room setting request from", client.RemoteAddr().String())
			onUpdateRoom(seq, p, client)
		case OnCloseResultWindow:
			log.Println("Recived a close resultWindow request from", client.RemoteAddr().String())
		case SetUserTeamRequest:
			log.Println("Recived a set user team request from", client.RemoteAddr().String())
		case GameStartCountdownRequest:
			log.Println("Recived a begin start game request from", client.RemoteAddr().String())
		default:
			log.Println("Recived a unknown room packet from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Recived a illegal room packet from", client.RemoteAddr().String())
	}
}

func praseRoomPacket(p packet, dest *inRoomPaket) bool {
	if p.datalen-HeaderLen < 2 {
		return false
	}
	(*dest).InRoomType = p.data[5]
	return true
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

func onNewRoom(seq *uint8, p packet, client net.Conn) {
	//检索房间数据报
	var roompkt InNewRoomPacket
	if !praseNewRoomQuest(p, &roompkt) {
		log.Println("Cannot prase a new room request !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr.userid <= 0 {
		log.Println("A user request a new room but not in server!")
		return
	}
	//检索玩家当前房间
	if uPtr.currentRoomId > 0 {
		log.Println("A user request a new room but already in a room!")
		uPtr.quitRoom()
		return
	}
	//创建房间
	rm := CreateRoom(roompkt, uPtr)
	if rm.id <= 0 {
		log.Println("Cannot create a new room !")
		return
	}
	//把房间加进服务器
	if !addChannelRoom(rm,
		uPtr.getUserChannelID(),
		uPtr.getUserChannelServerID()) {
		return
	}
	//修改用户相关信息
	uPtr.setUserRoom(rm.id)
	//生成返回数据报
	p.id = TypeRoom
	rst := append(BuildHeader(seq, p), OUTCreateAndJoin)
	rst = BytesCombine(rst, buildCreateAndJoin(rm))
	sendPacket(rst, client)
	log.Println("Sent a new room packet to", client.RemoteAddr().String())
	//生成房间设置数据包
	rst = BytesCombine(BuildHeader(seq, p), buildRoomSetting(rm))
	sendPacket(rst, client)
	log.Println("Sent a room setting packet to", client.RemoteAddr().String())
}

//创建新房间数据包
// func buildNewRoom(seq *uint8, p packet) []byte {
// 	var buf []byte
// 	p.id = TypeRoom
// 	buf = BytesCombine(BuildHeader(seq, p), []byte{OUTCreateAndJoin})

// 	return buf
// }
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

	WriteString(&buf, []byte(""), &offset)

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

func (u user) buildUserNetInfo() []byte {
	buf := make([]byte, 25)
	offset := 0
	WriteUint8(&buf, u.getUserTeam(), &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	cliadr := u.currentConnection.RemoteAddr().String()
	externalIPAddress, err := IPToUint32(cliadr[:slideIP(cliadr)])
	if err != nil {
		log.Fatalln("Prasing externalIpAddress error !")
		return buf
	}
	WriteUint32BE(&buf, externalIPAddress, &offset) //externalIpAddress
	WriteUint16(&buf, 0, &offset)                   //externalServerPort
	WriteUint16(&buf, 0, &offset)                   //externalClientPort
	WriteUint16(&buf, 0, &offset)                   //externalTvPort
	WriteUint32BE(&buf, 0, &offset)                 //localIpAddress
	WriteUint16(&buf, 0, &offset)                   //localServerPort
	WriteUint16(&buf, 0, &offset)                   //localClientPort
	WriteUint16(&buf, 0, &offset)                   //localTvPort
	return buf[:offset]
}

//创建房间设置数据包
func buildRoomSetting(room roomInfo) []byte {
	buf := make([]byte, 128+room.setting.lenOfName+ //实际计算是最大63字节+长度
		room.setting.lenOfunk09+
		room.setting.lenOfMultiMaps)
	offset := 0
	WriteUint8(&buf, OUTUpdateSettings, &offset)
	room.flags = getFlags(room)
	WriteUint64(&buf, room.flags, &offset)
	lowFlag := *(*uint32)(unsafe.Pointer(&room.flags))
	flags := room.flags >> 32
	highFlag := *(*uint32)(unsafe.Pointer(&flags))

	if lowFlag&0x1 != 0 {
		WriteString(&buf, room.setting.roomName, &offset)
	}
	if lowFlag&0x2 != 0 {
		WriteUint8(&buf, room.setting.unk00, &offset)
	}
	if lowFlag&0x4 != 0 {
		WriteUint8(&buf, room.setting.unk01, &offset)
		WriteUint32(&buf, room.setting.unk02, &offset)
		WriteUint32(&buf, room.setting.unk03, &offset)
	}
	if lowFlag&0x8 != 0 {
		WriteString(&buf, room.setting.unk09, &offset)
	}
	if lowFlag&0x10 != 0 {
		WriteUint16(&buf, room.setting.unk10, &offset)
	}
	if lowFlag&0x20 != 0 {
		WriteUint8(&buf, room.setting.forceCamera, &offset)
	}
	if lowFlag&0x40 != 0 {
		WriteUint8(&buf, room.setting.gameModeID, &offset)
	}
	if lowFlag&0x80 != 0 {
		WriteUint8(&buf, room.setting.mapID, &offset)
		WriteUint8(&buf, room.setting.unk13, &offset)
	}
	if lowFlag&0x100 != 0 {
		WriteUint8(&buf, room.setting.maxPlayers, &offset)
	}
	if lowFlag&0x200 != 0 {
		WriteUint8(&buf, room.setting.winLimit, &offset)
	}
	if lowFlag&0x400 != 0 {
		WriteUint16(&buf, room.setting.killLimit, &offset)
	}
	if lowFlag&0x800 != 0 {
		WriteUint8(&buf, room.setting.unk17, &offset)
	}
	if lowFlag&0x1000 != 0 {
		WriteUint8(&buf, room.setting.unk18, &offset)
	}
	if lowFlag&0x2000 != 0 {
		WriteUint8(&buf, room.setting.weaponRestrictions, &offset)
	}
	if lowFlag&0x4000 != 0 {
		WriteUint8(&buf, room.setting.status, &offset)
	}
	if lowFlag&0x8000 != 0 {
		WriteUint8(&buf, room.setting.unk21, &offset)
		WriteUint8(&buf, room.setting.mapCycleType, &offset)
		WriteUint8(&buf, room.setting.unk23, &offset)
		WriteUint8(&buf, room.setting.unk24, &offset)
	}
	if lowFlag&0x10000 != 0 {
		WriteUint8(&buf, room.setting.unk25, &offset)
	}
	if lowFlag&0x20000 != 0 {
		WriteUint8(&buf, room.setting.lenOfMultiMaps, &offset)
		for _, v := range room.setting.multiMaps {
			WriteUint8(&buf, v, &offset)
		}
	}
	if lowFlag&0x40000 != 0 {
		WriteUint8(&buf, room.setting.teamBalanceType, &offset)
	}
	if lowFlag&0x80000 != 0 {
		WriteUint8(&buf, room.setting.unk29, &offset)
	}
	if lowFlag&0x100000 != 0 {
		WriteUint8(&buf, room.setting.unk30, &offset)
	}
	if lowFlag&0x200000 != 0 {
		WriteUint8(&buf, room.setting.unk31, &offset)
	}
	if lowFlag&0x400000 != 0 {
		WriteUint8(&buf, room.setting.unk32, &offset)
	}
	if lowFlag&0x800000 != 0 {
		WriteUint8(&buf, room.setting.unk33, &offset)
	}
	if lowFlag&0x1000000 != 0 {
		WriteUint8(&buf, room.setting.areBotsEnabled, &offset)
		if room.setting.areBotsEnabled != 0 {
			WriteUint8(&buf, room.setting.botDifficulty, &offset)
			WriteUint8(&buf, room.setting.numCtBots, &offset)
			WriteUint8(&buf, room.setting.numTrBots, &offset)
		}
	}

	if lowFlag&0x2000000 != 0 {
		WriteUint8(&buf, room.setting.unk35, &offset)
	}

	if lowFlag&0x4000000 != 0 {
		WriteUint8(&buf, room.setting.unk36, &offset)
	}

	if lowFlag&0x8000000 != 0 {
		WriteUint8(&buf, room.setting.unk37, &offset)
	}

	if lowFlag&0x10000000 != 0 {
		WriteUint8(&buf, room.setting.unk38, &offset)
	}

	if lowFlag&0x20000000 != 0 {
		WriteUint8(&buf, room.setting.unk39, &offset)
	}

	if lowFlag&0x40000000 != 0 {
		WriteUint8(&buf, room.setting.isIngame, &offset)
	}

	if lowFlag&0x80000000 != 0 {
		WriteUint16(&buf, room.setting.startMoney, &offset)
	}

	if highFlag&0x1 != 0 {
		WriteUint8(&buf, room.setting.changeTeams, &offset)
	}

	if highFlag&0x2 != 0 {
		WriteUint8(&buf, room.setting.unk43, &offset)
	}

	if highFlag&0x4 != 0 {
		WriteUint8(&buf, room.setting.hltvEnabled, &offset)
	}

	if highFlag&0x8 != 0 {
		WriteUint8(&buf, room.setting.unk45, &offset)
	}

	if highFlag&0x10 != 0 {
		WriteUint8(&buf, room.setting.respawnTime, &offset)
	}
	return buf[:offset]
}

func getFlags(room roomInfo) uint64 {
	lowFlag := 0
	highFlag := 0

	/* tslint:disable: no-bitwise */
	if room.setting.roomName != nil {
		lowFlag |= 0x1
	}
	if room.setting.unk00 != 0 {
		lowFlag |= 0x2
	}
	if room.setting.unk01 != 0 &&
		room.setting.unk02 != 0 &&
		room.setting.unk03 != 0 {
		lowFlag |= 0x4
	}
	if room.setting.unk09 != nil {
		lowFlag |= 0x8
	}
	if room.setting.unk10 != 0 {
		lowFlag |= 0x10
	}
	if room.setting.forceCamera != 0 {
		lowFlag |= 0x20
	}
	if room.setting.gameModeID != 0 {
		lowFlag |= 0x40
	}
	if room.setting.mapID != 0 && room.setting.unk13 != 0 {
		lowFlag |= 0x80
	}
	if room.setting.maxPlayers != 0 {
		lowFlag |= 0x100
	}
	if room.setting.winLimit != 0 {
		lowFlag |= 0x200
	}
	if room.setting.killLimit != 0 {
		lowFlag |= 0x400
	}
	if room.setting.unk17 != 0 {
		lowFlag |= 0x800
	}
	if room.setting.unk18 != 0 {
		lowFlag |= 0x1000
	}
	if room.setting.weaponRestrictions != 0 {
		lowFlag |= 0x2000
	}
	if room.setting.status != 0 {
		lowFlag |= 0x4000
	}
	if room.setting.unk21 != 0 &&
		room.setting.mapCycleType != 0 &&
		room.setting.unk23 != 0 &&
		room.setting.unk24 != 0 {
		lowFlag |= 0x8000
	}
	if room.setting.unk25 != 0 {
		lowFlag |= 0x10000
	}
	if room.setting.multiMaps != nil {
		lowFlag |= 0x20000
	}
	if room.setting.teamBalanceType != 0 {
		lowFlag |= 0x40000
	}
	if room.setting.unk29 != 0 {
		lowFlag |= 0x80000
	}
	if room.setting.unk30 != 0 {
		lowFlag |= 0x100000
	}
	if room.setting.unk31 != 0 {
		lowFlag |= 0x200000
	}
	if room.setting.unk32 != 0 {
		lowFlag |= 0x400000
	}
	if room.setting.unk33 != 0 {
		lowFlag |= 0x800000
	}
	if room.setting.areBotsEnabled != 0 {
		lowFlag |= 0x1000000
	}

	if room.setting.unk35 != 0 {
		lowFlag |= 0x2000000
	}

	if room.setting.unk36 != 0 {
		lowFlag |= 0x4000000
	}

	if room.setting.unk37 != 0 {
		lowFlag |= 0x8000000
	}

	if room.setting.unk38 != 0 {
		lowFlag |= 0x10000000
	}

	if room.setting.unk39 != 0 {
		lowFlag |= 0x20000000
	}

	if room.setting.isIngame != 0 {
		lowFlag |= 0x40000000
	}

	if room.setting.startMoney != 0 {
		lowFlag |= 0x80000000
	}

	if room.setting.changeTeams != 0 {
		highFlag |= 0x1
	}

	if room.setting.unk43 != 0 {
		highFlag |= 0x2
	}

	if room.setting.hltvEnabled != 0 {
		highFlag |= 0x4
	}

	if room.setting.unk45 != 0 {
		highFlag |= 0x8
	}

	if room.setting.respawnTime != 0 {
		highFlag |= 0x10
	}
	/* tslint:enable: no-bitwise */

	flags := uint64(highFlag)
	flags = flags << 32
	return flags + uint64(lowFlag)
}

func onRoomList(seq *uint8, p *packet, client net.Conn) {
	var pkt roomsRequestPacket
	if praseChannelRequest(*p, &pkt) {
		uPtr := getUserFromConnection(client)
		if uPtr.userid <= 0 {
			log.Println("A unknow Client from", client.RemoteAddr().String(), "request a RoomList !")
			return
		}
		//发送频道请求返回包
		chlsrv := getChannelServerWithID(pkt.channelServerIndex)
		if chlsrv == nil {
			log.Println("Client from", client.RemoteAddr().String(), "request a unknown channelServer !")
			return
		}
		rst := BuildLobbyReply(seq, *p)
		WriteLen(&rst) //写入长度
		client.Write(rst)
		log.Println("Sent a lobbyReply packet to", client.RemoteAddr().String())
		//发送频道请求所得房间列表
		chl := getChannelWithID(pkt.channelIndex, *chlsrv)
		if chl == nil {
			log.Println("Client from", client.RemoteAddr().String(), "request a unknown channel !")
			return
		}
		rst = BuildRoomList(seq, *p, *chl)
		WriteLen(&rst) //写入长度
		client.Write(rst)
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
	rm.id = id
	rm.hostUserID = u.userid
	rm.hostUserName = u.username
	rm.users = append(rm.users, *u)
	rm.parentChannel = chl.channelID
	rm.countingDown = false
	rm.countdown = DefaultCountdownNum
	rm.numPlayers = 1
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

//getNewRoomID() 暂定
func getNewRoomID(chl channelInfo) uint16 {
	if chl.roomNum > MAXROOMNUMS {
		log.Fatalln("Room is too much ! Unable to create more !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXROOMNUMS + 2]uint16
	//哈希思想
	for i := 0; i < int(chl.roomNum); i++ {
		intbuf[chl.rooms[i].id] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXROOMNUMS+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint16(i)
		}
	}
	return 0
}

type upSettingReq struct {
	flags              uint64
	lenOfRoomName      uint8
	roomName           []byte
	unk00              uint8
	unk01              uint8
	unk02              uint32
	unk03              uint32
	lenOfunk09         uint8
	unk09              []byte
	unk10              uint16
	forceCamera        uint8
	gameModeID         uint8
	mapID              uint8
	unk13              uint8
	maxPlayers         uint8
	winLimit           uint8
	killLimit          uint16
	unk17              uint8
	unk18              uint8
	weaponRestrictions uint8
	status             uint8
	unk21              uint8
	mapCycleType       uint8
	unk23              uint8
	unk24              uint8
	unk25              uint8
	numOfMultiMaps     uint8
	multiMaps          []uint8
	teamBalanceType    uint8
	unk29              uint8
	unk30              uint8
	unk31              uint8
	unk32              uint8
	unk33              uint8
	botEnabled         uint8
	botDifficulty      uint8
	numCtBots          uint8
	numTrBots          uint8
	unk35              uint8
	unk36              uint8
	unk37              uint8
	unk38              uint8
	unk39              uint8
	isIngame           uint8
	startMoney         uint16
	changeTeams        uint8
	unk43              uint8
	hltvEnabled        uint8
	unk45              uint8
	respawnTime        uint8
}

//onUpdateRoom 房主更新房间信息
func onUpdateRoom(seq *uint8, p packet, client net.Conn) {
	//检索数据报
	var pkt upSettingReq
	if !praseUpdateRoomPacket(p, &pkt) {
		log.Println("Client from", client.RemoteAddr().String(), "sent a illegal packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr.userid <= 0 {
		log.Println("Client from", client.RemoteAddr().String(), "try to update room but not in server !")
		return
	}
	//检查用户是不是房主
	curroom := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if (*curroom).id <= 0 {
		log.Println("Client from", client.RemoteAddr().String(), "try to update a null room but in server !")
		delChannelRoom((*curroom).id, uPtr.getUserChannelID(), uPtr.getUserChannelServerID())
		return
	}
	if (*curroom).hostUserID != (uPtr).userid {
		log.Println("Client from", client.RemoteAddr().String(), "try to update a room but isn't host !")
		return
	}
	//检查用户所在房间
	if (*curroom).id != (uPtr).currentRoomId {
		log.Println("Client from", client.RemoteAddr().String(), "try to update a room but not in !")
		return
	}
	//检查当前是不是正在倒计时
	if (*curroom).isGlobalCountdownInProgress() {
		log.Println("Client from", client.RemoteAddr().String(), "try to update a room but is counting !")
		return
	}
	//更新房间设置
	curroom.toUpdateSetting(pkt)
	p.id = TypeRoom
	//向房间所有玩家发送更新报文
	for k, v := range (*curroom).users {
		rst := BytesCombine(BuildHeader(seq, p), buildRoomSetting(*curroom))
		sendPacket(rst, v.currentConnection)
		log.Println("["+strconv.Itoa(k+1)+"/"+strconv.Itoa(int((*curroom).numPlayers))+"] Updated room for", v.currentConnection.RemoteAddr().String(), "!")
	}
	log.Println("Hoster from", client.RemoteAddr().String(), "updated room !")
}

func praseUpdateRoomPacket(src packet, dest *upSettingReq) bool {
	//前5字节数据包通用头部+1字节房间数据包通用头部
	offset := 6
	//读取flag，标记要读的有哪些数据
	flags := ReadUint64(src.data, &offset)
	lowFlag := *(*uint32)(unsafe.Pointer(&flags))
	//右移32比特位
	flags = flags >> 32
	highFlag := *(*uint32)(unsafe.Pointer(&flags))
	if lowFlag&0x1 != 0 {
		(*dest).lenOfRoomName = ReadUint8(src.data, &offset)
		(*dest).roomName = ReadString(src.data, &offset, int((*dest).lenOfRoomName))
	}
	if lowFlag&0x2 != 0 {
		(*dest).unk00 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x4 != 0 {
		(*dest).unk01 = ReadUint8(src.data, &offset)
		(*dest).unk02 = ReadUint32(src.data, &offset)
		(*dest).unk03 = ReadUint32(src.data, &offset)
	}
	if lowFlag&0x8 != 0 {
		(*dest).lenOfunk09 = ReadUint8(src.data, &offset)
		(*dest).unk09 = ReadString(src.data, &offset, int((*dest).lenOfRoomName))
	}
	if lowFlag&0x10 != 0 {
		(*dest).unk10 = ReadUint16(src.data, &offset)
	}
	if lowFlag&0x20 != 0 {
		(*dest).forceCamera = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x40 != 0 {
		(*dest).gameModeID = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x80 != 0 {
		(*dest).mapID = ReadUint8(src.data, &offset)
		(*dest).unk13 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x100 != 0 {
		(*dest).maxPlayers = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x200 != 0 {
		(*dest).winLimit = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x400 != 0 {
		(*dest).killLimit = ReadUint16(src.data, &offset)
	}
	if lowFlag&0x800 != 0 {
		(*dest).unk17 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x1000 != 0 {
		(*dest).unk18 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x2000 != 0 {
		(*dest).weaponRestrictions = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x4000 != 0 {
		(*dest).status = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x8000 != 0 {
		(*dest).unk21 = ReadUint8(src.data, &offset)
		(*dest).mapCycleType = ReadUint8(src.data, &offset)
		(*dest).unk23 = ReadUint8(src.data, &offset)
		(*dest).unk24 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x10000 != 0 {
		(*dest).unk25 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x20000 != 0 {
		(*dest).numOfMultiMaps = ReadUint8(src.data, &offset)
		for i := 0; i < int((*dest).numOfMultiMaps); i++ {
			(*dest).multiMaps[i] = ReadUint8(src.data, &offset)
		}
	}
	if lowFlag&0x40000 != 0 {
		(*dest).teamBalanceType = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x80000 != 0 {
		(*dest).unk29 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x100000 != 0 {
		(*dest).unk30 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x200000 != 0 {
		(*dest).unk31 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x400000 != 0 {
		(*dest).unk32 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x800000 != 0 {
		(*dest).unk33 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x1000000 != 0 {
		(*dest).botEnabled = ReadUint8(src.data, &offset)
		if (*dest).botEnabled != 0 {
			(*dest).botDifficulty = ReadUint8(src.data, &offset)
			(*dest).numCtBots = ReadUint8(src.data, &offset)
			(*dest).numTrBots = ReadUint8(src.data, &offset)
		}
	}

	if lowFlag&0x2000000 != 0 {
		(*dest).unk35 = ReadUint8(src.data, &offset)
	}

	if lowFlag&0x4000000 != 0 {
		(*dest).unk36 = ReadUint8(src.data, &offset)
	}

	if lowFlag&0x8000000 != 0 {
		(*dest).unk37 = ReadUint8(src.data, &offset)
	}

	if lowFlag&0x10000000 != 0 {
		(*dest).unk38 = ReadUint8(src.data, &offset)
	}

	if lowFlag&0x20000000 != 0 {
		(*dest).unk39 = ReadUint8(src.data, &offset)
	}

	if lowFlag&0x40000000 != 0 {
		(*dest).isIngame = ReadUint8(src.data, &offset)
	}

	if lowFlag&0x80000000 != 0 {
		(*dest).startMoney = ReadUint16(src.data, &offset)
	}

	if highFlag&0x1 != 0 {
		(*dest).changeTeams = ReadUint8(src.data, &offset)
	}

	if highFlag&0x2 != 0 {
		(*dest).unk43 = ReadUint8(src.data, &offset)
	}

	if highFlag&0x4 != 0 {
		(*dest).hltvEnabled = ReadUint8(src.data, &offset)
	}

	if highFlag&0x8 != 0 {
		(*dest).unk45 = ReadUint8(src.data, &offset)
	}

	if highFlag&0x10 != 0 {
		(*dest).respawnTime = ReadUint8(src.data, &offset)
	}
	return true
}

func getRoomFromID(chlsrvID uint8, chlID uint8, roomID uint16) *roomInfo {
	if chlsrvID <= 0 ||
		chlID <= 0 ||
		roomID <= 0 {
		return nil
	}
	chlsrv := getChannelServerWithID(chlsrvID)
	if chlsrv.serverIndex <= 0 {
		return nil
	}
	chl := getChannelWithID(chlID, *chlsrv)
	if chl.channelID <= 0 || chl.roomNum <= 0 {
		return nil
	}
	for k, v := range chl.rooms {
		if v.id == roomID {
			return &chl.rooms[k]
		}
	}
	return nil
}

func (ri roomInfo) isGlobalCountdownInProgress() bool {
	return ri.countingDown
}

func (dest *roomInfo) toUpdateSetting(src upSettingReq) {
	flags := src.flags
	lowFlag := *(*uint32)(unsafe.Pointer(&flags))
	//右移32比特位
	flags = flags >> 32
	highFlag := *(*uint32)(unsafe.Pointer(&flags))
	if lowFlag&0x1 != 0 {
		(*dest).setting.lenOfName = src.lenOfRoomName
		(*dest).setting.roomName = src.roomName
	}
	if lowFlag&0x2 != 0 {
		(*dest).setting.unk00 = src.unk00
	}
	if lowFlag&0x4 != 0 {
		(*dest).setting.unk01 = src.unk01
		(*dest).setting.unk02 = src.unk02
		(*dest).setting.unk03 = src.unk03
	}
	if lowFlag&0x8 != 0 {
		(*dest).setting.lenOfunk09 = src.lenOfunk09
		(*dest).setting.unk09 = src.unk09
	}
	if lowFlag&0x10 != 0 {
		(*dest).setting.unk10 = src.unk10
	}
	if lowFlag&0x20 != 0 {
		(*dest).setting.forceCamera = src.forceCamera
	}
	if lowFlag&0x40 != 0 {
		(*dest).setting.gameModeID = src.gameModeID
	}
	if lowFlag&0x80 != 0 {
		(*dest).setting.mapID = src.mapID
		(*dest).setting.unk13 = src.unk13
	}
	if lowFlag&0x100 != 0 {
		(*dest).setting.maxPlayers = src.maxPlayers
	}
	if lowFlag&0x200 != 0 {
		(*dest).setting.winLimit = src.winLimit
	}
	if lowFlag&0x400 != 0 {
		(*dest).setting.killLimit = src.killLimit
	}
	if lowFlag&0x800 != 0 {
		(*dest).setting.unk17 = src.unk17
	}
	if lowFlag&0x1000 != 0 {
		(*dest).setting.unk18 = src.unk18
	}
	if lowFlag&0x2000 != 0 {
		(*dest).setting.weaponRestrictions = src.weaponRestrictions
	}
	if lowFlag&0x4000 != 0 {
		(*dest).setting.status = src.status
	}
	if lowFlag&0x8000 != 0 {
		(*dest).setting.unk21 = src.unk21
		(*dest).setting.mapCycleType = src.mapCycleType
		(*dest).setting.unk23 = src.unk23
		(*dest).setting.unk24 = src.unk24
	}
	if lowFlag&0x10000 != 0 {
		(*dest).setting.unk25 = src.unk25
	}
	if lowFlag&0x20000 != 0 {
		(*dest).setting.lenOfMultiMaps = src.numOfMultiMaps
		for i := 0; i < int((*dest).setting.lenOfMultiMaps); i++ {
			(*dest).setting.multiMaps[i] = src.multiMaps[i]
		}
	}
	if lowFlag&0x40000 != 0 {
		(*dest).setting.teamBalanceType = src.teamBalanceType
	}
	if lowFlag&0x80000 != 0 {
		(*dest).setting.unk29 = src.unk29
	}
	if lowFlag&0x100000 != 0 {
		(*dest).setting.unk30 = src.unk30
	}
	if lowFlag&0x200000 != 0 {
		(*dest).setting.unk31 = src.unk31
	}
	if lowFlag&0x400000 != 0 {
		(*dest).setting.unk32 = src.unk32
	}
	if lowFlag&0x800000 != 0 {
		(*dest).setting.unk33 = src.unk33
	}
	if lowFlag&0x1000000 != 0 {
		(*dest).setting.areBotsEnabled = src.botEnabled
		if (*dest).setting.areBotsEnabled != 0 {
			(*dest).setting.botDifficulty = src.botDifficulty
			(*dest).setting.numCtBots = src.numCtBots
			(*dest).setting.numTrBots = src.numTrBots
		}
	}

	if lowFlag&0x2000000 != 0 {
		(*dest).setting.unk35 = src.unk35
	}

	if lowFlag&0x4000000 != 0 {
		(*dest).setting.unk36 = src.unk36
	}

	if lowFlag&0x8000000 != 0 {
		(*dest).setting.unk37 = src.unk37
	}

	if lowFlag&0x10000000 != 0 {
		(*dest).setting.unk38 = src.unk38
	}

	if lowFlag&0x20000000 != 0 {
		(*dest).setting.unk39 = src.unk39
	}

	if lowFlag&0x40000000 != 0 {
		(*dest).setting.isIngame = src.isIngame
	}

	if lowFlag&0x80000000 != 0 {
		(*dest).setting.startMoney = src.startMoney
	}

	if highFlag&0x1 != 0 {
		(*dest).setting.changeTeams = src.changeTeams
	}

	if highFlag&0x2 != 0 {
		(*dest).setting.unk43 = src.unk43
	}

	if highFlag&0x4 != 0 {
		(*dest).setting.hltvEnabled = src.hltvEnabled
	}

	if highFlag&0x8 != 0 {
		(*dest).setting.unk45 = src.unk45
	}

	if highFlag&0x10 != 0 {
		(*dest).setting.respawnTime = src.respawnTime
	}
}
