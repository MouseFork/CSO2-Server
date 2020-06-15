package main

const (
	FavoriteSetLoadout   = 0
	FavoriteSetCosmetics = 1
)

func BuildCosmetics(inventory userInventory) []byte {
	buf := make([]byte, 55)
	offset := 0
	WriteUint8(&buf, FavoriteSetCosmetics, &offset)
	WriteUint8(&buf, 10, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint32(&buf, inventory.CTModel, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint32(&buf, inventory.TModel, &offset)
	WriteUint8(&buf, 2, &offset)
	WriteUint32(&buf, inventory.headItem, &offset)
	WriteUint8(&buf, 3, &offset)
	WriteUint32(&buf, inventory.gloveItem, &offset)
	WriteUint8(&buf, 4, &offset)
	WriteUint32(&buf, inventory.backItem, &offset)
	WriteUint8(&buf, 5, &offset)
	WriteUint32(&buf, inventory.stepsItem, &offset)
	WriteUint8(&buf, 6, &offset)
	WriteUint32(&buf, inventory.cardItem, &offset)
	WriteUint8(&buf, 7, &offset)
	WriteUint32(&buf, inventory.sprayItem, &offset)
	WriteUint8(&buf, 8, &offset)
	WriteUint32(&buf, 0, &offset)
	WriteUint8(&buf, 9, &offset)
	WriteUint32(&buf, 0, &offset)
	return buf[:offset]
}
