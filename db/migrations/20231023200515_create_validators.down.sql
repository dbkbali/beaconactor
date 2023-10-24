-- Drop tables
DROP TABLE IF EXISTS validators CASCADE;
DROP TABLE IF EXISTS validator_balances CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS i_validators_public_key;
DROP INDEX IF EXISTS i_validators_index;
DROP INDEX IF EXISTS i_validators_withdrawal;
DROP INDEX IF EXISTS i_validators_v_index;
