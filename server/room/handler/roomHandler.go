package handler

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/model/packet"
	. "github.com/KouKouChan/CSO2-Server/model/room"
)

func onRoomRequest(seq *uint8, p Packet, client net.Conn) {
	var pkt InRoomPaket
	if pkt.PraseRoomPacket(p) {
		switch pkt.InRoomType {
		case NewRoomRequest:
			//log.Println("Recived a new room request from", client.RemoteAddr().String())
			onNewRoom(seq, p, client)
		case JoinRoomRequest:
			//log.Println("Recived a join room request from", client.RemoteAddr().String())
			onJoinRoom(seq, p, client)
		case LeaveRoomRequest:
			//log.Println("Recived a leave room request from", client.RemoteAddr().String())
			onLeaveRoom(seq, p, client)
		case ToggleReadyRequest:
			//log.Println("Recived a ready request from", client.RemoteAddr().String())
			onToggleReady(seq, p, client)
		case GameStartRequest:
			//log.Println("Recived a start game request from", client.RemoteAddr().String())
			onGameStart(seq, p, client)
		case UpdateSettings:
			//log.Println("Recived a update room setting request from", client.RemoteAddr().String())
			onUpdateRoom(seq, p, client)
		case OnCloseResultWindow:
			//log.Println("Recived a close resultWindow request from", client.RemoteAddr().String())
			onCloseResultRequest(seq, p, client)
		case SetUserTeamRequest:
			//log.Println("Recived a set user team request from", client.RemoteAddr().String())
			onChangeTeam(seq, p, client)
		case GameStartCountdownRequest:
			//log.Println("Recived a begin start game request from", client.RemoteAddr().String())
			onGameStartCountdown(p, client)
		default:
			log.Println("Unknown room packet", pkt.InRoomType, "from", client.RemoteAddr().String())
		}
	} else {
		log.Println("Error : Recived a illegal room packet from", client.RemoteAddr().String())
	}
}
