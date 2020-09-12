package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnGameStartCountdown(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InRoomCountdownPacket
	if !p.PraseRoomCountdownPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error Countdown packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to start counting but not in server !")
		return
	}
	//检查用户是不是房主
	curroom := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if curroom == nil ||
		curroom.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to start counting in a null room !")
		return
	}
	if curroom.HostUserID != uPtr.Userid {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to start counting but is not the host !")
		return
	}
	//检查用户所在房间
	if curroom.Id != uPtr.CurrentRoomId {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to start counting but in another room !")
		return
	}
	//检查当前游戏模式
	if !curroom.CanStartGame() {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to start countdown but mode is illegal !")
		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, DialogBox,
			GAME_ROOM_COUNT_MODE_ERROR)
		return
	}
	//检查是否应该继续倒计时
	should := pkt.ShouldCountdown()
	if should {
		curroom.ProgressCountdown(pkt.Count)
		DebugInfo(2, "User", string(uPtr.Username), "countdown at", curroom.GetCountdown(), "host is", pkt.Count)
	} else {
		curroom.StopCountdown()
		DebugInfo(2, "User", string(uPtr.Username), "cancled room countdown")
	}
	//所有玩家发送倒计时数据
	build := BuildCountdown(pkt, should)
	for _, v := range curroom.Users {
		rst := BytesCombine(BuildHeader(v.CurrentSequence, p.Id), build)
		SendPacket(rst, v.CurrentConnection)
	}
}

func BuildCountdown(p InRoomCountdownPacket, should bool) []byte {
	buf := make([]byte, 3)
	offset := 0
	WriteUint8(&buf, OUTCountdown, &offset)
	WriteUint8(&buf, p.CountdownType, &offset)
	if should {
		WriteUint8(&buf, p.Count, &offset)
	}
	return buf[:offset]
}
