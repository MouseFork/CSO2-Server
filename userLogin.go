package main

import (
	"log"
	"net"
)

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
	netCafeID          uint32
	unknown02          uint32
	userSn             uint64
	lenOfUnknownString uint16
	UnknownString03    []byte
	unknown04          uint8
	isLeague           uint8
	lenOfString        uint8
	String             []byte
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
		log.Println("Error : User", string(pkt.gameUsername), "from", (*client).RemoteAddr().String(), "login failed !")
		(*client).Close()
	}
	u.currentConnection = *client
	//绑定seq
	u.currentSequence = seq
	//把用户加入用户管理器
	if !addUser(&u) {
		log.Println("Error : User", string(pkt.gameUsername), "from", (*client).RemoteAddr().String(), "login failed !")
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
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildUserInfo(info, u.userid, true))
	WriteLen(&rst)       //写入长度
	(*client).Write(rst) //发送UserInfo消息
	log.Println("Sent a user info packet to", (*client).RemoteAddr().String())
	//ServerList部分
	onServerList(seq, p, client)
	//Inventory部分
	pkt.BasePacket.id = TypeInventory_Create
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildInventoryInfo(u))
	sendPacket(rst, *client)
	pkt.BasePacket.id = TypeInventory_Add
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildInventoryInfo(u))
	sendPacket(rst, *client)
	//unlock
	pkt.BasePacket.id = 0x5a
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildUnlockReply())
	sendPacket(rst, *client)
	//偏好装备
	pkt.BasePacket.id = TypeFavorite
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildCosmetics(u.inventory))
	sendPacket(rst, *client)
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildLoadout(u.inventory))
	sendPacket(rst, *client)
	//购买菜单
	pkt.BasePacket.id = TypeOption
	rst = BytesCombine(BuildHeader(seq, pkt.BasePacket), BuildBuyMenu(u.inventory))
	sendPacket(rst, *client)
	log.Println("Sent a user inventory packet to", (*client).RemoteAddr().String())
	return true
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
	(*p).netCafeID = ReadUint32BE((*p).BasePacket.data, &offset)
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
