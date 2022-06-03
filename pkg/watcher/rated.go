package watcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
)

// Simplified version of the getValidationEffectiveness response.
//
// https://api.rated.network/docs#/default/get_effectiveness_v0_eth_validators__validator_index__effectiveness_get
type getValidationEffectiveness struct {
	Data []getValidatorEffectivenessData `json:"data"`
}

type getValidatorEffectivenessData struct {
	Uptime                 float64 `json:"uptime"`
	AvgCorrectness         float64 `json:"avgCorrectness"`
	AttesterEffectiveness  float64 `json:"attesterEffectiveness"`
	ProposerEffectiveness  float64 `json:"proposerEffectiveness"`
	ValidatorEffectiveness float64 `json:"validatorEffectiveness"`
}

// getValidationStatistics fetches data about a validator on rated network.
func getValidationStatistics(cfg *core.Config, key string) (*getValidatorEffectivenessData, error) {
	// Key here is already prefixed with '0x', this is required for rated
	// network API (which also supports indexes).
	url := fmt.Sprintf("%s/v0/eth/validators/%s/effectiveness?size=1", cfg.ApiEndpoint, key)

	log.WithFields(log.Fields{
		"url":            url,
		"validation-key": key,
	}).Info("fetching rated data for validation key")

	res, err := http.Get(url)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url":            url,
			"validation-key": key,
		}).Warn("unable to fetch data about validation key from rated network")

		return nil, err
	}
	defer res.Body.Close()

	log.WithFields(log.Fields{
		"url":            url,
		"validation-key": key,
	}).Info("reading response from rated  network")

	if res.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url":         url,
			"status-code": res.StatusCode,
		}).Warn("unable to fetch data about this validation key")

		return nil, fmt.Errorf("non 200-response code received from rated network")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url":            url,
			"validation-key": key,
		}).Warn("unable to read rated network http body")

		return nil, err
	}

	log.WithFields(log.Fields{
		"url":            url,
		"validation-key": key,
	}).Info("parsing response from rated network")

	var response getValidationEffectiveness
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url":            url,
			"validation-key": key,
		}).Warn("unable to parse rated network http body into expected response")

		return nil, err
	}

	if len(response.Data) != 1 {
		log.WithFields(log.Fields{
			"url":            url,
			"status-code":    res.StatusCode,
			"validation-key": key,
			"nb-entries":     len(response.Data),
		}).Warn("expected 1 entry of statistics for the validation key")

		return nil, fmt.Errorf("unexpected data response from rated network")
	}

	log.WithFields(log.Fields{
		"url":            url,
		"validation-key": key,
	}).Info("fetched validator statistics about key from rated network")

	return &response.Data[0], nil
}
