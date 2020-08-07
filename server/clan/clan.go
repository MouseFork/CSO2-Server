package main

import(
	
	. "github.com/KouKouChan/CSO2-Server/blademaster"
)

func NewClan() Clan {
	return Clan{
		0,
		NewNullString(),
		0,
	}
}
