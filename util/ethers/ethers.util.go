package ethers

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

func GetTxFromAddressWithoutCall(client *ethclient.Client, tx *types.Transaction) (common.Address, error) {
	ctx := context.Background()

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return common.Address{}, err
	}

	signer := types.LatestSignerForChainID(chainId)
	from, err := types.Sender(signer, tx)
	return from, nil
}

func GetBalanceOf(client *ethclient.Client, address common.Address) (balance, pendingBalance *big.Int, e error) {
	var w sync.WaitGroup
	ctx := context.Background()

	w.Add(2)
	go func() {
		b, err := client.BalanceAt(ctx, address, nil)
		balance = b
		if err != nil {
			e = err
		}
		w.Done()
	}()
	go func() {
		pb, err := client.PendingBalanceAt(ctx, address)
		pendingBalance = pb
		if err != nil {
			e = err
		}
		w.Done()
	}()
	w.Wait()
	return
}

func WeiToEther(wei *big.Int) (ether *big.Float) {
	floatBalance := new(big.Float)
	floatBalance.SetString(wei.String())
	ether = new(big.Float).Quo(floatBalance, big.NewFloat(math.Pow10(18)))
	return
}

func DecryptedWallet(password, fileName string) (address, publicKey, privateKey string, e error) {
	bytesFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", "", "", err
	}
	key, err := keystore.DecryptKey(bytesFile, password)
	if err != nil {
		return "", "", "", err
	}

	bytesPrivateData := crypto.FromECDSA(key.PrivateKey)
	privateKey = hexutil.Encode(bytesPrivateData)
	bytesPublicData := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	publicKey = hexutil.Encode(bytesPublicData)
	address = crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex()
	return
}

func GenerateEncryptedWallet(password string) (address string, e error) {
	// save the file
	key := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := key.NewAccount(password)
	if err != nil {
		return "", err
	}
	address = account.Address.Hex()
	return
}

func GenerateWallet() (address, publicKey, privateKey string, e error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", err
	}

	byesPrivateData := crypto.FromECDSA(key)
	privateKey = hexutil.Encode(byesPrivateData)
	bytesPublicData := crypto.FromECDSAPub(&key.PublicKey)
	publicKey = hexutil.Encode(bytesPublicData)
	address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	return
}

func GetBaseFee(client *ethclient.Client, chainConfig *params.ChainConfig) (baseFee *big.Int, e error) {
	ctx := context.Background()

	blockNumber, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	block, err := client.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNumber))
	if err != nil {
		return nil, err
	}

	fmt.Println(block.BaseFee())
	baseFee = misc.CalcBaseFee(chainConfig, block.Header())
	fmt.Println(baseFee)
	return
}

func EstimateGas(client *ethclient.Client, value *big.Int, data []byte) (uint64, error) {
	ctx := context.Background()

	callMsg := ethereum.CallMsg{
		From:     common.Address{},
		To:       new(common.Address),
		Gas:      0,
		GasPrice: big.NewInt(0),
		Value:    value,
		Data:     data,
	}

	estimateGas, err := client.EstimateGas(ctx, callMsg)
	if err != nil {
		return 0, err
	}

	return estimateGas, nil
}
