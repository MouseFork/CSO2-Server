package main
const (
	UdpPacketSignature = 0x57
)

type InUDPmsg struct {
	signature uint8
    userId uint32
    portId uint16
    ipAddress uint32
    port uint16
    packetData []byte
    curOffset  uint32	//可能32位
    parsedSuccessfully bool
}
