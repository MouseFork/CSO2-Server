package servermanager

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func DelUserWithConn(con net.Conn) bool {
	if UsersManager.UserNum == 0 {
		DebugInfo(2, "UsersManager Error : There is no online user !")
		return false
	}
	for k, v := range UsersManager.Users {
		if v.CurrentConnection == con {
			CheckErr(UpdateUserToDB(v))
			rm := GetRoomFromID(v.GetUserChannelServerID(),
				v.GetUserChannelID(),
				v.GetUserRoomID())
			if rm != nil &&
				rm.Id > 0 {
				rm.RoomRemoveUser(v.Userid)
				if rm.NumPlayers <= 0 {
					DelChannelRoom(rm.Id,
						v.GetUserChannelID(),
						v.GetUserChannelServerID())

				} else {
					//p.id = TypeRoom
					//SentUserLeaveMes(&v, rm, p)
				}
			}
			UsersManager.Lock.Lock()
			defer UsersManager.Lock.Unlock()
			delete(UsersManager.Users, k)
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
	//如果map中不存在该ID，则返回ID
	UsersManager.Lock.Lock()
	defer UsersManager.Lock.Unlock()
	for i := 1; i <= MAXUSERNUM; i++ {
		if _, ok := UsersManager.Users[uint32(i)]; !ok {
			return uint32(i)
		}
	}
	return 0
}

//getUserByLogin 假定nexonUsername是唯一
func GetUserByLogin(account, passwd []byte) *User {
	//查看是否有已经登陆的同名用户
	for _, v := range UsersManager.Users {
		if string(v.NexonUsername) == string(account) {
			return nil
		}
	}
	//查看数据库是否有该用户
	return GetUserFromDatabase(account, passwd)
}

//通过连接获取用户
func GetUserFromConnection(client net.Conn) *User {
	if UsersManager.UserNum <= 0 {
		return nil
	}
	for _, v := range UsersManager.Users {
		if v.CurrentConnection == client {
			return v
		}
	}
	return nil
}

//通过ID获取用户
func GetUserFromID(id uint32) *User {
	if UsersManager.UserNum <= 0 {
		return nil
	}
	if v, ok := UsersManager.Users[id]; ok {
		return v
	}
	return nil
}
