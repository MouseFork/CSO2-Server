package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostSetUserLoadout(p *PacketData, client net.Conn) {
	//检查数据包
	var pkt InHostSetLoadoutPacket
	if !p.PraseSetUserLoadoutPacket(&pkt) {
		DebugInfo(2, "Error : Cannot prase a send UserLoadout packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : A user request to send UserLoadout but not in server!")
		return
	}
	dest := GetUserFromID(pkt.UserID)
	if dest == nil ||
		dest.Userid <= 0 {
		DebugInfo(2, "Error : A user request to send UserLoadout but dest user is null!")
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to send UserLoadout but in a null room !")
		return
	}
	//是不是房主
	if rm.HostUserID != uPtr.Userid {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to send UserLoadout but isn't host !")
		return
	}
	//发送用户背包数据
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, p.Id), BuildSetUserLoadout(dest))
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Send User", string(dest.Username), "Loadout to host", string(uPtr.Username))
}

func BuildSetUserLoadout(u *User) []byte {
	buf := make([]byte, 6)
	offset := 0
	WriteUint8(&buf, SetLoadout, &offset)
	WriteUint32(&buf, u.Userid, &offset)
	WriteUint8(&buf, 8, &offset) //类型数量
	//当前8个类型的装备
	curItem := uint8(0)
	temp := WriteItem(u.Inventory.CTModel, &curItem)
	temp = BytesCombine(temp, WriteItem(u.Inventory.TModel, &curItem))
	temp = BytesCombine(temp, WriteItem(u.Inventory.HeadItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.Inventory.GloveItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.Inventory.BackItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.Inventory.StepsItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.Inventory.CardItem, &curItem))
	temp = BytesCombine(temp, WriteItem(u.Inventory.SprayItem, &curItem))
	buf = BytesCombine(buf[:offset], temp)
	buf = append(buf, uint8(len(u.Inventory.Loadouts)))
	for _, v := range u.Inventory.Loadouts {
		buf = append(buf, uint8(len(v.Items)))
		curItem = 0
		for _, j := range v.Items {
			buf = BytesCombine(buf, WriteItem(j, &curItem))
		}
	}
	buf = append(buf, 0)
	return buf
}
