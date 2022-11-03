package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zsoltggs/golang-example/port-domain-service/pkg/generated/github.com/zsoltggs/golang-example/pkg/pds"
	"io"
)

//go:generate mockgen -package=mockpds -destination=mockpds/mockpds.go github.com/zsoltggs/golang-example/port-domain-service/pkg/generated/github.com/zsoltggs/golang-example/pkg/pds ServiceClient
type Parser struct {
	pdsClient pds.ServiceClient
}

func New(pdsClient pds.ServiceClient) *Parser {
	return &Parser{
		pdsClient: pdsClient,
	}
}

type Port struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

// TODO Solve batched reader
func (p Parser) Parse(ctx context.Context, file io.Reader) error {
	decoder := json.NewDecoder(file)

	// Read the array open bracket
	/*
		_, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("unable to parse opening token: %w", err)
		}
	*/

	data := make(map[string]*pds.Port)
	i := 0
	for decoder.More() {
		err := decoder.Decode(&data)
		if err != nil {
			return fmt.Errorf("unable to decode data : %w", err)
		}
		err = p.processBatchOfData(ctx, data)
		if err != nil {
			return fmt.Errorf("unable to process %d batch: %w", i, err)
		}
		i++
	}

	return nil
}

func (p Parser) processBatchOfData(ctx context.Context, data map[string]*pds.Port) error {
	for id, port := range data {
		port.Id = id
		_, err := p.pdsClient.UpsertPort(ctx, &pds.UpsertPortRequest{Port: port})
		if err != nil {
			return fmt.Errorf("unable to process record: %w", err)
		}
	}
	return nil
}
