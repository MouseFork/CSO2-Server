package main

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

func onHostGameEnd(p packet, client net.Conn) {
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : A user request to send GameEnd but not in server!")
		return
	}
	//找到玩家的房间
	rm := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if rm == nil ||
		rm.id <= 0 {
		log.Println("Error : User", string(uPtr.username), "try to send GameEnd but in a null room !")
		return
	}
	//是不是房主
	if rm.hostUserID != uPtr.userid {
		log.Println("Error : User", string(uPtr.username), "try to send GameEnd but isn't host !")
		return
	}
	//修改房间信息
	rm.setStatus(StatusWaiting)
	for k, v := range rm.users {
		p.id = TypeRoom
		u := getUserFromID(v.userid)
		if u == nil ||
			u.userid <= 0 {
			continue
		}
		//修改用户状态
		u.setUserStatus(UserNotReady)
		(*rm).users[k].currentstatus = UserNotReady
		//发送房间状态
		rst := BytesCombine(BuildHeader(v.currentSequence, p), buildRoomSetting(*rm))
		sendPacket(rst, v.currentConnection)
		//检查是否还在游戏内
		if u.currentIsIngame {
			//发送游戏结束数据包
			p.id = TypeHost
			rst = BytesCombine(BuildHeader(v.currentSequence, p), BuildHostStop())
			sendPacket(rst, v.currentConnection)
			//发送游戏结果
			p.id = TypeRoom
			rst = BytesCombine(BuildHeader(v.currentSequence, p), BuildGameResult())
			sendPacket(rst, v.currentConnection)
			//修改用户状态
			(*rm).users[k].currentIsIngame = false
			u.setUserIngame(false)
		}
	}
	//给每个人发送房间内所有人的准备状态
	for _, v := range rm.users {
		rst := BuildUserReadyStatus(v)
		for _, k := range rm.users {
			rst = BytesCombine(BuildHeader(k.currentSequence, p), rst)
			sendPacket(rst, k.currentConnection)
		}
	}

}

func BuildHostStop() []byte {
	return []byte{HostStop}
}

func BuildGameResult() []byte {
	buf := make([]byte, 38)
	offset := 0
	WriteUint8(&buf, OUTSetGameResult, &offset)
	WriteUint8(&buf, 0, &offset)  //unk01
	WriteUint8(&buf, 0, &offset)  //unk02
	WriteUint8(&buf, 0, &offset)  //unk03
	WriteUint64(&buf, 0, &offset) //unk04
	WriteUint64(&buf, 0, &offset) //unk05
	WriteUint8(&buf, 0, &offset)  //unk06
	WriteString(&buf, []byte("unk07"), &offset)
	WriteString(&buf, []byte("unk08"), &offset)
	WriteUint8(&buf, 0, &offset) //unk09
	WriteUint8(&buf, 0, &offset) //unk10
	return buf[:offset]
}
