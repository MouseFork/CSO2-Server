package quick

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnQuickList(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InQuickList
	if !p.PraseInQuickListPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error QuickList packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request QuickList but not in server !")
		return
	}
	//发送房间数据,暂时发送空数据
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, p.Id), BuildQuickList(pkt))
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Sent a null quickList to User", string(uPtr.Username))
}

func BuildQuickList(pkt InQuickList) []byte {
	buf := make([]byte, 2)
	offset := 0
	WriteUint8(&buf, QuickList, &offset)
	WriteUint8(&buf, 0, &offset) //num of room

	return buf[:offset]
}
