package clientapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/migalabs/goteth/pkg/spec"
)

func (s *APIClient) RequestBlockRewards(slot phase0.Slot) (spec.BlockRewards, error) {
	parsedURL, _ := url.Parse(s.bnEndpoint)
	parsedURL.Path = fmt.Sprintf("/eth/v1/beacon/rewards/blocks/%d", slot)

	req, _ := http.NewRequestWithContext(s.ctx, http.MethodGet, parsedURL.String(), nil)
	if parsedURL.User != nil {
		password, _ := parsedURL.User.Password()
		req.SetBasicAuth(parsedURL.User.Username(), password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var rewards spec.BlockRewards
	err = json.Unmarshal(body, &rewards)
	if err != nil {
		log.Fatalf("error parsing block rewards response: %s", err)
	}

	return rewards, err
}
