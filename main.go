package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/holepunch"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/quick"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/room"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/user"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/version"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/database/redis"
	. "github.com/KouKouChan/CSO2-Server/database/sqlite"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/register"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
	"github.com/garyburd/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//SERVERVERSION 版本号
	SERVERVERSION = "v0.3.0"
	Redis         redis.Conn
)

func ReadHead(client net.Conn) ([]byte, bool) {
	head, curlen := make([]byte, HeaderLen), 0
	for {
		n, err := client.Read(head[curlen:])
		if err != nil {
			return head, false
		}
		curlen += n
		if curlen >= HeaderLen {
			break
		}
	}
	return head, true
}

func ReadData(client net.Conn, len uint16) ([]byte, bool) {
	data, curlen := make([]byte, len), 0
	for {
		n, err := client.Read(data[curlen:])
		if err != nil {
			return data, false
		}
		curlen += n
		if curlen >= int(len) {
			break
		}
	}
	return data, true
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("检测到异常")
			fmt.Println("error:", err)
			fmt.Println("异常结束")
		}
	}()
	fmt.Println("Counter-Strike Online 2 Server", SERVERVERSION)
	fmt.Println("Initializing process ...")

	//get server exe path
	path, err := GetExePath()
	if err != nil {
		panic(err)
	}

	//read configure
	Conf.InitConf(path)

	//set verbose
	Level = Conf.DebugLevel
	LogFile = Conf.LogFile
	IsConsole = Conf.EnableConsole

	//init Logger
	if LogFile != 0 {
		InitLoger(path, "CSO2-Server.log")
	}

	//init TCP
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", Conf.PORT))
	if err != nil {
		fmt.Println("Init tcp socket error !\n")
		panic(err)
	}
	defer server.Close()

	//init UDP
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Conf.HolePunchPort))
	if err != nil {
		fmt.Println("Init udp addr error !\n")
		panic(err)
	}
	holepunchserver, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Init udp socket error !\n")
		panic(err)
	}
	defer holepunchserver.Close()

	//Init Database
	if Conf.EnableDataBase != 0 {
		DB, err = InitDatabase(path + "\\cso2.db")
		if err != nil {
			fmt.Println("Init database failed !")
			Conf.EnableDataBase = 0
		} else {
			fmt.Println("Database connected !")
			defer DB.Close()
		}
	} else {
		DB = nil
	}

	//Init Redis
	if Conf.EnableRedis != 0 {
		Redis, err := InitRedis(Conf.RedisIP + ":" + strconv.Itoa(int(Conf.RedisPort)))
		if err != nil {
			fmt.Println("connect to redis server failed !")
			Conf.EnableRedis = 0
		} else {
			fmt.Println("Redis server connected !")
			defer Redis.Close()
		}
	}

	//Init MainServer Info
	MainServer = NewMainServer()

	//Start UDP Server
	go StartHolePunchServer(strconv.Itoa(int(Conf.HolePunchPort)), holepunchserver)

	//Start TCP Server
	go TCPServer(server)

	//Start Register Server
	if Conf.EnableRegister != 0 {
		go OnRegister()
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	_ = <-ch
}

func TCPServer(server net.Listener) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("TCP server suffered a fault !")
			fmt.Println("error:", err)
			fmt.Println("Fault end!")
		}
	}()

	fmt.Println("Server is running at", "[AnyAdapter]:"+strconv.Itoa(int(Conf.PORT)))
	for {
		client, err := server.Accept()
		if err != nil {
			DebugInfo(2, "Server Accept data error !\n")
			continue
		}
		DebugInfo(2, "Server accept a new connection request at", client.RemoteAddr().String())
		go RecvMessage(client)
	}
}

//RecvMessage 循环处理收到的包
func RecvMessage(client net.Conn) {
	var seq uint8 = 0

	defer client.Close() //关闭con
	defer func() {
		if err := recover(); err != nil {
			OnSendMessage(&seq, client, DialogBox, GAME_SERVER_ERROR)
			fmt.Println("Client", client.RemoteAddr().String(), "suffered a fault !")
			fmt.Println(err)
			fmt.Println("Fault end!")
			DelUserWithConn(client)
		}
	}()

	client.Write([]byte("~SERVERCONNECTED\n"))

	for {
		//读取4字节数据包头部
		headBytes, err := ReadHead(client)
		if !err {
			goto close
		}
		var headPacket PacketHeader
		headPacket.Data = headBytes
		headPacket.PraseHeadPacket()
		if !headPacket.IsGoodPacket {
			DebugInfo(2, "Recived a illegal head from", client.RemoteAddr().String())
			continue
		}

		//读取数据部分
		bytes, err := ReadData(client, headPacket.Length)
		if !err {
			goto close
		}
		dataPacket := PacketData{
			bytes,
			headPacket.Sequence,
			headPacket.Length,
			bytes[0],
			1,
		}

		//执行功能
		switch dataPacket.Id {
		case PacketTypeQuickJoin:
			OnQuick(&dataPacket, client)
		case PacketTypeVersion:
			OnVersionPacket(&seq, client)
		case PacketTypeLogin:
			OnLogin(&seq, &dataPacket, client)
		case PacketTypeRequestChannels:
			OnServerList(client)
		case PacketTypeRequestRoomList:
			OnRoomList(&dataPacket, client)
		case PacketTypeRoom:
			OnRoomRequest(&dataPacket, client)
		case PacketTypeHost:
		case PacketTypeFavorite:
		case PacketTypeOption:
		case PacketTypePlayerInfo:
		default:
			DebugInfo(2, "Unknown packet", dataPacket.Id, "from", client.RemoteAddr().String())
		}
	}

close:
	DebugInfo(1, "client", client.RemoteAddr().String(), "closed the connection")
	DelUserWithConn(client)
	return
}
