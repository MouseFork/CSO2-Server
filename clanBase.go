package main

//Clan 战队
type Clan struct {
	clanID   uint32
	clanName []byte
	clanMark uint32
}

func newClan() Clan {
	return Clan{
		0,
		newNullString(),
		0,
	}
}
