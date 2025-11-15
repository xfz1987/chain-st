package token

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	tokensimulation "chainst/x/token/simulation"
	"chainst/x/token/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	tokenGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&tokenGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgMintTokens          = "op_weight_msg_token"
		defaultWeightMsgMintTokens int = 100
	)

	var weightMsgMintTokens int
	simState.AppParams.GetOrGenerate(opWeightMsgMintTokens, &weightMsgMintTokens, nil,
		func(_ *rand.Rand) {
			weightMsgMintTokens = defaultWeightMsgMintTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintTokens,
		tokensimulation.SimulateMsgMintTokens(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgTransferTokens          = "op_weight_msg_token"
		defaultWeightMsgTransferTokens int = 100
	)

	var weightMsgTransferTokens int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferTokens, &weightMsgTransferTokens, nil,
		func(_ *rand.Rand) {
			weightMsgTransferTokens = defaultWeightMsgTransferTokens
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferTokens,
		tokensimulation.SimulateMsgTransferTokens(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
