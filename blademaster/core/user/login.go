package user

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnLogin(seq *uint8, dataPacket *PacketData, client net.Conn) {
	var pkt InLoginPacket
	if !dataPacket.PraseLoginPacket(&pkt) {
		DebugInfo(2, "Error : User from", client.RemoteAddr().String(), "Sent a illegal login packet !")
		return
	}

	u, result := GetUserByLogin(pkt.NexonUsername, pkt.PassWd)
	switch result {
	case USER_PASSWD_ERROR:
		DebugInfo(2, "Error : User", string(pkt.NexonUsername), "from", client.RemoteAddr().String(), "login failed with error password !")
		OnSendMessage(seq, client, DialogBox, GAME_LOGIN_BAD_PASSWORD)
		return
	case USER_ALREADY_LOGIN:
		DebugInfo(2, "Error : User", string(pkt.NexonUsername), "from", client.RemoteAddr().String(), "already logged in !")
		OnSendMessage(seq, client, DialogBox, GAME_LOGIN_ALREADY)
		return
	case USER_NOT_FOUND:
		DebugInfo(2, "Error : User", string(pkt.NexonUsername), "from", client.RemoteAddr().String(), "not registered !")
		OnSendMessage(seq, client, DialogBox, GAME_LOGIN_BAD_USERNAME)
		return
	case USER_UNKOWN_ERROR:
		DebugInfo(2, "Error : User", string(pkt.NexonUsername), "from", client.RemoteAddr().String(), "login but suffered a error !")
		OnSendMessage(seq, client, DialogBox, GAME_LOGIN_ERROR)
		return
	default:
	}

	u.CurrentConnection = client
	u.CurrentSequence = seq

	//把用户加入用户管理器
	if !UsersManager.AddUser(u) {
		DebugInfo(2, "Error : User", string(pkt.NexonUsername), "from", client.RemoteAddr().String(), "login failed !")
		return
	}

	//UserStart部分
	rst := BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeUserStart), BuildUserStart(u))
	SendPacket(rst, u.CurrentConnection)
	DebugInfo(1, "User", string(u.NexonUsername), "from", client.RemoteAddr().String(), "logged in !")

	//UserInfo部分
	rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(u), u.Userid, true))
	SendPacket(rst, u.CurrentConnection)

	//ServerList部分
	//OnServerList(u.CurrentConnection)

	//Inventory部分
	rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeInventory_Create),
		BuildInventoryInfo(u))
	SendPacket(rst, u.CurrentConnection)

	//unlock
	rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeUnlock), BuildUnlockReply())
	SendPacket(rst, u.CurrentConnection)

	//偏好装备
	rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeFavorite), BuildCosmetics(&u.Inventory))
	SendPacket(rst, u.CurrentConnection)
	rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeFavorite), BuildLoadout(&u.Inventory))
	SendPacket(rst, u.CurrentConnection)

	//购买菜单
	rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeOption), BuildBuyMenu(&u.Inventory))
	SendPacket(rst, u.CurrentConnection)

	//achievement

	//friends
}

//BuildUserStart 返回结构
// userId
// loginName
// userName
// unk00
// holepunchPort
func BuildUserStart(u *User) []byte {
	//暂时都取GameUsername
	userbuf := make([]byte, 9+int(len(u.NexonUsername))+int(len(u.Username)))
	offset := 0
	WriteUint32(&userbuf, u.Userid, &offset)
	WriteString(&userbuf, u.NexonUsername, &offset)
	WriteString(&userbuf, u.Username, &offset)
	WriteUint8(&userbuf, 1, &offset)
	WriteUint16(&userbuf, uint16(Conf.HolePunchPort), &offset)
	return userbuf
}
