package main

import (
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

//发送出去的包结构，其中一些未知，知道后会加入user里去
type UserInfo struct {
	flags                uint32 // should always be 0xFFFFFFFF for a full update
	unk00                uint64 // nexon id?
	userName             []byte
	level                uint16
	curExp               uint64
	maxExp               uint64
	unk03                uint32
	rank                 uint8
	rankFrame            uint8
	points               uint64
	playedMatches        uint32
	wins                 uint32
	kills                uint32
	headshots            uint32
	deaths               uint32
	assists              uint32
	accuracy             uint16
	secondsPlayed        uint32
	unk15                uint32
	unk16                uint32
	unk17                uint8
	unk18                uint64
	unk19                uint32
	unk20                uint32
	unk21                uint32
	unk22                uint32
	unk23                uint32
	unk24                uint32
	unk25                uint32
	unk26                []byte
	unk27                uint32
	unk28                uint32
	unk29                uint32
	unk30                uint32
	netCafeName          []byte
	cash                 uint32
	unk33                uint32
	unk34                uint32
	clanName             []byte
	clanMark             uint32
	unk37                uint8
	unk38                []uint32 // array size is always 5
	unk39                []uint32 // array size is always 5
	unk40                uint8
	worldRank            uint32
	unk42                uint32
	unk43                uint8
	unk44                uint16
	unk45                uint32
	mPoints              uint32
	unk47                uint64
	unk48                uint32
	Title                uint16
	unk50                uint16
	unlockedTitles       []uint8
	signature            []byte
	unk53                uint8
	unk54                uint8
	unk55                uint32
	bestGamemode         uint32
	bestMap              uint32
	unk58                uint16
	unlockedAchievements []uint8 // 128字节
	unk60                uint32
	avatars              uint16
	unk62                uint16
	unlockedAvatars      []uint8 // 128字节
	isVip                uint8
	vipLevel             uint8
	vipExp               uint32
	unk67                uint32
	skill_human_curxp    uint64
	skill_human_maxxp    uint64
	skill_human_points   uint8
	skill_zombie_curxp   uint64
	skill_zombie_maxxp   uint64
	skill_zombie_points  uint8
	unk74                uint32
	unk75                uint32
	unk76                uint32
	unk77                uint32
	unk78                uint32
	unk79                uint32
	unk80                uint32
	unk81                uint32
	unk82                uint8
	unk83                uint8
	unk84                uint8
}

func BuildUserInfo(info UserInfo, id uint32, needID bool) []byte {
	infobuf := make([]byte, 1024)
	// if err != nil {
	// 	log.Println("Server occurred an error while senting user info packet !")
	// 	return nil
	// }
	offset := 0
	if needID {
		WriteUint32(&infobuf, id, &offset)
		WriteUint32(&infobuf, info.flags, &offset)
	} else {
		WriteUint32(&infobuf, 0xFFFFFFFF, &offset)
	}
	WriteUint64(&infobuf, info.unk00, &offset)
	WriteString(&infobuf, info.userName, &offset)
	WriteUint16(&infobuf, info.level, &offset)
	WriteUint64(&infobuf, info.curExp, &offset)
	WriteUint64(&infobuf, info.maxExp, &offset)
	WriteUint32(&infobuf, info.unk03, &offset)
	WriteUint8(&infobuf, info.rank, &offset)
	WriteUint8(&infobuf, info.rankFrame, &offset)
	WriteUint64(&infobuf, info.points, &offset)
	WriteUint32(&infobuf, info.playedMatches, &offset)
	WriteUint32(&infobuf, info.wins, &offset)
	WriteUint32(&infobuf, info.kills, &offset)
	WriteUint32(&infobuf, info.headshots, &offset)
	WriteUint32(&infobuf, info.deaths, &offset)
	WriteUint32(&infobuf, info.assists, &offset)
	WriteUint16(&infobuf, info.accuracy, &offset)
	WriteUint32(&infobuf, info.secondsPlayed, &offset)
	WriteUint32(&infobuf, info.unk15, &offset)
	WriteUint32(&infobuf, info.unk16, &offset)
	WriteUint8(&infobuf, info.unk17, &offset)
	WriteUint64(&infobuf, info.unk18, &offset)
	WriteUint32(&infobuf, info.unk19, &offset)
	WriteUint32(&infobuf, info.unk20, &offset)
	WriteUint32(&infobuf, info.unk21, &offset)
	WriteUint32(&infobuf, info.unk22, &offset)
	WriteUint32(&infobuf, info.unk23, &offset)
	WriteUint32(&infobuf, info.unk24, &offset)
	WriteUint32(&infobuf, info.unk25, &offset)
	WriteString(&infobuf, info.unk26, &offset)
	WriteUint32(&infobuf, info.unk27, &offset)
	WriteUint32(&infobuf, info.unk28, &offset)
	WriteUint32(&infobuf, info.unk29, &offset)
	WriteUint32(&infobuf, info.unk30, &offset)
	WriteString(&infobuf, info.netCafeName, &offset)
	WriteUint32(&infobuf, info.cash, &offset)
	WriteUint32(&infobuf, info.unk33, &offset)
	WriteUint32(&infobuf, info.unk34, &offset)
	WriteString(&infobuf, info.clanName, &offset)
	WriteUint32(&infobuf, info.clanMark, &offset)
	WriteUint8(&infobuf, info.unk37, &offset)

	for _, v := range info.unk38 {
		WriteUint32(&infobuf, v, &offset)
	}
	for _, v := range info.unk39 {
		WriteUint32(&infobuf, v, &offset)
	}

	WriteUint8(&infobuf, info.unk40, &offset)
	WriteUint32(&infobuf, info.worldRank, &offset)
	WriteUint32(&infobuf, info.unk42, &offset)
	WriteUint8(&infobuf, info.unk43, &offset)
	WriteUint16(&infobuf, info.unk44, &offset)
	WriteUint32(&infobuf, info.unk45, &offset)
	WriteUint32(&infobuf, info.mPoints, &offset)
	WriteUint64(&infobuf, info.unk47, &offset)
	WriteUint32(&infobuf, info.unk48, &offset)
	WriteUint16(&infobuf, info.Title, &offset)
	WriteUint16(&infobuf, info.unk50, &offset)
	for _, v := range info.unlockedTitles {
		WriteUint8(&infobuf, v, &offset)
	}

	WriteString(&infobuf, info.signature, &offset)
	WriteUint8(&infobuf, info.unk53, &offset)
	WriteUint8(&infobuf, info.unk54, &offset)
	WriteUint32(&infobuf, info.unk55, &offset)
	WriteUint32(&infobuf, info.bestGamemode, &offset)
	WriteUint32(&infobuf, info.bestMap, &offset)
	WriteUint16(&infobuf, info.unk58, &offset)
	for _, v := range info.unlockedAchievements {
		WriteUint8(&infobuf, v, &offset)
	}
	WriteUint32(&infobuf, info.unk60, &offset)
	WriteUint16(&infobuf, info.avatars, &offset)
	WriteUint16(&infobuf, info.unk62, &offset)
	for _, v := range info.unlockedAvatars {
		WriteUint8(&infobuf, v, &offset)
	}
	WriteUint8(&infobuf, info.isVip, &offset)
	WriteUint8(&infobuf, info.vipLevel, &offset)
	WriteUint32(&infobuf, info.vipExp, &offset)
	WriteUint32(&infobuf, info.unk67, &offset)
	WriteUint64(&infobuf, info.skill_human_curxp, &offset)
	WriteUint64(&infobuf, info.skill_human_maxxp, &offset)
	WriteUint8(&infobuf, info.skill_human_points, &offset)
	WriteUint64(&infobuf, info.skill_zombie_curxp, &offset)
	WriteUint64(&infobuf, info.skill_zombie_maxxp, &offset)
	WriteUint8(&infobuf, info.skill_zombie_points, &offset)
	WriteUint32(&infobuf, info.unk74, &offset)
	WriteUint32(&infobuf, info.unk75, &offset)
	WriteUint32(&infobuf, info.unk76, &offset)
	WriteUint32(&infobuf, info.unk77, &offset)
	WriteUint32(&infobuf, info.unk78, &offset)
	WriteUint32(&infobuf, info.unk79, &offset)
	WriteUint32(&infobuf, info.unk80, &offset)
	WriteUint32(&infobuf, info.unk81, &offset)
	WriteUint8(&infobuf, info.unk82, &offset)
	WriteUint8(&infobuf, info.unk83, &offset)
	WriteUint8(&infobuf, info.unk84, &offset)
	return infobuf[0:offset]
}

func newUserInfo(u user) UserInfo {
	isvip := uint8(0)
	if u.isVIP() {
		isvip = 1
	}
	return UserInfo{
		0xFFFFFFFF,
		0x2241158F,
		u.username,
		u.level,
		u.currentExp,
		u.maxExp,
		0x0313,
		u.rank,
		u.rankFrame,
		u.points,
		u.playedMatches,
		u.wins,
		u.kills,
		u.headshots,
		u.deaths,
		u.assists,
		u.accuracy,
		u.secondsPlayed,
		0,
		0x32,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		newNullString(),
		0,
		0,
		0,
		0,
		u.netCafeName,
		u.cash,
		0,
		0,
		u.clanName,
		u.clanMark,
		0,
		[]uint32{0, 0, 0, 0, 0},
		[]uint32{0, 0, 0, 0, 0},
		0,
		u.worldRank,
		0,
		0,
		0xFF,
		0,
		u.mpoints,
		0,
		0,
		u.titleId,
		0,
		u.unlockedTitles,
		// []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		u.signature,
		0,
		0,
		7,
		u.bestGamemode,
		u.bestMap,
		0,
		u.unlockedAchievements,
		// []uint8{0x00, 0x00, 0x18, 0x08, 0x00, 0x00, 0x00, 0x00, 0x42, 0x02,
		// 	0x18, 0xC0, 0x09, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0xC0, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0xC8, 0xB7, 0x08, 0x00, 0x00, 0x04, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		0xA5C8,
		u.avatar,
		0,
		u.unlockedAvatars,
		// []uint8{0x00, 0x00, 0x18, 0x08, 0x00, 0x00, 0x00, 0x00, 0x42, 0x02,
		// 	0x18, 0xC0, 0x09, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0xC0, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0xC8, 0xB7, 0x08, 0x00, 0x00, 0x04, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00,
		// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3E, 0x00, 0x00},
		isvip,
		u.vipLevel,
		u.vipXp,
		0,
		u.skillHumanCurXp,
		u.skillHumanMaxXp,
		u.skillHumanPoints,
		u.skillZombieCurXp,
		u.skillZombieMaxXp,
		u.skillZombiePoints,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
	}
}

func (u user) buildUserNetInfo() []byte {
	buf := make([]byte, 25)
	offset := 0
	WriteUint8(&buf, u.getUserTeam(), &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint32BE(&buf, u.netInfo.ExternalIpAddress, &offset) //externalIpAddress
	WriteUint16(&buf, u.netInfo.ExternalServerPort, &offset)  //externalServerPort
	WriteUint16(&buf, u.netInfo.ExternalClientPort, &offset)  //externalClientPort
	WriteUint16(&buf, u.netInfo.ExternalTvPort, &offset)      //externalTvPort
	WriteUint32BE(&buf, u.netInfo.LocalIpAddress, &offset)    //localIpAddress
	WriteUint16(&buf, u.netInfo.LocalServerPort, &offset)     //localServerPort
	WriteUint16(&buf, u.netInfo.LocalClientPort, &offset)     //localClientPort
	WriteUint16(&buf, u.netInfo.LocalTvPort, &offset)         //localTvPort
	// log.Println(u.netInfo.ExternalIpAddress,
	// 	u.netInfo.ExternalServerPort,
	// 	u.netInfo.ExternalClientPort,
	// 	u.netInfo.ExternalTvPort,
	// 	u.netInfo.LocalIpAddress,
	// 	u.netInfo.LocalServerPort,
	// 	u.netInfo.LocalClientPort,
	// 	u.netInfo.LocalTvPort)
	return buf[:offset]
}
