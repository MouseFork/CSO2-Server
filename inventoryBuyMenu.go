package main

const (
	OptionSetBuyMenu = 1
)

//userBuyMenu 用户的购买菜单
type userBuyMenu struct {
	pistols     []uint32
	shotguns    []uint32
	smgs        []uint32
	rifles      []uint32
	snipers     []uint32
	machineguns []uint32
	melees      []uint32
	equipment   []uint32
}

func createNewUserBuyMenu() userBuyMenu {
	return userBuyMenu{
		[]uint32{5280, 5279, 5337, 5356, 5294, 5360, 5262, 103, 106},
		[]uint32{5130, 5293, 5306, 5261, 5242, 5264, 5265, 5230, 137},
		[]uint32{5251, 5295, 5238, 5320, 5285, 5347, 5310, 162, 105},
		[]uint32{46, 45, 5296, 5184, 5355, 113, 102, 161, 157},
		[]uint32{5133, 5118, 5206, 5241, 5225, 146, 125, 160, 163},
		[]uint32{5125, 5314, 5260, 87, 5332, 5366, 5276, 5233, 159},
		[]uint32{79, 5232, 84, 5221, 5304, 5330, 5253, 5231, 5353},
		[]uint32{36, 37, 23, 4, 8, 34, 0, 0, 0},
	}
}

func BuildBuyMenu(inventory userInventory) []byte {
	l := 6 * (len(inventory.buyMenu.pistols) +
		len(inventory.buyMenu.shotguns) +
		len(inventory.buyMenu.smgs) +
		len(inventory.buyMenu.rifles) +
		len(inventory.buyMenu.snipers) +
		len(inventory.buyMenu.machineguns) +
		len(inventory.buyMenu.melees) +
		len(inventory.buyMenu.equipment))
	buf := make([]byte, 4+l)
	offset := 0
	WriteUint8(&buf, OptionSetBuyMenu, &offset)
	WriteUint16(&buf, 369, &offset)
	WriteUint8(&buf, 2, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.pistols)), &offset)
	for k, v := range inventory.buyMenu.pistols {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}
	WriteUint8(&buf, uint8(len(inventory.buyMenu.shotguns)), &offset)
	for k, v := range inventory.buyMenu.shotguns {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.smgs)), &offset)
	for k, v := range inventory.buyMenu.smgs {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.rifles)), &offset)
	for k, v := range inventory.buyMenu.rifles {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.snipers)), &offset)
	for k, v := range inventory.buyMenu.snipers {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.machineguns)), &offset)
	for k, v := range inventory.buyMenu.machineguns {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.melees)), &offset)
	for k, v := range inventory.buyMenu.melees {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	WriteUint8(&buf, uint8(len(inventory.buyMenu.equipment)), &offset)
	for k, v := range inventory.buyMenu.equipment {
		WriteUint8(&buf, uint8(k), &offset)
		WriteUint32(&buf, v, &offset)
	}

	return buf[:offset]
}
