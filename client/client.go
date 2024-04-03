package client

import (
	"encoding/hex"
	"log"
	"math/big"
	"xk6-eth/testmetrics"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
	"go.k6.io/k6/js/modules"
)

const (
	DefaultPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
	DefaultAddress    = "0x85da99c8a7c2c95964c8efd687e95e632fc533d6"
)

var (
	DefaultWallet *wallet.Key
)

func init() {
	privateKey, err := hex.DecodeString(DefaultPrivateKey)
	if err != nil {
		log.Fatal("Couldn't decode default private key")
	}

	DefaultWallet, err = wallet.NewWalletFromPrivKey(privateKey)
	if err != nil {
		log.Fatal("Couldn't create default wallet")
	}
}

type Transaction struct {
	From      string
	To        string
	Input     []byte
	GasPrice  uint64
	GasFeeCap uint64
	GasTipCap uint64
	Gas       uint64
	Value     int64
	Nonce     uint64
	ChainId   int64
}

type Client struct {
	Client  *jsonrpc.Client
	VU      modules.VU
	Metrics testmetrics.Metrics
	Wallet  *wallet.Key
	Nonce   uint64
	ChainId *big.Int
}

func SendRawTransaction(rc *jsonrpc.Client, w *wallet.Key, tx ethgo.Transaction) (string, error) {
	gas, err := EstimateGas(rc, w, tx)
	if err != nil {
		return "", err
	}

	tx.Type = ethgo.TransactionLegacy
	tx.Gas = gas

	signer := wallet.NewEIP155Signer(tx.ChainID.Uint64())
	signedTx, err := signer.SignTx(&tx, w)
	if err != nil {
		return "", err
	}

	txRLP, err := signedTx.MarshalRLPTo(nil)
	if err != nil {
		return "", err
	}

	txHash, err := rc.Eth().SendRawTransaction(txRLP)

	return txHash.String(), err
}

func (client *Client) SendRawTransaction(tx Transaction) (string, error) {
	return "", nil
}

func EstimateGas(rc *jsonrpc.Client, w *wallet.Key, tx ethgo.Transaction) (uint64, error) {
	msg := &ethgo.CallMsg{
		From:     tx.From,
		To:       tx.To,
		Value:    tx.Value,
		Data:     tx.Input,
		GasPrice: tx.GasPrice,
	}

	gas, err := rc.Eth().EstimateGas(msg)
	if err != nil {
		return 0, err
	}

	return gas, nil
}

func (client *Client) EstimateGas(tx Transaction) (uint64, error) {
	return 0, nil
}
