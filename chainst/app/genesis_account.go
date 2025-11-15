package app

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var _ authtypes.GenesisAccount = (*GenesisAccount)(nil)

// GenesisAccount defines a type that implements the GenesisAccount interface
// to be used for simulation accounts in the genesis state.
// 定义创世账户结构
type GenesisAccount struct {
	*authtypes.BaseAccount

	// 锁仓账户字段
	// 初始锁仓金额
	OriginalVesting  sdk.Coins `json:"original_vesting" yaml:"original_vesting"`   // total vesting coins upon initialization
	// 已委托的自由币
	DelegatedFree    sdk.Coins `json:"delegated_free" yaml:"delegated_free"`       // delegated vested coins at time of delegation
	// 已委托的锁仓币
	DelegatedVesting sdk.Coins `json:"delegated_vesting" yaml:"delegated_vesting"` // delegated vesting coins at time of delegation
	// 锁仓开始时间
	StartTime        int64     `json:"start_time" yaml:"start_time"`               // vesting start time (UNIX Epoch time)
	// 锁仓结束时间
	EndTime          int64     `json:"end_time" yaml:"end_time"`                   // vesting end time (UNIX Epoch time)

	// 模块账户字段
	// 模块名称
	ModuleName        string   `json:"module_name" yaml:"module_name"`               // name of the module account
	// 模块权限
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // permissions of module account
}

// 验证账户有效性
func (sga GenesisAccount) Validate() error {
	if !sga.OriginalVesting.IsZero() {
		if sga.StartTime >= sga.EndTime {
			return errors.New("vesting start-time cannot be before end-time")
		}
	}

	if sga.ModuleName != "" {
		ma := authtypes.ModuleAccount{
			BaseAccount: sga.BaseAccount, Name: sga.ModuleName, Permissions: sga.ModulePermissions,
		}
		if err := ma.Validate(); err != nil {
			return err
		}
	}

	return sga.BaseAccount.Validate()
}
