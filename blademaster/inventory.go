package blademaster

import (
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type (
	UserInventory struct {
		NumOfItem uint16              //物品数量
		Items     []UserInventoryItem //物品

		CTModel   uint32 //当前的CT模型
		TModel    uint32 //当前的T模型
		HeadItem  uint32 //当前的头部装饰
		GloveItem uint32 //当前的手套
		BackItem  uint32 //当前的背部物品
		StepsItem uint32 //当前的脚步效果
		CardItem  uint32 //当前的卡片
		SprayItem uint32 //当前的喷漆

		BuyMenu  UserBuyMenu //购买菜单
		Loadouts []UserLoadout
	}
	UserInventoryItem struct {
		Id    uint32 //物品id
		Count uint16 //数量
	}
)

func CreateNewUserInventory() UserInventory {
	Inv := UserInventory{
		0,
		//createDeafaultInventoryItem(),
		CreateFullInventoryItem(),
		1047,
		1048,
		0,
		0,
		0,
		0,
		0,
		42001,
		CreateNewUserBuyMenu(),
		CreateNewLoadout(),
	}
	Inv.NumOfItem = uint16(len(Inv.Items))
	return Inv
}

func CreateDeafaultInventoryItem() []UserInventoryItem {
	items := []UserInventoryItem{}

	return items
}

func CreateFullInventoryItem() []UserInventoryItem {
	items := []UserInventoryItem{}
	var i uint32
	//用户角色
	for i = 1005; i <= 1058; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	//添加默认武器
	number := []uint32{2, 3, 4, 6, 8, 13, 14, 15, 18, 19, 21, 23, 27, 34, 36, 37, 80, 128, 101, 1001, 1002, 1003, 1004, 49009, 49004}
	for _, v := range number {
		items = append(items, UserInventoryItem{v, 1})
	}
	//解锁武器
	for i = 1; i <= 33; i++ {
		if IsIllegal(i) {
			continue
		}
		items = append(items, UserInventoryItem{i, 1})
	}
	for i = 44; i <= 163; i++ {
		if IsIllegal(i) {
			continue
		}
		items = append(items, UserInventoryItem{i, 1})
	}
	//僵尸技能
	items = append(items, UserInventoryItem{2019, 1})
	items = append(items, UserInventoryItem{3, 1})
	items = append(items, UserInventoryItem{2020, 1})
	items = append(items, UserInventoryItem{50, 1})
	for i = 2021; i <= 2023; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	//武器皮肤
	for i = 5042; i <= 5370; i++ {
		if IsIllegal(i) {
			continue
		}
		items = append(items, UserInventoryItem{i, 1})
	}
	items = append(items, UserInventoryItem{5997, 1})
	//帽子
	for i = 10001; i <= 10133; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	//背包
	for i = 20001; i <= 20107; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	//手套
	for i = 30001; i <= 30027; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	//脚部特效
	for i = 40001; i <= 40025; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	//喷漆
	for i = 42001; i <= 42020; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	//道具
	for i = 49001; i <= 49010; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	items = append(items, UserInventoryItem{49999, 1})
	//角色卡片
	for i = 60001; i <= 60004; i++ {
		items = append(items, UserInventoryItem{i, 1})
	}
	return items
}

func IsIllegal(num uint32) bool {
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
