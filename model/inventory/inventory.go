package inventory

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
