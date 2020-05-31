package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
)

var (
	//SERVERVERSION 版本号
	SERVERVERSION = "v0.1.3"
	//PORT 端口
	PORT = 30001
	//HOLEPUNCHPORT 端口
	HOLEPUNCHPORT = 30002
	//MainServer 主服务器
	MainServer = serverManager{
		0,
		[]channelServer{},
	}
	//UserManager 全局用户管理
	UserManager = userManager{
		0,
		[]user{},
	}
)

func main() {
	defer func() {
		fmt.Println("检测到异常")
		// 获取异常信息
		if err := recover(); err != nil {
			//  输出异常信息
			fmt.Println("error:", err)
		}
		fmt.Println("异常结束")
	}()
	fmt.Println("Counter-Strike Online 2 Server", SERVERVERSION)
	//初始化TCP
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatal("Init tcp socket error !\n")
		os.Exit(-1)
	}
	//初始化UDP
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatal("Init udp addr error !\n")
		os.Exit(-1)
	}
	holepunchserver, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Init udp socket error !\n")
		os.Exit(-1)
	}
	//延迟关闭
	defer server.Close()
	defer holepunchserver.Close()

	fmt.Println("Initializing process ...")
	//初始化主频道服务器
	MainServer = newMainServer()
	//开启UDP服务
	go startHolePunchServer(holepunchserver)
	//开启TCP服务
	fmt.Println("Server is running at", "[localhost]:"+strconv.Itoa(PORT))
	for {
		client, err := server.Accept()
		if err != nil {
			log.Fatal("Server Accept data error !\n")
			continue
		}
		log.Println("Server accept a new connection request at", client.RemoteAddr().String())
		go RecvMessage(client)
	}
}

func startHolePunchServer(server *(net.UDPConn)) {
	defer server.Close()
	fmt.Println("Server holepunch is running at", "[localhost]:"+strconv.Itoa(HOLEPUNCHPORT))
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := server.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("Server holepunch read error from", server.RemoteAddr().String())
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		// client, err := (*server).Accept()
		// if err != nil {
		// 	log.Fatal("Server Accept data error !\n")
		// 	continue
		// }
		// log.Println("Server accept a new connection request at", client.RemoteAddr().String())
		// go RecvHolePunchMessage(client)
	}
}

//RecvHolePunchMessage 处理收到的包
func RecvHolePunchMessage(holeclient net.Conn) {

}

//RecvMessage 处理收到的包
func RecvMessage(client net.Conn) {
	defer client.Close() //关闭con
	var seq uint8 = 0
	client.Write([]byte("~SERVERCONNECTED\n"))
	for {
		bytes := make([]byte, math.MaxUint16)
		n, err := client.Read(bytes) //读数据
		if err == nil {
			if n == 0 {
				continue
			}
			var pkt packet
			pkt.data = bytes
			log.Println("Prasing a packet from", client.RemoteAddr().String())
			pkt.PrasePacket()
			if !pkt.IsGoodPacket {
				log.Println("Recived a illegal packet from", client.RemoteAddr().String())
				continue
			}
			switch pkt.id {
			case TypeVersion:
				log.Println("Recived a client version packet from", client.RemoteAddr().String())
				client.Write(onVersionPacket(&seq, pkt))
			case TypeLogin:
				log.Println("Recived a login request packet from", client.RemoteAddr().String())
				if !onLoginPacket(&seq, &pkt, &client) {
					log.Println("Recived a illegal packet from", client.RemoteAddr().String())
				}
			case TypeRequestChannels:
				log.Println("Recived a ChannelList request packet from", client.RemoteAddr().String())
				onServerList(&seq, &pkt, &client)
			case TypeRequestRoomList:
				log.Println("Recived a RoomList request packet from", client.RemoteAddr().String())
				onRoomList(&seq, &pkt, client)
			case TypeRoom:
				//log.Println("Recived a Room request packet from", client.RemoteAddr().String())
				onRoomRequest(&seq, pkt, client)
			case TypeHost:
				log.Println("Recived a Host request packet from", client.RemoteAddr().String())
			case TypeFavorite:
				log.Println("Recived a favorite request packet from", client.RemoteAddr().String())
			default:
				log.Println("Recived a unknown packet from", client.RemoteAddr().String())
			}
		} else {
			log.Println("client", client.RemoteAddr().String(), "closed the connection")
			delUserWithConn(client)
			client.Close() //关闭con
			return
		}
	}
}

//GetIP 获取IP
func GetIP() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return "Error"
}
