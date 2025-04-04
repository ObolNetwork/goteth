package spec

const (
	MainnetGenesis               = 1606824023
	SepoliaGenesis               = 1655733600
	HoleskyGenesis               = 1695902400
	MainnetBeaconContractAddress = "0x00000000219ab540356cBB839Cbe05303d7705Fa"
	SepoliaBeaconContractAddress = "0x7f02C3E3c98b133055B8B348B2Ac625669Ed295D"
	HoleskyBeaconContractAddress = "0x4242424242424242424242424242424242424242"
	HoodiBeaconContractAddress   = MainnetBeaconContractAddress
	DepositEventTopic            = "0x649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c5"
	DepositEventDataLength       = 576
)

var BeaconContractAddresses = map[string]string{
	"mainnet": MainnetBeaconContractAddress,
	"sepolia": SepoliaBeaconContractAddress,
	"holesky": HoleskyBeaconContractAddress,
	"hoodi":   HoodiBeaconContractAddress,
}

/*
Phase0
*/

const (
	MaxEffectiveInc             = 32
	BaseRewardFactor            = 64
	BaseRewardPerEpoch          = 4
	EffectiveBalanceInc         = 1000000000
	SlotsPerEpoch               = 32
	ProposerRewardQuotient      = 8
	SlotsPerHistoricalRoot      = 8192
	SlotSeconds                 = 12
	EpochSlots                  = 32
	WhistleBlowerRewardQuotient = 512
	MinInclusionDelay           = 1

	AttSourceFlagIndex = 0
	AttTargetFlagIndex = 1
	AttHeadFlagIndex   = 2
)

/*
Altair
*/
const (
	// spec weight constants
	TimelySourceWeight = 14
	TimelyTargetWeight = 26
	TimelyHeadWeight   = 14

	SyncRewardWeight  = 2
	ProposerWeight    = 8
	WeightDenominator = 64
	SyncCommitteeSize = 512
)

var (
	ParticipatingFlagsWeight = [3]int{TimelySourceWeight, TimelyTargetWeight, TimelyHeadWeight}
)

type ModelType int8

const (
	BlockModel ModelType = iota
	BlockDropModel
	OrphanModel
	EpochModel
	EpochDropModel
	PoolSummaryModel
	ProposerDutyModel
	ProposerDutyDropModel
	ValidatorLastStatusModel
	ValidatorRewardsModel
	ValidatorRewardDropModel
	WithdrawalModel
	WithdrawalDropModel
	TransactionsModel
	TransactionDropModel
	ReorgModel
	FinalizedCheckpointModel
	HeadEventModel
	ValidatorRewardsAggregationModel
	SlashingModel
	BLSToExecutionChangeModel
	DepositModel
	ETH1DepositModel
)

type ValidatorStatus int8

const (
	QUEUE_STATUS ValidatorStatus = iota
	ACTIVE_STATUS
	EXIT_STATUS
	SLASHED_STATUS
	NUMBER_OF_STATUS // Add new status before this
)
