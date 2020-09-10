package holepunch

import (
	"fmt"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/udp"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func StartHolePunchServer(port string, server *net.UDPConn) {
	defer server.Close()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UDP server suffered a fault !")
			fmt.Println("error:", err)
			fmt.Println("Fault end!")
		}
	}()
	fmt.Println("Server UDPholepunch is running at", "[AnyAdapter]:"+port)
	for {
		data := make([]byte, 1024)
		n, ClientAddr, err := server.ReadFromUDP(data)
		if err != nil {
			DebugInfo(2, "UDP read error from", ClientAddr.String())
			continue
		}
		go RecvHolePunchMessage(data[:n], n, ClientAddr, server)
	}
}

//RecvHolePunchMessage 处理收到的包
func RecvHolePunchMessage(data []byte, len int, client *net.UDPAddr, server *net.UDPConn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recv UDP packet suffered a fault !")
			fmt.Println(err)
			fmt.Println("Fault end!")
			return
		}
	}()
	//分析数据包
	var p InUDPmsg
	if !p.PraseUDPpacket(data, len) {
		DebugInfo(2, "UDP had a illegal packet from", client.String())
		return
	}
	if p.IsHeartbeat() {
		return
	}
	cliadr := client.IP.To4().String()
	externalIPAddress, err := IPToUint32(cliadr)
	if err != nil {
		DebugInfo(2, "Error : Prasing externalIpAddress error !")
		return
	}
	//找到对应玩家
	uPtr := GetUserFromID(p.UserId)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		//log.Println("UDP had a packet from", client.String(), "but not logged in !")
		return
	}
	//更新netinfo
	index := uPtr.UpdateHolepunch(p.PortId, p.Port, uint16(client.Port))
	if index == 0xFFFF {
		DebugInfo(2, "Error : User", uPtr.Username, "update Holepunch failed !")
		return
	}
	uPtr.NetInfo.ExternalIpAddress = externalIPAddress
	uPtr.NetInfo.LocalIpAddress = p.IpAddress
	//发送返回数据
	rst := BuildUDPHolepunch(index)
	server.WriteToUDP(rst, client)
}
