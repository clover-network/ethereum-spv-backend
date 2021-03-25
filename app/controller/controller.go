package controller

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func GetDerive(txID string) {
	ctx := context.Background()
	client, err := ethclient.Dial("https://kovan.infura.io/v3/0f818347f004401992c6c63df1c85bf3")
	if err != nil {
		log.Fatal(err)
	}

	txHash := common.HexToHash(txID)

	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Fatal(err)
	}
	blockHash := receipt.BlockHash

	block, err := client.BlockByHash(ctx, blockHash)
	if err != nil {
		log.Fatal(err)
	}
	header := block.Header()
	transactionRoot, err := hex.DecodeString(header.TxHash.Hex()[2:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx Root: %x\n", transactionRoot)
	fmt.Println("tx Root: ", transactionRoot)

	txs := block.Transactions()
	trie := new(trie.Trie)

	h := types.DeriveSha(txs, trie)
	fmt.Println("Trie hash: ", h.Hex())
	fmt.Println("Tx 0: ", txs[0].Hash().Hex())

	key, err := rlp.EncodeToBytes(uint(1))
	if err != nil {
		log.Fatal(err)
	}
	v := trie.Get(key)
	fmt.Printf("value: %x \n", v)

	it := trie.NodeIterator(key)

	for it.Next(true) {
		if it.Leaf() {
			var t *types.Transaction
			err = rlp.DecodeBytes(it.LeafBlob(), &t)
			if err != nil {
				log.Fatal(err)
			}
			if t.Hash().Hex() == txID {
				fmt.Println("Leaf Blob: ", t.Hash().Hex())
				proof := it.LeafProof()
				fmt.Printf("Leaf Proof: %x", proof)
				fmt.Println("Proof length: ", len(proof))

				buf := new(bytes.Buffer)
				err := t.EncodeRLP(buf)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Data: %x \n", buf.Bytes())
				hash := crypto.Keccak256(buf.Bytes())
				fmt.Printf("Keccak: %x \n", hash)

				var val *[][]byte
				err = rlp.DecodeBytes(buf.Bytes(), &val)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Value: 0x%x \n", val)
				fmt.Println("############################### Merkle Path ###############################")

				for _, rawBytes := range proof {
					hash := crypto.Keccak256(rawBytes)
					fmt.Printf("%x: ", hash)

					fmt.Printf(" %x \n", rawBytes)
				}
			}
		}
	}
}
