package main

import . "github.com/KouKouChan/CSO2-Server/kerlong"

type userLoadout struct {
	items []uint32
}

func BuildLoadout(inventory userInventory) []byte {
	buf := make([]byte, 5+len(inventory.loadouts)*96)
	offset := 0
	WriteUint8(&buf, FavoriteSetLoadout, &offset)
	for i, v := range inventory.loadouts {
		for j, k := range v.items {
			WriteUint8(&buf, uint8(i), &offset)
			WriteUint8(&buf, uint8(j), &offset)
			WriteUint32(&buf, k, &offset)
		}
	}
	return buf[:offset]
}

func createNewLoadout() []userLoadout {
	return []userLoadout{
		{[]uint32{5336, 5356, 5330, 4, 23, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[]uint32{5285, 5294, 5231, 4, 23, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[]uint32{5206, 5356, 5365, 4, 23, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}
}
