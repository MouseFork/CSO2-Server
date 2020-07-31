package inventory

func createNewUserInventory() userInventory {
	Inv := userInventory{
		0,
		//createDeafaultInventoryItem(),
		createFullInventoryItem(),
		1047,
		1048,
		0,
		0,
		0,
		0,
		0,
		42001,
		createNewUserBuyMenu(),
		createNewLoadout(),
	}
	Inv.numOfItem = uint16(len(Inv.items))
	return Inv
}

func createDeafaultInventoryItem() []userInventoryItem {
	items := []userInventoryItem{}

	return items
}

func createFullInventoryItem() []userInventoryItem {
	items := []userInventoryItem{}
	var i uint32
	//用户角色
	for i = 1005; i <= 1058; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	//添加默认武器
	number := []uint32{2, 3, 4, 6, 8, 13, 14, 15, 18, 19, 21, 23, 27, 34, 36, 37, 80, 128, 101, 1001, 1002, 1003, 1004, 49009, 49004}
	for _, v := range number {
		items = append(items, userInventoryItem{v, 1})
	}
	//解锁武器
	for i = 1; i <= 33; i++ {
		if isIllegal(i) {
			continue
		}
		items = append(items, userInventoryItem{i, 1})
	}
	for i = 44; i <= 163; i++ {
		if isIllegal(i) {
			continue
		}
		items = append(items, userInventoryItem{i, 1})
	}
	//僵尸技能
	items = append(items, userInventoryItem{2019, 1})
	items = append(items, userInventoryItem{3, 1})
	items = append(items, userInventoryItem{2020, 1})
	items = append(items, userInventoryItem{50, 1})
	for i = 2021; i <= 2023; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	//武器皮肤
	for i = 5042; i <= 5370; i++ {
		if isIllegal(i) {
			continue
		}
		items = append(items, userInventoryItem{i, 1})
	}
	items = append(items, userInventoryItem{5997, 1})
	//帽子
	for i = 10001; i <= 10133; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	//背包
	for i = 20001; i <= 20107; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	//手套
	for i = 30001; i <= 30027; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	//脚部特效
	for i = 40001; i <= 40025; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	//喷漆
	for i = 42001; i <= 42020; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	//道具
	for i = 49001; i <= 49010; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	items = append(items, userInventoryItem{49999, 1})
	//角色卡片
	for i = 60001; i <= 60004; i++ {
		items = append(items, userInventoryItem{i, 1})
	}
	return items
}

func isIllegal(num uint32) bool {
	switch num {
	case 2:
		return true
	case 3:
		return true
	case 4:
		return true
	case 6:
		return true
	case 8:
		return true
	case 13:
		return true
	case 14:
		return true
	case 15:
		return true
	case 18:
		return true
	case 19:
		return true
	case 21:
		return true
	case 23:
		return true
	case 27:
		return true
	case 56:
		return true
	case 58:
		return true
	case 69:
		return true
	case 107:
		return true
	case 117:
		return true
	case 134:
		return true
	case 139:
		return true
	case 5172:
		return true
	case 5173:
		return true
	case 5174:
		return true
	case 5227:
		return true
	case 5228:
		return true
	case 5229:
		return true
	default:
		return false
	}
}

func WriteItem(num uint32, curitem *uint8) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint8(&buf, *curitem, &offset)
	(*curitem)++
	WriteUint32(&buf, num, &offset)
	return buf
}
