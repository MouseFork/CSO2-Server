package room

import (
	"unsafe"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

//创建房间设置数据包
func BuildRoomSetting(room *Room) []byte {
	buf := make([]byte, 128+room.Setting.LenOfName+ //实际计算是最大63字节+长度
		room.Setting.LenOfPassWd+
		room.Setting.LenOfMultiMaps)
	offset := 0
	WriteUint8(&buf, OUTUpdateSettings, &offset)
	var flags uint64
	if room.Flags != room.Lastflags {
		flags = room.Lastflags
	} else {
		flags = room.Flags
	}
	WriteUint64(&buf, flags, &offset)
	lowFlag := *(*uint32)(unsafe.Pointer(&flags))
	flags = flags >> 32
	highFlag := *(*uint32)(unsafe.Pointer(&flags))
	if lowFlag&0x1 != 0 {
		WriteString(&buf, room.Setting.RoomName, &offset)
	}
	if lowFlag&0x2 != 0 {
		WriteUint8(&buf, room.Setting.Unk00, &offset)
	}
	if lowFlag&0x4 != 0 {
		WriteUint8(&buf, room.Setting.Unk01, &offset)
		WriteUint32(&buf, room.Setting.Unk02, &offset)
		WriteUint32(&buf, room.Setting.Unk03, &offset)
	}
	if lowFlag&0x8 != 0 {
		WriteString(&buf, room.Setting.PassWd, &offset)
	}
	if lowFlag&0x10 != 0 {
		WriteUint16(&buf, room.Setting.Unk10, &offset)
	}
	if lowFlag&0x20 != 0 {
		WriteUint8(&buf, room.Setting.ForceCamera, &offset)
	}
	if lowFlag&0x40 != 0 {
		WriteUint8(&buf, room.Setting.GameModeID, &offset)
	}
	if lowFlag&0x80 != 0 {
		WriteUint8(&buf, room.Setting.MapID, &offset)
		WriteUint8(&buf, room.Setting.Unk13, &offset)
	}
	if lowFlag&0x100 != 0 {
		WriteUint8(&buf, room.Setting.MaxPlayers, &offset)
	}
	if lowFlag&0x200 != 0 {
		WriteUint8(&buf, room.Setting.WinLimit, &offset)
	}
	if lowFlag&0x400 != 0 {
		WriteUint16(&buf, room.Setting.KillLimit, &offset)
	}
	if lowFlag&0x800 != 0 {
		WriteUint8(&buf, room.Setting.Unk17, &offset)
	}
	if lowFlag&0x1000 != 0 {
		WriteUint8(&buf, room.Setting.Unk18, &offset)
	}
	if lowFlag&0x2000 != 0 {
		WriteUint8(&buf, room.Setting.WeaponRestrictions, &offset)
	}
	if lowFlag&0x4000 != 0 {
		WriteUint8(&buf, room.Setting.Status, &offset)
	}
	if lowFlag&0x8000 != 0 {
		WriteUint8(&buf, room.Setting.Unk21, &offset)
		WriteUint8(&buf, room.Setting.MapCycleType, &offset)
		WriteUint8(&buf, room.Setting.Unk23, &offset)
		WriteUint8(&buf, room.Setting.Unk24, &offset)
	}
	if lowFlag&0x10000 != 0 {
		WriteUint8(&buf, room.Setting.Unk21, &offset)
	}
	if lowFlag&0x20000 != 0 {
		WriteUint8(&buf, room.Setting.LenOfMultiMaps, &offset)
		for _, v := range room.Setting.MultiMaps {
			WriteUint8(&buf, v, &offset)
		}
	}
	if lowFlag&0x40000 != 0 {
		WriteUint8(&buf, room.Setting.TeamBalanceType, &offset)
	}
	if lowFlag&0x80000 != 0 {
		WriteUint8(&buf, room.Setting.Unk29, &offset)
	}
	if lowFlag&0x100000 != 0 {
		WriteUint8(&buf, room.Setting.Unk30, &offset)
	}
	if lowFlag&0x200000 != 0 {
		WriteUint8(&buf, room.Setting.Unk31, &offset)
	}
	if lowFlag&0x400000 != 0 {
		WriteUint8(&buf, room.Setting.Unk32, &offset)
	}
	if lowFlag&0x800000 != 0 {
		WriteUint8(&buf, room.Setting.Unk33, &offset)
	}
	if lowFlag&0x1000000 != 0 {
		WriteUint8(&buf, room.Setting.AreBotsEnabled, &offset)
		if room.Setting.AreBotsEnabled != 0 {
			WriteUint8(&buf, room.Setting.BotDifficulty, &offset)
			WriteUint8(&buf, room.Setting.NumCtBots, &offset)
			WriteUint8(&buf, room.Setting.NumTrBots, &offset)
		}
	}

	if lowFlag&0x2000000 != 0 {
		WriteUint8(&buf, room.Setting.Unk35, &offset)
	}

	if lowFlag&0x4000000 != 0 {
		WriteUint8(&buf, room.Setting.Unk36, &offset)
	}

	if lowFlag&0x8000000 != 0 {
		WriteUint8(&buf, room.Setting.Unk37, &offset)
	}

	if lowFlag&0x10000000 != 0 {
		WriteUint8(&buf, room.Setting.Unk38, &offset)
	}

	if lowFlag&0x20000000 != 0 {
		WriteUint8(&buf, room.Setting.Unk39, &offset)
	}

	if lowFlag&0x40000000 != 0 {
		WriteUint8(&buf, room.Setting.IsIngame, &offset)
	}

	if lowFlag&0x80000000 != 0 {
		WriteUint16(&buf, room.Setting.StartMoney, &offset)
	}

	if highFlag&0x1 != 0 {
		WriteUint8(&buf, room.Setting.ChangeTeams, &offset)
	}

	if highFlag&0x2 != 0 {
		WriteUint8(&buf, room.Setting.Unk43, &offset)
	}

	if highFlag&0x4 != 0 {
		WriteUint8(&buf, room.Setting.HltvEnabled, &offset)
	}

	if highFlag&0x8 != 0 {
		WriteUint8(&buf, room.Setting.Unk45, &offset)
	}

	if highFlag&0x10 != 0 {
		WriteUint8(&buf, room.Setting.RespawnTime, &offset)
	}
	return buf[:offset]
}

func GetFlags(room Room) uint64 {
	lowFlag := 0
	highFlag := 0

	if room.Setting.RoomName != nil {
		lowFlag |= 0x1
	}
	if room.Setting.Unk00 != 0 {
		lowFlag |= 0x2
	}
	if room.Setting.Unk01 != 0 &&
		room.Setting.Unk02 != 0 &&
		room.Setting.Unk03 != 0 {
		lowFlag |= 0x4
	}
	if room.Setting.PassWd != nil {
		lowFlag |= 0x8
	}
	if room.Setting.Unk10 != 0 {
		lowFlag |= 0x10
	}
	if room.Setting.ForceCamera != 0 {
		lowFlag |= 0x20
	}
	if room.Setting.GameModeID != 0 {
		lowFlag |= 0x40
	}
	if room.Setting.MapID != 0 && room.Setting.Unk13 != 0 {
		lowFlag |= 0x80
	}
	if room.Setting.MaxPlayers != 0 {
		lowFlag |= 0x100
	}
	if room.Setting.WinLimit != 0 {
		lowFlag |= 0x200
	}
	if room.Setting.KillLimit != 0 {
		lowFlag |= 0x400
	}
	if room.Setting.Unk17 != 0 {
		lowFlag |= 0x800
	}
	if room.Setting.Unk18 != 0 {
		lowFlag |= 0x1000
	}
	if room.Setting.WeaponRestrictions != 0 {
		lowFlag |= 0x2000
	}
	if room.Setting.Status != 0 {
		lowFlag |= 0x4000
	}
	if room.Setting.Unk21 != 0 &&
		room.Setting.MapCycleType != 0 &&
		room.Setting.Unk23 != 0 &&
		room.Setting.Unk24 != 0 {
		lowFlag |= 0x8000
	}
	if room.Setting.Unk25 != 0 {
		lowFlag |= 0x10000
	}
	if room.Setting.MultiMaps != nil {
		lowFlag |= 0x20000
	}
	if room.Setting.TeamBalanceType != 0 {
		lowFlag |= 0x40000
	}
	if room.Setting.Unk29 != 0 {
		lowFlag |= 0x80000
	}
	if room.Setting.Unk30 != 0 {
		lowFlag |= 0x100000
	}
	if room.Setting.Unk31 != 0 {
		lowFlag |= 0x200000
	}
	if room.Setting.Unk32 != 0 {
		lowFlag |= 0x400000
	}
	if room.Setting.Unk33 != 0 {
		lowFlag |= 0x800000
	}
	if room.Setting.AreBotsEnabled != 0 {
		lowFlag |= 0x1000000
	}

	if room.Setting.Unk35 != 0 {
		lowFlag |= 0x2000000
	}

	if room.Setting.Unk36 != 0 {
		lowFlag |= 0x4000000
	}

	if room.Setting.Unk37 != 0 {
		lowFlag |= 0x8000000
	}

	if room.Setting.Unk38 != 0 {
		lowFlag |= 0x10000000
	}

	if room.Setting.Unk39 != 0 {
		lowFlag |= 0x20000000
	}

	if room.Setting.IsIngame != 0 {
		lowFlag |= 0x40000000
	}

	if room.Setting.StartMoney != 0 {
		lowFlag |= 0x80000000
	}

	if room.Setting.ChangeTeams != 0 {
		highFlag |= 0x1
	}

	if room.Setting.Unk43 != 0 {
		highFlag |= 0x2
	}

	if room.Setting.HltvEnabled != 0 {
		highFlag |= 0x4
	}

	if room.Setting.Unk45 != 0 {
		highFlag |= 0x8
	}

	if room.Setting.RespawnTime != 0 {
		highFlag |= 0x10
	}

	flags := uint64(highFlag)
	flags = flags << 32
	return flags + uint64(lowFlag)
}
