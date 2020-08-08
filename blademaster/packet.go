package blademaster

import (
	"math"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type (
	Packet struct {
		Data      []byte
		Sequence  uint8
		Length    uint16
		Id        uint8
		CurOffset int
	}
	Header struct {
		Data         []byte
		IsGoodPacket bool
		Sequence     uint8
		Length       uint16
	}
	//房间请求
	InRoomPaket struct {
		InRoomType uint8
	}

	//房间所属频道，用于请求频道
	RoomsRequestPacket struct {
		ChannelServerIndex uint8
		ChannelIndex       uint8
	}

	InFavoritePacket struct {
		PacketType uint8
	}

	InHostPacket struct {
		InHostType uint8
	}
)

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
	TypePlayerInfo       = 69
	TypeUdp              = 70
	TypeBan              = 74
	TypeOption           = 76
	TypeFavorite         = 77
	TypeQuickJoin        = 80
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

func (p *Header) PraseHeadPacket() {
	if p.Data[0] != TypeSignature {
		p.IsGoodPacket = false
		return
	}
	p.IsGoodPacket = true
	p.Sequence = ReadUint8(p.Data, &p.CurOffset)
	p.Length = ReadUint16(p.Data, &p.CurOffset)
}

func (p *Packet) PraseRoomPacket(dest *InRoomPaket) bool {
	if p.Length < 2 {
		return false
	}
	dest.InRoomType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *Packet) PraseQuickPacket(dest *InQuickPacket) bool {
	if p.Length < 2 {
		return false
	}
	dest.inQuickType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *Packet) PraseInQuickListPacket(dest *InQuickList) bool {
	if p.Length < 4 ||
		dest == nil {
		return false
	}
	dest.gameModID = ReadUint8(p.Data, &p.CurOffset)
	dest.IsEnableBot = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *Packet) praseFavoritePacket(dest *InFavoritePacket) bool {
	if p.Length < 2 {
		return false
	}
	dest.packetType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *Packet) PraseFavoriteSetCosmeticsPacket(dest *InFavoriteSetCosmetics) bool {
	if p.Length < 7 {
		return false
	}
	dest.slot = ReadUint8(p.Data, &p.CurOffset)
	dest.itemId = ReadUint32(p.Data, &p.CurOffset)
	return true
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
func BuildHeader(seq *uint8, p Packet) []byte {
	header := make([]byte, 5)
	header[0] = TypeSignature
	header[1] = GetNextSeq(seq)
	header[2] = 0
	header[3] = 0
	header[4] = p.Id
	return header
}

//WriteLen 写入数据长度到数据包通用头部
func WriteLen(data *[]byte) {
	headerL := uint16(len(*data)) - HeaderLen
	(*data)[2] = uint8(headerL)
	(*data)[3] = uint8(headerL >> 8)
}

//NewNullString 新建空的字符串
func NewNullString() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

//SendPacket 发送数据包
func SendPacket(data []byte, client net.Conn) {
	WriteLen(&data)
	client.Write(data)
}
