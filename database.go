package main

import (
	"crypto/md5"
	"errors"
	"log"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

//从数据库中读取用户数据
//如果是新用户则保存到数据库中
func getUserFromDatabase(p loginPacket) user {
	if DB != nil {
		query, err := DB.Prepare("SELECT * FROM userinfo WHERE LoginName = ?")
		if err == nil {
			defer query.Close()
			u := getNewUser()
			var inventory []byte
			var clanID uint32
			err = query.QueryRow(p.nexonUsername).Scan(&u.loginName, &u.username, &u.password, &u.level, &u.rank,
				&u.rankFrame, &u.points, &u.currentExp, &u.playedMatches, &u.wins, &u.kills,
				&u.headshots, &u.deaths, &u.assists, &u.accuracy, &u.secondsPlayed, &u.netCafeName,
				&u.cash, &clanID, &u.worldRank, &u.mpoints, &u.titleId, &u.unlockedTitles, &u.signature,
				&u.bestGamemode, &u.bestMap, &u.unlockedAchievements, &u.avatar, &u.unlockedAvatars,
				&u.vipLevel, &u.vipXp, &u.skillHumanCurXp, &u.skillHumanPoints, &u.skillZombieCurXp,
				&u.skillZombiePoints, &inventory)
			if err != nil {
				log.Println("Suffered a error while getting User", string(p.nexonUsername)+"'s data !", err)
				u = getNewUser()
				u.setID(getNewUserID())
				u.setUserName(p)
				u.password = p.PassWd
				CheckErr(AddUserToDB(u))
				return u
			}
			//检查密码
			str := md5.Sum(p.PassWd)
			for i := 0; i < 16; i++ {
				if str[i] != u.password[i] {
					u = getNewUser()
					u.setID(getNewUserID())
					u.setUserName(p)
					u.password = p.PassWd
					log.Println("password error!", str, u.password)
					return u
				}
			}
			//设置仓库
			u.inventory = praseInventory(inventory)
			//设置战队...
			log.Println("User", string(u.username), "data found !")
			u.setID(getNewUserID())
			return u
			// u.setID(getNewUserID())
			// u.setUserName(p)
			// u.password = p.PassWd
			// CheckErr(AddUserToDB(u))
			// return u
		} else { //出错
			u := getNewUser()
			u.setID(getNewUserID())
			u.setUserName(p)
			u.password = p.PassWd
			log.Println(err)
			CheckErr(AddUserToDB(u))
			return u
		}
	}
	u := getNewUser()
	u.setID(getNewUserID())
	u.setUserName(p)
	u.password = p.PassWd
	return u
}

func praseInventory(inventory []byte) userInventory {
	var inv userInventory
	offset := 0
	inv.numOfItem = ReadUint16(inventory, &offset)
	for i := 0; i < int(inv.numOfItem); i++ {
		var it userInventoryItem
		it.id = ReadUint32(inventory, &offset)
		it.count = ReadUint16(inventory, &offset)
		inv.items = append(inv.items, it)
	}
	inv.CTModel = ReadUint32(inventory, &offset)
	inv.TModel = ReadUint32(inventory, &offset)
	inv.headItem = ReadUint32(inventory, &offset)
	inv.gloveItem = ReadUint32(inventory, &offset)
	inv.backItem = ReadUint32(inventory, &offset)
	inv.stepsItem = ReadUint32(inventory, &offset)
	inv.cardItem = ReadUint32(inventory, &offset)
	inv.sprayItem = ReadUint32(inventory, &offset)
	//buymenu
	len := ReadUint8(inventory, &offset)
	inv.buyMenu.pistols = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.buyMenu.shotguns = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.buyMenu.smgs = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.buyMenu.rifles = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.buyMenu.snipers = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.buyMenu.machineguns = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.buyMenu.melees = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.buyMenu.equipment = ReadUint32Array(inventory, &offset, int(len))
	//loadouts
	len = ReadUint8(inventory, &offset)
	for i := 0; i < int(len); i++ {
		var ul userLoadout
		l := ReadUint8(inventory, &offset)
		ul.items = ReadUint32Array(inventory, &offset, int(l))
		inv.loadouts = append(inv.loadouts, ul)
	}
	return inv
}

func InventoryToBytes(inventory userInventory) []byte {
	buf := make([]byte, 8096)
	offset := 0
	WriteUint16(&buf, inventory.numOfItem, &offset)
	for i := 0; i < int(inventory.numOfItem); i++ {
		WriteUint32(&buf, inventory.items[i].id, &offset)
		WriteUint16(&buf, inventory.items[i].count, &offset)
	}
	WriteUint32(&buf, inventory.CTModel, &offset)
	WriteUint32(&buf, inventory.TModel, &offset)
	WriteUint32(&buf, inventory.headItem, &offset)
	WriteUint32(&buf, inventory.gloveItem, &offset)
	WriteUint32(&buf, inventory.backItem, &offset)
	WriteUint32(&buf, inventory.stepsItem, &offset)
	WriteUint32(&buf, inventory.cardItem, &offset)
	WriteUint32(&buf, inventory.sprayItem, &offset)
	//buymenu
	WriteUint8(&buf, uint8(len(inventory.buyMenu.pistols)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.pistols, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.shotguns)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.shotguns, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.smgs)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.smgs, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.rifles)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.rifles, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.snipers)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.snipers, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.machineguns)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.machineguns, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.melees)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.melees, &offset)
	WriteUint8(&buf, uint8(len(inventory.buyMenu.equipment)), &offset)
	WriteUint32Array(&buf, inventory.buyMenu.equipment, &offset)
	//loadouts
	WriteUint8(&buf, uint8(len(inventory.loadouts)), &offset)
	for i := 0; i < len(inventory.loadouts); i++ {
		WriteUint8(&buf, uint8(len(inventory.loadouts[i].items)), &offset)
		WriteUint32Array(&buf, inventory.loadouts[i].items, &offset)
	}
	return buf[:offset]
}

func AddUserToDB(u user) error {
	if DB == nil {
		return errors.New("DataBase not connected")
	}
	stmt, err := DB.Prepare(`INSERT INTO userinfo(LoginName, UserName, PassWord,  
		Level, Rank, RankFrame, Points, CurrentExp, PlayedMatches, Wins, Kills,	
		HeadShots, Deathes, Assists, accuracy, SecondsPlayed, netCafeName,	
		Cash, ClanID, WorldRank, Mpoints, TitleID, UnlockefTitleID, signature,	
		bestGamemode, bestMap, unlockedAchievements, avatar, unlockedAvatars,	
		viplevel, vipXp, skillHumanCurXp, skillHumanPoints, skillZombieCurXp,	
		skillZombiePoints, Inventory) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?	
		   ,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`) //36个
	if err != nil {
		return err
	}
	defer stmt.Close()
	pass := md5.Sum(u.password)
	_, err = stmt.Exec(u.loginName, u.username, pass[:], u.level, u.rank,
		u.rankFrame, u.points, u.currentExp, u.playedMatches, u.wins, u.kills,
		u.headshots, u.deaths, u.assists, u.accuracy, u.secondsPlayed, u.netCafeName,
		u.cash, 0, u.worldRank, u.mpoints, u.titleId, u.unlockedTitles, u.signature,
		u.bestGamemode, u.bestMap, u.unlockedAchievements, u.avatar, u.unlockedAvatars,
		u.vipLevel, u.vipXp, u.skillHumanCurXp, u.skillHumanPoints, u.skillZombieCurXp,
		u.skillZombiePoints, InventoryToBytes(u.inventory))
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserToDB(u user) error {
	if DB == nil {
		return errors.New("DataBase not connected")
	}
	stmt, err := DB.Prepare(`Update userinfo set Level=?, 
		Rank=?, RankFrame=?, Points=?, CurrentExp=?, PlayedMatches=?, Wins=?, Kills=?,	
		HeadShots=?, Deathes=?, Assists=?, accuracy=?, SecondsPlayed=?, netCafeName=?,	
		Cash=?, ClanID=?, WorldRank=?, Mpoints=?, TitleID=?, UnlockefTitleID=?, signature=?,	
		bestGamemode=?, bestMap=?, unlockedAchievements=?, avatar=?, unlockedAvatars=?,	
		viplevel=?, vipXp=?, skillHumanCurXp=?, skillHumanPoints=?, skillZombieCurXp=?,	
		skillZombiePoints=?, Inventory=? WHERE LoginName=? `) //36个
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.level, u.rank,
		u.rankFrame, u.points, u.currentExp, u.playedMatches, u.wins, u.kills,
		u.headshots, u.deaths, u.assists, u.accuracy, u.secondsPlayed, u.netCafeName,
		u.cash, 0, u.worldRank, u.mpoints, u.titleId, u.unlockedTitles, u.signature,
		u.bestGamemode, u.bestMap, u.unlockedAchievements, u.avatar, u.unlockedAvatars,
		u.vipLevel, u.vipXp, u.skillHumanCurXp, u.skillHumanPoints, u.skillZombieCurXp,
		u.skillZombiePoints, InventoryToBytes(u.inventory), u.loginName)
	if err != nil {
		return err
	}
	return nil
}
