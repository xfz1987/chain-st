package keeper

import (
	"context"
	"fmt"

	"chainst/x/token/types"

	"cosmossdk.io/math"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) MintTokens(ctx context.Context, msg *types.MsgMintTokens) (*types.MsgMintTokensResponse, error) {
	// 验证创建者地址
	creatorAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid creator address")
	}

	// 验证金额
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "amount must be greater than 0")
	}

	// 创建币种
	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, math.NewInt(int64(msg.Amount))))

	// 使用 bank keeper 铸造代币到模块账户
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to mint coins")
	}

	// 从模块账户转账到创建者账户
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins to creator")
	}

	// 获取 SDK 上下文并触发事件
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"mint_tokens",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("amount", fmt.Sprintf("%d", msg.Amount)),
			sdk.NewAttribute("denom", msg.Denom),
		),
	)

	return &types.MsgMintTokensResponse{}, nil
}
