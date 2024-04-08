package ethereum

import (
	"log"
	"sync"
	"time"
	"xk6-eth/client"

	"github.com/umbracle/ethgo"
	"go.k6.io/k6/metrics"
)

var (
	selected = -1
	once     sync.Once
)

func selection(client *client.Client) {
	once.Do(func() {
		selected = int(client.VU.State().VUID)
	})

	if selected != int(client.VU.State().VUID) {
		return
	}

	go polling(client)
}

func polling(client *client.Client) {
	previousBlockNumber, err := client.Client.Eth().BlockNumber()
	if err != nil {
		log.Fatal("Couldn't get block number")
	}

	previousBlock, err := client.Client.Eth().GetBlockByNumber(ethgo.BlockNumber(previousBlockNumber), false)
	if err != nil {
		log.Fatal("Couldn't get block")
	}

	for {
		blockNumber, err := client.Client.Eth().BlockNumber()
		if err != nil {
			log.Fatal("Couldn't get block number")
		}

		if previousBlockNumber < blockNumber {
			block, err := client.Client.Eth().GetBlockByNumber(ethgo.BlockNumber(blockNumber), false)
			if err != nil {
				log.Fatal("Couldn't get block")
			}

			timestamp := block.Timestamp - previousBlock.Timestamp
			TPS := float64(len(block.TransactionsHashes)) / float64(timestamp)

			metrics.PushIfNotDone(client.VU.Context(), client.VU.State().Samples, metrics.ConnectedSamples{
				Samples: []metrics.Sample{
					{
						TimeSeries: metrics.TimeSeries{
							Metric: client.Metrics.TPS,
						},
						Value: TPS,
						Time:  time.Now(),
					},
					{
						TimeSeries: metrics.TimeSeries{
							Metric: client.Metrics.Block,
						},
						Value: 1,
						Time:  time.Now(),
					},
				},
			})

			previousBlock, previousBlockNumber = block, blockNumber
		}

		time.Sleep(time.Second)
	}
}
