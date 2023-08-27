package blockchain

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		txIDStr    string
		assetIDStr string
		amount     *big.Int
		txType     string
	}{
		{
			txIDStr:    "05ed06557eb93800fa2e70143da255857c6a9ca80206bffed3f96a39a8507c7d",
			assetIDStr: "deadbeef7eb93800fa2e70143da255857c6a9ca80206bffed3f96a39a8507c7d",
			amount:     big.NewInt(4242424242),
			txType:     "settlement",
		},
	}
	fileName := "testSQLite.db"
	os.Remove(fileName)

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		t.Fatal(err)
	}
	transactionRepository := NewSQLiteRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := transactionRepository.Migrate(ctx); err != nil {
		t.Fatal(err)
	}

	for _, tC := range testCases {
		t.Run(tC.txIDStr, func(t *testing.T) {
			txID, err := hex.DecodeString(tC.txIDStr)
			if err != nil {
				t.Fatal(err)
			}
			assetID, err := hex.DecodeString(tC.assetIDStr)
			if err != nil {
				t.Fatal(err)
			}

			tx := Transaction{
				Type:    tC.txType,
				TxID:    txID,
				Amount:  tC.amount,
				AssetID: assetID,
			}
			saved, err := transactionRepository.Create(ctx, tx)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%#v\n", saved)
			read, err := transactionRepository.GetByTxid(ctx, txID)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("retrieved: %#v\n", read)
		})
	}
}
