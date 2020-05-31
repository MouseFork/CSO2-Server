package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"math"
	"net"
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
	(*p).length = getUint16((*p).data[2:4])
	(*p).id = (*p).data[4]
}

func ReadUint8(b []byte, offset *int) uint8 {
	buf := bytes.NewBuffer(b[(*offset):])
	var i uint8
	binary.Read(buf, binary.BigEndian, &i)
	(*offset)++
	return i
}

func ReadUint16(b []byte, offset *int) uint16 {
	buf := bytes.NewBuffer(b[(*offset):])
	var i uint16
	binary.Read(buf, binary.BigEndian, &i)
	(*offset) += 2
	return i
}

func ReadUint32(b []byte, offset *int) uint32 {
	buf := bytes.NewBuffer(b[(*offset):])
	var i uint32
	binary.Read(buf, binary.BigEndian, &i)
	(*offset) += 4
	return i
}

func ReadUint64(b []byte, offset *int) uint64 {
	buf := bytes.NewBuffer(b[(*offset):])
	var i uint64
	binary.Read(buf, binary.BigEndian, &i)
	(*offset) += 8
	return i
}

func ReadString(b []byte, offset *int, len int) []byte {
	(*offset) += len
	return b[(*offset)-len : (*offset)]
}

func getUint16(b []byte) uint16 {
	buf := bytes.NewBuffer(b)
	var i uint16
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

func getUint32(b []byte) uint32 {
	buf := bytes.NewBuffer(b)
	var i uint32
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

func getUint64(b []byte) uint64 {
	buf := bytes.NewBuffer(b)
	var i uint64
	binary.Read(buf, binary.BigEndian, &i)
	return i
}

func GetNextSeq(seq *uint8) uint8 {
	if *seq > MAXSEQUENCE {
		*seq = 0
		return 0
	}
	(*seq)++
	return *seq
}

func WriteUint8(b *[]byte, i uint8, offset *int) {
	(*b)[*offset] = i
	(*offset)++
}

func WriteUint16(b *[]byte, i uint16, offset *int) {
	(*b)[*offset] = uint8(i)
	(*b)[*offset+1] = uint8(i >> 8)
	(*offset) += 2
}

func WriteUint32(b *[]byte, i uint32, offset *int) {
	(*b)[*offset] = uint8(i)
	(*b)[*offset+1] = uint8(i >> 8)
	(*b)[*offset+2] = uint8(i >> 16)
	(*b)[*offset+3] = uint8(i >> 24)
	(*offset) += 4
}

func WriteUint64(b *[]byte, i uint64, offset *int) {
	(*b)[*offset] = uint8(i)
	(*b)[*offset+1] = uint8(i >> 8)
	(*b)[*offset+2] = uint8(i >> 16)
	(*b)[*offset+3] = uint8(i >> 24)
	(*b)[*offset+4] = uint8(i >> 32)
	(*b)[*offset+5] = uint8(i >> 40)
	(*b)[*offset+6] = uint8(i >> 48)
	(*b)[*offset+7] = uint8(i >> 56)
	(*offset) += 8
}

func WriteString(dest *[]byte, src []byte, offset *int) int {
	l := len(src)
	(*dest)[*offset] = uint8(l)
	(*offset)++
	for i := 0; i < l; i++ {
		(*dest)[*offset] = src[i]
		(*offset)++
	}
	return l + 1
}

func BuildHeader(seq *uint8, p packet) []byte {
	header := make([]byte, 5)
	header[0] = TypeSignature
	header[1] = GetNextSeq(seq)
	header[2] = 0
	header[3] = 0
	header[4] = p.id
	return header
}

func WriteLen(data *[]byte) {
	headerL := uint16(len(*data)) - HeaderLen
	(*data)[2] = uint8(headerL)
	(*data)[3] = uint8(headerL >> 8)
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

//结构体转数组
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

func newNullString() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

func sendPacket(data []byte, client net.Conn) {
	WriteLen(&data)
	client.Write(data)
}
