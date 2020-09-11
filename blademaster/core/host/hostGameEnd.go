package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/room"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostGameEnd(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : A user request to send GameEnd but not in server!")
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to send GameEnd but in a null room !")
		return
	}
	//是不是房主
	if rm.HostUserID != uPtr.Userid {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to send GameEnd but isn't host !")
		return
	}
	//修改房间信息
	rm.SetStatus(StatusWaiting)
	header := BuildGameResultHeader(*rm)
	for _, v := range rm.Users {
		//修改用户状态
		v.SetUserStatus(UserNotReady)
		//发送房间状态
		rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), BuildRoomSetting(rm, 0xFFFFFFFFFFFFFFFF))
		SendPacket(rst, v.CurrentConnection)
		//检查是否还在游戏内
		if v.CurrentIsIngame {
			//发送游戏结束数据包
			rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeHost), BuildHostStop())
			SendPacket(rst, v.CurrentConnection)
			//发送游戏战绩
			rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), header, BuildGameResult(v))
			SendPacket(rst, v.CurrentConnection)
			DebugInfo(2, "Sent game result to User", string(v.Username))
			//修改用户状态
			v.SetUserIngame(false)
		}
	}
	//给每个人发送房间内所有人的准备状态
	for _, v := range rm.Users {
		rst := BuildUserReadyStatus(v)
		for _, k := range rm.Users {
			rst = BytesCombine(BuildHeader(k.CurrentSequence, PacketTypeRoom), rst)
			SendPacket(rst, k.CurrentConnection)
		}
	}
	rm.ResetRoomKillNum()
	rm.ResetRoomScore()
	rm.ResetRoomWinner()
}

func BuildHostStop() []byte {
	return []byte{HostStop}
}
func BuildGameResult(u *User) []byte {
	buf := make([]byte, 40)
	offset := 0
	WriteUint64(&buf, u.CurrentExp, &offset) //now total EXP
	WriteUint64(&buf, u.Points, &offset)     //now total point
	WriteUint8(&buf, 0, &offset)             //unk18
	WriteString(&buf, []byte("unk19"), &offset)
	WriteString(&buf, []byte("unk20"), &offset)
	WriteUint16(&buf, 0, &offset) //unk21 ，maybe 2 bytes
	WriteUint16(&buf, 0, &offset) //unk22 ，maybe 2 bytes
	return buf[:offset]
}

func BuildGameResultHeader(rm Room) []byte {
	buf := make([]byte, 30)
	offset := 0
	WriteUint8(&buf, OUTSetGameResult, &offset)
	WriteUint8(&buf, 0, &offset)                     //unk01
	WriteUint8(&buf, rm.Setting.GameModeID, &offset) //game mod？ 0x02 0x01
	switch rm.Setting.GameModeID {
	case ModeOriginal, ModePig:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint8(&buf, 0, &offset)             //上半局CT winNum，开启阵营互换情况
		WriteUint8(&buf, 0, &offset)             //上半局TR winNum
		//WriteUint16(&buf, 0, &offset)            //unk00
	case ModeStealth:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint8(&buf, 0, &offset)             //上半场CT winNum?
		WriteUint8(&buf, 0, &offset)             //上半场TR winNum?
	case ModeDeathmatch, ModeTeamdeath, ModeTeamdeath_mutation:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint32(&buf, rm.CtKillNum, &offset) //CT killnum
		WriteUint32(&buf, rm.TrKillNum, &offset) //TR killnum
		WriteUint64(&buf, 0, &offset)            //unk02
	case ModeGhost:
		WriteUint8(&buf, 0, &offset)
		WriteUint32(&buf, 0, &offset)
	case ModeZombie, ModeZombiecraft, ModeZombie_commander, ModeZombie_prop, ModeZombie_zeta, ModeZ_scenario, ModeZ_scenario_side, ModeHeroes:

	default:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint16(&buf, 0, &offset)            //unk00
	}
	WriteUint8(&buf, uint8(rm.GetNumOfRealReadyPlayers()), &offset) //usernum？
	buf = buf[:offset]
	for k, v := range rm.Users {
		if v.CurrentIsIngame {
			temp := make([]byte, 100)
			offset = 0
			WriteUint32(&temp, v.Userid, &offset)           //userid
			WriteUint8(&temp, 0, &offset)                   //unk01
			WriteUint64(&temp, 0, &offset)                  //unk03
			WriteUint64(&temp, 0, &offset)                  //unk04
			WriteUint16(&temp, v.CurrentKillNum, &offset)   //killnum
			WriteUint16(&temp, v.CurrentAssistNum, &offset) //assistnum
			WriteUint16(&temp, v.CurrentDeathNum, &offset)  //deathnum
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
			WriteUint8(&temp, uint8(v.Level), &offset)      //current level ?
			WriteUint8(&temp, uint8(v.Level), &offset)      //next level ？
			WriteUint8(&temp, 0, &offset)                   //unk15
			WriteUint8(&temp, uint8(k+1), &offset)          //rank
			WriteUint16(&temp, v.CurrentKillNum, &offset)   //连续击杀数
			WriteUint32(&temp, 0, &offset)                  //unk16 ，maybe 4 bytes
			WriteUint8(&temp, v.CurrentTeam, &offset)       //user team
			switch rm.Setting.GameModeID {
			case ModeOriginal, ModePig:
				WriteUint32(&temp, 0, &offset) //unk17
			case ModeDeathmatch, ModeTeamdeath, ModeTeamdeath_mutation:
			case ModeStealth:
				WriteUint16(&temp, 0, &offset) //unk17
			case ModeZombie, ModeZombiecraft, ModeZombie_commander, ModeZombie_prop, ModeZombie_zeta, ModeGhost, ModeZ_scenario, ModeZ_scenario_side, ModeHeroes:
			default:
				WriteUint32(&temp, 0, &offset) //unk17,貌似有时候不用
			}
			buf = BytesCombine(buf, temp[:offset])
		}
	}
	return buf
}
