package sync

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"github.com/qizikd/EthInfo/db"
	"math/big"
	"time"
)

var Host = "http://127.0.0.1:8545"

//var Host = "http://47.244.176.129:8545"

func Start() {
	client, err := ethclient.Dial(Host)
	if err != nil {
		glog.Error("连接infura节点失败", err)
		err = errors.New("连接infura节点失败")
		return
	}
	defer client.Close()
	lastblocknum, err := db.GetCoinLastblocknum("eth")
	if err != nil {
		lastblocknum = 0
	}
	for {
		fmt.Printf("block:%d\n", lastblocknum)
		blocknum := big.NewInt(lastblocknum)
		block, err := client.BlockByNumber(context.Background(), blocknum)
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		fmt.Printf("txnum:%d\n", len(block.Transactions()))
		//go func() {
		sync(block, client)
		//}()
		lastblocknum++
		db.SetCoinLastblocknum("eth", lastblocknum)
	}
}

func UpdateEthGasUsed() {
	client, err := ethclient.Dial(Host)
	if err != nil {
		glog.Error("连接infura节点失败", err)
		err = errors.New("连接infura节点失败")
		return
	}
	defer client.Close()
	lastid := -1
	for {
		txs, err := db.GetEthTxs(lastid, 100)
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err.Error())
			continue
		}
		if len(txs) == 0 {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("eth id:%d txid:%s\n", txs[0].Id, txs[0].TxId)
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			rec, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.TxId))
			if err != nil {
				fmt.Printf("获取交易接受详情失败(%s):%s\n", common.HexToHash(tx.TxId), err.Error())
				continue
			}
			gasUsed := big.NewInt(int64(rec.GasUsed))
			status := rec.Status
			err = db.UpdateEthGasused(tx.Id, gasUsed.Int64(), int64(status))
			if err != nil {
				fmt.Printf("获取交易接受详情失败(%s):%s\n", common.HexToHash(tx.TxId), err.Error())
				continue
			}
			lastid = tx.Id
		}
	}
}

func UpdateErc20GasUsed() {
	client, err := ethclient.Dial(Host)
	if err != nil {
		glog.Error("连接infura节点失败", err)
		err = errors.New("连接infura节点失败")
		return
	}
	defer client.Close()
	lastid := -1
	for {
		txs, err := db.GetEthTxs(lastid, 100)
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err.Error())
			continue
		}
		if len(txs) == 0 {
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("erc20 id:%d txid:%s\n", txs[0].Id, txs[0].TxId)
		for i := 0; i < len(txs); i++ {
			tx := txs[i]
			rec, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.TxId))
			if err != nil {
				fmt.Printf("获取交易接受详情失败(%s):%s\n", common.HexToHash(tx.TxId), err.Error())
				continue
			}
			gasUsed := big.NewInt(int64(rec.GasUsed))
			status := rec.Status
			err = db.UpdateEthGasused(tx.Id, gasUsed.Int64(), int64(status))
			if err != nil {
				fmt.Printf("获取交易接受详情失败(%s):%s\n", common.HexToHash(tx.TxId), err.Error())
				continue
			}
			lastid = tx.Id
		}
	}
}

func sync(block *types.Block, client *ethclient.Client) {
	for i := 0; i < len(block.Transactions()); i++ {
		tx := block.Transactions()[i]
		if tx.To() == nil {
			//创建合约交易
			continue
		}
		to := tx.To().Hex()
		singer := types.NewEIP155Signer(tx.ChainId())
		_from, err := singer.Sender(tx)
		from := _from.Hex()
		if err != nil {
			from = tx.From.Hex()
			fmt.Println("解析from错误:", err.Error(), tx.Hash().Hex())
			//continue
		}
		gaslimt := big.NewInt(int64(tx.Gas()))
		//rec, err := client.TransactionReceipt(context.Background(),tx.Hash())
		//if err != nil {
		//	fmt.Print("获取交易接受详情失败(%s):%s\n",tx.Hash().Hex(),err.Error())
		//	continue
		//}
		gasUsed := gaslimt //big.NewInt(int64(rec.GasUsed))
		status := int64(0)
		value := tx.Value()
		if value.Int64() == 0 {
			data := hex.EncodeToString(tx.Data())
			if len(data) < 72 {
				continue
			}
			//fmt.Println("合约交易:", data)
			token := to
			if data[0:8] == "a9059cbb" {
				//
				to = data[8+24 : 8+64]
				value.SetString(data[8+64:len(data)], 16)
			} else if data[0:8] == "23b872dd" {
				//
				from = data[8+24 : 8+64]
				to = data[8+64+24 : 8+64+64]
				value.SetString(data[8+64+64:len(data)], 16)
			} else {
				continue
			}
			err = db.InserErc20tx(block.Number().Int64(), block.Hash().Hex(), tx.Hash().Hex(), from, to, token, value.Int64(),
				tx.GasPrice().Int64(), gaslimt.Int64(), gasUsed.Int64(), time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"), status)
		} else {
			err = db.InserEthtx(block.Number().Int64(), block.Hash().Hex(), tx.Hash().Hex(), from, to, value.Int64(),
				tx.GasPrice().Int64(), gaslimt.Int64(), gasUsed.Int64(), time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"), status)
		}
		if err != nil {
			fmt.Println("插入交易失败：", err.Error())
			continue
		}
	}
}
