package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnCloseResultRequest(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to close result but not in server !")
		return
	}
	//发送数据
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeHost), BuildCloseResultWindow())
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "User", string(uPtr.Username), "closed game result window from room id", uPtr.CurrentRoomId)

	//在线时间奖励或进度等
	//...
}

func BuildCloseResultWindow() []byte {
	buf := make([]byte, 1)
	buf[0] = LeaveResultWindow
	return buf
}
