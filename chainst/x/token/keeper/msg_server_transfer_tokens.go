package keeper

import (
	"context"
	"fmt"

	"chainst/x/token/types"

	"cosmossdk.io/math"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferTokens(ctx context.Context, msg *types.MsgTransferTokens) (*types.MsgTransferTokensResponse, error) {
	// 验证发送者地址
	fromAddr, err := k.addressCodec.StringToBytes(msg.From)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid from address")
	}

	// 验证接收者地址
	toAddr, err := k.addressCodec.StringToBytes(msg.To)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid to address")
	}

	// 验证金额
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "amount must be greater than 0")
	}

	// 创建币种
	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, math.NewInt(int64(msg.Amount))))

	// 使用 bank keeper 转账
	if err := k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, coins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to send coins")
	}

	// 获取 SDK 上下文并触发事件
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"transfer_tokens",
			sdk.NewAttribute("from", msg.From),
			sdk.NewAttribute("to", msg.To),
			sdk.NewAttribute("amount", fmt.Sprintf("%d", msg.Amount)),
			sdk.NewAttribute("denom", msg.Denom),
		),
	)

	return &types.MsgTransferTokensResponse{}, nil
}
