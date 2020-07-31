package blademaster

import (
	"sync"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

//房间信息
type (
	RoomInfo struct {
		Id                 uint16
		Lastflags          uint64
		Flags              uint64
		RoomNumber         uint8
		PasswordProtected  uint8
		Unk08              uint8
		HostUserID         uint32
		HostUserName       []byte
		Unk11              uint8
		Unk12              uint8
		Unk13              uint32
		Unk14              uint16
		Unk15              uint16
		Unk16              uint32
		Unk17              uint16
		Unk18              uint16
		Unk19              uint8
		Unk20              uint8
		Unk21              uint8
		Unk24              uint8
		Unk26              uint8
		Unk27              []uint8
		Unk28              uint8
		Unk29              uint8
		Unk30              uint64
		Unk31              uint8
		Unk35              uint8
		AreFlashesDisabled uint8
		CanSpec            uint8
		IsVipRoom          uint8
		VipRoomLevel       uint8

		//设置
		Setting       RoomSettings
		CountingDown  bool
		Countdown     uint8
		NumPlayers    uint8
		UserIDs       []uint32
		ParentChannel uint8
		CtScore       uint8
		TrScore       uint8
		CtKillNum     uint32
		TrKillNum     uint32
		WinnerTeam    uint8

		RoomMutex *sync.Mutex
	}

	//未知，用于请求频道
	LobbyJoinRoom struct {
		Unk00 uint8
		Unk01 uint8
		Unk02 uint8
	}
)

const (
	//房间操作，加入、暂停等

	GameStart         = 0
	HostJoin          = 1
	HostStop          = 3
	LeaveResultWindow = 4

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
	MAXROOMNUMS         = 0xFF
	DefaultCountdownNum = 7
)

func (rm RoomInfo) IsGlobalCountdownInProgress() bool {
	return rm.CountingDown
}

func (rm RoomInfo) RoomGetUser(id uint32) *User {
	if id <= 0 ||
		rm.Id <= 0 ||
		rm.NumPlayers <= 0 {
		return nil
	}
	for _, v := range rm.UserIDs {
		u := GetUserFromID(v)
		if u != nil {
			return u
		}
	}
	return nil
}

func (rm *RoomInfo) StopCountdown() {
	if rm == nil {
		return
	}
	rm.Countdown = DefaultCountdownNum
	rm.CountingDown = false
}

func (rm *RoomInfo) SetStatus(status uint8) {
	if rm == nil {
		return
	}
	if status == 1 ||
		status == 2 {
		rm.Setting.Status = status
		rm.Setting.IsIngame = status - 1
	}
}

func (rm RoomInfo) CanStartGame() bool {
	switch rm.Setting.GameModeID {
	case deathmatch, original, originalmr, casualbomb, casualoriginal, eventmod01, eventmod02, diy, campaign1, campaign2, campaign3, campaign4, campaign5, tdm_small, de_small, madcity, madcity_team, gunteamdeath, gunteamdeath_re, stealth, teamdeath, teamdeath_mutation, pig:
		if rm.GetNumOfReadyPlayers() < 2 {
			return false
		}
	case giant, hide, hide2, hide_match, hide_origin, hide_Item, hide_multi, ghost, tag, zombie, zombiecraft, zombie_commander, zombie_prop, zombie_zeta:
		if rm.GetNumOfRealReadyPlayers() < 2 {
			return false
		}
	}
	return true
}

func (rm *RoomInfo) ProgressCountdown(num uint8) {
	if rm.Countdown > DefaultCountdownNum ||
		rm.Countdown < 0 {
		rm.Countdown = 0
	}
	if rm.CountingDown == false {
		rm.CountingDown = true
		rm.Countdown = DefaultCountdownNum
	}
	rm.Countdown--
	if rm.Countdown != num {
		DebugInfo(2, "Error : Host is counting", num, "but room is", rm.Countdown)
	}
}

func (rm *RoomInfo) GetCountdown() uint8 {
	if rm.CountingDown == false {
		DebugInfo(2, "Error : tried to get countdown without counting down")
		return 0
	}
	if rm.Countdown > DefaultCountdownNum ||
		rm.Countdown < 0 {
		rm.Countdown = DefaultCountdownNum
	}
	return rm.Countdown
}

func (rm RoomInfo) GetAllCtNum() int {
	num := 0
	for _, v := range rm.UserIDs {
		u := GetUserFromID(v)
		if u != nil && u.GetUserTeam() == CounterTerrorist {
			num++
		}
	}
	return num
}

func (rm RoomInfo) GetAllTrNum() int {
	num := 0
	for _, v := range rm.UserIDs {
		u := GetUserFromID(v)
		if u != nil && u.GetUserTeam() == Terrorist {
			num++
		}
	}
	return num
}

func (rm RoomInfo) GetFreeSlots() int {
	return int(rm.Setting.MaxPlayers - rm.NumPlayers)
}

func (rm *RoomInfo) JoinUser(u *User) bool {
	destTeam := rm.FindDesirableTeam()
	if destTeam <= 0 {
		DebugInfo(2, "Error : Cant add User", string(u.Username), "to room", string(rm.Setting.RoomName))
		return false
	}
	rm.NumPlayers++
	u.CurrentTeam = uint8(destTeam)
	u.SetUserStatus(UserNotReady)
	u.SetUserRoom(rm.Id)
	u.SetUserIngame(false)
	rm.UserIDs = append(rm.UserIDs, u.Userid)
	return true
}

func (rm RoomInfo) FindDesirableTeam() int {
	trNum := 0
	ctNum := 0
	for _, v := range rm.UserIDs {
		u := GetUserFromID(v)
		if u.GetUserTeam() == Terrorist {
			trNum++
		} else if u.GetUserTeam() == CounterTerrorist {
			ctNum++
		} else {
			DebugInfo(2, "Error : User", string(u.Username), "is in Unknown team in room", string(rm.Setting.RoomName))
			return 0
		}
	}
	if rm.Setting.AreBotsEnabled != 0 {
		u := rm.RoomGetUser(rm.HostUserID)
		if u == nil ||
			u.Userid <= 0 {
			return 0
		}
		botsInHostTeam := 0
		if u.GetUserTeam() == CounterTerrorist {
			botsInHostTeam = int(rm.Setting.NumCtBots)
			if botsInHostTeam > 0 {
				return CounterTerrorist
			}
		} else if u.GetUserTeam() == Terrorist {
			botsInHostTeam = int(rm.Setting.NumTrBots)
			if botsInHostTeam > 0 {
				return Terrorist
			}
		} else {
			DebugInfo(2, "Error : Host", string(u.Username), "is in Unknown team in room", string(rm.Setting.RoomName))
			return 0
		}
	}
	if trNum < ctNum {
		return Terrorist
	} else {
		return CounterTerrorist
	}
}

func (rm *RoomInfo) CheckIngameStatus() {
	if rm == nil {
		return
	}
	if rm.NumPlayers <= 0 {
		rm.SetStatus(StatusWaiting)
		return
	}
	for _, v := range rm.UserIDs {
		u := GetUserFromID(v)
		if u != nil && u.CurrentIsIngame {
			rm.SetStatus(StatusIngame)
			return
		}
	}
	rm.SetStatus(StatusWaiting)
}

func (rm RoomInfo) GetNumOfRealReadyPlayers() int {
	num := 0
	for _, v := range rm.UserIDs {
		u := GetUserFromID(v)
		if u != nil && (u.IsUserReady() ||
			u.Userid == rm.HostUserID) {
			num++
		}
	}
	return num
}
func (rm RoomInfo) GetNumOfReadyPlayers() int {
	botPlayers := int(rm.Setting.NumCtBots + rm.Setting.NumTrBots)
	if rm.Setting.TeamBalanceType == WithBots {
		numCts := rm.GetAllCtNum()
		numTrs := rm.GetAllTrNum()
		requiredBalanceBots := IntAbs(numCts - numTrs)
		botPlayers = Ternary(botPlayers > requiredBalanceBots, botPlayers, requiredBalanceBots).(int)
	}
	return botPlayers + rm.GetNumOfRealReadyPlayers()
}

func (rm *RoomInfo) SetRoomScore(ctScore uint8, trScore uint8) {
	if rm == nil {
		return
	}
	rm.CtScore = ctScore
	rm.TrScore = trScore
}

func (rm *RoomInfo) ResetRoomScore() {
	if rm == nil {
		return
	}
	rm.CtScore = 0
	rm.TrScore = 0
}
func (rm *RoomInfo) SetRoomWinner(Winner uint8) {
	if rm == nil {
		return
	}
	rm.WinnerTeam = Winner
}

func (rm *RoomInfo) ResetRoomWinner() {
	if rm == nil {
		return
	}
	rm.WinnerTeam = 0
}
func (rm *RoomInfo) CountRoomCtKill() {
	if rm == nil {
		return
	}
	rm.CtKillNum++
}

func (rm *RoomInfo) CountRoomTrKill() {
	if rm == nil {
		return
	}
	rm.TrKillNum++
}

func (rm *RoomInfo) ResetRoomKillNum() {
	if rm == nil {
		return
	}
	rm.CtKillNum = 0
	rm.TrKillNum = 0
}
