package datastore

import (
	"context"
	"fmt"
	"sync"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

var (
	DbConn *DB
	pgOnce sync.Once
)

func NewDb(ctx context.Context, connString string) (*DB, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			fmt.Println("unable to create connection pool: %w", err)
		}
		fmt.Println("connected to database")

		DbConn = &DB{db}
	})

	return DbConn, nil
}

func (db *DB) DbExists(ctx context.Context, dbName string) (bool, error) {
	query := `select exists(SELECT * FROM pg_catalog.pg_database WHERE datname = $1);`
	result, err := db.Query(ctx, query, dbName)
	if err != nil {
		return false, err
	}
	return result.Next(), nil
}

// type Validator struct {
// 	Index                      uint64 `json:"index"`
// 	Balance                    uint64 `json:"balance"`
// 	Status                     string `json:"status"`
// 	Pubkey                     string `json:"pubkey"`
// 	WithdrawalCredentials      string `json:"withdrawal_credentials"`
// 	EffectiveBalance           uint64 `json:"effective_balance"`
// 	Slashed                    bool   `json:"slashed"`
// 	ActivationEligibilityEpoch uint64 `json:"activation_eligibility_epoch"`
// 	ActivationEpoch            uint64 `json:"activation_epoch"`
// 	ExitEpoch                  uint64 `json:"exit_epoch"`
// 	WithdrawableEpoch          uint64 `json:"withdrawable_epoch"`
// }

type Validator struct {
	Index     phase0.ValidatorIndex
	Balance   phase0.Gwei
	Status    v1.ValidatorState
	Validator *phase0.Validator
}

func (db *DB) CreateValidatorTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS validators (
		index bigint PRIMARY KEY,
		balance bigint,
		status text,
		pubkey text,
		withdrawal_credentials text,
		effective_balance bigint,
		slashed boolean,
		activation_eligibility_epoch bigint,
		activation_epoch bigint,
		exit_epoch bigint,
		withdrawable_epoch bigint
	);`
	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateTrackedValidator(ctx context.Context, validator Validator) error {
	query := `INSERT INTO validators (index, balance, status, pubkey, withdrawal_credentials, effective_balance, slashed, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
	_, err := db.Exec(ctx, query, validator.Index, validator.Balance, validator.Status, validator.Validator)
	if err != nil {
		return err
	}
	return nil
}
