package typestruct

import "unsafe"

//每个房间的设置数据
type RoomSettings struct {
	LenOfName          uint8
	RoomName           []byte
	Unk00              uint8
	Unk01              uint8
	Unk02              uint32
	Unk03              uint32
	LenOfPassWd        uint8
	PassWd             []byte
	Unk10              uint16
	ForceCamera        uint8
	GameModeID         uint8
	MapID              uint8
	Unk13              uint8
	MaxPlayers         uint8
	WinLimit           uint8
	KillLimit          uint16
	Unk17              uint8
	Unk18              uint8
	WeaponRestrictions uint8
	Status             uint8
	Unk21              uint8
	MapCycleType       uint8
	Unk23              uint8
	Unk24              uint8
	Unk25              uint8
	LenOfMultiMaps     uint8
	MultiMaps          []byte
	TeamBalanceType    uint8
	Unk29              uint8
	Unk30              uint8
	Unk31              uint8
	Unk32              uint8
	Unk33              uint8
	AreBotsEnabled     uint8
	BotDifficulty      uint8
	NumCtBots          uint8
	NumTrBots          uint8
	Unk35              uint8
	Unk36              uint8
	Unk37              uint8
	Unk38              uint8
	Unk39              uint8
	StartMoney         uint16
	ChangeTeams        uint8
	Unk43              uint8
	HltvEnabled        uint8
	Unk45              uint8
	RespawnTime        uint8
	NextMapEnabled     uint8
	Difficulty         uint8
	IsIngame           uint8
}

func (dest *Room) ToUpdateSetting(src *InUpSettingReq) {
	flags := src.Flags
	lowFlag := *(*uint32)(unsafe.Pointer(&flags))
	//右移32比特位
	flags = flags >> 32
	highFlag := *(*uint32)(unsafe.Pointer(&flags))
	if lowFlag&0x1 != 0 {
		dest.Setting.LenOfName = src.LenOfRoomName
		dest.Setting.RoomName = src.RoomName
	}
	if lowFlag&0x2 != 0 {
		dest.Setting.Unk00 = src.Unk00
	}
	if lowFlag&0x4 != 0 {
		dest.Setting.Unk01 = src.Unk01
		dest.Setting.Unk02 = src.Unk02
		dest.Setting.Unk03 = src.Unk03
	}
	if lowFlag&0x8 != 0 {
		dest.Setting.LenOfPassWd = src.LenOfpasswd
		dest.Setting.PassWd = src.Passwd
		if dest.Setting.LenOfPassWd > 0 {
			dest.PasswordProtected = 1
		} else {
			dest.PasswordProtected = 0
		}
	}
	if lowFlag&0x10 != 0 {
		dest.Setting.Unk10 = src.Unk10
	}
	if lowFlag&0x20 != 0 {
		dest.Setting.ForceCamera = src.ForceCamera
	}
	if lowFlag&0x40 != 0 {
		dest.Setting.GameModeID = src.GameModeID
	}
	if lowFlag&0x80 != 0 {
		dest.Setting.MapID = src.MapID
		dest.Setting.Unk13 = src.Unk13
	}
	if lowFlag&0x100 != 0 {
		dest.Setting.MaxPlayers = src.MaxPlayers
	}
	if lowFlag&0x200 != 0 {
		dest.Setting.WinLimit = src.WinLimit
	}
	if lowFlag&0x400 != 0 {
		dest.Setting.KillLimit = src.KillLimit
	}
	if lowFlag&0x800 != 0 {
		dest.Setting.Unk17 = src.Unk17
	}
	if lowFlag&0x1000 != 0 {
		dest.Setting.Unk18 = src.Unk18
	}
	if lowFlag&0x2000 != 0 {
		dest.Setting.WeaponRestrictions = src.WeaponRestrictions
	}
	if lowFlag&0x4000 != 0 {
		dest.Setting.Status = src.Status
	}
	if lowFlag&0x8000 != 0 {
		dest.Setting.Unk21 = src.Unk21
		dest.Setting.MapCycleType = src.MapCycleType
		dest.Setting.Unk23 = src.Unk23
		dest.Setting.Unk24 = src.Unk24
	}
	if lowFlag&0x10000 != 0 {
		dest.Setting.Unk25 = src.Unk25
	}
	if lowFlag&0x20000 != 0 {
		dest.Setting.LenOfMultiMaps = src.NumOfMultiMaps
		dest.Setting.MultiMaps = make([]byte, src.NumOfMultiMaps)
		for i := 0; i < int(dest.Setting.LenOfMultiMaps); i++ {
			dest.Setting.MultiMaps[i] = src.MultiMaps[i]
		}
	}
	if lowFlag&0x40000 != 0 {
		dest.Setting.TeamBalanceType = src.TeamBalanceType
	}
	if lowFlag&0x80000 != 0 {
		dest.Setting.Unk29 = src.Unk29
	}
	if lowFlag&0x100000 != 0 {
		dest.Setting.Unk30 = src.Unk30
	}
	if lowFlag&0x200000 != 0 {
		dest.Setting.Unk31 = src.Unk31
	}
	if lowFlag&0x400000 != 0 {
		dest.Setting.Unk32 = src.Unk32
	}
	if lowFlag&0x800000 != 0 {
		dest.Setting.Unk33 = src.Unk33
	}
	if lowFlag&0x1000000 != 0 {
		dest.Setting.AreBotsEnabled = src.BotEnabled
		if dest.Setting.AreBotsEnabled != 0 {
			dest.Setting.BotDifficulty = src.BotDifficulty
			dest.Setting.NumCtBots = src.NumCtBots
			dest.Setting.NumTrBots = src.NumTrBots
		}
	}

	if lowFlag&0x2000000 != 0 {
		dest.Setting.Unk35 = src.Unk35
	}

	if lowFlag&0x4000000 != 0 {
		dest.Setting.Unk36 = src.Unk36
	}

	if lowFlag&0x8000000 != 0 {
		dest.Setting.Unk37 = src.Unk37
	}

	if lowFlag&0x10000000 != 0 {
		dest.Setting.Unk38 = src.Unk38
	}

	if lowFlag&0x20000000 != 0 {
		dest.Setting.Unk39 = src.Unk39
	}

	if lowFlag&0x40000000 != 0 {
		dest.Setting.IsIngame = src.IsIngame
	}

	if lowFlag&0x80000000 != 0 {
		dest.Setting.StartMoney = src.StartMoney
	}

	if highFlag&0x1 != 0 {
		dest.Setting.ChangeTeams = src.ChangeTeams
	}

	if highFlag&0x2 != 0 {
		dest.Setting.Unk43 = src.Unk43
	}

	if highFlag&0x4 != 0 {
		dest.Setting.HltvEnabled = src.HltvEnabled
	}

	if highFlag&0x8 != 0 {
		dest.Setting.Unk45 = src.Unk45
	}

	if highFlag&0x10 != 0 {
		dest.Setting.RespawnTime = src.RespawnTime
	}
	dest.Lastflags = src.Flags
}
