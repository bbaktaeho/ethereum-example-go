package ethersutil

import (
	"io/ioutil"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func WeiToEther(wei *big.Int) (ether *big.Float) {
	floatBalance := new(big.Float)
	floatBalance.SetString(wei.String())
	ether = new(big.Float).Quo(floatBalance, big.NewFloat(math.Pow10(18)))
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
