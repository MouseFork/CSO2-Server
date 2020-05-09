package main

func onVersionPacket(seq *uint8, p packet) []byte {
	header := BuildHeader(seq, p)
	header[1] = 0
	*seq = 0
	IsBadHash := make([]byte, 1)
	IsBadHash[0] = 0
	hash := []byte("6246015df9a7d1f7311f888e7e861f18")
	rst := BytesCombine(header, IsBadHash, hash)
	WriteLen(&rst)
	return rst
}
