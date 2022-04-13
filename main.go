package main

import (
	"context"
	cusutil "ethereum-client-example-go/util"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RPC_RUL = "https://ropsten.infura.io/v3/7ed60e8c0f574496b5297cdc767ac2c0"
)

func GetBalanceOf(client *ethclient.Client, address common.Address) (balance, pendingBalance *big.Int) {
	defer cusutil.TimeTrack(time.Now(), "GetBalanceOf")
	var w sync.WaitGroup
	ctx := context.Background()

	w.Add(2)
	go func() {
		b, err := client.BalanceAt(ctx, address, nil)
		balance = b
		if err != nil {
			log.Fatal(err)
		}
		w.Done()
	}()
	go func() {
		pb, err := client.PendingBalanceAt(ctx, address)
		pendingBalance = pb
		if err != nil {
			log.Fatal(err)
		}
		w.Done()
	}()
	w.Wait()
	return
}

// legacy transaction
// todo: goroutine
func TransferEtherTransaction(client *ethclient.Client, from, to string, amount *big.Int, privateHexKey string) (receipt *types.Receipt) {
	ctx := context.Background()
	privateKey, err := crypto.HexToECDSA(privateHexKey)
	if err != nil {
		log.Fatal(err)
	}

	if !common.IsHexAddress(from) || !common.IsHexAddress(to) {
		log.Fatal("Is not address")
	}

	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(from))
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	callMsg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Gas:      0,
		GasPrice: big.NewInt(0),
		Value:    amount,
		Data:     nil,
	}

	gasUsed, err := client.EstimateGas(ctx, callMsg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("gasLimit:", gasUsed)

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("gasPrice:", gasPrice)

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(to),
		amount,
		gasUsed,
		gasPrice,
		nil,
	)

	chainId, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("transactionHash:", tx.Hash().String())

	receipt, err = bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	ctx := context.Background()

	client, _ := ethclient.DialContext(ctx, RPC_RUL)
	defer client.Close()

	from := "0x2E14747bF40385F958F4C5265A312504935a2afE"
	privateKey := "43f035f47cae3a368377d4418dac7c830f6a40a029947cb9f746fefc70537d72" // todo
	to := "0xdDF36eCdf1fA200a1dFF510544CA1B948d9e7Fec"
	amount := big.NewInt(int64(math.Pow10(18)))
	receipt := TransferEtherTransaction(client, from, to, amount, privateKey)
	fmt.Println("gasUsed:", receipt.GasUsed)
	fmt.Println("txHash:", receipt.TxHash)
}
