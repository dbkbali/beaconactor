package server

import (
	"beacon-actor/datastore"
	"context"
	"fmt"
	"log"
	"net/http"

	eth2client "github.com/attestantio/go-eth2-client"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/labstack/echo/v4"
)

type ServerConfig struct {
	Logger        *log.Logger
	ListenAddress string
}

type Server struct {
	ServerConfig
	db           *datastore.DB
	beaconClient *eth2client.Service
}

type TrackValidatorRequest struct {
	Index uint64 `json:"index" form:"index" query:"index"`
}

func NewServer(config *ServerConfig, dbConn *datastore.DB, beaconClient *eth2client.Service) *Server {
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
	// e.GET("/tracked_validators", s.handleTrackedValidators)
	// e.GET("/tracked_validators/:address", s.handleTrackedValidator)
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
	stateID := "head"

	validatorResponse, err := (*s.beaconClient).(eth2client.ValidatorsProvider).Validators(c.Request().Context(), stateID, validatorIndexes)
	if err != nil {
		return err
	}

	validator := validatorResponse[phase0.ValidatorIndex(vIndex.Index)]

	if err = s.db.CreateTrackedValidator(context.Background(), datastore.Validator(*validator)); err != nil {
		log.Printf("error creating tracked validator: %v", err)
	}

	// validator := fetchValidatorResponse.Validator
	return c.JSON(http.StatusCreated, validatorResponse[phase0.ValidatorIndex(vIndex.Index)])
}
