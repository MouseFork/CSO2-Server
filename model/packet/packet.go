package packet

import (
	"math"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type Packet struct {
	Data         []byte
	Datalen      int
	IsGoodPacket bool
	Sequence     uint8
	Length       uint16
	Id           uint8
	CurOffset    int
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

func (p *Packet) PrasePacket() {
	if p.Data[0] != TypeSignature {
		p.IsGoodPacket = false
		return
	}
	p.CurOffset = 1
	p.IsGoodPacket = true
	p.Sequence = ReadUint8(p.Data, &p.CurOffset)
	p.Length = ReadUint16(p.Data, &p.CurOffset)
	p.Id = ReadUint8(p.Data, &p.CurOffset)
	p.Datalen = int(p.Length) + HeaderLen
	if len(p.Data) >= p.Datalen {
		p.Data = p.Data[:p.Datalen]
	}
}
