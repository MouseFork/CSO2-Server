package main

import (
	"unsafe"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

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

//创建房间设置数据包
func buildRoomSetting(room roomInfo) []byte {
	buf := make([]byte, 128+room.setting.lenOfName+ //实际计算是最大63字节+长度
		room.setting.lenOfunk09+
		room.setting.lenOfMultiMaps)
	offset := 0
	WriteUint8(&buf, OUTUpdateSettings, &offset)
	//room.flags = getFlags(room)
	room.flags = 0xFFFFFFFFFFFFFFFF
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
		(*dest).setting.multiMaps = make([]byte, src.numOfMultiMaps)
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
