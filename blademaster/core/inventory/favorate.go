package inventory

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	FavoriteSetLoadout   = 0
	FavoriteSetCosmetics = 1

	OptionSetBuyMenu = 1
)

var (
	DeafaultInventoryItem = BuildDefaultInventoryInfo()
)

func OnFavorite(p *PacketData, client net.Conn) {
	var pkt InFavoritePacket
	if !p.PraseFavoritePacket(&pkt) {
		DebugInfo(2, "Error : Recived a illegal favorite packet from", client.RemoteAddr().String())
		return
	}
	switch pkt.PacketType {
	case FavoriteSetLoadout:
		//log.Println("Recived a favorite SetLoadout packet from", client.RemoteAddr().String())
		OnFavoriteSetLoadout(p, client)
	case FavoriteSetCosmetics:
		//log.Println("Recived a favorite SetCosmetics packet from", client.RemoteAddr().String())
		OnFavoriteSetCosmetics(p, client)
	default:
		DebugInfo(2, "Unknown favorite packet", pkt.PacketType, "from", client.RemoteAddr().String())
	}
}

func BuildCosmetics(inventory *UserInventory) []byte {
	buf := make([]byte, 2)
	offset := 0
	curItem := uint8(0)
	WriteUint8(&buf, FavoriteSetCosmetics, &offset)
	WriteUint8(&buf, 10, &offset)
	temp := WriteItem(inventory.CTModel, &curItem)
	temp = BytesCombine(temp, WriteItem(inventory.TModel, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.HeadItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.GloveItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.BackItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.StepsItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.CardItem, &curItem))
	temp = BytesCombine(temp, WriteItem(inventory.SprayItem, &curItem))
	temp = BytesCombine(temp, WriteItem(0, &curItem))
	temp = BytesCombine(temp, WriteItem(0, &curItem))
	buf = BytesCombine(buf[:offset], temp)
	return buf
}

func BuildBuyMenu(inventory *UserInventory) []byte {
	l := 6 * (len(inventory.BuyMenu.Pistols) +
		len(inventory.BuyMenu.Shotguns) +
		len(inventory.BuyMenu.Smgs) +
		len(inventory.BuyMenu.Rifles) +
		len(inventory.BuyMenu.Snipers) +
		len(inventory.BuyMenu.Machineguns) +
		len(inventory.BuyMenu.Melees) +
		len(inventory.BuyMenu.Equipment))
	buf := make([]byte, 4+l)
	offset := 0
	WriteUint8(&buf, OptionSetBuyMenu, &offset)
	WriteUint16(&buf, 369, &offset)
	WriteUint8(&buf, 2, &offset)
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
