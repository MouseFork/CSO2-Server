package main

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

const (
	Congratulate    = 11
	SystemImportant = 20
	DialogBox       = 21
	System          = 22
	DialogBoxExit   = 60
)

var (
	GAME_ROOM_JOIN_FAILED_CLOSED         = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_CLOSED")
	GAME_ROOM_JOIN_FAILED_FULL           = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_FULL")
	GAME_ROOM_JOIN_FAILED_BAD_PASSWORD   = []byte("#CSO2_POPUP_ROOM_JOIN_FAILED_INVALID_PASSWD")
	GAME_ROOM_CHANGETEAM_FAILED          = []byte("#CSO2_POPUP_ROOM_CHANGETEAM_FAILED")
	GAME_ROOM_COUNTDOWN_FAILED_NOENEMIES = []byte("#CSO2_UI_ROOM_COUNTDOWN_FAILED_NOENEMY")
	GAME_LOGIN_BAD_USERNAME              = []byte("#CSO2_LoginAuth_Certify_NoPassport")
	GAME_LOGIN_BAD_PASSWORD              = []byte("#CSO2_LoginAuth_WrongPassword")
	GAME_LOGIN_INVALID_USERINFO          = []byte("#CSO2_ServerMessage_INVALID_USERINFO")
)

func onSendMessage(seq *uint8, client net.Conn, tp uint8, msg []byte) {
	var p packet
	p.id = TypeChat
	rst := BuildHeader(seq, p)
	rst = append(rst, tp)
	rst = BytesCombine(rst, BuildMessage(msg, tp))
	sendPacket(rst, client)
}

func BuildMessage(msg []byte, tp uint8) []byte {
	if tp == Congratulate {
		buf := make([]byte, 1)
		buf[0] = 0
		return BytesCombine(buf, BuildString(msg))
	}
	return BuildLongString(msg)
}
