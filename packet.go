package main

import (
	"math"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type packet struct {
	data         []byte
	datalen      int
	IsGoodPacket bool
	sequence     uint8
	length       uint16
	id           uint8
}

const (
	TypeSignature        = 0x55
	TypeVersion          = 0
	TypeReply            = 1
	TypeLogin            = 3
	TypeServerList       = 5
	TypeCharacter        = 6
	TypeRequestRoomList  = 7
	TypeRequestChannels  = 10
	TypeRoom             = 65
	TypeChat             = 67
	TypeHost             = 68
	TypeUdp              = 70
	TypeBan              = 74
	TypeOption           = 76
	TypeFavorite         = 77
	TypeQuickStart       = 86
	TypeAutomatch        = 88
	TypeFriend           = 89
	TypeUnlock           = 90
	TypeGZ               = 95
	TypeAchievement      = 96
	TypeConfigInfo       = 106
	TypeLobby            = 107
	TypeUserStart        = 150
	TypeRoomList         = 151
	TypeInventory_Add    = 152
	TypeInventory_Create = 154
	TypeUserInfo         = 157

	MINSEQUENCE = 0
	MAXSEQUENCE = math.MaxUint8
	HeaderLen   = 4
)

func (p *packet) PrasePacket() {
	(*p).datalen = len((*p).data)
	if (*p).data[0] != TypeSignature || (*p).datalen < 5 {
		(*p).IsGoodPacket = false
		return
	}
	(*p).IsGoodPacket = true
	(*p).sequence = (*p).data[1]
	(*p).length = GetUint16((*p).data[2:4])
	(*p).id = (*p).data[4]
}

//GetNextSeq 获取下一次的seq数据包序号
func GetNextSeq(seq *uint8) uint8 {
	if *seq > MAXSEQUENCE {
		*seq = 0
		return 0
	}
	(*seq)++
	return *seq
}

//BuildHeader 建立数据包通用头部
func BuildHeader(seq *uint8, p packet) []byte {
	header := make([]byte, 5)
	header[0] = TypeSignature
	header[1] = GetNextSeq(seq)
	header[2] = 0
	header[3] = 0
	header[4] = p.id
	return header
}

//WriteLen 写入数据长度到数据包通用头部
func WriteLen(data *[]byte) {
	headerL := uint16(len(*data)) - HeaderLen
	(*data)[2] = uint8(headerL)
	(*data)[3] = uint8(headerL >> 8)
}

//newNullString 新建空的字符串
func newNullString() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

//sendPacket 发送数据包
func sendPacket(data []byte, client net.Conn) {
	WriteLen(&data)
	client.Write(data)
}
