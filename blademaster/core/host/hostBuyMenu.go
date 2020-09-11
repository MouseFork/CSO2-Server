package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostSetUserBuyMenu(p *PacketData, client net.Conn) {
	//检查数据包
	var pkt InHostSetBuyMenu
	if !p.PraseSetBuyMenuPacket(&pkt) {
		DebugInfo(2, "Error : Cannot prase a send BuyMenu packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : A user request to send BuyMenu but not in server!")
		return
	}
	dest := GetUserFromID(pkt.Userid)
	if dest == nil ||
		dest.Userid <= 0 {
		DebugInfo(2, "Error : A user request to send BuyMenu but dest user is null!")
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to send BuyMenu but in a null room !")
		return
	}
	destRm := GetRoomFromID(dest.GetUserChannelServerID(),
		dest.GetUserChannelID(),
		dest.GetUserRoomID())
	if destRm == nil ||
		destRm.Id <= 0 {
		DebugInfo(2, "Error : User", string(dest.Username), "try to send BuyMenu but in a null room !")
		return
	}
	if rm.Id != destRm.Id {
		DebugInfo(2, "Error : User", string(dest.Username), "try to send BuyMenu to", string(dest.Username), "but not a same room !")
		return
	}
	//是不是房主
	if rm.HostUserID != uPtr.Userid {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to send BuyMenu but isn't host !")
		return
	}
	//发送数据包
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, p.Id), BuildSetBuyMenu(dest.Userid, &dest.Inventory))
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Send User", string(dest.Username), "BuyMenu to host", string(uPtr.Username))

}

func BuildSetBuyMenu(id uint32, inventory *UserInventory) []byte {
	l := 6 * (len(inventory.BuyMenu.Pistols) +
		len(inventory.BuyMenu.Shotguns) +
		len(inventory.BuyMenu.Smgs) +
		len(inventory.BuyMenu.Rifles) +
		len(inventory.BuyMenu.Snipers) +
		len(inventory.BuyMenu.Machineguns) +
		len(inventory.BuyMenu.Melees) +
		len(inventory.BuyMenu.Equipment))
	buf := make([]byte, 8+l)
	offset := 0
	WriteUint8(&buf, SetBuyMenu, &offset)
	WriteUint32(&buf, id, &offset)
	WriteUint16(&buf, 369, &offset) //buyMenuByteLength
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Pistols)), &offset)
	for k, v := range inventory.BuyMenu.Pistols {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Shotguns)), &offset)
	for k, v := range inventory.BuyMenu.Shotguns {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Smgs)), &offset)
	for k, v := range inventory.BuyMenu.Smgs {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Rifles)), &offset)
	for k, v := range inventory.BuyMenu.Rifles {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Snipers)), &offset)
	for k, v := range inventory.BuyMenu.Snipers {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Machineguns)), &offset)
	for k, v := range inventory.BuyMenu.Machineguns {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Melees)), &offset)
	for k, v := range inventory.BuyMenu.Melees {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Equipment)), &offset)
	for k, v := range inventory.BuyMenu.Equipment {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}
	return buf[:offset]
}
