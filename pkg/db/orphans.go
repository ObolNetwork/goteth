package db

import (
	"strings"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/migalabs/goteth/pkg/spec"
	"github.com/migalabs/goteth/pkg/utils"
)

var (
	orphansTable      = "t_orphans"
	insertOrphanQuery = `
	INSERT INTO %s (
		f_timestamp,
		f_epoch, 
		f_slot,
		f_graffiti,
		f_proposer_index,
		f_proposed,
		f_attestations,
		f_deposits,
		f_proposer_slashings,
		f_attester_slashings,
		f_voluntary_exits,
		f_sync_bits,
		f_el_fee_recp,
		f_el_gas_limit,
		f_el_gas_used,
		f_el_base_fee_per_gas,
		f_el_block_hash,
		f_el_transactions,
		f_el_block_number,
		f_ssz_size_bytes,
		f_snappy_size_bytes,
		f_compression_time_ms,
		f_decompression_time_ms,
		f_payload_size_bytes)
		VALUES`
)

func orphansInput(blocks []spec.AgnosticBlock) proto.Input {
	// one object per column
	var (
		f_timestamp             proto.ColUInt64
		f_epoch                 proto.ColUInt64
		f_slot                  proto.ColUInt64
		f_graffiti              proto.ColStr
		f_proposer_index        proto.ColUInt64
		f_proposed              proto.ColBool
		f_attestations          proto.ColUInt64
		f_deposits              proto.ColUInt64
		f_proposer_slashings    proto.ColUInt64
		f_attester_slashings    proto.ColUInt64
		f_voluntary_exits       proto.ColUInt64
		f_sync_bits             proto.ColUInt64
		f_el_fee_recp           proto.ColStr
		f_el_gas_limit          proto.ColUInt64
		f_el_gas_used           proto.ColUInt64
		f_el_base_fee_per_gas   proto.ColUInt64
		f_el_block_hash         proto.ColStr
		f_el_transactions       proto.ColUInt64
		f_el_block_number       proto.ColUInt64
		f_payload_size_bytes    proto.ColUInt64
		f_ssz_size_bytes        proto.ColFloat32
		f_snappy_size_bytes     proto.ColFloat32
		f_compression_time_ms   proto.ColFloat32
		f_decompression_time_ms proto.ColFloat32
	)

	for _, block := range blocks {
		f_timestamp.Append(uint64(block.ExecutionPayload.Timestamp))
		f_epoch.Append(uint64(block.Slot / spec.SlotsPerEpoch))
		f_slot.Append(uint64(block.Slot))

		graffiti := strings.ToValidUTF8(string(block.Graffiti[:]), "?")
		graffiti = strings.ReplaceAll(graffiti, "\u0000", "")
		f_graffiti.Append(graffiti)

		f_proposer_index.Append(uint64(block.ProposerIndex))
		f_proposed.Append(block.Proposed)
		if block.Attestations != nil {
			f_attestations.Append(uint64(len(block.Attestations)))
		} else if block.ElectraAttestations != nil {
			f_attestations.Append(uint64(len(block.ElectraAttestations)))
		}
		f_deposits.Append(uint64(len(block.Deposits)))
		f_proposer_slashings.Append(uint64(len(block.ProposerSlashings)))
		f_attester_slashings.Append(uint64(len(block.AttesterSlashings)))
		f_voluntary_exits.Append(uint64(len(block.VoluntaryExits)))
		f_sync_bits.Append(uint64(block.SyncAggregate.SyncCommitteeBits.Count()))

		// Execution Payload
		f_el_fee_recp.Append(block.ExecutionPayload.FeeRecipient.String())
		f_el_gas_limit.Append(uint64(block.ExecutionPayload.GasLimit))
		f_el_gas_used.Append(uint64(block.ExecutionPayload.GasUsed))
		f_el_base_fee_per_gas.Append(uint64(block.ExecutionPayload.BaseFeePerGas))
		f_el_block_hash.Append(block.ExecutionPayload.BlockHash.String())
		f_el_transactions.Append(uint64(len(block.ExecutionPayload.Transactions)))
		f_el_block_number.Append(uint64(block.ExecutionPayload.BlockNumber))

		// Size stats
		f_payload_size_bytes.Append(uint64(block.ExecutionPayload.PayloadSize))
		f_ssz_size_bytes.Append(float32(block.SSZsize))
		f_snappy_size_bytes.Append(float32(block.SnappySize))
		f_compression_time_ms.Append(float32(utils.DurationToFloat64Millis(block.CompressionTime)))
		f_decompression_time_ms.Append(float32(utils.DurationToFloat64Millis(block.DecompressionTime)))

	}

	return proto.Input{
		{Name: "f_timestamp", Data: f_timestamp},
		{Name: "f_epoch", Data: f_epoch},
		{Name: "f_slot", Data: f_slot},
		{Name: "f_graffiti", Data: f_graffiti},
		{Name: "f_proposer_index", Data: f_proposer_index},
		{Name: "f_proposed", Data: f_proposed},
		{Name: "f_attestations", Data: f_attestations},
		{Name: "f_deposits", Data: f_deposits},
		{Name: "f_proposer_slashings", Data: f_proposer_slashings},
		{Name: "f_attester_slashings", Data: f_attester_slashings},
		{Name: "f_voluntary_exits", Data: f_voluntary_exits},
		{Name: "f_sync_bits", Data: f_sync_bits},
		{Name: "f_el_fee_recp", Data: f_el_fee_recp},
		{Name: "f_el_gas_limit", Data: f_el_gas_limit},
		{Name: "f_el_gas_used", Data: f_el_gas_used},
		{Name: "f_el_base_fee_per_gas", Data: f_el_base_fee_per_gas},
		{Name: "f_el_block_hash", Data: f_el_block_hash},
		{Name: "f_el_transactions", Data: f_el_transactions},
		{Name: "f_el_block_number", Data: f_el_block_number},
		{Name: "f_ssz_size_bytes", Data: f_ssz_size_bytes},
		{Name: "f_snappy_size_bytes", Data: f_snappy_size_bytes},
		{Name: "f_compression_time_ms", Data: f_compression_time_ms},
		{Name: "f_decompression_time_ms", Data: f_decompression_time_ms},
		{Name: "f_payload_size_bytes", Data: f_payload_size_bytes},
	}
}

func (p *DBService) PersistOrphans(data []spec.AgnosticBlock) error {
	persistObj := PersistableObject[spec.AgnosticBlock]{
		input: orphansInput,
		table: orphansTable,
		query: insertOrphanQuery,
	}

	for _, item := range data {
		persistObj.Append(item)
	}

	err := p.Persist(persistObj.ExportPersist())
	if err != nil {
		log.Errorf("error persisting orphans: %s", err.Error())
	}
	return err
}
