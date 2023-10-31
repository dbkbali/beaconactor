package server

import (
	"beacon-actor/datastore/postgresql"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	eth2client "github.com/dbkbali/go-eth2-client"
	v1 "github.com/dbkbali/go-eth2-client/api/v1"
	"github.com/dbkbali/go-eth2-client/spec/phase0"
	"github.com/labstack/echo/v4"
)

type ServerConfig struct {
	Logger        *log.Logger
	ListenAddress string
}

type Server struct {
	ServerConfig
	db           *postgresql.DB
	beaconClient *eth2client.Service
}

type TrackValidatorRequest struct {
	Index uint64 `json:"index" form:"index" query:"index"`
}

func NewServer(config *ServerConfig, dbConn *postgresql.DB, beaconClient *eth2client.Service) *Server {
	return &Server{
		ServerConfig: *config,
		db:           dbConn,
		beaconClient: beaconClient,
	}
}

func (s *Server) Start() error {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	e.GET("/tracked_validators", s.handleTrackedValidators)
	e.GET("/tracked_validators/:index", s.handleTrackedValidator)
	e.POST("/tracked_validators", s.handleCreateTrackedValidator)
	e.Logger.Fatal(e.Start(s.ListenAddress))
	return nil
}

func (s *Server) handleCreateTrackedValidator(c echo.Context) error {
	vIndex := new(TrackValidatorRequest)
	if err := c.Bind(vIndex); err != nil {
		return err
	}

	// Use the index variable of type phase0.ValidatorIndex
	fmt.Printf("validator index: %+v\n", vIndex)
	validatorIndexes := []phase0.ValidatorIndex{phase0.ValidatorIndex(vIndex.Index)}
	stateID := "finalized"

	validatorResponse, err := (*s.beaconClient).(eth2client.ValidatorsProvider).Validators(c.Request().Context(), stateID, validatorIndexes)
	if err != nil {
		return err
	}

	epochResponse, err := (*s.beaconClient).(eth2client.FinalityProvider).Finality(c.Request().Context(), stateID)
	if err != nil {
		return err
	}
	finalizedEpoch := epochResponse.Finalized.Epoch
	fmt.Printf("epoch response: %+v\n", finalizedEpoch)
	validator := validatorResponse[phase0.ValidatorIndex(vIndex.Index)]

	if err = s.db.CreateTrackedValidator(context.Background(), v1.Validator(*validator), finalizedEpoch); err != nil {
		log.Printf("error creating tracked validator: %v", err)
	}

	activationEpoch := validator.Validator.ActivationEpoch
	fmt.Printf("activation epoch: %+v\n", activationEpoch)
	validatorPubkey := validator.Validator.PublicKey
	// iterate from activation epoch to finalized epoch and fetcj attestation rewards
	for i := activationEpoch; i <= finalizedEpoch; i++ {
		fmt.Printf("epoch: %+v\n", i)
		beaconRewards, err := (*s.beaconClient).(eth2client.BeaconAttestationRewardsProvider).BeaconAttestationRewards(c.Request().Context(), i, validatorIndexes)
		if err != nil {
			fmt.Printf("error fetching attestation rewards: %v epoch %d\n", err, i)
		}
		// set pubkey
		beaconRewards[phase0.ValidatorIndex(vIndex.Index)].Pubkey = validatorPubkey
		fmt.Printf("beacon rewards: %+v\n", beaconRewards)
		if err = s.db.CreateAttestationReward(context.Background(), v1.BeaconAttestationRewards(*beaconRewards[phase0.ValidatorIndex(vIndex.Index)])); err != nil {
			log.Printf("error creating attestation reward: %v", err)
		}

	}

	// validator := fetchValidatorResponse.Validator
	return c.JSON(http.StatusCreated, validatorResponse[phase0.ValidatorIndex(vIndex.Index)])
}

func (s *Server) handleTrackedValidators(c echo.Context) error {
	validators, err := s.db.GetTrackedValidators(context.Background())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, validators)
}

func (s *Server) handleTrackedValidator(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.ParseUint(indexStr, 10, 64)
	if err != nil {
		return err
	}
	validator, err := s.db.GetTrackedValidator(context.Background(), index)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, validator)
}
