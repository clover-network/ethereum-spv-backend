package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/clover-network/ethereum-spv-backend/app/merkle"
)

func main() {
	ctx := context.Background()
	client, err := ethclient.Dial("https://eth-02.dccn.ankr.com")
	if err != nil {
		log.Fatal(err)
	}
	txHashToVerify := "0x64dae5116265a07299ad260c3897ea9596db69a73bb92a7042e60460df167fe2"
	var txIndex uint

	txHash := common.HexToHash(txHashToVerify)
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		log.Fatal("Unable to get transaction. ", err)
	}

	txJSON, err := tx.MarshalJSON()
	if err != nil {
		log.Fatal()
	}
	fmt.Println("Tx JSON:")
	fmt.Println(string(txJSON))

	fmt.Println(tx.Hash().Hex()) // 0x55fc64d1acae0dc87798657da57bbfd0a835ec65f0a962371ab3eb1f27f24575
	fmt.Println(isPending)       // false

	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		log.Fatal(err)
	}
	blockHash := receipt.BlockHash
	txCount, err := client.TransactionCount(ctx, blockHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transcation count: ", txCount)
	fmt.Println(receipt.BlockNumber)

	trie := merkle.NewTrie()

	for i := uint(0); i < txCount; i++ {
		t, err := client.TransactionInBlock(ctx, blockHash, i)
		if err != nil {
			log.Fatal(err)
		}
		if t.Hash().Hex() == txHashToVerify {
			txIndex = i
		}

		key, err := rlp.EncodeToBytes(i)
		if err != nil {
			log.Fatal(err)
		}
		rlp, err := rlp.EncodeToBytes(t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(i)
		trie.Put(key, rlp)
	}

	block, err := client.BlockByHash(ctx, blockHash)
	if err != nil {
		log.Fatal(err)
	}
	header := block.Header()
	fmt.Println(header.TxHash.Hex()[2:])
	transactionRoot, err := hex.DecodeString(header.TxHash.Hex()[2:])
	if err != nil {
		log.Fatal(err)
	}

	key, err := rlp.EncodeToBytes(txIndex)
	if err != nil {
		log.Fatal(err)
	}
	proof, found := trie.Prove(key)
	if !found {
		log.Fatal("proof not found")
	}
	fmt.Println(proof)

	txRLP, err := merkle.VerifyProof(transactionRoot, key, proof)
	if err != nil {
		log.Fatal(err)
	}

	// verify that if the verification passes, it returns the RLP encoded transaction
	r, err := merkle.FromEthTransaction(tx).GetRLP()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(txRLP)
	fmt.Println(r)

}
