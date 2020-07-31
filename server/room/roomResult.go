package handler

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func onCloseResultRequest(seq *uint8, p packet, client net.Conn) {
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to close result but not in server !")
		return
	}
	//发送数据
	p.id = TypeHost
	rst := BytesCombine(BuildHeader(uPtr.currentSequence, p), BuildCloseResultWindow())
	sendPacket(rst, uPtr.currentConnection)
	DebugInfo(2, "User", string(uPtr.username), "closed game result window from room id", uPtr.currentRoomId)
}

func BuildCloseResultWindow() []byte {
	buf := make([]byte, 1)
	buf[0] = LeaveResultWindow
	return buf
}
