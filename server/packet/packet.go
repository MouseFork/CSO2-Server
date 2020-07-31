package packet

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/model/packet"
)

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
