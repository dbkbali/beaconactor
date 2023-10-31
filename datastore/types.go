package datastore

import (
	"github.com/dbkbali/go-eth2-client/spec/phase0"
)

type Validator struct {
	PublicKey                  phase0.BLSPubKey
	Index                      phase0.ValidatorIndex
	EffectiveBalance           phase0.Gwei
	Slashed                    bool
	ActivationEligibilityEpoch phase0.Epoch
	ActivationEpoch            phase0.Epoch
	ExitEpoch                  phase0.Epoch
	WithdrawableEpoch          phase0.Epoch
	WithdrawalCredentials      [32]byte
}
