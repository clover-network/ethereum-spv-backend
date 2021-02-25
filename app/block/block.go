package app

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SubcribeToHead() {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/0f818347f004401992c6c63df1c85bf3")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println("New block: ")
			fmt.Println("Header Hash: ", header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Block Hash: ", block.Hash().Hex())                    // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
			fmt.Println("Block Number: ", block.Number().Uint64())             // 3477413
			fmt.Println("Block Time: ", block.Time())                          // 1529525947
			fmt.Println("Block Nonce: ", block.Nonce())                        // 130524141876765836
			fmt.Println("Amount of transactions: ", len(block.Transactions())) // 7
		}
	}
}
