package main

import (
	"context"
	"database/sql"
	"db-connector/blockchain"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

const fileName = "sqlite.db"

func main() {
	os.Remove(fileName)

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	transactionRepository := blockchain.NewSQLiteRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := transactionRepository.Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	txid, err := hex.DecodeString("05ed06557eb93800fa2e70143da255857c6a9ca80206bffed3f96a39a8507c7d")
	if err != nil {
		log.Fatal(err)
	}

	assetID, err := hex.DecodeString("dead06557eb93800fa2e70143da255857c6a9ca80206bffed3f96a39a850beef")
	if err != nil {
		log.Fatal(err)
	}

	tx := blockchain.Transaction{
		Type:    "settlement",
		TxID:    txid,
		Amount:  big.NewInt(42),
		AssetID: assetID,
	}
	saved, err := transactionRepository.Create(ctx, tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", saved)

	readTx, err := transactionRepository.GetByTxid(ctx, txid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("retreived: %#v\n", readTx)
}
