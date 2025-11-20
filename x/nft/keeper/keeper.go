package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/errors"
	accumulatortypes "github.com/cosmos/cosmos-sdk/x/accumulator/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		paramstore        paramtypes.Subspace
		bankKeeper        types.BankKeeper
		stakingKeeper     types.StakingKeeper
		accumulatorKeeper types.AccumulatorKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	accumulatorKeeper types.AccumulatorKeeper,
) *Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		memKey:            memKey,
		paramstore:        ps,
		bankKeeper:        bankKeeper,
		stakingKeeper:     stakingKeeper,
		accumulatorKeeper: accumulatorKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) DistributeToAddress(ctx sdk.Context, amount sdk.Coins, owner sdk.AccAddress, nftAddress sdk.AccAddress) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, nftAddress, types.ModuleName, amount)
	if err != nil {
		err = errors.Wrap(err, "failed to distribute tokens to nft module")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, amount)
	if err != nil {
		err = errors.Wrap(err, "failed to distribute tokens to nft module")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	return nil
}

func (k Keeper) IsDelegated(ctx sdk.Context, nftAddress sdk.AccAddress) bool {
	return len(k.stakingKeeper.GetAllDelegatorDelegations(ctx, nftAddress)) > 0
}

func (k Keeper) IsModuleAdmin(ctx sdk.Context, address string) bool {
	params := k.GetParams(ctx)

	return params.ModuleAdmin == address
}

func (k *Keeper) CreateNft(ctx sdk.Context, owner string, startVestingBlock int64, vestingPeriodReward *big.Int) (*types.NFT, uint64, error) {
	var (
		nftAddress string
		err        error
	)

	sequence := k.GetNftSequence(ctx) + 1

	// loop is used to avoid sequence collisions
	for {
		nftAddress, err = sdk.Bech32ifyAddressBytes(
			k.GetPrefix(ctx),
			address.Derive(
				authtypes.NewModuleAddress(accumulatortypes.ModuleName),
				[]byte(strconv.FormatUint(sequence, 10)),
			),
		)

		if err != nil {
			return nil, 0, errorsmod.Wrap(err, "failed to retrieve NFT address")
		}

		_, ok := k.GetNFT(ctx, nftAddress)
		if ok {
			sequence++
			continue
		}

		break
	}

	newNft := types.NFT{
		Address:             nftAddress,
		Owner:               owner,
		RewardPerPeriod:     sdk.NewCoin(k.GetBondDenom(ctx), sdk.NewIntFromBigInt(vestingPeriodReward)),
		VestingPeriodsCount: 0,
		AvailableToWithdraw: sdk.NewCoin(k.GetBondDenom(ctx), sdk.ZeroInt()),
		Denom:               k.GetBondDenom(ctx),
		StartVestingBlock:   startVestingBlock,
	}

	return &newNft, sequence, nil
}
