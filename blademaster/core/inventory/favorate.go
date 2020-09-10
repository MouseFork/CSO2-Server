package inventory

import (
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

const (
	FavoriteSetLoadout   = 0
	FavoriteSetCosmetics = 1

	OptionSetBuyMenu = 1
)

var (
	DeafaultInventoryItem = BuildDefaultInventoryInfo()
)

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
