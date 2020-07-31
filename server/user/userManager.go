package user

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/server/database"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func DelUserWithConn(con net.Conn) bool {
	if UsersManager.UserNum == 0 {
		DebugInfo(2, "Error : There is no online user !")
		return false
	}
	for i, v := range UsersManager.Users {
		if v.CurrentConnection == con {
			CheckErr(UpdateUserToDB(v))
			rm := GetRoomFromID(v.GetUserChannelServerID(),
				v.GetUserChannelID(),
				v.GetUserRoomID())
			if rm != nil &&
				rm.id > 0 {
				rm.roomRemoveUser(v)
				if rm.numPlayers <= 0 {
					delChannelRoom(rm.id,
						v.GetUserChannelID(),
						v.GetUserChannelServerID())

				} else {
					var p Packet
					p.id = TypeRoom
					sentUserLeaveMes(&v, rm, p)
				}
			}
			UsersManager.Users = append(UsersManager.Users[:i], UsersManager.Users[i+1:]...)
			UsersManager.UserNum--
			return true
		}
	}
	return false
}

func GetNewUserID() uint32 {
	if UsersManager.UserNum > MAXUSERNUM {
		DebugInfo(2, "Online users is too much , unable to get a new id !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXUSERNUM + 2]uint32
	//哈希思想
	for i := 0; i < int(UsersManager.UserNum); i++ {
		intbuf[UsersManager.Users[i].Userid] = 1
	}
	//找到空闲的ID
	for i := 1; i < int(MAXUSERNUM+2); i++ {
		if intbuf[i] == 0 {
			//找到了空闲ID
			return uint32(i)
		}
	}
	return 0
}

//getUserByLogin 假定nexonUsername是唯一
func GetUserByLogin(p loginPacket) User {
	//查看是否有已经登陆的同名用户，待定
	for _, v := range UsersManager.Users {
		if string(v.Username) == string(p.nexonUsername) {
			return v
		}
	}
	//查看数据库是否有该用户
	return getUserFromDatabase(p)
}
