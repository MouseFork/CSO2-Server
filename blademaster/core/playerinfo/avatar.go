package playerinfo

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
)

func OnSetAvatar(p *PacketData, client net.Conn) {
	var pkt InSetAvatarPacket
	if !p.PraseSetAvatarPacket(&pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal SetAvatar packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to SetAvatar but not in server !")
		return
	}
	//修改数据
	uPtr.SetAvatar(pkt.AvatarId)
	//发送数据包
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildSetAvatarPacket(uPtr.Userid, pkt.AvatarId))
	SendPacket(rst, uPtr.CurrentConnection)
	log.Println("User", string(uPtr.Username), "Set Avatar to", pkt.AvatarId)
	//如果是在房间内
}

func BuildSetAvatarPacket(id uint32, avatar uint16) []byte {
	buf := make([]byte, 10)
	offset := 0
	WriteUint32(&buf, id, &offset)
	WriteUint32(&buf, 0x800000, &offset)
	WriteUint16(&buf, avatar, &offset)
	return buf[:offset]
}
