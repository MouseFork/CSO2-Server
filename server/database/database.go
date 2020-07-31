package database

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"log"

	. "github.com/KouKouChan/CSO2-Server/blademaster"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/server/user"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	DB *sql.DB
)

//从数据库中读取用户数据
//如果是新用户则保存到数据库中
func getUserFromDatabase(loginname, passwd []byte) User {
	if DB != nil {
		query, err := DB.Prepare("SELECT * FROM userinfo WHERE LoginName = ?")
		if err == nil {
			defer query.Close()
			u := GetNewUser()
			var inventory []byte
			var clanID uint32
			err = query.QueryRow(loginname).Scan(&u.LoginName, &u.Username, &u.Password, &u.Level, &u.Rank,
				&u.RankFrame, &u.Points, &u.CurrentExp, &u.PlayedMatches, &u.Wins, &u.Kills,
				&u.Headshots, &u.Deaths, &u.Assists, &u.Accuracy, &u.SecondsPlayed, &u.NetCafeName,
				&u.Cash, &clanID, &u.WorldRank, &u.Mpoints, &u.TitleId, &u.UnlockedTitles, &u.Signature,
				&u.BestGamemode, &u.BestMap, &u.UnlockedAchievements, &u.Avatar, &u.UnlockedAvatars,
				&u.VipLevel, &u.VipXp, &u.SkillHumanCurXp, &u.SkillHumanPoints, &u.SkillZombieCurXp,
				&u.SkillZombiePoints, &inventory)
			if err != nil {
				DebugInfo(1, "Suffered a error while getting User", string(loginname)+"'s data !", err)
				u = GetNewUser()
				u.SetID(GetNewUserID())
				u.SetUserName(loginname, loginname)
				u.Password = passwd
				CheckErr(AddUserToDB(u))
				return u
			}
			//检查密码
			str := md5.Sum(passwd)
			for i := 0; i < 16; i++ {
				if str[i] != u.Password[i] {
					u = GetNewUser()
					u.SetID(GetNewUserID())
					u.SetUserName(loginname, loginname)
					u.Password = passwd
					DebugInfo(1, "password error!", str, u.Password)
					return u
				}
			}
			//设置仓库
			u.Inventory = praseInventory(inventory)
			//设置战队...
			DebugInfo(1, "User", string(u.Username), "data found !")
			u.SetID(GetNewUserID())
			return u
			// u.setID(getNewUserID())
			// u.setUserName(p)
			// u.password = p.PassWd
			// CheckErr(AddUserToDB(u))
			// return u
		} else { //出错
			u := GetNewUser()
			u.SetID(GetNewUserID())
			u.SetUserName(loginname, loginname)
			u.Password = passwd
			log.Println(err)
			CheckErr(AddUserToDB(u))
			return u
		}
	}
	u := GetNewUser()
	u.SetID(GetNewUserID())
	u.SetUserName(loginname, loginname)
	u.Password = passwd
	return u
}

func praseInventory(inventory []byte) UserInventory {
	var inv UserInventory
	offset := 0
	inv.NumOfItem = ReadUint16(inventory, &offset)
	for i := 0; i < int(inv.NumOfItem); i++ {
		var it UserInventoryItem
		it.Id = ReadUint32(inventory, &offset)
		it.Count = ReadUint16(inventory, &offset)
		inv.Items = append(inv.Items, it)
	}
	inv.CTModel = ReadUint32(inventory, &offset)
	inv.TModel = ReadUint32(inventory, &offset)
	inv.HeadItem = ReadUint32(inventory, &offset)
	inv.GloveItem = ReadUint32(inventory, &offset)
	inv.BackItem = ReadUint32(inventory, &offset)
	inv.StepsItem = ReadUint32(inventory, &offset)
	inv.CardItem = ReadUint32(inventory, &offset)
	inv.SprayItem = ReadUint32(inventory, &offset)
	//buymenu
	len := ReadUint8(inventory, &offset)
	inv.BuyMenu.Pistols = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Shotguns = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Smgs = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Rifles = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Snipers = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Machineguns = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Melees = ReadUint32Array(inventory, &offset, int(len))
	len = ReadUint8(inventory, &offset)
	inv.BuyMenu.Equipment = ReadUint32Array(inventory, &offset, int(len))
	//loadouts
	len = ReadUint8(inventory, &offset)
	for i := 0; i < int(len); i++ {
		var ul UserLoadout
		l := ReadUint8(inventory, &offset)
		ul.Items = ReadUint32Array(inventory, &offset, int(l))
		inv.Loadouts = append(inv.Loadouts, ul)
	}
	return inv
}

func InventoryToBytes(inventory UserInventory) []byte {
	buf := make([]byte, 8096)
	offset := 0
	WriteUint16(&buf, inventory.NumOfItem, &offset)
	for i := 0; i < int(inventory.NumOfItem); i++ {
		WriteUint32(&buf, inventory.Items[i].Id, &offset)
		WriteUint16(&buf, inventory.Items[i].Count, &offset)
	}
	WriteUint32(&buf, inventory.CTModel, &offset)
	WriteUint32(&buf, inventory.TModel, &offset)
	WriteUint32(&buf, inventory.HeadItem, &offset)
	WriteUint32(&buf, inventory.GloveItem, &offset)
	WriteUint32(&buf, inventory.BackItem, &offset)
	WriteUint32(&buf, inventory.StepsItem, &offset)
	WriteUint32(&buf, inventory.CardItem, &offset)
	WriteUint32(&buf, inventory.SprayItem, &offset)
	//buymenu
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Pistols)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Pistols, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Shotguns)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Shotguns, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Smgs)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Smgs, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Rifles)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Rifles, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Snipers)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Snipers, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Machineguns)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Machineguns, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Melees)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Melees, &offset)
	WriteUint8(&buf, uint8(len(inventory.BuyMenu.Equipment)), &offset)
	WriteUint32Array(&buf, inventory.BuyMenu.Equipment, &offset)
	//loadouts
	WriteUint8(&buf, uint8(len(inventory.Loadouts)), &offset)
	for i := 0; i < len(inventory.Loadouts); i++ {
		WriteUint8(&buf, uint8(len(inventory.Loadouts[i].Items)), &offset)
		WriteUint32Array(&buf, inventory.Loadouts[i].Items, &offset)
	}
	return buf[:offset]
}

func AddUserToDB(u User) error {
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
	pass := md5.Sum(u.Password)
	_, err = stmt.Exec(u.LoginName, u.Username, pass[:], u.Level, u.Rank,
		u.RankFrame, u.Points, u.CurrentExp, u.PlayedMatches, u.Wins, u.Kills,
		u.Headshots, u.Deaths, u.Assists, u.Accuracy, u.SecondsPlayed, u.NetCafeName,
		u.Cash, 0, u.WorldRank, u.Mpoints, u.TitleId, u.UnlockedTitles, u.Signature,
		u.BestGamemode, u.BestMap, u.UnlockedAchievements, u.Avatar, u.UnlockedAvatars,
		u.VipLevel, u.VipXp, u.SkillHumanCurXp, u.SkillHumanPoints, u.SkillZombieCurXp,
		u.SkillZombiePoints, InventoryToBytes(u.Inventory))
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserToDB(u User) error {
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
	_, err = stmt.Exec(u.Level, u.Rank,
		u.RankFrame, u.Points, u.CurrentExp, u.PlayedMatches, u.Wins, u.Kills,
		u.Headshots, u.Deaths, u.Assists, u.Accuracy, u.SecondsPlayed, u.NetCafeName,
		u.Cash, 0, u.WorldRank, u.Mpoints, u.TitleId, u.UnlockedTitles, u.Signature,
		u.BestGamemode, u.BestMap, u.UnlockedAchievements, u.Avatar, u.UnlockedAvatars,
		u.VipLevel, u.VipXp, u.SkillHumanCurXp, u.SkillHumanPoints, u.SkillZombieCurXp,
		u.SkillZombiePoints, InventoryToBytes(u.Inventory), u.LoginName)
	if err != nil {
		return err
	}
	return nil
}
