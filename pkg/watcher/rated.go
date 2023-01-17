package watcher

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"

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

type authenticationError struct {
	Detail    string `json:"detail"`
}

func getValidationStatistics(cfg *core.Config, key string) (*getValidatorEffectivenessData, error) {

	url := fmt.Sprintf("%s/v0/eth/validators/%s/effectiveness?size=1", cfg.ApiEndpoint, key)

	log.WithFields(log.Fields{
		"url":            url,
		"validation-key": key,
	}).Info("Fetching rated data for validation key")

	maxRetries := 10

	client := new(http.Client)
	bearerToken := fmt.Sprintf("Bearer %s", cfg.ApiAccessToken)

	for r := 0; r <= maxRetries; r++ {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"url":            url,
				"validation-key": key,
			}).Warn("Unable to create http request. Aborting")

			// Abort on unknown error
			return nil, err
		}

		req.Header.Add("X-Rated-Network", cfg.Network)
		req.Header.Add("Authorization", bearerToken)
		res, err := client.Do(req)

		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"url":            url,
				"network":        cfg.Network,
				"validation-key": key,
			}).Warn("Unable to fetch validator data. Aborting")

			// Abort on unknown error
			return nil, err
		}

		defer res.Body.Close()

		log.WithFields(log.Fields{
			"url":            url,
			"network":        cfg.Network,
			"validation-key": key,
		}).Info("Reading response from rated network")

		if res.StatusCode == 500 {
			log.WithFields(log.Fields{
				"url":           url,
				"network":       cfg.Network,
				"status-code":   res.StatusCode,
				"validator-key": key,
			}).Error("Unable to fetch validator data due to internal server error. Aborting")

			return nil, fmt.Errorf("Internal server error received from rated network")

		} else if res.StatusCode == 404 {
			log.WithFields(log.Fields{
				"url":            url,
				"network":        cfg.Network,
				"status-code":    res.StatusCode,
				"validation-key": key,
			}).Warn("Validator not found")

			return nil, fmt.Errorf("Validator not found")

		} else if res.StatusCode == 401 {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"url":            url,
					"network":        cfg.Network,
					"validation-key": key,
				}).Warn("Unable to read rated network http body")

				return nil, err
			}

			var response authenticationError
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"url":            url,
					"network":        cfg.Network,
					"validation-key": key,
				}).Warn("Unable to read rated network http body")

				return nil, err
			}

			return nil, fmt.Errorf("Authentication failure: %q.", response.Detail)

		} else if res.StatusCode == 200 {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"url":            url,
					"network":        cfg.Network,
					"validation-key": key,
				}).Warn("Unable to read rated network http body")

				return nil, err
			}

			log.WithFields(log.Fields{
				"url":            url,
				"network":        cfg.Network,
				"validation-key": key,
			}).Info("Parsing response from rated network")

			var response getValidationEffectiveness
			err = json.Unmarshal(body, &response)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"url":            url,
					"network":        cfg.Network,
					"validation-key": key,
				}).Warn("Unable to parse rated network http body into expected response")

				return nil, err
			}

			if len(response.Data) != 1 {
				log.WithFields(log.Fields{
					"url":            url,
					"network":        cfg.Network,
					"status-code":    res.StatusCode,
					"validation-key": key,
					"nb-entries":     len(response.Data),
				}).Warn("Expected 1 entry of statistics for the validation key")

				return nil, fmt.Errorf("unexpected data response from rated network")
			} else {
				log.WithFields(log.Fields{
					"url":            url,
					"network":        cfg.Network,
					"validation-key": key,
				}).Info("Successfully fetched statistics for validator")

				return &response.Data[0], nil
			}

		} else if res.StatusCode == 429 {
			min := 0.95
			max := 1.05
			sleepFor := math.Pow(2, float64(r)) * ((rand.Float64() * (max - min)) + min)

			log.WithFields(log.Fields{
				"url":         url,
				"network":     cfg.Network,
				"status-code": res.StatusCode,
			}).Warn(fmt.Sprintf("Rate limit exceeded. Retrying in %v seconds.", sleepFor))

			time.Sleep(time.Duration(sleepFor) * time.Second)

		} else {

			log.WithFields(log.Fields{
				"url":            url,
				"network":        cfg.Network,
				"status-code":    res.StatusCode,
				"validation-key": key,
			}).Warn("Unknown status code received. Aborting")

			return nil, fmt.Errorf("Unknown status code received from rated network")
		}
	}

	return nil, fmt.Errorf(fmt.Sprintf("Failed to fetch Validator data after %v attempts", maxRetries))

}
