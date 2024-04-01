package client

import (
	"fmt"
	"xk6-eth/testmetrics"

	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
	"go.k6.io/k6/js/modules"
)

const (
	DefaultPrivateKey = "42b6e34dc21598a807dc19d7784c71b2a7a01f6480dc6f58258f78e539f1a1fa"
)

type Client struct {
	Client  *jsonrpc.Client
	VU      modules.VU
	Metrics testmetrics.Metrics
	Wallet  *wallet.Key
}

func (client *Client) Print() {
	fmt.Println(client.VU.State().VUID)
}
