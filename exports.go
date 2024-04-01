package ethereum

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"xk6-eth/client"

	"github.com/dop251/goja"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

// GOJA runtime constructor for the Client object (https://github.com/dop251/goja?tab=readme-ov-file#native-constructors)
func (module *Module) NewClient(call goja.ConstructorCall) *goja.Object {
	runTime := module.vu.Runtime()

	privateKey, err := hex.DecodeString(client.DefaultPrivateKey)
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

	client := &client.Client{
		Client:  rpcClient,
		VU:      module.vu,
		Metrics: module.metrics,
		Wallet:  wallet,
	}

	fmt.Println("New client successfully created!")

	return runTime.ToValue(client).ToObject(runTime)
}

// Premine provides the initial accounts funding with the coins
func (module *Module) Premine(number int) *goja.Object {
	runTime := module.vu.Runtime()

	fmt.Println("Mining...")

	for i := range 5 {
		fmt.Println("-mining ", i+1)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done!")

	accounts := []struct{}{}

	return runTime.ToValue(accounts).ToObject(runTime)
}
