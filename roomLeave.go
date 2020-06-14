package main

import (
	"log"
	"net"
)

func onLeaveRoom(seq *uint8, p packet, client net.Conn) {
	//找到玩家
	uPtr := getUserFromConnection(client)
	if uPtr.userid <= 0 {
		log.Println("Client from", client.RemoteAddr().String(), "try to leave room but not in server !")
		return
	}
	//找到玩家的房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm.id <= 0 {
		log.Println("Client from", client.RemoteAddr().String(), "try to leave a null room !")
		return
	}
	//检查玩家游戏状态，准备情况下并且开始倒计时了，那么就不允许离开房间
	if uPtr.isUserReady() &&
		rm.isGlobalCountdownInProgress() {
		log.Println("Client from", client.RemoteAddr().String(), "try to leave room but is started !")
		return
	}
	//房间移除玩家
	rm.roomRemoveUser(*uPtr)
	//设置玩家状态
	p.datalen = 7
	p.data[5] = uPtr.getUserChannelServerID()
	p.data[6] = uPtr.getUserChannelID()
	uPtr.quitRoom()
	//发送房间列表给玩家
	onRoomList(seq, &p, client)
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
