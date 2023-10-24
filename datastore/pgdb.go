package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

var (
	DbConn         *DB
	pgOnce         sync.Once
	farFutureEpoch = phase0.Epoch(0xffffffffffffffff)
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

//	type Validator struct {
//		Index                      uint64 `json:"index"`
//		Balance                    uint64 `json:"balance"`
//		Status                     string `json:"status"`
//		Pubkey                     string `json:"pubkey"`
//		WithdrawalCredentials      string `json:"withdrawal_credentials"`
//		EffectiveBalance           uint64 `json:"effective_balance"`
//		Slashed                    bool   `json:"slashed"`
//		ActivationEligibilityEpoch uint64 `json:"activation_eligibility_epoch"`
//		ActivationEpoch            uint64 `json:"activation_epoch"`
//		ExitEpoch                  uint64 `json:"exit_epoch"`
//		WithdrawableEpoch          uint64 `json:"withdrawable_epoch"`
//	}

type Validator struct {
	Index     phase0.ValidatorIndex
	Balance   phase0.Gwei
	Status    v1.ValidatorState
	Validator *phase0.Validator
}

func (db *DB) CreateTrackedValidator(ctx context.Context, validator Validator) error {
	publicKey := pgtype.Bytea{}
	publicKey.Set(validator.Validator.PublicKey[:])
	withdrawalCreds := pgtype.Bytea{}
	withdrawalCreds.Set(validator.Validator.WithdrawalCredentials[:])

	var activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch sql.NullInt64

	activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch = processEpochTypes(validator)

	query := `INSERT INTO validators 
		(public_key, index, vStatus, slashed, activation_eligibility_epoch, activation_epoch, exit_epoch, 
			effective_balance, withdrawal_credentials, withdrawable_epoch) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT DO NOTHING RETURNING index;`
	_, err := db.Exec(ctx, query,
		publicKey,
		validator.Index,
		validator.Status,
		validator.Validator.Slashed,
		activationEligibilityEpoch,
		activationEpoch,
		exitEpoch,
		validator.Validator.EffectiveBalance,
		withdrawalCreds,
		withdrawableEpoch,
	)
	if err != nil {
		return err
	}

	query = `INSERT INTO validator_balances (v_index, balance, vStatus, created_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP) ON CONFLICT DO NOTHING RETURNING v_index;`
	_, err = db.Exec(ctx, query, validator.Index, validator.Balance, validator.Status)
	if err != nil {
		return err
	}
	return nil
}

func processEpochTypes(validator Validator) (activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch sql.NullInt64) {

	if validator.Validator.ActivationEligibilityEpoch != farFutureEpoch {
		activationEligibilityEpoch.Valid = true
		activationEligibilityEpoch.Int64 = (int64)(validator.Validator.ActivationEligibilityEpoch)
	}
	if validator.Validator.ActivationEpoch != farFutureEpoch {
		activationEpoch.Valid = true
		activationEpoch.Int64 = (int64)(validator.Validator.ActivationEpoch)
	}
	if validator.Validator.ExitEpoch != farFutureEpoch {
		exitEpoch.Valid = true
		exitEpoch.Int64 = (int64)(validator.Validator.ExitEpoch)
	}
	if validator.Validator.WithdrawableEpoch != farFutureEpoch {
		withdrawableEpoch.Valid = true
		withdrawableEpoch.Int64 = (int64)(validator.Validator.WithdrawableEpoch)
	}
	return
}
