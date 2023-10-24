
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
    balance BIGINT NOT NULL,
    vStatus VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS i_validators_v_index ON validator_balances(v_index, vstatus);

-- Path: database/migrations/20231023200515_create_validators.down.sql