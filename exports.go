package ethereum

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"xk6-eth/client"

	"github.com/dop251/goja"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

// GOJA runtime constructor for the Client object (https://github.com/dop251/goja?tab=readme-ov-file#native-constructors)
func (module *Module) NewClient(call goja.ConstructorCall) *goja.Object {
	runTime := module.vu.Runtime()

	privateKey, err := hex.DecodeString(call.Arguments[0].String())
	if err != nil {
		log.Fatal("Couldn't decode private key")
	}

	wallet, err := wallet.NewWalletFromPrivKey(privateKey)
	if err != nil {
		log.Fatal("Couldn't create new wallet")
	}

	rpcClient, err := jsonrpc.NewClient("http://localhost:10002")
	if err != nil {
		log.Fatal("Couldn't create new RPC client")
	}

	nonce, err := rpcClient.Eth().GetNonce(wallet.Address(), ethgo.Pending)
	if err != nil {
		log.Fatal("Couldn't get nonce")
	}

	chainId, err := rpcClient.Eth().ChainID()
	if err != nil {
		log.Fatal("Couldn't get chain id")
	}

	client := &client.Client{
		Client:  rpcClient,
		VU:      module.vu,
		Metrics: module.metrics,
		Wallet:  wallet,
		Nonce:   nonce,
		ChainId: chainId,
	}

	fmt.Println("New client (" + wallet.Address().String() + ") successfully created!")

	return runTime.ToValue(client).ToObject(runTime)
}

// Premine provides the initial accounts funding
//
// The number of funded accounts is equal to the number of VUs
func (module *Module) Premine() *goja.Object {
	runTime := module.vu.Runtime()

	var accounts []struct {
		PrivateKey string
		Address    string
	}

	premineTxs := make(map[string]bool)

	rpcClient, err := jsonrpc.NewClient("http://localhost:10002")
	if err != nil {
		log.Fatal("Couldn't create new RPC client")
	}
	defer rpcClient.Close()

	chainId, err := rpcClient.Eth().ChainID()
	if err != nil {
		log.Fatal("Couldn't get chain id")
	}

	gasPrice, err := rpcClient.Eth().GasPrice()
	if err != nil {
		log.Fatal("Couldn't get current gas price")
	}

	nonce, err := rpcClient.Eth().GetNonce(client.DefaultWallet.Address(), ethgo.Pending)
	if err != nil {
		log.Fatal("Couldn't get nonce")
	}

	fmt.Println("Pre-mining...")

	for i := range module.vu.State().Options.VUs.Int64 {
		fmt.Println("- pre-mining", i+1)

		wallet, err := wallet.GenerateKey()
		if err != nil {
			log.Fatal("Couldn't create new wallet")
		}

		privateKey, err := wallet.MarshallPrivateKey()
		if err != nil {
			log.Fatal("Couldn't serialize private key")
		}

		account := struct {
			PrivateKey string
			Address    string
		}{
			hex.EncodeToString(privateKey),
			wallet.Address().String(),
		}

		to := wallet.Address()
		value, _ := big.NewInt(0).SetString("5000000000000000000", 10)

		txHash, err := client.SendRawTransaction(rpcClient, client.DefaultWallet, ethgo.Transaction{
			From:     client.DefaultWallet.Address(),
			To:       &to,
			Value:    value,
			GasPrice: gasPrice,
			Nonce:    nonce,
			ChainID:  chainId,
		})

		if err != nil {
			log.Fatal("Couldn't send pre-mining transaction")
		}

		fmt.Println("Pre-mining transaction (" + txHash + ") was successfully sent")

		premineTxs[txHash] = false

		nonce++

		accounts = append(accounts, account)
	}

	fmt.Println("Waiting for the confirmation of the pre-mining transactions...")

	confirmation(rpcClient, premineTxs)

	fmt.Println("Pre-mining done!")

	return runTime.ToValue(accounts).ToObject(runTime)
}

func confirmation(rc *jsonrpc.Client, premineTxs map[string]bool) {
	for {
		done := true

		for txHash, mined := range premineTxs {
			if mined {
				continue
			}

			receipt, err := rc.Eth().GetTransactionReceipt(ethgo.HexToHash(txHash))
			if err != nil {
				log.Fatal("Couldn't get receipt for the pre-mining transaction")
			}

			if receipt == nil {
				done = false
				continue
			}

			fmt.Println("Transaction (" + txHash + ") confirmed in the block " + strconv.FormatUint(receipt.BlockNumber, 10))

			premineTxs[txHash] = true
		}

		if done {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
}
