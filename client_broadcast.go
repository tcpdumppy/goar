package goar

import (
	"errors"
	"fmt"

	"github.com/everFinance/sandy_log/log"
)

func (c *Client) GetTxDataFromPeers(txId string) ([]byte, error) {
	peers, err := c.GetPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range peers {
		pNode := NewClient("http://" + peer)
		data, err := pNode.GetTransactionData(txId)
		if err != nil {
			log.Error("get tx data failed", "error", err, "peer", peer)
			continue
		}
		return data, nil
	}

	return nil, errors.New("get tx data from peers failed")
}

func (c *Client) BroadcastData(txId string, data []byte, numOfNodes int64) error {
	peers, err := c.GetPeers()
	if err != nil {
		return err
	}

	count := int64(0)
	for _, peer := range peers {

		fmt.Printf("upload peer: %s, count: %d\n", peer, count)
		arNode := NewClient("http://" + peer)
		uploader, err := CreateUploader(arNode, txId, data)
		if err != nil {
			continue
		}

	Loop:
		for !uploader.IsComplete() {
			if err := uploader.UploadChunk(); err != nil {
				break Loop
			}
			if uploader.LastResponseStatus != 200 {
				break Loop
			}
		}
		if uploader.IsComplete() { // upload success
			count++
		}
		if count >= numOfNodes {
			return nil
		}
	}

	return fmt.Errorf("upload tx data to peers failed, txId: %s", txId)
}