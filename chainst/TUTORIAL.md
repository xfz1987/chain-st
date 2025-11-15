# ChainST - Cosmos SDK 区块链开发教程

本教程将指导你使用 Cosmos SDK 开发一个完整的区块链项目。

## 项目功能

1. ✅ 代币铸造 (Token Minting)
2. ✅ 用户账户创建
3. ✅ 代币转账
4. ✅ 验证者和共识机制 (PoS)
5. ✅ 区块链浏览器 API
6. ✅ 简单的 Web 区块链浏览器

## 开发环境要求

- Go 1.21+ (已安装: go1.25.4)
- Ignite CLI v29.6.0+ (已安装: v29.6.1-dev)
- Node.js 18+ (用于前端浏览器)

## 项目结构

```
chainst/
├── app/                    # 应用程序配置和初始化
├── cmd/                    # 命令行入口
├── proto/                  # Protocol Buffers 定义
├── x/                      # 自定义模块
│   └── token/             # 代币模块
│       ├── keeper/        # 状态管理器
│       ├── types/         # 类型定义
│       └── module.go      # 模块定义
├── config.yml             # Ignite 配置文件
├── go.mod                 # Go 依赖
└── README.md              # 项目说明
```

## 步骤 1: 初始化项目 (已完成)

项目已使用以下命令初始化:

```bash
cd ~/work/web3/chain-st
ignite scaffold chain chainst --no-module --address-prefix chainst
```

## 步骤 2: 创建代币模块

### 2.1 创建 token 模块

```bash
cd chainst
ignite scaffold module token --dep bank
```

### 2.2 添加代币铸造消息

```bash
ignite scaffold message mint-tokens amount:uint denom:string --module token --signer creator
```

### 2.3 添加代币转账消息

```bash
ignite scaffold message transfer-tokens to:string amount:uint denom:string --module token --signer from
```

### 2.4 添加代币信息查询

```bash
ignite scaffold query token-info denom:string --module token
```

## 步骤 3: 实现代币铸造逻辑

编辑文件: `x/token/keeper/msg_server_mint_tokens.go`

```go
package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"chainst/x/token/types"
)

func (k msgServer) MintTokens(goCtx context.Context, msg *types.MsgMintTokens) (*types.MsgMintTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 验证金额
	if msg.Amount == 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// 创建币种
	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, sdk.NewIntFromUint64(msg.Amount)))

	// 铸造代币到创建者地址
	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	// 使用 bank keeper 铸造代币
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, err
	}

	// 转账到创建者
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins); err != nil {
		return nil, err
	}

	// 触发事件
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"mint_tokens",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("amount", fmt.Sprintf("%d", msg.Amount)),
			sdk.NewAttribute("denom", msg.Denom),
		),
	)

	return &types.MsgMintTokensResponse{}, nil
}
```

## 步骤 4: 实现代币转账逻辑

编辑文件: `x/token/keeper/msg_server_transfer_tokens.go`

```go
package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"chainst/x/token/types"
)

func (k msgServer) TransferTokens(goCtx context.Context, msg *types.MsgTransferTokens) (*types.MsgTransferTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 验证金额
	if msg.Amount == 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// 获取地址
	fromAddr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	toAddr, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}

	// 创建币种
	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, sdk.NewIntFromUint64(msg.Amount)))

	// 转账
	if err := k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, coins); err != nil {
		return nil, err
	}

	// 触发事件
	ctx.EventManager().EmitEvent(
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
```

## 步骤 5: 构建和启动区块链

### 5.1 安装依赖

```bash
go mod tidy
```

### 5.2 启动区块链

```bash
ignite chain serve
```

这个命令会:
- 编译区块链代码
- 初始化本地测试网
- 创建验证者节点
- 启动区块链

默认端口:
- RPC: http://localhost:26657
- API: http://localhost:1317
- gRPC: localhost:9090

## 步骤 6: 测试区块链功能

### 6.1 创建用户账户

账户会在启动时自动创建。查看账户:

```bash
chainstd keys list
```

### 6.2 铸造代币

```bash
chainstd tx token mint-tokens 1000000 mytoken --from alice --chain-id chainst --yes
```

### 6.3 查询余额

```bash
chainstd query bank balances $(chainstd keys show alice -a)
```

### 6.4 转账代币

```bash
chainstd tx token transfer-tokens $(chainstd keys show bob -a) 100 mytoken --from alice --chain-id chainst --yes
```

### 6.5 查询区块信息

```bash
# 查询最新区块
chainstd query block

# 查询指定高度区块
chainstd query block 100
```

## 步骤 7: 区块链浏览器 API

Cosmos SDK 自带 REST API,可以通过以下端点访问:

### 主要 API 端点:

1. **获取最新区块**
   ```
   GET http://localhost:1317/cosmos/base/tendermint/v1beta1/blocks/latest
   ```

2. **获取指定高度区块**
   ```
   GET http://localhost:1317/cosmos/base/tendermint/v1beta1/blocks/{height}
   ```

3. **查询账户余额**
   ```
   GET http://localhost:1317/cosmos/bank/v1beta1/balances/{address}
   ```

4. **查询交易**
   ```
   GET http://localhost:1317/cosmos/tx/v1beta1/txs/{hash}
   ```

5. **查询验证者列表**
   ```
   GET http://localhost:1317/cosmos/staking/v1beta1/validators
   ```

## 步骤 8: 创建简单的区块链浏览器

接下来的步骤将创建一个 Web 界面来显示区块链信息。

## Cosmos SDK 的工作原理

### 用户管理
在 Cosmos SDK 中,用户账户是通过密钥对自动管理的:
- 每个账户有一个地址 (address)
- 账户可以持有代币
- 使用 `chainstd keys` 命令管理账户

### 挖矿机制 (验证者)
Cosmos 使用 Proof of Stake (PoS) 而不是传统的 PoW (工作量证明):
- 验证者通过质押代币获得出块权
- 使用 Tendermint 共识算法
- 验证者轮流提议和验证区块
- 质押越多,获得出块权的机会越大
- 本项目配置为每 12 秒产生一个新区块

### 区块产生流程
1. 验证者提议新区块
2. 其他验证者验证并投票
3. 达成共识后,区块被添加到链上
4. 新的区块高度生成

## 常用命令

```bash
# 启动区块链
ignite chain serve

# 查看账户列表
chainstd keys list

# 查询区块链状态
chainstd status

# 查询账户余额
chainstd query bank balances <address>

# 发送交易
chainstd tx bank send <from> <to> <amount> --chain-id chainst

# 查询交易
chainstd query tx <hash>
```

## 下一步

现在你已经有了一个功能完整的 Cosmos SDK 区块链!你可以:

1. 添加更多自定义功能
2. 部署到测试网络
3. 集成前端应用
4. 添加更多验证者节点
5. 实现跨链功能 (IBC)

## 故障排除

### 问题: 端口被占用
```bash
lsof -i :26657
kill -9 <PID>
```

### 问题: 重置区块链数据
```bash
ignite chain clean
```

### 问题: 依赖问题
```bash
go clean -modcache
go mod tidy
```

## 资源链接

- Cosmos SDK 官方文档: https://docs.cosmos.network/
- Tendermint 文档: https://docs.tendermint.com/
- Ignite CLI 文档: https://docs.ignite.com/

---

祝你开发愉快!如有问题,请查阅官方文档或社区支持。
