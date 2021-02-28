package controller

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/clover-network/ethereum-spv-backend/app/merkle"
)

func VerifyTransaction(txID string) (bool, map[string][]byte, error) {
	ctx := context.Background()
	client, err := ethclient.Dial("https://kovan.infura.io/v3/0f818347f004401992c6c63df1c85bf3")
	if err != nil {
		return false, nil, err
	}

	txHash := common.HexToHash(txID)
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return false, nil, err
	}

	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return false, nil, err
	}
	blockHash := receipt.BlockHash
	txCount, err := client.TransactionCount(ctx, blockHash)
	if err != nil {
		return false, nil, err
	}
	fmt.Println("Transcation count: ", txCount)
	fmt.Println(receipt.BlockNumber)

	trie := merkle.NewTrie()
	var txIndex uint

	for i := uint(0); i < txCount; i++ {
		t, err := client.TransactionInBlock(ctx, blockHash, i)
		if err != nil {
			return false, nil, err
		}
		if t.Hash().Hex() == txID {
			txIndex = i
		}

		key, err := rlp.EncodeToBytes(i)
		if err != nil {
			return false, nil, err
		}
		rlp, err := rlp.EncodeToBytes(t)
		if err != nil {
			return false, nil, err
		}
		fmt.Println(i)
		trie.Put(key, rlp)
	}

	key, err := rlp.EncodeToBytes(txIndex)
	if err != nil {
		return false, nil, err
	}

	proof, merklePath, found := trie.Prove(key)
	if !found {
		return false, nil, nil
	}

	block, err := client.BlockByHash(ctx, blockHash)
	if err != nil {
		return false, nil, err
	}
	header := block.Header()
	transactionRoot, err := hex.DecodeString(header.TxHash.Hex()[2:])
	if err != nil {
		return false, nil, err
	}
	txRLP, err := merkle.VerifyProof(transactionRoot, key, proof)
	if err != nil {
		return false, nil, err
	}

	// verify that if the verification passes, it returns the RLP encoded transaction
	rlp, err := merkle.FromEthTransaction(tx).GetRLP()
	if err != nil {
		return false, nil, err
	}

	result := bytes.Compare(txRLP, rlp)
	if result == 0 {
		return true, merklePath, nil
	}
	return false, merklePath, nil
}
