package types

const (
	// ModuleName defines the module name
	ModuleName = "nft"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_nft"

	// ModuleAdminKey defines module admin key
	ParamModuleAdminKey          = "moduleAdmin"
	ParamBondDenomKey            = "bondDenom"
	ParamPrefixKey               = "prefix"
	ParamNftSequenceKey          = "nftSequence"
	ParamsNftTotalVestingTimeKey = "totalVestingTime"
	ParamsNftVestingPeriodKey    = "nftVestingPeriod"
	ParamVestingPeriodsLimitKey  = "vestingPeriodsLimitKey"
	ParamNftTokenAmountKey       = "nftTokenAmount"
	ParamBatchIndexKey           = "batchIndex"
	ParamBatchSizeKey            = "batchSize"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
