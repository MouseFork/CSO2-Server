package user

import (
	"net"
	"sync"

	. "github.com/KouKouChan/CSO2-Server/model/inventory"
)

type (
	User struct {
		//个人信息
		Userid               uint32
		LoginName            []byte
		Username             []byte
		Password             []byte
		Level                uint16
		Rank                 uint8
		RankFrame            uint8
		Points               uint64
		CurrentExp           uint64
		MaxExp               uint64
		PlayedMatches        uint32
		Wins                 uint32
		Kills                uint32
		Headshots            uint32
		Deaths               uint32
		Assists              uint32
		Accuracy             uint16
		SecondsPlayed        uint32
		NetCafeName          []byte
		Cash                 uint32
		ClanName             []byte
		ClanMark             uint32
		WorldRank            uint32
		Mpoints              uint32
		TitleId              uint16
		UnlockedTitles       []byte
		Signature            []byte
		BestGamemode         uint32
		BestMap              uint32
		UnlockedAchievements []byte
		Avatar               uint16
		UnlockedAvatars      []byte
		VipLevel             uint8
		VipXp                uint32
		SkillHumanCurXp      uint64
		SkillHumanMaxXp      uint64
		SkillHumanPoints     uint8
		SkillZombieCurXp     uint64
		SkillZombieMaxXp     uint64
		SkillZombiePoints    uint8
		//连接
		CurrentConnection net.Conn
		//频道房间信息
		CurrentChannelServerIndex uint8
		CurrentChannelIndex       uint8
		CurrentRoomId             uint16
		CurrentTeam               uint8
		Currentstatus             uint8
		CurrentIsIngame           bool
		CurrentSequence           *uint8
		CurrentKillNum            uint16
		CurrentDeathNum           uint16
		CurrentAssistNum          uint16
		NetInfo                   UserNetInfo
		//仓库信息
		Inventory UserInventory

		UserMutex *sync.Mutex
	}

	UserNetInfo struct {
		ExternalIpAddress  uint32
		ExternalClientPort uint16
		ExternalServerPort uint16
		ExternalTvPort     uint16

		LocalIpAddress  uint32
		LocalClientPort uint16
		LocalServerPort uint16
		LocalTvPort     uint16
	}
)

const (
	//MAXUSERNUM 最大用户数
	MAXUSERNUM = 1024 //房间状态

	UserNotReady = 0
	UserIngame   = 1
	UserReady    = 2
)

func (u user) isVIP() bool {
	if u.vipLevel <= 0 {
		return false
	}
	return true
}

func (u *user) setID(id uint32) {
	if u == nil {
		return
	}
	(*u).userid = id
}

func (u *user) setUserName(p loginPacket) {
	if u == nil {
		return
	}
	(*u).loginName = p.nexonUsername
	(*u).username = p.gameUsername
}

func (u *user) setUserChannelServer(id uint8) {
	if u == nil {
		return
	}
	(*u).currentChannelServerIndex = id
}

func (u *user) setUserChannel(id uint8) {
	if u == nil {
		return
	}
	(*u).currentChannelIndex = id
}

func (u *user) setUserRoom(id uint16) {
	if u == nil {
		return
	}
	(*u).currentRoomId = id
}

func (u *user) quitChannel() {
	if u == nil {
		return
	}
	(*u).currentChannelIndex = 0
}

func (u *user) quitRoom() {
	if u == nil {
		return
	}
	(*u).currentRoomId = 0
	(*u).currentTeam = Unknown
	(*u).currentstatus = UserNotReady
	(*u).currentIsIngame = false
}

func (u *user) setUserStatus(status uint8) {
	if u == nil {
		return
	}
	if status <= 2 &&
		status >= 0 {
		(*u).currentstatus = status
	}
}

//获取用户所在分区服务器ID
func (u user) getUserChannelServerID() uint8 {
	if u.userid <= 0 {
		return 0
	}
	return u.currentChannelServerIndex
}

//获取用户所在频道ID
func (u user) getUserChannelID() uint8 {
	if u.userid <= 0 {
		return 0
	}
	return u.currentChannelIndex
}

//获取用户所在房间ID
func (u user) getUserRoomID() uint16 {
	if u.userid <= 0 {
		return 0
	}
	return u.currentRoomId
}

func (u user) getUserTeam() uint8 {
	return u.currentTeam
}

func (u user) isUserReady() bool {
	return u.currentstatus == UserReady
}

func (u *user) setUserIngame(ingame bool) {
	if u == nil {
		return
	}
	(*u).currentIsIngame = ingame
	if ingame {
		(*u).currentstatus = UserIngame
	} else {
		(*u).currentstatus = UserNotReady
	}

}

func (u *user) updateHolepunch(portId uint16, localPort uint16, externalPort uint16) uint16 {
	if u == nil {
		return 0xFFFF
	}
	switch portId {
	case UDPTypeClient:
		(*u).netInfo.LocalClientPort = localPort
		(*u).netInfo.ExternalClientPort = externalPort
		return 0
	case UDPTypeServer:
		(*u).netInfo.LocalServerPort = localPort
		(*u).netInfo.ExternalServerPort = externalPort
		return 1
	case UDPTypeSourceTV:
		(*u).netInfo.LocalTvPort = localPort
		(*u).netInfo.ExternalTvPort = externalPort
		return 2
	default:
		return 0xFFFF
	}
}

func (u *user) CountKillNum(num uint16) {
	if u == nil {
		return
	}
	(*u).currentKillNum = num
}

func (u *user) CountDeadNum(num uint16) {
	if u == nil {
		return
	}
	(*u).currentDeathNum = num
}
func (u *user) CountAssistNum() {
	if u == nil {
		return
	}
	(*u).currentAssistNum++
}
func (u *user) ResetAssistNum() {
	if u == nil {
		return
	}
	(*u).currentAssistNum = 0
}
func (u *user) ResetKillNum() {
	if u == nil {
		return
	}
	(*u).currentKillNum = 0
}

func (u *user) ResetDeadNum() {
	if u == nil {
		return
	}
	(*u).currentDeathNum = 0
}

func (u *user) SetSignature(sig []byte) {
	if u == nil {
		return
	}
	(*u).signature = sig
}

func (u *user) SetAvatar(id uint16) {
	if u == nil {
		return
	}
	(*u).avatar = id
}

func (u *user) SetTitle(id uint16) {
	if u == nil {
		return
	}
	(*u).titleId = id
}
