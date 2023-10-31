-- Drop tables
DROP TABLE IF EXISTS validators CASCADE;
DROP TABLE IF EXISTS validator_balances CASCADE;
DROP TABLE IF EXISTS attestation_rewards CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS i_validators_public_key;
DROP INDEX IF EXISTS i_validators_index;
DROP INDEX IF EXISTS i_validators_withdrawal;
DROP INDEX IF EXISTS i_validators_v_index;
DROP INDEX IF EXISTS i_validators_f_epoch;
DROP INDEX IF EXISTS i_attestation_v_index_epoch;
DROP INDEX IF EXISTS i_attestation_pubkey_epoch;

