package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"

	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/database/redis"
	. "github.com/KouKouChan/CSO2-Server/database/sqlite"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/model/user"
	. "github.com/KouKouChan/CSO2-Server/model/usermanager"
	. "github.com/KouKouChan/CSO2-Server/server"
	. "github.com/KouKouChan/CSO2-Server/verbose"
	"github.com/garyburd/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//SERVERVERSION 版本号
	SERVERVERSION = "v0.3.0"
	DB            *sql.DB
	Redis         redis.Conn
	Conf          CSO2Conf
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
	//获取路径
	path, err := GetExePath()
	if err != nil {
		panic(err)
	}
	//读取配置
	Conf.InitConf(path)
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
		log.Println("Init tcp socket error !\n")
		panic(err)
	}
	defer server.Close()
	//初始化UDP
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Conf.HolePunchPort))
	if err != nil {
		log.Println("Init udp addr error !\n")
		panic(err)
	}
	holepunchserver, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Println("Init udp socket error !\n")
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
		bytes := make([]byte, math.MaxUint16)
		n, err := client.Read(bytes) //读数据
		if err == nil {
			if n == 0 {
				continue
			}
			var pkt packet
			pkt.data = bytes
			//log.Println("Prasing a packet from", client.RemoteAddr().String())
			pkt.PrasePacket()
			if !pkt.IsGoodPacket {
				DebugInfo(2, "Recived a illegal packet from", client.RemoteAddr().String())
				continue
			}
			switch pkt.id {
			case TypeQuickJoin:
				onQuick(&seq, pkt, client)
			case TypeVersion:
				onVersionPacket(&seq, pkt, client)
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
		} else {
			DebugInfo(1, "client", client.RemoteAddr().String(), "closed the connection")
			delUserWithConn(client)
			client.Close() //关闭con
			return
		}
	}
}
