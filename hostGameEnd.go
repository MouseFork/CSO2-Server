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
	header := BuildGameResultHeader(*rm)
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
			//发送游戏战绩
			p.id = TypeRoom
			rst = BytesCombine(BuildHeader(v.currentSequence, p), header, BuildGameResult(*u))
			sendPacket(rst, v.currentConnection)
			log.Println("Sent game result to User", string(u.username))
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
	rm.resetRoomKillNum()
	rm.resetRoomScore()
	rm.resetRoomWinner()
}

func BuildHostStop() []byte {
	return []byte{HostStop}
}
func BuildGameResult(u user) []byte {
	buf := make([]byte, 40)
	offset := 0
	WriteUint64(&buf, u.currentExp, &offset) //now total EXP
	WriteUint64(&buf, u.points, &offset)     //now total point
	WriteUint8(&buf, 0, &offset)             //unk18
	WriteString(&buf, []byte("unk19"), &offset)
	WriteString(&buf, []byte("unk20"), &offset)
	WriteUint16(&buf, 0, &offset) //unk21 ，maybe 2 bytes
	WriteUint16(&buf, 0, &offset) //unk22 ，maybe 2 bytes
	return buf[:offset]
}

func BuildGameResultHeader(rm roomInfo) []byte {
	buf := make([]byte, 30)
	offset := 0
	WriteUint8(&buf, OUTSetGameResult, &offset)
	WriteUint8(&buf, 0, &offset)                     //unk01
	WriteUint8(&buf, rm.setting.gameModeID, &offset) //game mod？ 0x02 0x01
	switch rm.setting.gameModeID {
	case original, pig, stealth:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint16(&buf, 0, &offset)            //unk00
	case deathmatch, teamdeath, teamdeath_mutation:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint32(&buf, rm.CtKillNum, &offset) //CT killnum
		WriteUint32(&buf, rm.TrKillNum, &offset) //TR killnum
		WriteUint64(&buf, 0, &offset)            //unk02
	case ghost:
		WriteUint8(&buf, 0, &offset)
		WriteUint32(&buf, 0, &offset)
	case zombie, zombiecraft, zombie_commander, zombie_prop, zombie_zeta:

	default:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint16(&buf, 0, &offset)            //unk00
	}
	WriteUint8(&buf, uint8(rm.getNumOfRealReadyPlayers()), &offset) //usernum？
	buf = buf[:offset]
	for k, v := range rm.users {
		if v.currentIsIngame {
			u := getUserFromID(v.userid)
			if u == nil ||
				u.userid <= 0 {
				continue
			}
			temp := make([]byte, 100)
			offset = 0
			WriteUint32(&temp, u.userid, &offset)           //userid
			WriteUint8(&temp, 0, &offset)                   //unk01
			WriteUint64(&temp, 0, &offset)                  //unk03
			WriteUint64(&temp, 0, &offset)                  //unk04
			WriteUint16(&temp, u.currentKillNum, &offset)   //killnum
			WriteUint16(&temp, u.currentAssistNum, &offset) //assistnum
			WriteUint16(&temp, u.currentDeathNum, &offset)  //deathnum
			WriteUint16(&temp, 0, &offset)                  //unk05 ，maybe 2 bytes
			WriteUint16(&temp, 0, &offset)                  //unk06 ，maybe 2 bytes 0x56 = 86
			WriteUint16(&temp, 0, &offset)                  //unk07 ，maybe 2 bytes 0x2b = 43
			WriteUint64(&temp, 100, &offset)                //gained EXP
			WriteUint32(&temp, 0, &offset)                  //unk08 ，maybe 4 bytes
			WriteUint16(&temp, 0, &offset)                  //unk09 ，maybe 2 bytes
			WriteUint8(&temp, 0, &offset)                   //unk10 ，maybe 1 bytes
			WriteUint64(&temp, 100, &offset)                //gained point
			WriteUint32(&temp, 0, &offset)                  //unk11 ，maybe 4 bytes
			WriteUint16(&temp, 0, &offset)                  //unk12 ，maybe 2 bytes
			WriteUint8(&temp, 0, &offset)                   //unk13 ，maybe 1 bytes
			WriteUint8(&temp, uint8(u.level), &offset)      //current level ?
			WriteUint8(&temp, uint8(u.level), &offset)      //next level ？
			WriteUint8(&temp, 0, &offset)                   //unk15
			WriteUint8(&temp, uint8(k+1), &offset)          //rank
			WriteUint16(&temp, u.currentKillNum, &offset)   //连续击杀数
			WriteUint32(&temp, 0, &offset)                  //unk16 ，maybe 4 bytes
			WriteUint8(&temp, u.currentTeam, &offset)       //user team
			switch rm.setting.gameModeID {
			case original, pig, stealth:
				WriteUint32(&temp, 0, &offset) //unk17,貌似有时候不用
			case deathmatch, teamdeath, teamdeath_mutation:
			case zombie, zombiecraft, zombie_commander, zombie_prop, zombie_zeta, ghost:
			default:
				WriteUint32(&temp, 0, &offset) //unk17,貌似有时候不用
			}
			buf = BytesCombine(buf, temp[:offset])
		}
	}
	return buf
}
