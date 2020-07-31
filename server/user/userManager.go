package user

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/model/user"
	. "github.com/KouKouChan/CSO2-Server/model/usermanager"
)

func delUserWithConn(con net.Conn) bool {
	if UsersManager.UserNum == 0 {
		log.Println("Error : There is no online user !")
		return false
	}
	for i, v := range UsersManager.Users {
		if v.CurrentConnection == con {
			CheckErr(UpdateUserToDB(v))
			rm := getRoomFromID(v.getUserChannelServerID(),
				v.getUserChannelID(),
				v.getUserRoomID())
			if rm != nil &&
				rm.id > 0 {
				rm.roomRemoveUser(v)
				if rm.numPlayers <= 0 {
					delChannelRoom(rm.id,
						v.getUserChannelID(),
						v.getUserChannelServerID())

				} else {
					var p packet
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

func getNewUserID() uint32 {
	if UserManager.userNum > MAXUSERNUM {
		log.Println("Online users is too much , unable to get a new id !")
		//ID=0 是非法的
		return 0
	}
	var intbuf [MAXUSERNUM + 2]uint32
	//哈希思想
	for i := 0; i < int(UserManager.userNum); i++ {
		intbuf[UserManager.users[i].userid] = 1
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
func getUserByLogin(p loginPacket) user {
	//查看是否有已经登陆的同名用户，待定
	for _, v := range UserManager.users {
		if string(v.username) == string(p.nexonUsername) {
			return v
		}
	}
	//查看数据库是否有该用户
	return getUserFromDatabase(p)
}

//通过连接获取用户
func getUserFromConnection(client net.Conn) *user {
	if UserManager.userNum <= 0 {
		return nil
	}
	for k, v := range UserManager.users {
		if v.currentConnection == client {
			return &UserManager.users[k]
		}
	}
	return nil
}

//通过ID获取用户
func getUserFromID(id uint32) *user {
	if UserManager.userNum <= 0 {
		return nil
	}
	for k, v := range UserManager.users {
		if v.userid == id {
			return &UserManager.users[k]
		}
	}
	return nil
}
