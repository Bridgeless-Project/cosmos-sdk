package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/nft/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/pkg/errors"
)

func (m msgServer) BecameValidator(goCtx context.Context, msg *types.MsgBecameValidator) (*types.MsgBecameValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := m.stakingKeeper.GetValidator(ctx, valAddr); found {
		return nil, stakingTypes.ErrValidatorOwnerExists
	}

	pk, ok := msg.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	if _, found := m.stakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk)); found {
		return nil, stakingTypes.ErrValidatorPubKeyExists
	}

	minDelegation, ok := sdk.NewIntFromString(m.stakingKeeper.MinDelegationAmount(ctx))
	if !ok {
		return nil, errors.New("invalid min delegation amount")
	}

	bondDenom := m.GetBondDenom(ctx)
	amount := sdk.NewInt(0)
	nfts := make([]types.NFT, 0)

	for _, nftAddress := range msg.NftAddresses {
		nft, ok := m.GetNFT(ctx, nftAddress)
		if !ok {
			return nil, types.ErrNFTNotFound
		}

		addr, err := sdk.AccAddressFromBech32(nft.Address)
		if err != nil {
			return nil, errors.Wrap(err, "invalid nft address")
		}

		nftbalance := m.bankKeeper.GetBalance(ctx, addr, bondDenom)
		amount = amount.Add(nftbalance.Amount)
		nfts = append(nfts, nft)
	}

	if amount.LT(minDelegation) {
		return nil, errors.Wrapf(
			types.ErrMinSelfDelegation,
			"provided delegation: %s%s, required: %s%s.",
			amount.String(),
			bondDenom,
			minDelegation.String(),
			bondDenom,
		)
	}

	if len(nfts) == 0 {
		return nil, types.ErrNoNFTsProvided
	}

	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		pkType := pk.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return nil, sdkerrors.Wrapf(
				stakingTypes.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pk.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	validator, err := stakingTypes.NewValidator(valAddr, pk, msg.Description)
	if err != nil {
		return nil, err
	}

	commission := stakingTypes.NewCommissionWithTime(
		msg.Commission.Rate, msg.Commission.MaxRate,
		msg.Commission.MaxChangeRate, ctx.BlockHeader().Time,
	)

	validator, err = validator.SetInitialCommission(commission)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set initial commission")
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get delegator address from bech32")
	}

	m.stakingKeeper.SetValidator(ctx, validator)
	err = m.stakingKeeper.SetValidatorByConsAddr(ctx, validator)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set validator by cons addr")
	}
	m.stakingKeeper.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	if err = m.stakingKeeper.AfterValidatorCreated(ctx, validator.GetOperator()); err != nil {
		return nil, errors.Wrap(err, "failed in after validator created hook")
	}

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	for _, nft := range nfts {
		addr, err := sdk.AccAddressFromBech32(nft.Address)
		if err != nil {
			return nil, errors.Wrap(err, "invalid nft address")
		}

		err = m.DelegateNFT(ctx, addr, delegatorAddress, valAddr, sdk.NewInt(0), true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to delegate NFTs to validator")
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingTypes.EventTypeCreateValidator,
			sdk.NewAttribute(stakingTypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, stakingTypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return new(types.MsgBecameValidatorResponse), nil
}
