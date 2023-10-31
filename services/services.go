package services

import (
	"context"
	"fmt"

	eth2client "github.com/dbkbali/go-eth2-client"
	"github.com/dbkbali/go-eth2-client/spec/phase0"
)

type Eth2Service struct {
	client *eth2client.Service
}

func NewEth2Service(client *eth2client.Service) *Eth2Service {
	return &Eth2Service{
		client: client,
	}
}

func (service *Eth2Service) UpdateBeaconAttestationRewards(ctx context.Context, epoch phase0.Epoch, validatorIndex phase0.ValidatorIndex) error {
	// get the tracked validator from the database
	// get the attestation rewards not yet stored in the db from the activation date to the current epoch
	// store the attestation rewards in the db for that validator -
	// 1. get the attestation rewards for the validator
	// 2. store the attestation rewards in the db

	beaconRewards, err := (*service.client).(eth2client.BeaconAttestationRewardsProvider).BeaconAttestationRewards(ctx, phase0.Epoch(0), []phase0.ValidatorIndex{validatorIndex})
	if err != nil {
		return err
	}
	fmt.Printf("beacon rewards: %+v\n", beaconRewards)
	return nil
}
