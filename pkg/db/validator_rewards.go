package db

import (
	"log"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/jackc/pgx/v5"

	"github.com/migalabs/goteth/pkg/spec"
)

// Postgres intregration variables
var (
	InsertValidator = `
	INSERT INTO t_validator_rewards_summary (	
		f_val_idx, 
		f_epoch, 
		f_balance_eth, 
		f_reward, 
		f_max_reward,
		f_max_att_reward,
		f_max_sync_reward,
		f_att_slot,
		f_base_reward,
		f_in_sync_committee,
		f_missing_source,
		f_missing_target,
		f_missing_head,
		f_status,
		f_block_api_reward,
		f_block_experimental_reward)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	ON CONFLICT ON CONSTRAINT t_validator_rewards_summary_pkey
		DO NOTHING;
	`

	DropValidatorRewardsInEpochQuery = `
		DELETE FROM t_validator_rewards_summary
		WHERE f_epoch = $1;
	`

	DropValidatorRewardsQuery = `
	DELETE FROM t_validator_rewards_summary
	WHERE f_epoch = $1 and f_val_idx = $2;
`
)

func insertValidator(inputValidator spec.ValidatorRewards) (string, []interface{}) {
	return InsertValidator, inputValidator.ToArray()
}

func ValidatorOperation(inputValidator spec.ValidatorRewards) (string, []interface{}) {

	q, args := insertValidator(inputValidator)
	return q, args
}

func (p *PostgresDBService) CopyValRewards(rowSrc [][]interface{}) int64 {

	startTime := time.Now()

	connection, err := p.psqlPool.Acquire(p.ctx)
	if err != nil {
		log.Fatalf("Error while acquiring connection from the database pool for val_rewards: %s!!", err.Error())
	}
	defer connection.Release()

	count, err := connection.CopyFrom(
		p.ctx,
		pgx.Identifier{"t_validator_rewards_summary"},
		[]string{"f_val_idx",
			"f_epoch",
			"f_balance_eth",
			"f_reward",
			"f_max_reward",
			"f_max_att_reward",
			"f_max_sync_reward",
			"f_att_slot",
			"f_base_reward",
			"f_in_sync_committee",
			"f_missing_source",
			"f_missing_target",
			"f_missing_head",
			"f_status",
			"f_block_api_reward",
			"f_block_experimental_reward"},
		pgx.CopyFromRows(rowSrc))

	if err != nil {
		wlog.Warnf("could not copy val_rewards rows into db, they probably already exist in the given epoch: %s", err.Error())
	} else {
		metrics := PersistMetrics{}
		metrics.Rows = uint64(count)
		metrics.PersistTime = time.Since(startTime)
		metrics.RatePersisted = float64(count) / float64(time.Since(startTime).Seconds())
		p.metrics["val_rewards"] = metrics
	}

	wlog.Infof("persisted val_rewards %d rows in %f seconds", count, time.Since(startTime).Seconds())

	return count
}

type ValidatorRewardsDropType phase0.Epoch

func (s ValidatorRewardsDropType) Type() spec.ModelType {
	return spec.ValidatorRewardDropModel
}

func DropValidatorRewards(epoch ValidatorRewardsDropType) (string, []interface{}) {
	resultArgs := make([]interface{}, 0)
	resultArgs = append(resultArgs, epoch)
	return DropValidatorRewardsInEpochQuery, resultArgs
}
