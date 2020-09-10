package typestruct

//UserBuyMenu 用户的购买菜单
type UserBuyMenu struct {
	Pistols     []uint32
	Shotguns    []uint32
	Smgs        []uint32
	Rifles      []uint32
	Snipers     []uint32
	Machineguns []uint32
	Melees      []uint32
	Equipment   []uint32
}

func CreateNewUserBuyMenu() UserBuyMenu {
	return UserBuyMenu{
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
