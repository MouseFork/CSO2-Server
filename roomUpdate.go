package main

import (
	"log"
	"net"
	"strconv"
	"unsafe"
)

type upSettingReq struct {
	flags              uint64
	lenOfRoomName      uint8
	roomName           []byte
	unk00              uint8
	unk01              uint8
	unk02              uint32
	unk03              uint32
	lenOfunk09         uint8
	unk09              []byte
	unk10              uint16
	forceCamera        uint8
	gameModeID         uint8
	mapID              uint8
	unk13              uint8
	maxPlayers         uint8
	winLimit           uint8
	killLimit          uint16
	unk17              uint8
	unk18              uint8
	weaponRestrictions uint8
	status             uint8
	unk21              uint8
	mapCycleType       uint8
	unk23              uint8
	unk24              uint8
	unk25              uint8
	numOfMultiMaps     uint8
	multiMaps          []uint8
	teamBalanceType    uint8
	unk29              uint8
	unk30              uint8
	unk31              uint8
	unk32              uint8
	unk33              uint8
	botEnabled         uint8
	botDifficulty      uint8
	numCtBots          uint8
	numTrBots          uint8
	unk35              uint8
	unk36              uint8
	unk37              uint8
	unk38              uint8
	unk39              uint8
	isIngame           uint8
	startMoney         uint16
	changeTeams        uint8
	unk43              uint8
	hltvEnabled        uint8
	unk45              uint8
	respawnTime        uint8
}

//onUpdateRoom 房主更新房间信息
func onUpdateRoom(seq *uint8, p packet, client net.Conn) {
	//检索数据报
	var pkt upSettingReq
	if !praseUpdateRoomPacket(p, &pkt) {
		log.Println("Error : Client from", client.RemoteAddr().String(), "sent a illegal packet !")
		return
	}
	//找到对应用户
	uPtr := getUserFromConnection(client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : Client from", client.RemoteAddr().String(), "try to toggle ready status but not in server !")
		return
	}
	//检查用户是不是房主
	curroom := getRoomFromID(uPtr.getUserChannelServerID(),
		uPtr.getUserChannelID(),
		uPtr.getUserRoomID())
	if curroom == nil ||
		curroom.id <= 0 {
		log.Println("Error : User:", string(uPtr.username), "try to update a null room but in server !")
		return
	}
	if curroom.hostUserID != uPtr.userid {
		log.Println("Error : User:", string(uPtr.username), "try to update a room but isn't host !")
		return
	}
	//检查用户所在房间
	if curroom.id != uPtr.currentRoomId {
		log.Println("Error : User:", string(uPtr.username), "try to update a room but not in !")
		return
	}
	//检查当前是不是正在倒计时
	if (*curroom).isGlobalCountdownInProgress() {
		log.Println("Error : User:", string(uPtr.username), "try to update a room but is counting !")
		return
	}
	//更新房间设置
	curroom.toUpdateSetting(pkt)
	p.id = TypeRoom
	//向房间所有玩家发送更新报文
	for k, v := range (*curroom).users {
		rst := BytesCombine(BuildHeader(v.currentSequence, p), buildRoomSetting(*curroom))
		sendPacket(rst, v.currentConnection)
		log.Println("["+strconv.Itoa(k+1)+"/"+strconv.Itoa(int((*curroom).numPlayers))+"] Updated room for", v.currentConnection.RemoteAddr().String(), "!")
	}
	log.Println("Host from", client.RemoteAddr().String(), "updated room !")
}

func praseUpdateRoomPacket(src packet, dest *upSettingReq) bool {
	//前5字节数据包通用头部+1字节房间数据包通用头部
	offset := 6
	//读取flag，标记要读的有哪些数据
	flags := ReadUint64(src.data, &offset)
	lowFlag := *(*uint32)(unsafe.Pointer(&flags))
	//右移32比特位
	flags = flags >> 32
	highFlag := *(*uint32)(unsafe.Pointer(&flags))
	if lowFlag&0x1 != 0 {
		(*dest).lenOfRoomName = ReadUint8(src.data, &offset)
		(*dest).roomName = ReadString(src.data, &offset, int((*dest).lenOfRoomName))
	}
	if lowFlag&0x2 != 0 {
		(*dest).unk00 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x4 != 0 {
		(*dest).unk01 = ReadUint8(src.data, &offset)
		(*dest).unk02 = ReadUint32(src.data, &offset)
		(*dest).unk03 = ReadUint32(src.data, &offset)
	}
	if lowFlag&0x8 != 0 {
		(*dest).lenOfunk09 = ReadUint8(src.data, &offset)
		(*dest).unk09 = ReadString(src.data, &offset, int((*dest).lenOfRoomName))
	}
	if lowFlag&0x10 != 0 {
		(*dest).unk10 = ReadUint16(src.data, &offset)
	}
	if lowFlag&0x20 != 0 {
		(*dest).forceCamera = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x40 != 0 {
		(*dest).gameModeID = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x80 != 0 {
		(*dest).mapID = ReadUint8(src.data, &offset)
		(*dest).unk13 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x100 != 0 {
		(*dest).maxPlayers = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x200 != 0 {
		(*dest).winLimit = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x400 != 0 {
		(*dest).killLimit = ReadUint16(src.data, &offset)
	}
	if lowFlag&0x800 != 0 {
		(*dest).unk17 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x1000 != 0 {
		(*dest).unk18 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x2000 != 0 {
		(*dest).weaponRestrictions = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x4000 != 0 {
		(*dest).status = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x8000 != 0 {
		(*dest).unk21 = ReadUint8(src.data, &offset)
		(*dest).mapCycleType = ReadUint8(src.data, &offset)
		(*dest).unk23 = ReadUint8(src.data, &offset)
		(*dest).unk24 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x10000 != 0 {
		(*dest).unk25 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x20000 != 0 {
		(*dest).numOfMultiMaps = ReadUint8(src.data, &offset)
		for i := 0; i < int((*dest).numOfMultiMaps); i++ {
			(*dest).multiMaps[i] = ReadUint8(src.data, &offset)
		}
	}
	if lowFlag&0x40000 != 0 {
		(*dest).teamBalanceType = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x80000 != 0 {
		(*dest).unk29 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x100000 != 0 {
		(*dest).unk30 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x200000 != 0 {
		(*dest).unk31 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x400000 != 0 {
		(*dest).unk32 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x800000 != 0 {
		(*dest).unk33 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x1000000 != 0 {
		(*dest).botEnabled = ReadUint8(src.data, &offset)
		if (*dest).botEnabled != 0 {
			(*dest).botDifficulty = ReadUint8(src.data, &offset)
			(*dest).numCtBots = ReadUint8(src.data, &offset)
			(*dest).numTrBots = ReadUint8(src.data, &offset)
		}
	}
	if lowFlag&0x2000000 != 0 {
		(*dest).unk35 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x4000000 != 0 {
		(*dest).unk36 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x8000000 != 0 {
		(*dest).unk37 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x10000000 != 0 {
		(*dest).unk38 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x20000000 != 0 {
		(*dest).unk39 = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x40000000 != 0 {
		(*dest).isIngame = ReadUint8(src.data, &offset)
	}
	if lowFlag&0x80000000 != 0 {
		(*dest).startMoney = ReadUint16(src.data, &offset)
	}
	if highFlag&0x1 != 0 {
		(*dest).changeTeams = ReadUint8(src.data, &offset)
	}
	if highFlag&0x2 != 0 {
		(*dest).unk43 = ReadUint8(src.data, &offset)
	}
	if highFlag&0x4 != 0 {
		(*dest).hltvEnabled = ReadUint8(src.data, &offset)
	}
	if highFlag&0x8 != 0 {
		(*dest).unk45 = ReadUint8(src.data, &offset)
	}
	if highFlag&0x10 != 0 {
		(*dest).respawnTime = ReadUint8(src.data, &offset)
	}
	return true
}
