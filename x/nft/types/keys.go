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
	ParamModuleAdminKey       = "module_admin"
	ParamBondDenomKey         = "bond_denom"
	ParamPrefixKey            = "prefix"
	ParamNftSequenceKey       = "nft_sequence"
	ParamsNftVestingTimeKey   = "vesting_time"
	ParamsNftVestingPeriodKey = "nft_vesting_period"
	ParamVestingCountKey      = "vesting_count"
	ParamNftCostKey           = "nft_cost"
	ParamBatchIndexKey        = "batch_index"
	ParamBatchSizeKey         = "batch_size"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
