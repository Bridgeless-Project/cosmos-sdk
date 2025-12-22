package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrSample                     = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidAmount              = sdkerrors.Register(ModuleName, 1101, "invalid amount")
	ErrValidatorNotFound          = sdkerrors.Register(ModuleName, 1102, "validator not found")
	ErrFailedToDelegate           = sdkerrors.Register(ModuleName, 1103, "failed to delegate nft")
	ErrFailedToSendTokenAmount    = sdkerrors.Register(ModuleName, 1104, "failed to send token amount")
	ErrNFTIsDelegated             = sdkerrors.Register(ModuleName, 1105, "NFT is delegated")
	ErrNFTNotFound                = sdkerrors.Register(ModuleName, 1106, "NFT not found")
	ErrNFTInvalidOwner            = sdkerrors.Register(ModuleName, 1107, "NFT's owner is invalid")
	ErrFailedToBeginRedeligation  = sdkerrors.Register(ModuleName, 1108, "failed to begin redelegation")
	ErrAddressISNFT               = sdkerrors.Register(ModuleName, 1109, "address is an NFT")
	ErrUnauthorizedNFTCreator     = sdkerrors.Register(ModuleName, 1110, "unauthorized NFT creator")
	ErrInvalidSequence            = sdkerrors.Register(ModuleName, 1111, "invalid sequence")
	ErrInvalidPrefix              = sdkerrors.Register(ModuleName, 1112, "invalid prefix")
	ErrInvalidBondDenom           = sdkerrors.Register(ModuleName, 1113, "invalid bond denom")
	ErrInvalidBalance             = sdkerrors.Register(ModuleName, 1114, "invalid balance")
	ErrInvalidVestingPeriod       = sdkerrors.Register(ModuleName, 1115, "invalid vesting period")
	ErrInvalidNftCost             = sdkerrors.Register(ModuleName, 1117, "invalid nft cost")
	ErrInvalidBatchSize           = sdkerrors.Register(ModuleName, 1118, "invalid batch size")
	ErrNoNFTsProvided             = sdkerrors.Register(ModuleName, 1119, "nfts not provided")
	ErrMinSelfDelegation          = sdkerrors.Register(ModuleName, 1120, "min self delegation not met")
	ErrEmptyNftMetadataUri        = sdkerrors.Register(ModuleName, 1121, "empty nft metadata uri")
	ErrInvalidVestingPeriodsLimit = sdkerrors.Register(ModuleName, 1122, "invalid vesting period limit")
)
