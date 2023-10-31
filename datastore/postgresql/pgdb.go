package postgresql

import (
	"beacon-actor/datastore"
	"context"
	"database/sql"
	"fmt"
	"sync"

	v1 "github.com/dbkbali/go-eth2-client/api/v1"
	"github.com/dbkbali/go-eth2-client/spec/phase0"
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

func (db *DB) CreateTrackedValidator(ctx context.Context, validator v1.Validator, finalizeEpoch phase0.Epoch) error {
	publicKey := pgtype.Bytea{}
	publicKey.Set(validator.Validator.PublicKey[:])
	withdrawalCreds := pgtype.Bytea{}
	withdrawalCreds.Set(validator.Validator.WithdrawalCredentials[:])

	var activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch sql.NullInt64

	activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch = processEpochTypes(validator)
	fmt.Printf("validator %v\n", validator.Validator.String())
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

	query = `INSERT INTO validator_balances (v_index, f_epoch, balance, effective_balance, vStatus, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) ON CONFLICT DO NOTHING RETURNING v_index;`
	_, err = db.Exec(ctx, query, validator.Index, finalizeEpoch, validator.Balance, validator.Validator.EffectiveBalance, validator.Status)
	if err != nil {
		return err
	}
	return nil
}

func processEpochTypes(validator v1.Validator) (activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch sql.NullInt64) {

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

func (db *DB) GetTrackedValidators(ctx context.Context) ([]datastore.Validator, error) {
	panic("not implemented")
	// query := `SELECT index, balance, vStatus, public_key, withdrawal_credentials, effective_balance, slashed,
	// 	activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch
	// 	FROM validators ORDER BY index ASC;`
	// rows, err := db.Query(ctx, query)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var validators []Validator
	// for rows.Next() {
	// 	var validator Validator
	// 	var publicKey, withdrawalCreds pgtype.Bytea
	// 	var activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch sql.NullInt64

	// 	err := rows.Scan(&validator.Index, &validator.Balance, &validator.Status, &publicKey, &withdrawalCreds,
	// 		&validator.Validator.EffectiveBalance, &validator.Validator.Slashed, &activationEligibilityEpoch,
	// 		&activationEpoch, &exitEpoch, &withdrawableEpoch)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	var pk phase0.BLSPubKey
	// 	copy(pk[:], publicKey.Bytes)
	// 	validator.Validator.PublicKey = pk

	// 	var wc phase0.BLSPubKey
	// 	copy(wc[:], withdrawalCreds.Bytes)
	// 	validator.Validator.WithdrawalCredentials = wc

	// 	validator.Validator.PublicKey = publicKey.Bytes
	// 	validator.Validator.WithdrawalCredentials = withdrawalCreds.Bytes

	// 	validator.Validator.ActivationEligibilityEpoch = phase0.Epoch(activationEligibilityEpoch.Int64)
	// 	validator.Validator.ActivationEpoch = phase0.Epoch(activationEpoch.Int64)
	// 	validator.Validator.ExitEpoch = phase0.Epoch(exitEpoch.Int64)
	// 	validator.Validator.WithdrawableEpoch = phase0.Epoch(withdrawableEpoch.Int64)

	// 	validators = append(validators, validator)
	// }
	// return validators, nil
}

func (db *DB) GetTrackedValidator(ctx context.Context, index uint64) (datastore.Validator, error) {
	panic("not implemented")
	// query := `SELECT index, balance, vStatus, public_key, withdrawal_credentials, effective_balance, slashed,
	// 	activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch
	// 	FROM validators WHERE index = $1;`
	// row := db.QueryRow(ctx, query, index)

	// var validator Validator
	// var publicKey, withdrawalCreds pgtype.Bytea
	// var activationEligibilityEpoch, activationEpoch, exitEpoch, withdrawableEpoch sql.NullInt64

	// err := row.Scan(&validator.Index, &validator.Balance, &validator.Status, &publicKey, &withdrawalCreds,
	// 	&validator.Validator.EffectiveBalance, &validator.Validator.Slashed, &activationEligibilityEpoch,
	// 	&activationEpoch, &exitEpoch, &withdrawableEpoch)
	// if err != nil {
	// 	return validator, err
	// }

	// validator.Validator.PublicKey = publicKey.Bytes
	// validator.Validator.WithdrawalCredentials = withdrawalCreds.Bytes

	// validator.Validator.ActivationEligibilityEpoch = phase0.Epoch(activationEligibilityEpoch.Int64)
	// validator.Validator.ActivationEpoch = phase0.Epoch(activationEpoch.Int64)
	// validator.Validator.ExitEpoch = phase0.Epoch(exitEpoch.Int64)
	// validator.Validator.WithdrawableEpoch = phase0.Epoch(withdrawableEpoch.Int64)

	// return validator, nil
}

func (db *DB) CreateAttestationReward(ctx context.Context, rewards v1.BeaconAttestationRewards) error {
	query := `INSERT INTO attestation_rewards (v_index, pubkey, epoch, head, ctarget, source, inactivity, total_rewards) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT DO NOTHING RETURNING v_index;`
	_, err := db.Exec(ctx, query, rewards.Index, rewards.Pubkey, rewards.Epoch, rewards.Head, rewards.Source, rewards.Target, rewards.Inactivity, rewards.Total)
	if err != nil {
		return err
	}
	return nil
}
