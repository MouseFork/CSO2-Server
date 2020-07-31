package usermanager

import (
	"log"

	. "github.com/KouKouChan/CSO2-Server/model/lock"
	. "github.com/KouKouChan/CSO2-Server/model/user"
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

func (dest *UserManager) addUser(src *User) bool {
	if dest == nil {
		return false
	}
	if src.Userid == 0 {
		log.Println("Error : ID of User", (*src).Username, "is illegal !")
		return false
	}
	if dest.UserNum > MAXUSERNUM {
		log.Println("Error : Online users is too more to login !")
		return false
	}
	for _, v := range dest.Users {
		if v.Userid == src.Userid {
			log.Println("Error : User is already logged in !")
			return false
		}
	}
	UserManagerMutex.Lock()
	defer UserManagerMutex.Unlock()
	dest.UserNum++
	dest.Users = append(dest.Users, *src)
	return true
}

func (dest *UserManager) delUser(src *User) bool {
	if dest == nil {
		return false
	}
	if src.Userid == 0 {
		log.Println("Error : ID of User", (*src).Username, "is illegal !")
		return false
	}
	if dest.UserNum == 0 {
		log.Println("Error : There is no online user !")
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
