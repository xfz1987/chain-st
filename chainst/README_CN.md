# ChainST - 基于 Cosmos SDK 的区块链项目

这是一个功能完整的区块链项目,基于 Cosmos SDK 开发,包含代币铸造、转账、区块浏览器等功能。

## 功能特性

✅ **代币铸造** - 支持自定义代币的铸造
✅ **代币转账** - 账户间的代币转移
✅ **用户账户管理** - 创建和管理区块链账户
✅ **PoS 共识** - 基于 Tendermint 的权益证明机制
✅ **区块链浏览器** - Web 界面查看区块和交易
✅ **REST API** - 完整的 HTTP API 接口

## 快速开始

### 1. 启动区块链

在项目根目录执行:

```bash
ignite chain serve
```

这个命令会自动:
- 编译区块链代码
- 初始化本地测试网络
- 创建验证者节点
- 启动区块链服务

启动后,区块链将在以下端口运行:
- **RPC**: http://localhost:26657
- **API**: http://localhost:1317
- **gRPC**: localhost:9090

### 2. 打开区块链浏览器

在浏览器中打开:

```
file:///Users/autonomic/work/web3/chain-st/explorer/index.html
```

或者使用 Python 启动一个简单的 HTTP 服务器:

```bash
cd ../explorer
python3 -m http.server 8000
# 然后访问 http://localhost:8000
```

## 使用教程

### 查看账户列表

```bash
chainstd keys list
```

默认会创建两个测试账户: `alice` 和 `bob`

### 查看账户地址

```bash
chainstd keys show alice -a
```

### 查询账户余额

```bash
chainstd query bank balances $(chainstd keys show alice -a)
```

### 铸造代币

使用 token 模块铸造新代币:

```bash
chainstd tx token mint-tokens 1000000 mytoken \
  --from alice \
  --chain-id chainst \
  --yes
```

参数说明:
- `1000000`: 代币数量
- `mytoken`: 代币名称
- `--from alice`: 使用 alice 账户签名
- `--chain-id chainst`: 链 ID

### 转账代币

将代币从一个账户转移到另一个账户:

```bash
chainstd tx token transfer-tokens \
  $(chainstd keys show bob -a) \
  100000 \
  mytoken \
  --from alice \
  --chain-id chainst \
  --yes
```

参数说明:
- 第一个参数: 接收者地址
- `100000`: 转账数量
- `mytoken`: 代币名称

### 使用原生 bank 模块转账

也可以使用 Cosmos SDK 内置的 bank 模块:

```bash
chainstd tx bank send \
  $(chainstd keys show alice -a) \
  $(chainstd keys show bob -a) \
  1000stake \
  --chain-id chainst \
  --yes
```

### 查询最新区块

```bash
chainstd query block
```

### 查询指定高度区块

```bash
chainstd query block 100
```

### 查询交易

```bash
chainstd query tx <交易哈希>
```

### 查询节点状态

```bash
chainstd status
```

## API 使用

### REST API 端点

区块链提供了完整的 REST API,以下是一些常用端点:

**1. 获取最新区块**
```bash
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/blocks/latest
```

**2. 获取指定高度区块**
```bash
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/blocks/100
```

**3. 查询账户余额**
```bash
curl http://localhost:1317/cosmos/bank/v1beta1/balances/<地址>
```

**4. 查询交易**
```bash
curl http://localhost:1317/cosmos/tx/v1beta1/txs/<交易哈希>
```

**5. 查询验证者列表**
```bash
curl http://localhost:1317/cosmos/staking/v1beta1/validators
```

**6. 查询节点信息**
```bash
curl http://localhost:26657/status
```

### gRPC API

区块链也提供 gRPC 接口,监听在 `localhost:9090`

## 项目结构

```
chainst/
├── app/                    # 应用程序配置
│   ├── app.go             # 应用主文件
│   └── ...
├── cmd/                    # 命令行入口
│   └── chainstd/          # 节点守护进程
├── proto/                  # Protocol Buffers 定义
│   ├── chainst/           # 自定义模块的 proto 文件
│   └── ...
├── x/                      # 自定义模块
│   └── token/             # Token 模块
│       ├── keeper/        # 状态管理器
│       ├── types/         # 类型定义
│       └── module.go      # 模块定义
├── config.yml             # Ignite 配置
├── go.mod                 # Go 依赖
└── README_CN.md           # 本文件
```

## 核心概念

### 1. 账户系统

Cosmos SDK 使用账户模型,每个账户有:
- **地址**: 类似于 `chainst1...` 的 Bech32 格式地址
- **余额**: 可以持有多种代币
- **序号 (Sequence)**: 防止交易重放攻击

账户通过公私钥对管理,私钥由用户保管。

### 2. 共识机制

ChainST 使用 Tendermint 共识引擎:
- **Proof of Stake (PoS)**: 权益证明,而非工作量证明
- **验证者**: 质押代币的节点参与共识
- **BFT**: 拜占庭容错,可容忍 1/3 的恶意节点
- **快速确认**: 12 秒出块时间

### 3. 挖矿 vs 验证

与传统的 PoW 不同:
- **没有矿工**: 使用验证者替代
- **没有算力竞争**: 通过质押代币获得出块权
- **节能环保**: 不需要大量计算

成为验证者的步骤:
```bash
# 创建验证者
chainstd tx staking create-validator \
  --amount=1000000stake \
  --pubkey=$(chainstd tendermint show-validator) \
  --moniker="我的验证者" \
  --chain-id=chainst \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --from=alice \
  --yes
```

### 4. Token 模块

自定义的 token 模块提供:
- `MintTokens`: 铸造新代币
- `TransferTokens`: 转账代币
- `QueryTokenInfo`: 查询代币信息

## 开发和扩展

### 添加新的消息类型

```bash
ignite scaffold message <消息名> <字段1>:<类型1> <字段2>:<类型2> --module token
```

### 添加新的查询

```bash
ignite scaffold query <查询名> <参数>:<类型> --module token
```

### 添加新的模块

```bash
ignite scaffold module <模块名> --dep <依赖模块>
```

### 编译代码

```bash
ignite chain build
```

### 测试

```bash
go test ./...
```

## 故障排除

### 端口被占用

如果遇到端口冲突:

```bash
# 查找占用进程
lsof -i :26657
lsof -i :1317

# 杀死进程
kill -9 <PID>
```

### 重置区块链数据

```bash
ignite chain clean
```

### 清理 Go 缓存

```bash
go clean -modcache
go mod tidy
```

### 查看日志

启动时的日志会显示在终端,查看详细日志:

```bash
chainstd start --log_level info
```

## 部署到生产环境

### 1. 配置 systemd 服务

创建 `/etc/systemd/system/chainstd.service`:

```ini
[Unit]
Description=ChainST Node
After=network.target

[Service]
Type=simple
User=<用户名>
ExecStart=/usr/local/bin/chainstd start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl enable chainstd
sudo systemctl start chainstd
```

### 2. 配置防火墙

开放必要的端口:
```bash
sudo ufw allow 26656  # P2P
sudo ufw allow 26657  # RPC
sudo ufw allow 1317   # API
```

### 3. 配置 Nginx 反向代理

```nginx
server {
    listen 80;
    server_name api.chainst.example.com;

    location / {
        proxy_pass http://localhost:1317;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 高级功能

### IBC (跨链通信)

Cosmos SDK 支持 IBC 协议,可以与其他 Cosmos 链交互:

```bash
# 配置 IBC relayer
# 详见: https://github.com/cosmos/relayer
```

### 升级链

使用治理提案进行链上升级:

```bash
chainstd tx gov submit-proposal software-upgrade <升级名> \
  --title="升级到 v2.0" \
  --description="添加新功能" \
  --upgrade-height=100000 \
  --from=alice \
  --deposit=10000000stake
```

## 资源链接

- **Cosmos SDK 文档**: https://docs.cosmos.network/
- **Tendermint 文档**: https://docs.tendermint.com/
- **Ignite CLI 文档**: https://docs.ignite.com/
- **Cosmos 开发者门户**: https://tutorials.cosmos.network/

## 贡献

欢迎提交 Issue 和 Pull Request!

## 许可证

MIT License

---

## 完整操作示例

下面是一个完整的使用流程示例:

```bash
# 1. 启动区块链
ignite chain serve

# 2. 新开一个终端,查看账户
chainstd keys list

# 3. 查看 alice 的余额
chainstd query bank balances $(chainstd keys show alice -a)

# 4. 铸造新代币
chainstd tx token mint-tokens 1000000 mytoken \
  --from alice \
  --chain-id chainst \
  --yes

# 5. 等待几秒,然后查询余额
chainstd query bank balances $(chainstd keys show alice -a)

# 6. 转账给 bob
chainstd tx token transfer-tokens \
  $(chainstd keys show bob -a) \
  50000 \
  mytoken \
  --from alice \
  --chain-id chainst \
  --yes

# 7. 查询 bob 的余额
chainstd query bank balances $(chainstd keys show bob -a)

# 8. 查询最新区块
chainstd query block

# 9. 查询节点状态
chainstd status
```

祝你使用愉快!
