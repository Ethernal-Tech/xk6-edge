package ethereum

import (
	"log"
	"sync"
	"time"
	"xk6-eth/testmetrics"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

var (
	selected = -1
	once     sync.Once
)

func selection(VUID int, rpcClient *jsonrpc.Client, metrics testmetrics.Metrics) {
	once.Do(func() {
		selected = VUID
	})

	if selected != VUID {
		return
	}

	go polling(rpcClient, metrics)
}

func polling(rpcClient *jsonrpc.Client, metrics testmetrics.Metrics) {
	previousBlockNumber, err := rpcClient.Eth().BlockNumber()
	if err != nil {
		log.Fatal("Couldn't get block number")
	}

	previousBlock, err := rpcClient.Eth().GetBlockByNumber(ethgo.BlockNumber(previousBlockNumber), false)
	if err != nil {
		log.Fatal("Couldn't get block")
	}

	for {
		blockNumber, err := rpcClient.Eth().BlockNumber()
		if err != nil {
			log.Fatal("Couldn't get block number")
		}

		if previousBlockNumber < blockNumber {
			block, err := rpcClient.Eth().GetBlockByNumber(ethgo.BlockNumber(blockNumber), false)
			if err != nil {
				log.Fatal("Couldn't get block")
			}

			timestamp := block.Timestamp - previousBlock.Timestamp
			TPS := float64(len(block.TransactionsHashes)) / float64(timestamp)

			_ = TPS

			previousBlock, previousBlockNumber = block, blockNumber
		}

		time.Sleep(time.Second)
	}
}
