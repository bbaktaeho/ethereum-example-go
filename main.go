package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RPC_RUL = "https://ropsten.infura.io/v3/7ed60e8c0f574496b5297cdc767ac2c0"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func weiToEther(wei *big.Int) (ether *big.Float) {
	floatBalance := new(big.Float)
	floatBalance.SetString(wei.String())
	ether = new(big.Float).Quo(floatBalance, big.NewFloat(math.Pow10(18)))
	return
}

func GetBalanceOf(client *ethclient.Client, address common.Address) (balance, pendingBalance *big.Int) {
	defer timeTrack(time.Now(), "GetBalanceOf")
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

func GenerateWallet() (address, publicKey, privateKey string) {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	byesPrivateData := crypto.FromECDSA(key)
	privateKey = hexutil.Encode(byesPrivateData)
	bytesPublicData := crypto.FromECDSAPub(&key.PublicKey)
	publicKey = hexutil.Encode(bytesPublicData)
	address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	return
}

func GenerateEncryptedWallet(password string) (address string) {
	// save the file
	key := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := key.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	address = account.Address.Hex()
	return
}

func DecryptedWallet(password, fileName string) (address, publicKey, privateKey string) {
	bytesFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey(bytesFile, password)
	if err != nil {
		log.Fatal(err)
	}

	bytesPrivateData := crypto.FromECDSA(key.PrivateKey)
	privateKey = hexutil.Encode(bytesPrivateData)
	bytesPublicData := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	publicKey = hexutil.Encode(bytesPublicData)
	address = crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex()
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
	fmt.Println("gasUsed:", gasUsed)

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

	// balance, pendingBalance := GetBalanceOf(client, common.HexToAddress("0x81b7e08f65bdf5648606c89998a9cc8164397647"))
	// fmt.Println("wei:", balance, pendingBalance)
	// fmt.Println("ether:", weiToEther(balance).String(), weiToEther(pendingBalance).String())

	// address, publicKey, privateKey := GenerateWallet()
	// fmt.Println("GeneratedWallet:", address, publicKey, privateKey)

	// // if ok, newAddress := GenerateEncryptedWallet("1234"); ok {
	// // 	fmt.Println(newAddress)
	// // }

	// address, publicKey, privateKey = DecryptedWallet("1234", "./wallet/UTC--2022-03-27T03-36-43.952947000Z--0ef1516c77bd10dcedcd646f5308dc6870fda2e2")
	// fmt.Println("DecryptedWallet:", address, publicKey, privateKey)

	from := "0x606Ef56B359c70A287045D015571C23546260178"
	privateKey := "" // todo
	to := "0x9250788f25aBDc2128a6Eb1B4a1B3393b09bDbb1"
	amount := big.NewInt(int64(math.Pow10(18)))
	receipt := TransferEtherTransaction(client, from, to, amount, privateKey)
	fmt.Println(receipt.TxHash)
}
