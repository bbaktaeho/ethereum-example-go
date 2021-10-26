package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var url = "http://localhost:8545"

// fee collector
var address = "0xaC84E3540BC6866eeF4d2735462a05b090Ea4f10"

func main() {
	client, err := ethclient.DialContext(context.Background(), url)
	if err != nil {
		log.Fatalf("Error to create a ether client:%v", err)
	}
	defer client.Close()
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error to get a block:%v", err)
	}
	fmt.Println("block number:", block.Number())

	// balance := GetBalance(address, client)
	// fmt.Println(balance)

	// publicKey, privateKey, address := GenerateAddress()
	// fmt.Println("public key:", publicKey)
	// fmt.Println("private key:", privateKey)
	// fmt.Println("address:", address)

	// fmt.Println()

	// GenerateKeystore()
	// ReadKeystore()
	GetEther()
}

func GetBalance(addr string, client *ethclient.Client) *big.Int {
	address := common.HexToAddress(addr)
	// wei (big int)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Error to get the balance:%v", err)
	}
	// fBalance := new(big.Float)
	// fBalance.SetString(balance.String())
	// // wei (bit float)
	// fmt.Println(fBalance)
	// // ether
	// value := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
	// fmt.Println(value)
	return balance
}

func GetEther() {
	url := "https://rinkeby.infura.io/v3/38ff0015e86549b28e5ad2d68f876a23"
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	a1 := common.HexToAddress("96435b7ac04116a667e7b6761436e429ae7dc885")
	a2 := common.HexToAddress("56f91f4c0ea9cda4ff24c4643c5c1cf2b23e2fef")

	b1, err := client.BalanceAt(context.Background(), a1, nil)
	if err != nil {
		log.Fatal(err)
	}

	b2, err := client.BalanceAt(context.Background(), a2, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b1)
	fmt.Println(b2)
}
