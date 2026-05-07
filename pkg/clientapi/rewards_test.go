package clientapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/migalabs/goteth/pkg/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fetchBlockRewards performs the HTTP call and validation extracted from RequestBlockRewards.
// This allows testing without needing a full APIClient with a real beacon node connection.
func fetchBlockRewards(baseURL string, slot uint64) (spec.BlockRewards, error) {
	uri := baseURL + "/eth/v1/beacon/rewards/blocks/" + fmt.Sprintf("%d", slot)
	resp, err := http.Get(uri)
	if err != nil {
		return spec.BlockRewards{}, fmt.Errorf("block rewards request failed for slot %d: %w", slot, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return spec.BlockRewards{}, fmt.Errorf("block rewards read body failed for slot %d: %w", slot, err)
	}

	if resp.StatusCode != http.StatusOK {
		return spec.BlockRewards{}, fmt.Errorf("block rewards API returned status %d for slot %d: %s", resp.StatusCode, slot, string(body))
	}

	var rewards spec.BlockRewards
	err = json.Unmarshal(body, &rewards)
	if err != nil {
		return spec.BlockRewards{}, fmt.Errorf("block rewards parse failed for slot %d: %w", slot, err)
	}

	if rewards.Data.Total == 0 && rewards.Data.Attestations == 0 && rewards.Data.SyncAggregate == 0 {
		return spec.BlockRewards{}, fmt.Errorf("block rewards response has all zero fields for slot %d, likely invalid", slot)
	}

	return rewards, nil
}

func TestBlockRewards_ValidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/eth/v1/beacon/rewards/blocks/13687200", r.URL.Path)
		resp := spec.BlockRewards{
			ExecutionOptimistic: false,
			Finalized:           true,
			Data: spec.BlockRewardsContent{
				ProposerIndex:     757954,
				Total:             49024174,
				Attestations:      47313857,
				SyncAggregate:     1710317,
				ProposerSlashings: 0,
				AttesterSlashings: 0,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	rewards, err := fetchBlockRewards(server.URL, 13687200)

	require.NoError(t, err)
	assert.Equal(t, uint64(49024174), rewards.Data.Total)
	assert.Equal(t, uint64(47313857), rewards.Data.Attestations)
	assert.Equal(t, uint64(1710317), rewards.Data.SyncAggregate)
}

func TestBlockRewards_HTTP500ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":500,"message":"Internal error"}`))
	}))
	defer server.Close()

	rewards, err := fetchBlockRewards(server.URL, 13687200)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 500")
	assert.Equal(t, uint64(0), rewards.Data.Total)
}

func TestBlockRewards_AllZeroFieldsReturnsError(t *testing.T) {
	// Simulates the exact production bug: Lighthouse returns HTTP 200 with
	// a valid JSON structure but all zero values during state reconstruction.
	// Before the fix, goteth would silently write these zeros to t_block_rewards.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"execution_optimistic":false,"finalized":true,"data":{"proposer_index":"0","total":"0","attestations":"0","sync_aggregate":"0","proposer_slashings":"0","attester_slashings":"0"}}`))
	}))
	defer server.Close()

	rewards, err := fetchBlockRewards(server.URL, 13687200)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "all zero fields")
	assert.Equal(t, uint64(0), rewards.Data.Total)
}

func TestBlockRewards_MalformedJSONReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	_, err := fetchBlockRewards(server.URL, 13687200)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse failed")
}

func TestBlockRewards_HTTP404ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"code":404,"message":"Block not found"}`))
	}))
	defer server.Close()

	_, err := fetchBlockRewards(server.URL, 99999999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")
}

func TestBlockRewards_OnlySlashingsNonZeroIsValid(t *testing.T) {
	// A block with only proposer slashings reward and zero attestations/sync
	// should still be valid (rare but legitimate)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := spec.BlockRewards{
			Data: spec.BlockRewardsContent{
				ProposerIndex:     100,
				Total:             500000,
				Attestations:      0,
				SyncAggregate:     0,
				ProposerSlashings: 500000,
				AttesterSlashings: 0,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	rewards, err := fetchBlockRewards(server.URL, 13687200)

	// Total > 0 so it should pass validation even though attestations and sync are 0
	require.NoError(t, err)
	assert.Equal(t, uint64(500000), rewards.Data.Total)
}
