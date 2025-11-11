package keeper

import (
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft/types"
)

//func distributeNft(owner string, commonValue, nftValue *big.Int, vperiod, vperiodCount int64, vrewards sdk.Coin, denom string, startId int64) []nfttypes.NFT {
//	nfts := make([]nfttypes.NFT, 0)
//	id := startId
//	for commonValue.Cmp(nftValue) != -1 {
//		nftAddress, err := sdk.Bech32ifyAddressBytes(prefix, address.Derive(authtypes.NewModuleAddress(acumulatortypes.ModuleName), []byte(strconv.FormatInt(id, 10))))
//		if err != nil {
//			panic(err)
//		}
//		newNft := nfttypes.NFT{
//			Address:             nftAddress,
//			Owner:               owner,
//			VestingPeriod:       vperiod,
//			RewardPerPeriod:     vrewards,
//			VestingPeriodsCount: vperiodCount,
//			AvailableToWithdraw: sdk.NewCoin(denom, sdk.NewInt(0)),
//			Denom:               Denom,
//		}
//
//		nfts = append(nfts, newNft)
//		id++
//		commonValue = commonValue.Sub(commonValue, nftValue)
//	}
//
//	return nfts
//}
