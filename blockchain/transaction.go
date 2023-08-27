package blockchain

import "math/big"

type Transaction struct {
	ID      int64    `json:"id"`
	TxID    []byte   `json:"tx_id"`
	Amount  *big.Int `json:"amount"`
	AssetID []byte   `json:"asset_id"`
	Type    string   `json:"type"`
}
