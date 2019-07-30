package main

import "github.com/qizikd/EthInfo/sync"

func main() {
	go sync.UpdateEthGasUsed()
	go sync.UpdateErc20GasUsed()
	sync.Start()
	//s := "a9059cbb00000000000000000000000075e89d5979e4f6fba9f97c104c2f0afb3f1dcb880000000000000000000000000211f3cedbef3143223d3acf0e589747933e852700000000000000000000000000000000000000000000000000000186b3ec24b8"
	//fmt.Println(s[8+24:8+64])
	//fmt.Println(s[8+64+24:8+64+64])
	//value := big.Int{}
	//value.SetString(s[8+64+64:len(s)],16)
	//fmt.Println(value.Int64())
}
