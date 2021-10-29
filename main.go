package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	localUrl = "http://localhost:7545"
	vmUrl    = "http://192.168.66.111:8545"
	address  = "0xaC84E3540BC6866eeF4d2735462a05b090Ea4f10"
)

type SendParams struct {
	From       common.Address
	To         common.Address
	PrivateKey *ecdsa.PrivateKey
	Amount     *big.Int
}

func main() {
	client := GetClientByUrl(vmUrl)
	defer client.Close()

	// feeCollector := common.HexToAddress(address)
	balance := GetBalance(client, common.HexToAddress("0x85F4c4C2e04CcAbbd32c7833D3b364921b0E3663"))
	fmt.Printf("%s's balance: %v wei\n", "0x85F4c4C2e04CcAbbd32c7833D3b364921b0E3663", balance)

	// key := GetKeyInKeystore("./wallet/UTC--2021-10-27T08-59-12.868625000Z--5cbfc05ec8a802e1be308798a61687b432b83ded.json", "1234")
	amount := new(big.Int)
	amount, _ = amount.SetString("1000000000000000000", 10)

	sendParams := SendParams{
		From:       common.HexToAddress("0x85F4c4C2e04CcAbbd32c7833D3b364921b0E3663"),
		To:         common.HexToAddress("0x2D6edcAb374812A03038369eC28472103FEe8F7d"),
		PrivateKey: GetPrivateKey("ca6abc2a50d13b3b4c1212f697cddba93d5370067c58e44ffe382a24c0ea7fff"),
		Amount:     amount,
	}
	SeadEth(client, sendParams)

}

func GetBalance(client *ethclient.Client, address common.Address) *big.Int {
	// wei (big int)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Error to get the balance: %v", err)
		balance = big.NewInt(0)
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

func GetPrivateKey(privateStr string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(privateStr)
	if err != nil {
		log.Fatal(err)
	}
	// publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("error casting public key to ECDSA")
	// }

	return privateKey
}

func GetClientByUrl(networkUrl string) *ethclient.Client {
	client, err := ethclient.DialContext(context.TODO(), networkUrl)
	if err != nil {
		log.Fatalf("Error to create a ether client:%v", err)
	}

	chainId, _ := client.ChainID(context.TODO())
	block, _ := client.BlockNumber(context.TODO())
	networkId, _ := client.NetworkID(context.TODO())
	fmt.Println("chain id:", chainId.String())
	fmt.Println("block number:", block)
	fmt.Println("network id:", networkId.String())

	return client
}

func GenerateAddress() (publicKey, privateKey, address string) {
	pvk, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal()
	}

	pData := crypto.FromECDSA(pvk)
	privateKey = hexutil.Encode(pData)

	puData := crypto.FromECDSAPub(&pvk.PublicKey)
	publicKey = hexutil.Encode(puData)

	address = crypto.PubkeyToAddress(pvk.PublicKey).Hex()
	return
}

func GenerateKeystore() {
	key := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "1234"
	_, err := key.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
}

func GetKeyInKeystore(file, password string) *keystore.Key {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey(b, password)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func GetAddressByPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address
}

func SeadEth(client *ethclient.Client, sendParams SendParams) {
	nonce, err := client.PendingNonceAt(context.TODO(), sendParams.From)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nonce:", nonce)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("gasPrice:", gasPrice, ", gasLimit:", gasLimit)
	var data []byte
	tx := types.NewTransaction(nonce, sendParams.To, sendParams.Amount, gasLimit, gasPrice, data)

	chainId, err := client.NetworkID(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	tx, err = types.SignTx(tx, types.NewEIP155Signer(chainId), sendParams.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())
}
