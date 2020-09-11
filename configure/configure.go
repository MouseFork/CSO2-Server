package configure

import (
	"fmt"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

var (
	Conf CSO2Conf
)

type CSO2Conf struct {
	PORT             uint32
	HolePunchPort    uint32
	EnableRedis      uint32
	EnableDataBase   uint32
	UseJson          uint32
	MaxUsers         uint32
	EnableShop       uint32
	UnlockAllWeapons uint32
	UnlockAllSkills  uint32
	RedisIP          string
	RedisPort        uint32
	DebugLevel       uint32
	LogFile          uint32
	EnableConsole    uint32
	EnableRegister   uint32
	EnableMail       uint32
	REGPort          uint32
	REGEmail         string
	REGPassWord      string
	REGSMTPaddr      string
}

func (conf *CSO2Conf) InitConf(path string) {
	if conf == nil {
		return
	}
	ini_parser := IniParser{}
	file := path + "\\server.conf"
	if err := ini_parser.LoadIni(file); err != nil {
		fmt.Printf("try load config file error[%s]\n", err.Error())
		fmt.Printf("Using default data ...\n")
		conf.EnableRedis = 0
		conf.EnableDataBase = 1
		conf.UseJson = 0
		conf.MaxUsers = 0
		conf.EnableShop = 0
		conf.UnlockAllWeapons = 1
		conf.UnlockAllSkills = 1
		conf.PORT = 30001
		conf.HolePunchPort = 30002
		conf.RedisIP = "127.0.0.1"
		conf.RedisPort = 6379
		conf.DebugLevel = 2
		conf.LogFile = 1
		conf.EnableConsole = 0
		conf.EnableRegister = 1
		conf.EnableMail = 0
		return
	}
	conf.EnableRedis = ini_parser.IniGetUint32("Database", "EnableRedis")
	conf.EnableDataBase = ini_parser.IniGetUint32("Database", "EnableDataBase")
	conf.UseJson = ini_parser.IniGetUint32("Database", "UseJson")
	conf.MaxUsers = ini_parser.IniGetUint32("Server", "MaxUsers")
	if conf.MaxUsers < 0 {
		conf.MaxUsers = 0
	}
	conf.EnableShop = ini_parser.IniGetUint32("Server", "EnableShop")
	conf.UnlockAllWeapons = ini_parser.IniGetUint32("Server", "UnlockAllWeapons")
	conf.UnlockAllSkills = ini_parser.IniGetUint32("Server", "UnlockAllSkills")
	conf.PORT = ini_parser.IniGetUint32("Server", "TCPPort")
	conf.HolePunchPort = ini_parser.IniGetUint32("Server", "UDPPort")
	conf.RedisIP = ini_parser.IniGetString("Server", "RedisIP")
	conf.RedisPort = ini_parser.IniGetUint32("Server", "RedisPort")
	conf.DebugLevel = ini_parser.IniGetUint32("Debug", "DebugLevel")
	if conf.DebugLevel > 2 || conf.DebugLevel < 0 {
		conf.DebugLevel = 2
	}
	conf.LogFile = ini_parser.IniGetUint32("Debug", "LogFile")
	conf.EnableConsole = ini_parser.IniGetUint32("Debug", "EnableConsole")
	conf.EnableRegister = ini_parser.IniGetUint32("Register", "EnableRegister")
	conf.EnableMail = ini_parser.IniGetUint32("Register", "EnableEmail")
	conf.REGPort = ini_parser.IniGetUint32("Register", "REGPort")
	conf.REGEmail = ini_parser.IniGetString("Register", "REGEmail")
	conf.REGPassWord = ini_parser.IniGetString("Register", "REGPassWord")
	conf.REGSMTPaddr = ini_parser.IniGetString("Register", "REGSMTPaddr")
}
