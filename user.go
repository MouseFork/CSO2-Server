/*struct of user packet {
	basePacket 			4 bytes
	type 				1 byte
	lenOfNexonUsername 	1 byte
	nexonUsername 		len bytes
	lenOfGameUsername 	1 byte
	gameUsername 		len bytes
	unknown01			1 byte
	lenOfPassWd			2 bytes
	PassWd				len bytes
	HddHwid				16 bytes
	netCafeId 			4 bytes
	unknown02			4 bytes
	userSn 				8 bytes
	lenOfUnknownString	2 bytes
	UnknownString03		len bytes
	unknown04 			1 byte
	isLeague 			1 byte
	{ always null ignore it /  包里好像不存在，忽视
		unk04 				1 byte
		lenOfUnknown05 		1 byte
		UnknownString05		len bytes
		lenOfUnknown06 		1 byte
		UnknownString06		len bytes
	}
	lenOfString			1 bytes
	String				len bytes
}*/

package main

import (
	"log"
	"net"
)

const (
	//MAXUSERNUM 最大用户数
	MAXUSERNUM = 1024
)

//全局用户管理
type userManager struct {
	userNum int
	users   []user
}

type loginPacket struct {
	//BasePacket 	   5 bytes
	BasePacket         packet
	lenOfNexonUsername uint8
	nexonUsername      []byte //假定nexonUsername是唯一
	lenOfGameUsername  uint8
	gameUsername       []byte
	unknown01          uint8
	lenOfPassWd        uint16
	PassWd             []byte
	//HddHwid 	   	   16 bytes
	HddHwid []byte
	//netCafeID 	   4 bytes
	netCafeID          []byte
	unknown02          uint32
	userSn             uint64
	lenOfUnknownString uint16
	UnknownString03    []byte
	unknown04          uint8
	isLeague           uint8
	lenOfString        uint8
	String             []byte
}

type user struct {
	//个人信息
	userid               uint32
	loginName            []byte
	username             []byte
	password             []byte
	level                uint16
	rank                 uint8
	rankFrame            uint8
	points               uint64
	currentExp           uint64
	maxExp               uint64
	playedMatches        uint32
	wins                 uint32
	kills                uint32
	headshots            uint32
	deaths               uint32
	assists              uint32
	accuracy             uint16
	secondsPlayed        uint32
	netCafeName          []byte
	cash                 uint32
	clanName             []byte
	clanMark             uint32
	worldRank            uint32
	mpoints              uint32
	titleId              uint16
	unlockedTitles       []byte
	signature            []byte
	bestGamemode         uint32
	bestMap              uint32
	unlockedAchievements []byte
	avatar               uint16
	unlockedAvatars      []byte
	vipLevel             uint8
	vipXp                uint32
	skillHumanCurXp      uint64
	skillHumanMaxXp      uint64
	skillHumanPoints     uint8
	skillZombieCurXp     uint64
	skillZombieMaxXp     uint64
	skillZombiePoints    uint8
	//连接
	currentConnection net.Conn
	//频道房间信息
	currentChannelServerIndex uint8
	currentChannelIndex       uint8
	currentRoomId             uint8
	//仓库信息
	//inventory userInventory
}

//发送出去的包结构，其中一些未知，知道后会加入user里去
type UserInfo struct {
	flags         uint32 // should always be 0xFFFFFFFF for a full update
	unk00         uint64 // nexon id?
	userName      []byte
	level         uint16
	curExp        uint64
	maxExp        uint64
	unk03         uint32
	rank          uint8
	rankFrame     uint8
	points        uint64
	playedMatches uint32
	wins          uint32
	kills         uint32
	headshots     uint32
	deaths        uint32
	assists       uint32
	accuracy      uint16
	secondsPlayed uint32
	unk15         uint32
	unk16         uint32
	unk17         uint8
	unk18         uint64
	unk19         uint32
	unk20         uint32
	unk21         uint32
	unk22         uint32
	unk23         uint32
	unk24         uint32
	unk25         uint32
	unk26         []byte
	unk27         uint32
	unk28         uint32
	unk29         uint32
	unk30         uint32
	unk31         []byte
	cash          uint32
	unk33         uint32
	unk34         uint32
	clanName      []byte
	clanMark      uint32
	unk37         uint8
	unk38         []uint32 // array size is always 5
	unk39         []uint32 // array size is always 5
	unk40         uint8
	worldRank     uint32
	unk42         uint32
	unk43         uint8
	unk44         uint16
	unk45         uint32
	unk46         uint32
	unk47         uint64
	unk48         uint32
	unk49         uint16
	unk50         uint16
	unk51         []uint8
	unk52         []byte
	unk53         uint8
	unk54         uint8
	unk55         uint32
	unk56         uint32
	unk57         uint32
	unk58         uint16
	unk59         []uint8 // it must always be 0x80 long
	unk60         uint32
	icon          uint16
	unk62         uint16
	unk63         []uint8
	isVip         uint8
	vipLevel      uint8
	vipExp        uint32
	unk67         uint32
	unk68         uint64
	unk69         uint64
	unk70         uint8
	unk71         uint64
	unk72         uint64
	unk73         uint8
	unk74         uint32
	unk75         uint32
	unk76         uint32
	unk77         uint32
	unk78         uint32
	unk79         uint32
	unk80         uint32
	unk81         uint32
	unk82         uint8
	unk83         uint8
	unk84         uint8
}

func onLoginPacket(seq *uint8, p *packet, client *(net.Conn)) bool {
	var pkt loginPacket
	pkt.BasePacket = *p
	PraseLoginPacket(&pkt)            //分析收到的用户数据
	if !pkt.BasePacket.IsGoodPacket { //如果包损坏或非法
		(*p).IsGoodPacket = false
		return false
	}
	//获得用户数据，待定
	u := getUserByLogin(pkt)
	if u.userid <= 0 {
		log.Println("User", string(pkt.gameUsername), "from", (*client).RemoteAddr().String(), "login failed !")
		(*client).Close()
	}
	u.currentConnection = *client
	//把用户加入用户管理器
	if !addUser(&u) {
		log.Println("User", string(pkt.gameUsername), "from", (*client).RemoteAddr().String(), "login failed !")
		(*client).Close()
	}
	//UserStart部分
	pkt.BasePacket.id = TypeUserStart //UserStart消息标识
	rst := BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildUserStart(u))
	WriteLen(&rst)       //写入长度
	(*client).Write(rst) //发送UserStart消息
	log.Println("User", string(u.loginName), "from", (*client).RemoteAddr().String(), "logged in !")
	log.Println("Sent a user start packet to", (*client).RemoteAddr().String())
	//UserInfo部分
	pkt.BasePacket.id = TypeUserInfo //发送UserInfo消息
	info := newUserInfo(u)
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildUserInfo(info, u.userid))
	WriteLen(&rst)       //写入长度
	(*client).Write(rst) //发送UserInfo消息
	log.Println("Sent a user info packet to", (*client).RemoteAddr().String())
	//ServerList部分
	onServerList(seq, p, client)
	/*(*p).id = TypeServerList
	rst1 := BytesCombine(BuildHeader(seq, *p), BuildServerList())
	WriteLen(&rst1) //写入长度
	rst = BytesCombine(rst, rst1)
	(*client).Write(rst) //发送UserStart消息
	log.Println("Sent a server list packet to", (*client).RemoteAddr().String())*/
	//Inventory部分

	return true
}

//BuildUserStart 返回结构
// userId
// loginName
// userName
// unk00
// holepunchPort
func BuildUserStart(u user) []byte {
	//暂时都取GameUsername
	userbuf := make([]byte, 9+int(len(u.loginName))+int(len(u.username)))
	offset := 0
	WriteUint32(&userbuf, u.userid, &offset)
	WriteString(&userbuf, u.loginName, &offset)
	WriteString(&userbuf, u.username, &offset)
	WriteUint8(&userbuf, 1, &offset)
	WriteUint16(&userbuf, uint16(HOLEPUNCHPORT), &offset)
	return userbuf
}

func BuildUserInfo(info UserInfo, id uint32) []byte {
	infobuf := make([]byte, 1024)
	// if err != nil {
	// 	log.Println("Server occurred an error while senting user info packet !")
	// 	return nil
	// }
	offset := 0
	WriteUint32(&infobuf, id, &offset)
	WriteUint32(&infobuf, info.flags, &offset)
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
	WriteString(&infobuf, info.unk31, &offset)
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
	WriteUint32(&infobuf, info.unk46, &offset)
	WriteUint64(&infobuf, info.unk47, &offset)
	WriteUint32(&infobuf, info.unk48, &offset)
	WriteUint16(&infobuf, info.unk49, &offset)
	WriteUint16(&infobuf, info.unk50, &offset)
	for _, v := range info.unk51 {
		WriteUint8(&infobuf, v, &offset)
	}

	WriteString(&infobuf, info.unk52, &offset)
	WriteUint8(&infobuf, info.unk53, &offset)
	WriteUint8(&infobuf, info.unk54, &offset)
	WriteUint32(&infobuf, info.unk55, &offset)
	WriteUint32(&infobuf, info.unk56, &offset)
	WriteUint32(&infobuf, info.unk57, &offset)
	WriteUint16(&infobuf, info.unk58, &offset)
	for _, v := range info.unk59 {
		WriteUint8(&infobuf, v, &offset)
	}
	WriteUint32(&infobuf, info.unk60, &offset)
	WriteUint16(&infobuf, info.icon, &offset)
	WriteUint16(&infobuf, info.unk62, &offset)
	for _, v := range info.unk63 {
		WriteUint8(&infobuf, v, &offset)
	}
	WriteUint8(&infobuf, info.isVip, &offset)
	WriteUint8(&infobuf, info.vipLevel, &offset)
	WriteUint32(&infobuf, info.vipExp, &offset)
	WriteUint32(&infobuf, info.unk67, &offset)
	WriteUint64(&infobuf, info.unk68, &offset)
	WriteUint64(&infobuf, info.unk69, &offset)
	WriteUint8(&infobuf, info.unk70, &offset)
	WriteUint64(&infobuf, info.unk71, &offset)
	WriteUint64(&infobuf, info.unk72, &offset)
	WriteUint8(&infobuf, info.unk73, &offset)
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

func PraseLoginPacket(p *loginPacket) {
	if (*p).BasePacket.datalen < 50 {
		(*p).BasePacket.IsGoodPacket = false
		return
	}
	lenOfData := (*p).BasePacket.datalen
	offset := 5

	(*p).lenOfNexonUsername = (*p).BasePacket.data[offset]
	offset++

	(*p).nexonUsername = (*p).BasePacket.data[offset : offset+int((*p).lenOfNexonUsername)]
	offset += int((*p).lenOfNexonUsername)
	if offset > lenOfData {
		(*p).BasePacket.IsGoodPacket = false
		return
	}

	(*p).lenOfGameUsername = (*p).BasePacket.data[offset]
	offset++

	(*p).gameUsername = (*p).BasePacket.data[offset : offset+int((*p).lenOfGameUsername)]
	offset += int((*p).lenOfGameUsername)

	(*p).unknown01 = (*p).BasePacket.data[offset]
	offset++

	(*p).lenOfPassWd = getUint16((*p).BasePacket.data[offset : offset+2])
	offset += 2

	(*p).PassWd = (*p).BasePacket.data[offset : offset+int((*p).lenOfPassWd)]
	offset += int((*p).lenOfPassWd)

	(*p).HddHwid = (*p).BasePacket.data[offset : offset+16]
	offset += 16
	if offset > lenOfData {
		(*p).BasePacket.IsGoodPacket = false
		return
	}

	(*p).netCafeID = (*p).BasePacket.data[offset : offset+4]
	offset += 4

	(*p).unknown02 = getUint32((*p).BasePacket.data[offset : offset+4])
	offset += 4

	(*p).userSn = getUint64((*p).BasePacket.data[offset : offset+8])
	offset += 8

	(*p).lenOfUnknownString = getUint16((*p).BasePacket.data[offset : offset+2])
	offset += 2

	(*p).UnknownString03 = (*p).BasePacket.data[offset : offset+int((*p).lenOfUnknownString)]
	offset += int((*p).lenOfUnknownString)

	(*p).unknown04 = (*p).BasePacket.data[offset]
	offset++

	(*p).isLeague = (*p).BasePacket.data[offset]
	offset++

	(*p).lenOfString = (*p).BasePacket.data[offset]
	offset++
	if offset+int((*p).lenOfString) > lenOfData {
		(*p).BasePacket.IsGoodPacket = false
		return
	}

	(*p).String = (*p).BasePacket.data[offset : offset+int((*p).lenOfUnknownString)]
	//...
}

func newUserInfo(u user) UserInfo {
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
		u.isVIP(),
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

func addUser(src *user) bool {
	if (*src).userid == 0 {
		log.Fatalln("ID of User", (*src).username, "is illegal !")
		return false
	}
	if UserManager.userNum > MAXUSERNUM {
		log.Fatalln("Online users is too more to login !")
		return false
	}
	for _, v := range UserManager.users {
		if v.userid == (*src).userid {
			log.Fatalln("User is already logged in !")
			return false
		}
	}
	UserManager.userNum++
	UserManager.users = append(UserManager.users, *src)
	return true
}

func delUser(src *user) bool {
	if (*src).userid == 0 {
		log.Fatalln("ID of User", (*src).username, "is illegal !")
		return false
	}
	if UserManager.userNum == 0 {
		log.Fatalln("There is no online user !")
		return false
	}
	for i, v := range UserManager.users {
		if v.userid == (*src).userid {
			UserManager.users = append(UserManager.users[:i], UserManager.users[i+1:]...)
			UserManager.userNum--
			return true
		}
	}
	return false
}

func delUserWithConn(con net.Conn) bool {
	if UserManager.userNum == 0 {
		log.Fatalln("There is no online user !")
		return false
	}
	for i, v := range UserManager.users {
		if v.currentConnection == con {
			UserManager.users = append(UserManager.users[:i], UserManager.users[i+1:]...)
			UserManager.userNum--
			return true
		}
	}
	return false
}
func getNewUserID() uint32 {
	if UserManager.userNum > MAXUSERNUM {
		log.Fatalln("Online users is too much , unable to get a new id !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXUSERNUM + 2]uint32
	//哈希思想
	for i := 0; i < int(UserManager.userNum); i++ {
		intbuf[UserManager.users[i].userid] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXUSERNUM+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint32(i)
		}
	}
	return 0
}

//假定nexonUsername是唯一
func getUserByLogin(p loginPacket) user {
	u := findOnlineUserByName(p.gameUsername)
	if u.userid <= 0 {
		return getUserFromDatabase(p)
	}
	return u
}

func findOnlineUserByName(name []byte) user {
	l := len(name)
	if l <= 0 {
		log.Fatalln("User name is illegal !")
		return getNewUser()
	}
	for _, v := range UserManager.users {
		if string(v.username) == string(name) {
			return v
		}
	}
	return getNewUser()
}

func (u user) isVIP() uint8 {
	if u.vipLevel <= 0 {
		return 0
	}
	return 1
}

func (u *user) setID(id uint32) {
	(*u).userid = id
}

func (u *user) setUserName(name []byte) {
	(*u).loginName = name
	(*u).username = name
}

func getNewUser() user {
	return user{
		0,
		[]byte{},        //loginname
		[]byte{},        //username,looks can change it to another name
		[]byte{},        //passwd
		1,               //level
		0,               //rank
		0,               //rankframe
		0x7AF3,          //points
		0,               //curEXP
		1000,            //maxEXP
		0x0A,            //playermatchs
		0,               //wins
		0,               //kills
		0,               //headshots
		0,               //deaths
		0,               //assists
		0x0A,            // accuracy
		0x290C,          // secondsPlayed
		newNullString(), // netCafeName
		0,               // cash
		newNullString(), // clanName
		0,               // clanMark
		0,               // worldRank
		0,               // mpoints
		0,               // titleId
		[]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // unlockedTitles
		newNullString(), // signature
		5,               // bestGamemode
		9,               // bestMap
		[]uint8{0x00, 0x00, 0x18, 0x08, 0x00, 0x00, 0x00, 0x00, 0x42, 0x02,
			0x18, 0xC0, 0x09, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0xC0, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0xC8, 0xB7, 0x08, 0x00, 0x00, 0x04, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // unlockedAchievements
		1006, // avatar
		[]uint8{0x00, 0x00, 0x18, 0x08, 0x00, 0x00, 0x00, 0x00, 0x42, 0x02,
			0x18, 0xC0, 0x09, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0xC0, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0xC8, 0xB7, 0x08, 0x00, 0x00, 0x04, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3E, 0x00, 0x00}, // unlockedAvatars
		0,      //viplevel
		0,      //vipXp
		0x02FB, //skillHumanCurXp
		0x19AC, //skillHumanMaxXp
		0,      //skillHumanPoints
		0,      //skillZombieCurXp
		0x16F6, //skillZombieMaxXp
		0,      //skillZombiePoints
		nil,    //connection
		1,      //serverid
		0,      //channelid
		0,      //roomid
	}
}

//暂定功能
//从数据库中读取用户数据
//如果是新用户则保存到数据库中
func getUserFromDatabase(p loginPacket) user {
	u := getNewUser()
	u.setID(getNewUserID())
	u.setUserName(p.gameUsername)
	u.password = p.PassWd
	return u
}

func setUserChannel() {

}

func getUserFromConnection(client net.Conn) {

}
