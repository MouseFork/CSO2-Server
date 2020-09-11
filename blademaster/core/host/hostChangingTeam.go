package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnChangingTeam(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InHostTeamChangingPacket
	if !p.PraseInTeamChangingPacket(&pkt) {
		DebugInfo(2, "Error : User from", client.RemoteAddr().String(), "sent a error TeamChanging packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : User from", client.RemoteAddr().String(), "sent TeamChanging but not in server !")
		return
	}
	destUser := GetUserFromID(uint32(pkt.UserId))
	if destUser == nil ||
		destUser.Userid <= 0 {
		DebugInfo(2, "Error : User from", client.RemoteAddr().String(), "sent TeamChanging but dester is not in server !")
		return
	}
	//找到房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.CurrentRoomId)
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User from", client.RemoteAddr().String(), "sent TeamChanging but is not host !")
		return
	}
	//更新数据
	//log.Println(p.data)
	destUser.SetUserTeam(pkt.NewTeam)
	result := BuildChangingTeam(destUser.Userid, destUser.CurrentTeam)
	for _, v := range rm.Users {
		rst := BytesCombine(BuildHeader(v.CurrentSequence, p.Id), result)
		SendPacket(rst, v.CurrentConnection)
		// if rm.setting.areBotsEnabled != 0 {
		// 	//(*u).currentTeam = pkt.newTeam
		// 	//log.Println("User", string(u.username), "changed team", u.currentTeam)
		// 	//(*rm).users[k].currentTeam = pkt.newTeam
		// 	//向每名玩家表示该玩家更换了队伍
		// 	rst1 := BuildChangingTeam(u.userid, pkt.newTeam)
		// 	for _, i := range rm.users {
		// 		rst = BytesCombine(BuildHeader(i.currentSequence, p), rst1)
		// 		sendPacket(rst, i.currentConnection)
		// 	}
		// }
	}
}
func BuildChangingTeam(id uint32, team uint8) []byte {
	buf := make([]byte, 7)
	offset := 0
	WriteUint8(&buf, OUTsetUserTeam, &offset)
	WriteUint8(&buf, 1, &offset) //numberOfUser
	WriteUint32(&buf, id, &offset)
	WriteUint8(&buf, team, &offset)
	return buf[:offset]
}
