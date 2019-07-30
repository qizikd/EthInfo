package db

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/qizikd/EthInfo/core"
	"github.com/qizikd/EthInfo/core/mysql"
	"math/big"
)

func InserEthtx(blocknum int64, blockhash, txid, from, to string, value, gasprice, gaslimit, gasuse int64, time string, status int64) (err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT txid FROM eth WHERE txid = '%s'", txid)
	rows, err := db.Query(sql)
	defer rows.Close()
	if rows.Next() {
		return
	}
	sql = fmt.Sprintf("INSERT INTO eth(blocknum,blockhash,txid,`from`,`to`,`value`,gasprice,gaslimit,gasuse,createtime, status) VALUES(%d,'%s','%s','%s','%s',%d,%d,%d,%d,'%s',%d)",
		blocknum, blockhash, txid, from, to, value, gasprice, gaslimit, gasuse, time, status)

	_, err = db.Exec(sql)
	return
}

func InserErc20tx(blocknum int64, blockhash, txid, from, to, token string, value, gasprice, gaslimit, gasuse int64, time string, status int64) (err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT txid FROM erc20 WHERE txid = '%s'", txid)
	rows, err := db.Query(sql)
	defer rows.Close()
	if rows.Next() {
		return
	}
	sql = fmt.Sprintf("INSERT INTO erc20(blocknum,blockhash,txid,`from`,`to`,token,`value`,gasprice,gaslimit,gasuse,createtime,status) VALUES(%d,'%s','%s','%s','%s','%s',%d,%d,%d,%d,'%s', %d)",
		blocknum, blockhash, txid, from, to, token, value, gasprice, gaslimit, gasuse, time, status)

	_, err = db.Exec(sql)
	return
}

func GetEthTxs(lastid int, limit int) (txs []core.TxInfo, err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT id,txid FROM eth WHERE id > %d  and update_status = 0 order by id asc limit %d", lastid, limit)
	//fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		glog.Error("查询失败 ", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		tx := core.TxInfo{}
		err = rows.Scan(&tx.Id, &tx.TxId)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	return
}

func GetErc20Txs(lastid int, limit int) (txs []core.TxInfo, err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT id,txid FROM erc20 WHERE id > %d and update_status = 0 order by id asc limit %d", lastid, limit)
	//fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		glog.Error("查询失败 ", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		tx := core.TxInfo{}
		err = rows.Scan(&tx.Id, &tx.TxId)
		if err != nil {
			return
		}
		txs = append(txs, tx)
	}
	return
}

func UpdateEthGasused(id int, Gasused int64, status int64) (err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("UPDATE eth SET gasuse = %d, `status` = %d WHERE id = %d ", Gasused, status, id))
	if err != nil {
		glog.Error("更新失败 ", err.Error())
		return
	}
	return
}

func UpdateErc20Gasused(id int, Gasused int64, status int64) (err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("UPDATE erc20 SET gasuse = %d, `status` = %d WHERE id = %d ", Gasused, status, id))
	if err != nil {
		glog.Error("更新失败 ", err.Error())
		return
	}
	return
}

func GetCoinLastblocknum(coin string) (num int64, err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT blocknum FROM coin WHERE coin = '%s'", coin)
	rows, err := db.Query(sql)
	defer rows.Close()
	if !rows.Next() {
		num = 0
		return
	}
	err = rows.Scan(&num)
	return
}

func SetCoinLastblocknum(coin string, num int64) (err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("UPDATE coin SET blocknum = %d WHERE coin = '%s' ", num, coin))
	if err != nil {
		glog.Error("更新失败 ", err.Error())
		return
	}
	return
}

func GetEthtxsByaddress(address string, limit int) (txs []core.TxInfo, err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT txid,`from`,`to`,`value`,gasprice,gaslimit,gasuse,ispengding,createtime FROM eth WHERE `from` = '%s' or `to` = '%s' order by id desc limit %d", address, address, limit)
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		glog.Error("查询失败 ", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		tx := core.TxInfo{}
		err = rows.Scan(&tx.TxId, &tx.From, &tx.To, &tx.Value, &tx.GasPrice, &tx.Gas, &tx.GasUse, &tx.IsPengding, &tx.TimeStamp)
		if err != nil {
			return
		}
		tx.Status = 1
		price := big.NewInt(tx.GasPrice)
		use := big.NewInt(tx.GasUse)
		fee := price.Mul(price, use)
		tx.Fee = fee.Int64()
		txs = append(txs, tx)
	}
	return
}

func GetErc20txsByaddress(address string, limit int) (txs []core.TxInfo, err error) {
	db, err := mysql.GetDbConn()
	if err != nil {
		glog.Error("连接数据库失败 ", err.Error())
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT txid,`from`,`to`,token,`value`,gasprice,gaslimit,gasuse,ispengding,createtime FROM erc20 WHERE `from` = '%s' or `to` = '%s' order by id desc limit %d", address, address, limit)
	fmt.Println(sql)
	rows, err := db.Query(sql)
	if err != nil {
		glog.Error("查询失败 ", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		tx := core.TxInfo{}
		err = rows.Scan(&tx.TxId, &tx.From, &tx.To, &tx.Token, &tx.Value, &tx.GasPrice, &tx.Gas, &tx.GasUse, &tx.IsPengding, &tx.TimeStamp)
		if err != nil {
			return
		}
		tx.Status = 1
		price := big.NewInt(tx.GasPrice)
		use := big.NewInt(tx.GasUse)
		fee := price.Mul(price, use)
		tx.Fee = fee.Int64()
		txs = append(txs, tx)
	}
	return
}
