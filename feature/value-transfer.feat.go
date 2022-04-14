package feature

import (
	"context"
	"ethereum-client-example-go/util/ethers"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func createLegacyTx(nonce, gasLimit uint64, gasPrice, value *big.Int, to *common.Address, data []byte) types.TxData {
	return &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       to,
		Value:    value,
		Data:     data,
	}
}

func createEIP1559Tx() {}

func TransferValue(client *ethclient.Client, privateHexKey string, to *common.Address, value *big.Int, txType int64) (*types.Receipt, error) {
	var txData types.TxData

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA(privateHexKey)
	if err != nil {
		return nil, err
	}

	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return nil, err
	}

	gasLimit, err := ethers.EstimateGas(client, value, nil)
	if err != nil {
		return nil, err
	}

	if txType == types.DynamicFeeTxType {
		txData = nil
	} else {
		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		txData = createLegacyTx(nonce, gasLimit, gasPrice, value, to, nil)
	}

	tx := types.NewTx(txData)
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKey)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return nil, err
	}
	return receipt, err
}
