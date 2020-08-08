package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	. "github.com/KouKouChan/CSO2-Server/blademaster"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/database/redis"
	. "github.com/KouKouChan/CSO2-Server/database/sqlite"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/server"
	. "github.com/KouKouChan/CSO2-Server/server/channel"
	. "github.com/KouKouChan/CSO2-Server/server/database"
	. "github.com/KouKouChan/CSO2-Server/server/favorate"
	. "github.com/KouKouChan/CSO2-Server/server/host"
	. "github.com/KouKouChan/CSO2-Server/server/inventory"
	. "github.com/KouKouChan/CSO2-Server/server/message"
	. "github.com/KouKouChan/CSO2-Server/server/playerInfo"
	. "github.com/KouKouChan/CSO2-Server/server/quick"
	. "github.com/KouKouChan/CSO2-Server/server/room"
	. "github.com/KouKouChan/CSO2-Server/server/user"
	. "github.com/KouKouChan/CSO2-Server/server/version"
	. "github.com/KouKouChan/CSO2-Server/verbose"
	"github.com/garyburd/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//SERVERVERSION 版本号
	SERVERVERSION = "v0.3.0"
	Redis         redis.Conn
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
	fmt.Println("Initializing process ...")
	//获取exe路径
	path, err := GetExePath()
	if err != nil {
		panic(err)
	}
	//读取配置
	Conf.InitConf(path)
	//设置verbose值
	Level = Conf.DebugLevel
	LogFile = Conf.LogFile
	IsConsole = Conf.EnableConsole
	//初始化日志记录器
	if LogFile != 0 {
		InitLoger(path)
	}
	//初始化TCP
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", Conf.PORT))
	if err != nil {
		fmt.Println("Init tcp socket error !\n")
		panic(err)
	}
	defer server.Close()
	//初始化UDP
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
	//初始化数据库
	if Conf.EnableDataBase != 0 {
		DB, err = InitDatabase(path + "\\cso2.db")
		if err != nil {
			fmt.Println("Init database failed !")
			Conf.EnableDataBase = 0
		} else {
			fmt.Println("Database connected !")
			defer DB.Close()
		}
	}
	//初始化Redis
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
	//初始化主频道服务器
	MainServer = NewMainServer()
	//开启UDP服务
	go StartHolePunchServer(strconv.Itoa(int(Conf.HolePunchPort)), holepunchserver)
	//开启TCP服务
	go TCPServer(server)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	_ = <-ch
}

func TCPServer(server net.Listener) {
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
	defer client.Close() //关闭con
	var seq uint8 = 0
	client.Write([]byte("~SERVERCONNECTED\n"))
	for {
		//读取4字节数据包头部
		headBytes := ReadHead(client)
		var headPacket Header
		headPacket.Data = headBytes
		headPacket.PraseHeadPacket()
		if !headPacket.IsGoodPacket {
			DebugInfo(2, "Recived a illegal head from", client.RemoteAddr().String())
			continue
		}
		//读取数据部分
		bytes := ReadData(client, headPacket.Length)
		dataPacket := Packet{
			bytes,
			headPacket.Sequence,
			headPacket.Length,
			bytes[0],
			1,
		}
		//执行功能
		switch dataPacket.id {
		case TypeQuickJoin:
			onQuick(&seq, pkt, client)
		case TypeVersion:
			OnVersionPacket(&seq, client)
		case TypeLogin:
			onLoginPacket(&seq, &pkt, &client)
		case TypeRequestChannels:
			onServerList(&seq, &pkt, &client)
		case TypeRequestRoomList:
			onRoomList(&seq, &pkt, client)
		case TypeRoom:
			onRoomRequest(&seq, pkt, client)
		case TypeHost:
			onHost(&seq, pkt, client)
		case TypeFavorite:
			onFavorite(&seq, pkt, client)
		case TypeOption:
			onOption(pkt, client)
		case TypePlayerInfo:
			onPlayerInfo(pkt, client)
		default:
			DebugInfo(2, "Unknown packet", pkt.id, "from", client.RemoteAddr().String())
		}
	}
end:
	DebugInfo(1, "client", client.RemoteAddr().String(), "closed the connection")
	delUserWithConn(client)
	client.Close() //关闭con
	return
}

func ReadHead(client net.Conn) []byte {
	head, curlen := make([]byte, HeaderLen), 0
	for {
		n, err := client.Read(bytes[curlen:])
		if err != nil {
			goto end
		}
		curlen += n
		if curlen >= HeaderLen {
			break
		}
	}
	return head
}

func ReadData(client net.Conn, len uint16) []byte {
	data, curlen := make([]byte, len), 0
	for {
		n, err := client.Read(data[curlen:])
		if err != nil {
			goto end
		}
		curlen += n
		if curlen >= len {
			break
		}
	}
	return data
}
