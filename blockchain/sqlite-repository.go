package blockchain

import (
	"context"
	"database/sql"
	"errors"
	"math/big"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

// Migrate() error
func (s *SQLiteRepository) Migrate(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS transactions(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        txid BINARY(32) NOT NULL,
        asset_id BINARY(32) NOT NULL,
		amount BLOB,
        type TEXT NOT NULL
    );
    `
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SQLiteRepository) Create(ctx context.Context, transaction Transaction) (*Transaction, error) {
	res, err := s.db.ExecContext(
		ctx,
		"INSERT INTO transactions(txid, asset_id, amount, type) values(?,?,?,?)",
		transaction.TxID,
		transaction.AssetID,
		transaction.Amount.Bytes(),
		transaction.Type,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	transaction.ID = id

	return &transaction, nil
}

func (s *SQLiteRepository) All(ctx context.Context) ([]Transaction, error) {
	rows, err := s.db.Query("SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Transaction
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.ID, &tx.TxID, &tx.Amount, &tx.AssetID, &tx.Type); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotExists
			}
			return nil, err
		}
		all = append(all, tx)
	}
	return all, nil
}

func (s *SQLiteRepository) GetByTxid(ctx context.Context, txid []byte) (*Transaction, error) {
	// Retrieve the row from the SQL query
	row := s.db.QueryRow("SELECT id, txid, amount, asset_id, type FROM transactions WHERE txid = ?", txid)
	// Declare temporary variables
	var id int64
	var txID []byte
	var amountBytes []byte
	var assetID []byte
	var transactionType string
	// Scan the row into the temporary variables
	err := row.Scan(&id, &txID, &amountBytes, &assetID, &transactionType)
	if err != nil {
		return nil, err
	}
	amount := new(big.Int).SetBytes(amountBytes)
	return &Transaction{
		ID:      id,
		TxID:    txID,
		Amount:  amount,
		AssetID: assetID,
		Type:    transactionType,
	}, nil

	// row := s.db.QueryRow("SELECT * FROM transactions WHERE txid = ?", txid)
	// tx := Transaction{}
	// if err := row.Scan(
	// 	&tx.ID,
	// 	&tx.TxID,
	// 	&tx.Amount,
	// 	&tx.AssetID,
	// 	&tx.Type,
	// ); err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return nil, ErrNotExists
	// 	}
	// 	return nil, err
	// }
	// return &tx, nil
}

func (s *SQLiteRepository) Update(ctx context.Context, id int64, updated Transaction) (*Transaction, error) {
	if id < 1 {
		return nil, errors.New("invalid update ID")
	}
	res, err := s.db.ExecContext(
		ctx,
		"INSERT INTO transactions(txid, asset_id, amount, type) values(?,?,?,?) WHERE id = ?",
		updated.TxID,
		updated.AssetID,
		updated.Amount.Bytes(),
		updated.Type,
		id,
	)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (s *SQLiteRepository) Delete(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return nil
}
