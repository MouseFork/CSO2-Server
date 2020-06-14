package main

import (
	"log"
	"net"
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
			onLeaveRoom(seq, p, client)
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
