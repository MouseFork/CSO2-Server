package playerinfo

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
)

func OnSetTitle(p *PacketData, client net.Conn) {
	var pkt InSetTitlePacket
	if !p.PraseSetTitlePacket(&pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal SetTitle packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to SetTitle but not in server !")
		return
	}
	//修改数据
	uPtr.SetTitle(pkt.TitleId)
	//发送数据包
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildSetTitlePacket(uPtr.Userid, pkt.TitleId))
	SendPacket(rst, uPtr.CurrentConnection)
	log.Println("User", string(uPtr.Username), "Set Title to", pkt.TitleId)
	//如果是在房间内
}

func BuildSetTitlePacket(id uint32, Title uint16) []byte {
	buf := make([]byte, 10)
	offset := 0
	WriteUint32(&buf, id, &offset)
	WriteUint32(&buf, 0x8000, &offset)
	WriteUint16(&buf, Title, &offset)
	return buf[:offset]
}
