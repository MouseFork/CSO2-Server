package typestruct

import (
	"sync"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

//房间信息
type (
	Room struct {
		Id uint16
		//Flags              uint64
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
		Users         map[uint32]*User
		ParentChannel uint8
		CtScore       uint8
		TrScore       uint8
		CtKillNum     uint32
		TrKillNum     uint32
		WinnerTeam    uint8

		RoomMutex *sync.Mutex
	}
)

const (
	//host操作，加入、暂停等

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
	ModeOriginal            = 1
	ModeTeamdeath           = 2
	ModeZombie              = 3
	ModeStealth             = 4
	ModeGunteamdeath        = 5
	ModeTutorial            = 6
	ModeHide                = 7
	ModePig                 = 8
	ModeAnimationtest_vcd   = 9
	ModeGz_survivor         = 10
	ModeDevtest             = 11
	ModeOriginalmr          = 12
	ModeOriginalmrdraw      = 13
	ModeCasualbomb          = 14
	ModeDeathmatch          = 15
	ModeScenario_test       = 16
	ModeGz                  = 17
	ModeGz_intro            = 18
	ModeGz_tour             = 19
	ModeGz_pve              = 20
	ModeEventmod01          = 21
	ModeDuel                = 22
	ModeGz_ZB               = 23
	ModeHeroes              = 24
	ModeEventmod02          = 25
	ModeZombiecraft         = 26
	ModeCampaign1           = 27
	ModeCampaign2           = 28
	ModeCampaign3           = 29
	ModeCampaign4           = 30
	ModeCampaign5           = 31
	ModeCampaign6           = 32
	ModeCampaign7           = 33
	ModeCampaign8           = 34
	ModeCampaign9           = 35
	ModeZ_scenario          = 36
	ModeZombie_prop         = 37
	ModeGhost               = 38
	ModeTag                 = 39
	ModeHide_match          = 40
	ModeHide_ice            = 41
	ModeDiy                 = 42
	ModeHide_Item           = 43
	ModeZd_boss1            = 44
	ModeZd_boss2            = 45
	ModeZd_boss3            = 46
	ModePractice            = 47
	ModeZombie_commander    = 48
	ModeCasualoriginal      = 49
	ModeHide2               = 50
	ModeGunball             = 51
	ModeZombie_zeta         = 53
	ModeTdm_small           = 54
	ModeDe_small            = 55
	ModeGunteamdeath_re     = 56
	ModeEndless_wave        = 57
	ModeRankmatch_original  = 58
	ModeRankmatch_teamdeath = 59
	ModePlay_ground         = 60
	ModeMadcity             = 61
	ModeHide_origin         = 62
	ModeTeamdeath_mutation  = 63
	ModeGiant               = 64
	ModeZ_scenario_side     = 65
	ModeHide_multi          = 66
	ModeMadcity_team        = 67
	ModeRankmatch_stealth   = 68

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

	DefaultCountdownNum = 7

	//Start CountDown
	InProgress = 0
	Stop       = 1
)

func (rm Room) IsGlobalCountdownInProgress() bool {
	return rm.CountingDown
}

func (rm Room) RoomGetUser(id uint32) *User {
	if id <= 0 ||
		rm.Id <= 0 ||
		rm.NumPlayers <= 0 {
		return nil
	}
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	if v, ok := rm.Users[id]; ok {
		return v
	}
	return nil
}

func (rm *Room) StopCountdown() {
	if rm == nil {
		return
	}
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	rm.Countdown = DefaultCountdownNum
	rm.CountingDown = false
}

func (rm *Room) SetStatus(status uint8) {
	if rm == nil {
		return
	}
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	if status == 1 ||
		status == 2 {
		rm.Setting.Status = status
		rm.Setting.IsIngame = status - 1
	}
}

func (rm Room) CanStartGame() bool {
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	switch rm.Setting.GameModeID {
	case ModeDeathmatch, ModeOriginal, ModeOriginalmr, ModeCasualbomb,
		ModeCasualoriginal, ModeEventmod01, ModeEventmod02, ModeDiy,
		ModeCampaign1, ModeCampaign2, ModeCampaign3, ModeCampaign4,
		ModeCampaign5, ModeTdm_small, ModeDe_small, ModeMadcity, ModeMadcity_team,
		ModeGunteamdeath, ModeGunteamdeath_re, ModeStealth, ModeTeamdeath,
		ModeTeamdeath_mutation, ModePig:
		if rm.GetNumOfReadyPlayers() < 2 {
			return false
		}
	case ModeGiant, ModeHide, ModeHide2, ModeHide_match, ModeHide_origin,
		ModeHide_Item, ModeHide_multi, ModeGhost, ModeTag, ModeZombie, ModeZombiecraft,
		ModeZombie_commander, ModeZombie_prop, ModeZombie_zeta:
		if rm.GetNumOfRealReadyPlayers() < 2 {
			return false
		}
	}
	return true
}

func (rm *Room) ProgressCountdown(num uint8) {
	if rm == nil {
		return
	}
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
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

func (rm *Room) GetCountdown() uint8 {
	if rm == nil {
		return 0
	}
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
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

func (rm Room) GetAllCtNum() int {
	num := 0
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	for _, v := range rm.Users {
		if v != nil && v.GetUserTeam() == UserForceCounterTerrorist {
			num++
		}
	}
	return num
}

func (rm Room) GetAllTrNum() int {
	num := 0
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	for _, v := range rm.Users {
		if v != nil && v.GetUserTeam() == UserForceTerrorist {
			num++
		}
	}
	return num
}

func (rm Room) GetFreeSlots() int {
	return int(rm.Setting.MaxPlayers - rm.NumPlayers)
}

func (rm *Room) JoinUser(u *User) bool {
	destTeam := rm.FindDesirableTeam()
	if destTeam <= 0 {
		DebugInfo(2, "Error : Cant add User", string(u.Username), "to room", string(rm.Setting.RoomName))
		return false
	}
	rm.RoomMutex.Lock()
	rm.NumPlayers++
	rm.RoomMutex.Unlock()
	u.CurrentTeam = uint8(destTeam)
	u.SetUserStatus(UserNotReady)
	u.SetUserRoom(rm.Id)
	u.SetUserIngame(false)
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	if _, ok := rm.Users[u.Userid]; !ok {
		rm.Users[u.Userid] = u
		return true
	}
	return false
}

func (rm Room) FindDesirableTeam() int {
	trNum := 0
	ctNum := 0
	rm.RoomMutex.Lock()
	for _, v := range rm.Users {
		if v.GetUserTeam() == UserForceTerrorist {
			trNum++
		} else if v.GetUserTeam() == UserForceCounterTerrorist {
			ctNum++
		} else {
			DebugInfo(2, "Error : User", string(v.Username), "is in Unknown team in room", string(rm.Setting.RoomName))
			return 0
		}
	}
	rm.RoomMutex.Unlock()
	if rm.Setting.AreBotsEnabled != 0 {
		u := rm.RoomGetUser(rm.HostUserID)
		if u == nil ||
			u.Userid <= 0 {
			return 0
		}
		botsInHostTeam := 0
		if u.GetUserTeam() == UserForceCounterTerrorist {
			botsInHostTeam = int(rm.Setting.NumCtBots)
			if botsInHostTeam > 0 {
				return UserForceCounterTerrorist
			}
		} else if u.GetUserTeam() == UserForceTerrorist {
			botsInHostTeam = int(rm.Setting.NumTrBots)
			if botsInHostTeam > 0 {
				return UserForceTerrorist
			}
		} else {
			DebugInfo(2, "Error : Host", string(u.Username), "is in Unknown team in room", string(rm.Setting.RoomName))
			return 0
		}
	}
	if trNum < ctNum {
		return UserForceTerrorist
	} else {
		return UserForceCounterTerrorist
	}
}

func (rm *Room) CheckIngameStatus() {
	if rm == nil {
		return
	}
	if rm.NumPlayers <= 0 {
		rm.SetStatus(StatusWaiting)
		return
	}
	for _, v := range rm.Users {
		if v != nil && v.CurrentIsIngame {
			rm.SetStatus(StatusIngame)
			return
		}
	}
	rm.SetStatus(StatusWaiting)
}

func (rm Room) GetNumOfRealReadyPlayers() int {
	num := 0
	for _, v := range rm.Users {
		if v != nil && (v.IsUserReady() ||
			v.Userid == rm.HostUserID) {
			num++
		}
	}
	return num
}
func (rm Room) GetNumOfReadyPlayers() int {
	botPlayers := int(rm.Setting.NumCtBots + rm.Setting.NumTrBots)
	if rm.Setting.TeamBalanceType == WithBots {
		numCts := rm.GetAllCtNum()
		numTrs := rm.GetAllTrNum()
		requiredBalanceBots := IntAbs(numCts - numTrs)
		botPlayers = Ternary(botPlayers > requiredBalanceBots, botPlayers, requiredBalanceBots).(int)
	}
	return botPlayers + rm.GetNumOfRealReadyPlayers()
}

func (rm *Room) SetRoomScore(ctScore uint8, trScore uint8) {
	if rm == nil {
		return
	}
	rm.CtScore = ctScore
	rm.TrScore = trScore
}

func (rm *Room) ResetRoomScore() {
	if rm == nil {
		return
	}
	rm.CtScore = 0
	rm.TrScore = 0
}
func (rm *Room) SetRoomWinner(Winner uint8) {
	if rm == nil {
		return
	}
	rm.WinnerTeam = Winner
}

func (rm *Room) ResetRoomWinner() {
	if rm == nil {
		return
	}
	rm.WinnerTeam = 0
}
func (rm *Room) CountRoomCtKill() {
	if rm == nil {
		return
	}
	rm.CtKillNum++
}

func (rm *Room) CountRoomTrKill() {
	if rm == nil {
		return
	}
	rm.TrKillNum++
}

func (rm *Room) ResetRoomKillNum() {
	if rm == nil {
		return
	}
	rm.CtKillNum = 0
	rm.TrKillNum = 0
}

func (rm *Room) RoomRemoveUser(id uint32) {
	if rm.NumPlayers <= 0 {
		return
	}
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	//找到玩家,玩家数-1，删除房间玩家
	if _, ok := rm.Users[id]; ok {
		delete(rm.Users, id)
		rm.NumPlayers--
	}
}
