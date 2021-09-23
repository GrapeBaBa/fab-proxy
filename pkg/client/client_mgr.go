package client

import (
	"github.com/Grapebaba/fab-proxy/pkg/config"
	"github.com/hyperledger/fabric/gossip/util"
)

type Client struct {
	//endorserClients  []peer.EndorserClient
	broadcastClients     []*BroadcastClientWrapper
	deliverClient        *DeliverClientWrapper
	initBroadcastClient  *BroadcastClientWrapper
	readBroadcastClients []*BroadcastClientWrapper
}

func New(concurrency int, nodes []config.Node) *Client {
	bClients := make([]*BroadcastClientWrapper, 0)
	rBClients := make([]*BroadcastClientWrapper, 0)
	for _, node := range nodes {
		for i := 0; i < concurrency; i++ {
			bc, _ := CreateBroadcastClient(node.Addr, node.RootCerts[0], node.OverrideHostname)
			bClients = append(bClients, bc)
		}

	}

	for i := 0; i < concurrency; i++ {
		rbc, _ := CreateBroadcastClient(nodes[0].Addr, nodes[0].RootCerts[0], nodes[0].OverrideHostname)
		rBClients = append(rBClients, rbc)
	}

	dClient, _ := CreateDeliverClient(nodes[0].Addr, nodes[0].RootCerts[0], nodes[0].OverrideHostname)
	iBClient, _ := CreateBroadcastClient(nodes[0].Addr, nodes[0].RootCerts[0], nodes[0].OverrideHostname)
	mgr := &Client{
		broadcastClients:     bClients,
		deliverClient:        dClient,
		initBroadcastClient:  iBClient,
		readBroadcastClients: rBClients,
	}
	return mgr
}

func (cm *Client) GetBroadcastClient() *BroadcastClientWrapper {
	return cm.broadcastClients[util.RandomInt(len(cm.broadcastClients))]
}

func (cm *Client) GetDeliverClient() *DeliverClientWrapper {
	return cm.deliverClient
}

func (cm *Client) GetBroadcastClients() []*BroadcastClientWrapper {
	return cm.broadcastClients
}

func (cm *Client) GetInitBroadcastClient() *BroadcastClientWrapper {
	return cm.initBroadcastClient
}

func (cm *Client) GetRBroadcastClient() *BroadcastClientWrapper {
	return cm.readBroadcastClients[util.RandomInt(len(cm.readBroadcastClients))]
}

func (cm *Client) GetRBroadcastClients() []*BroadcastClientWrapper {
	return cm.readBroadcastClients
}
