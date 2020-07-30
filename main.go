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
	"github.com/garyburd/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//SERVERVERSION 版本号
	SERVERVERSION = "v0.3.0"
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
	DB    *sql.DB
	Redis redis.Conn
	Conf  CSO2Conf
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
	//初始化TCP
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", Conf.PORT))
	if err != nil {
		log.Fatal("Init tcp socket error !\n")
		panic(err)
	}
	defer server.Close()
	//初始化UDP
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", Conf.HolePunchPort))
	if err != nil {
		log.Fatal("Init udp addr error !\n")
		panic(err)
	}
	holepunchserver, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Init udp socket error !\n")
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
	MainServer = newMainServer()
	//开启UDP服务
	go startHolePunchServer(holepunchserver)
	//开启TCP服务
	fmt.Println("Server is running at", "[AnyAdapter]:"+strconv.Itoa(int(Conf.PORT)))
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
				log.Println("Recived a illegal packet from", client.RemoteAddr().String())
				continue
			}
			switch pkt.id {
			case TypeQuickJoin:
				onQuick(&seq, pkt, client)
			case TypeVersion:
				//log.Println("Recived a client version packet from", client.RemoteAddr().String())
				onVersionPacket(&seq, pkt, client)
			case TypeLogin:
				//log.Println("Recived a login request packet from", client.RemoteAddr().String())
				//if ! {
				//log.Println("Recived a illegal packet from", client.RemoteAddr().String())
				//}
				onLoginPacket(&seq, &pkt, &client)
			case TypeRequestChannels:
				//log.Println("Recived a ChannelList request packet from", client.RemoteAddr().String())
				onServerList(&seq, &pkt, &client)
			case TypeRequestRoomList:
				//log.Println("Recived a RoomList request packet from", client.RemoteAddr().String())
				onRoomList(&seq, &pkt, client)
			case TypeRoom:
				//log.Println("Recived a Room request packet from", client.RemoteAddr().String())
				onRoomRequest(&seq, pkt, client)
			case TypeHost:
				//log.Println("Recived a Host request packet from", client.RemoteAddr().String())
				onHost(&seq, pkt, client)
			case TypeFavorite:
				//log.Println("Recived a favorite request packet from", client.RemoteAddr().String())
				onFavorite(&seq, pkt, client)
			case TypeOption:
				//log.Println("Recived a favorite request packet from", client.RemoteAddr().String())
				onOption(pkt, client)
			case TypePlayerInfo:
				onPlayerInfo(pkt, client)
			default:
				log.Println("Unknown packet", pkt.id, "from", client.RemoteAddr().String())
			}
		} else {
			log.Println("client", client.RemoteAddr().String(), "closed the connection")
			delUserWithConn(client)
			client.Close() //关闭con
			return
		}
	}
}
