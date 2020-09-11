package inventory

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnFavoriteSetLoadout(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InFavoriteSetLoadout
	if !p.PraseFavoriteSetLoadoutPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error SetLoadout packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to SetLoadout but not in server !")
		return
	}
	//设置武器
	if pkt.Loadout > 2 ||
		pkt.WeaponSlot > 6 {
		DebugInfo(2, "Error : User", string(uPtr.Username), "try to SetLoadout with invalid data !")
		return
	}
	uPtr.Inventory.Loadouts[pkt.Loadout].Items[pkt.WeaponSlot] = pkt.ItemId
	DebugInfo(1, "Setting User", string(uPtr.Username), "new weapon", pkt.ItemId, "to slot", pkt.WeaponSlot, "in loadout", pkt.Loadout)
}
func BuildLoadout(inventory *UserInventory) []byte {
	buf := make([]byte, 5+len(inventory.Loadouts)*96)
	offset := 0
	WriteUint8(&buf, FavoriteSetLoadout, &offset)
	WriteUint8(&buf, uint8(len(inventory.Loadouts))*16, &offset)
	for i, v := range inventory.Loadouts {
		for j, k := range v.Items {
			WriteUint8(&buf, uint8(i), &offset)
			WriteUint8(&buf, uint8(j), &offset)
			WriteUint32(&buf, k, &offset)
		}
	}
	return buf[:offset]
}
