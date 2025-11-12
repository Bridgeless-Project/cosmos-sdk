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
	ParamModuleAdminKey             = "module_admin"
	ParamBondDenomKey               = "bond_denom"
	ParamPrefixKey                  = "prefix"
	ParamNftSequenceKey             = "nft_sequence"
	ParamsNftVestingPeriodKey       = "nft_vesting_period"
	ParamsNftVestingPeriodRewardKey = "nft_vesting_period_reward"
	ParamVestingCountKey            = "vesting_count"
	ParamNftCostKey                 = "nft_cost"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
