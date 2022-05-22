package watcher

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
)

// Simplified version of the getStateValidation response, we are only interested
// in the index of the validation here.
//
// https://ethereum.github.io/beacon-APIs/#/Beacon/getStateValidation
type getStateValidation struct {
	Data getStateValidationData `json:"data"`
}

type getStateValidationData struct {
	Index string `json:"index"`
}

// GetValidationIndex fetches the index of a validation on the blockchain.
func getValidationIndex(cfg *core.Config, validationKey string) (int, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/states/head/validations/%s", cfg.BeaconEndpoint, validationKey)

	log.WithFields(log.Fields{
		"url": url,
		"validation-key": validationKey,
	}).Info("fetching index of validation key")

        res, err := http.Get(url)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url": url,
			"validation-key": validationKey,
		}).Warn("unable to fetch data about validation key")
                return 0, err
        }
	defer res.Body.Close()

	log.WithFields(log.Fields{
		"url": url,
		"validation-key": validationKey,
	}).Info("reading response from the beacon")

	if res.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url":    url,
			"status-code": res.StatusCode,
		}).Warn("unable to fetch data about this validation")
		return 0, fmt.Errorf("non 200-response code received from beacon")
	}

        body, err := io.ReadAll(res.Body)
        if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url": url,
			"validation-key": validationKey,
		}).Warn("unable to read http body from beacon")
                return 0, err
        }

	log.WithFields(log.Fields{
		"url": url,
		"validation-key": validationKey,
	}).Info("parsing response from the beacon")

	var response getStateValidation
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url": url,
			"validation-key": validationKey,
		}).Warn("unable to parse http body into expected response")
                return 0, err
	}

	log.WithFields(log.Fields{
		"url": url,
		"validation-key": validationKey,
		"index-str": response.Data.Index,
	}).Info("converting validation index string into a numeric")

	index, err := strconv.Atoi(response.Data.Index)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"url": url,
			"index-str": response.Data.Index,
			"validation-key": validationKey,
		}).Warn("unable to convert index into a numeric value")
		return 0, err
	}

	log.WithFields(log.Fields{
		"url": url,
		"validation-key": validationKey,
		"index-str": response.Data.Index,
		"index-numeric": index,
	}).Info("properly fetch validation index for the given validation key")

	return index, nil
}
