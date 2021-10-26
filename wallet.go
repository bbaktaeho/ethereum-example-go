package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

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

func ReadKeystore() {
	password := "1234"
	b, err := ioutil.ReadFile("./wallet/UTC--2021-10-26T08-01-02.752907000Z--96435b7ac04116a667e7b6761436e429ae7dc885")
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey(b, password)
	if err != nil {
		log.Fatal(err)
	}
	pData := crypto.FromECDSA(key.PrivateKey)
	fmt.Println(hexutil.Encode(pData))
	puData := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	fmt.Println(hexutil.Encode(puData))
	address := crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex()
	fmt.Println(address)
}
