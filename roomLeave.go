package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

func onLeaveRoom(seq *uint8, p packet, client net.Conn) {
	//找到玩家
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to leave room but not in server !")
		return
	}
	//找到玩家的房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User:", string(uPtr.username), "try to leave a null room !")
		return
	}
	//检查玩家游戏状态，准备情况下并且开始倒计时了，那么就不允许离开房间
	if uPtr.isUserReady() &&
		rm.isGlobalCountdownInProgress() {
		log.Println("Error : User:", string(uPtr.username), "try to leave room but is started !")
		return
	}
	//房间移除玩家
	rm.roomRemoveUser(*uPtr)
	//检查房间是否为空
	if rm.numPlayers <= 0 {
		delChannelRoom(rm.id,
			uPtr.getUserChannelID(),
			uPtr.getUserChannelServerID())

	} else {
		//向其他玩家发送离开信息
		//如果玩家是房主
		if rm.hostUserID == uPtr.userid {
			(*rm).hostUserID = rm.users[0].userid
			for _, v := range rm.users {
				rst1 := append(BuildHeader(v.currentSequence, p), OUTPlayerLeave)
				rst1 = BytesCombine(rst1, BuildUserLeave(uPtr.userid))
				rst2 := append(BuildHeader(v.currentSequence, p), OUTSetHost)
				rst2 = BytesCombine(rst2, BuildSetHost(rm.hostUserID))
				sendPacket(rst1, v.currentConnection)
				sendPacket(rst2, v.currentConnection)
			}
			log.Println("Sent a set roomHost packet to other users")
		} else {
			for _, v := range rm.users {
				rst1 := append(BuildHeader(v.currentSequence, p), OUTPlayerLeave)
				rst1 = BytesCombine(rst1, BuildUserLeave(uPtr.userid))
				sendPacket(rst1, v.currentConnection)
			}
		}
		log.Println("Sent a leave room packet to other users")
	}
	//设置玩家状态
	p.datalen = 7
	p.data[5] = uPtr.getUserChannelServerID()
	p.data[6] = uPtr.getUserChannelID()
	uPtr.quitRoom()
	//发送房间列表给玩家
	onRoomList(seq, &p, client)
	log.Println("User:", string(uPtr.username), "left room")
}

func BuildUserLeave(id uint32) []byte {
	buf := make([]byte, 4)
	offset := 0
	WriteUint32(&buf, id, &offset)
	return buf
}
func BuildSetHost(id uint32) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint32(&buf, id, &offset)
	buf[4] = 0
	return buf
}
func (rm *roomInfo) roomRemoveUser(u user) {
	if rm.numPlayers <= 0 {
		return
	}
	//找到玩家,玩家数-1，删除房间玩家
	for k, v := range rm.users {
		if v.userid == u.userid {
			(*rm).users = append(rm.users[:k], rm.users[k+1:]...)
			(*rm).numPlayers--
			return
		}
	}
}
