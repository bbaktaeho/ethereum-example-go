package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"

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

// func CreateEIP1559Transaction(client *ethclient.Client, to, privateHexKey string, amount *big.Int) (receipt *types.Receipt, e error) {
// 	ctx := context.Background()
// 	privateKey, err := crypto.HexToECDSA(privateHexKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	from := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

// 	if !common.IsHexAddress(to) {
// 		e = errors.New("invalid to address")
// 	}

// 	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(from))
// 	if err != nil {
// 		return nil, err
// 	}
// 	block.BaseFee()

// 	// tx := types.NewTx(&types.DynamicFeeTx{
// 	// 	Nonce: nonce,
// 	// })
// }

// legacy transaction
// todo: goroutine
func TransferEtherTransaction(client *ethclient.Client, from, to string, amount *big.Int, privateHexKey string) (receipt *types.Receipt, e error) {
	ctx := context.Background()
	privateKey, err := crypto.HexToECDSA(privateHexKey)
	if err != nil {
		e = err
	}

	if !common.IsHexAddress(from) || !common.IsHexAddress(to) {
		e = errors.New("invalid addres")
	}

	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(from))
	if err != nil {
		e = err
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
		e = err
	}
	log.Println("gasLimit:", gasUsed)

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		e = err
	}
	log.Println("gasPrice:", gasPrice)

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
		e = err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		e = err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		e = err
	}
	log.Println("transactionHash:", tx.Hash().String())

	receipt, err = bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		e = err
	}
	return
}

func main() {
	ctx := context.Background()

	client, _ := ethclient.DialContext(ctx, RPC_RUL)
	defer client.Close()

	// from := "0x2E14747bF40385F958F4C5265A312504935a2afE"
	// privateKey := "43f035f47cae3a368377d4418dac7c830f6a40a029947cb9f746fefc70537d72" // todo
	// toHexAddress := "0xdDF36eCdf1fA200a1dFF510544CA1B948d9e7Fec"
	// to := common.HexToAddress(toHexAddress)
	// value := big.NewInt(int64(math.Pow10(18)))

	// receipt, err := feature.TransferValue(client, privateKey, &to, value, 0)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	txHash := common.HexToHash("0xd6499e8217a72c71b49a8967263ce8f9b8b82d3c72df492008114590b7b54da6")
	receipt, _ := client.TransactionReceipt(ctx, txHash)
	obj, _ := json.Marshal(receipt)
	log.Printf("%+v\n", obj)
}
