package blademaster

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

//全局用户管理
type UserManager struct {
	UserNum int
	Users   []User
}

var (
	//UserManager 全局用户管理
	UsersManager = UserManager{
		0,
		[]User{},
	}
)

func (dest *UserManager) AddUser(src *User) bool {
	if dest == nil {
		return false
	}
	if src.Userid == 0 {
		DebugInfo(2, "Error : ID of User", (*src).Username, "is illegal !")
		return false
	}
	if dest.UserNum > MAXUSERNUM {
		DebugInfo(2, "Error : Online users is too more to login !")
		return false
	}
	for _, v := range dest.Users {
		if v.Userid == src.Userid {
			DebugInfo(2, "Error : User is already logged in !")
			return false
		}
	}
	UserManagerMutex.Lock()
	defer UserManagerMutex.Unlock()
	dest.UserNum++
	dest.Users = append(dest.Users, *src)
	return true
}

func (dest *UserManager) DelUser(src *User) bool {
	if dest == nil {
		return false
	}
	if src.Userid == 0 {
		DebugInfo(2, "Error : ID of User", (*src).Username, "is illegal !")
		return false
	}
	if dest.UserNum == 0 {
		DebugInfo(2, "Error : There is no online user !")
		return false
	}
	for i, v := range dest.Users {
		if v.Userid == src.Userid {
			UserManagerMutex.Lock()
			defer UserManagerMutex.Unlock()
			dest.Users = append(dest.Users[:i], dest.Users[i+1:]...)
			dest.UserNum--
			return true
		}
	}
	return false
}

//通过连接获取用户
func GetUserFromConnection(client net.Conn) *User {
	if UsersManager.UserNum <= 0 {
		return nil
	}
	for k, v := range UsersManager.Users {
		if v.CurrentConnection == client {
			return &UsersManager.Users[k]
		}
	}
	return nil
}

//通过ID获取用户
func GetUserFromID(id uint32) *User {
	if UsersManager.UserNum <= 0 {
		return nil
	}
	for k, v := range UsersManager.Users {
		if v.Userid == id {
			return &UsersManager.Users[k]
		}
	}
	return nil
}
