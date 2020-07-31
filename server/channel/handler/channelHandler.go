package handler

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/model/channel"
	. "github.com/KouKouChan/CSO2-Server/model/packet"
	. "github.com/KouKouChan/CSO2-Server/server/packet"
)

//处理请求服务器列表
func onServerList(seq *uint8, p *Packet, client *(net.Conn)) {
	uPtr := getUserFromConnection(*client)
	if uPtr == nil ||
		uPtr.userid <= 0 {
		log.Println("Error : A unknow Client from", (*client).RemoteAddr().String(), "request a ChannelList !")
		return
	}
	(*p).Id = TypeServerList
	rst := BuildHeader(seq, *p)
	rst = append(rst, MainServer.ServerNum)
	for i := 0; i < int(MainServer.ServerNum); i++ {
		rst = BytesCombine(rst, BuildChannelServerList(MainServer.Servers[i]))
	}
	WriteLen(&rst)       //写入长度
	(*client).Write(rst) //发送UserStart消息
	log.Println("Sent a server list packet to", (*client).RemoteAddr().String())
	uPtr.setUserChannelServer(0)
	uPtr.quitChannel()
}

//建立某个频道服务器数据包
func BuildChannelServerList(chlsrv ChannelServer) []byte {
	var list []byte
	l := len(chlsrv.ServerName)
	list = append(list, chlsrv.ServerIndex,
		chlsrv.ServerStatus,
		chlsrv.ServerType,
		uint8(l),
	)
	for i := 0; i < l; i++ {
		list = append(list, chlsrv.ServerName[i])
	}
	list = append(list, chlsrv.ChannelCount)
	chl := BuildChannelList(chlsrv.ChannelCount, chlsrv.Channels)
	for i := 0; i < len(chl); i++ {
		list = append(list, chl[i])
	}
	return list
}

func BuildChannelList(num uint8, channels []ChannelInfo) []byte {
	var list []byte
	for i := 0; i < int(num); i++ {
		temp := make([]byte, 9+len(channels[i].ChannelName))
		offset := 0
		WriteUint8(&temp, channels[i].ChannelID, &offset)
		WriteString(&temp, channels[i].ChannelName, &offset)
		WriteUint16(&temp, channels[i].Unk00, &offset)
		WriteUint16(&temp, channels[i].Unk01, &offset)
		WriteUint8(&temp, channels[i].Unk02, &offset)
		WriteUint8(&temp, channels[i].Unk03, &offset)
		WriteUint8(&temp, channels[i].Unk04, &offset)
		for j := 0; j < len(temp); j++ {
			list = append(list, temp[j])
		}
	}
	return list
}
