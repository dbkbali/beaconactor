
-- These are the validators we are tracking
CREATE TABLE IF NOT EXISTS validators (
    public_key BYTEA NOT NULL,
    index BIGINT NOT NULL PRIMARY KEY,
    vStatus VARCHAR(255) NOT NULL,
    slashed BOOLEAN NOT NULL,
    activation_eligibility_epoch BIGINT,
    activation_epoch BIGINT,
    exit_epoch BIGINT,
    effective_balance BIGINT NOT NULL,
    withdrawal_credentials BYTEA NOT NULL,
    withdrawable_epoch BIGINT
);

CREATE UNIQUE INDEX IF NOT EXISTS i_validators_public_key ON validators (public_key);
CREATE UNIQUE INDEX IF NOT EXISTS i_validators_index ON validators (index);
CREATE INDEX IF NOT EXISTS i_validators_withdrawal ON validators (withdrawal_credentials);

-- Validator balances 
CREATE TABLE IF NOT EXISTS validator_balances (
    v_index BIGINT NOT NULL PRIMARY KEY REFERENCES validators(index),
    f_epoch BIGINT NOT NULL,
    balance BIGINT NOT NULL,
    effective_balance BIGINT NOT NULL,
    vStatus VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS i_validators_v_index ON validator_balances(v_index, vstatus);
CREATE UNIQUE INDEX IF NOT EXISTS i_validators_f_epoch ON validator_balances(v_index, f_epoch);
-- Path: database/migrations/20231023200515_create_validators.down.sql

CREATE TABLE IF NOT EXISTS attestation_rewards (
    v_index BIGINT NOT NULL REFERENCES validators(index),
    pubkey BYTEA NOT NULL,
    epoch BIGINT NOT NULL,
    head BIGINT NOT NULL,
    ctarget BIGINT NOT NULL,
    source BIGINT NOT NULL,
    inactivity BIGINT NOT NULL,
    total_rewards BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS i_attestation_v_index_epoch ON attestation_rewards(v_index, epoch);
CREATE UNIQUE INDEX IF NOT EXISTS i_attestation_pubkey_epoch ON attestation_rewards(pubkey, epoch);

