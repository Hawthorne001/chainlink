package memory

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	"github.com/stretchr/testify/require"

	solTestUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/testutils"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
)

type EVMChain struct {
	Backend     *simulated.Backend
	DeployerKey *bind.TransactOpts
	Users       []*bind.TransactOpts
}

type SolanaChain struct {
	Client      *solRpc.Client
	DeployerKey *solana.PrivateKey
}

func fundAddress(t *testing.T, from *bind.TransactOpts, to common.Address, amount *big.Int, backend *simulated.Backend) {
	ctx := tests.Context(t)
	nonce, err := backend.Client().PendingNonceAt(ctx, from.From)
	require.NoError(t, err)
	gp, err := backend.Client().SuggestGasPrice(ctx)
	require.NoError(t, err)
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      21000,
		To:       &to,
		Value:    amount,
	})
	signedTx, err := from.Signer(from.From, rawTx)
	require.NoError(t, err)
	err = backend.Client().SendTransaction(ctx, signedTx)
	require.NoError(t, err)
	backend.Commit()
}

func GenerateChains(t *testing.T, numChains int, numUsers int) map[uint64]EVMChain {
	chains := make(map[uint64]EVMChain)
	for i := 0; i < numChains; i++ {
		chainID := chainsel.TEST_90000001.EvmChainID + uint64(i)
		chains[chainID] = evmChain(t, numUsers)
	}
	return chains
}

func getTestSolanaChainSelectors() []uint64 {
	result := []uint64{}
	for _, x := range chainsel.SolanaALL {
		if x.Name == x.ChainID {
			result = append(result, x.Selector)
		}
	}
	return result
}

func GenerateChainsSol(t *testing.T, numChains int) map[uint64]SolanaChain {
	testSolanaChainSelectors := getTestSolanaChainSelectors()
	if len(testSolanaChainSelectors) < numChains {
		t.Fatalf("not enough test solana chain selectors available")
	}
	chains := make(map[uint64]SolanaChain)
	for i := 0; i < numChains; i++ {
		chainID := testSolanaChainSelectors[i]
		url, _ := solTestUtil.SetupLocalSolNodeWithFlags(t)
		admin, gerr := solana.NewRandomPrivateKey()
		solTestUtil.FundTestAccounts(t, []solana.PublicKey{admin.PublicKey()}, url)
		require.NoError(t, gerr)
		chains[chainID] = SolanaChain{
			Client:      solRpc.New(url),
			DeployerKey: &admin,
		}
	}
	return chains
}

func GenerateChainsWithIds(t *testing.T, chainIDs []uint64, numUsers int) map[uint64]EVMChain {
	chains := make(map[uint64]EVMChain)
	for _, chainID := range chainIDs {
		chains[chainID] = evmChain(t, numUsers)
	}
	return chains
}

func evmChain(t *testing.T, numUsers int) EVMChain {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	owner, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	genesis := types.GenesisAlloc{
		owner.From: {Balance: assets.Ether(1_000_000).ToInt()}}
	// create a set of user keys
	var users []*bind.TransactOpts
	for j := 0; j < numUsers; j++ {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		require.NoError(t, err)
		users = append(users, user)
		genesis[user.From] = types.Account{Balance: assets.Ether(1_000_000).ToInt()}
	}
	// there have to be enough initial funds on each chain to allocate for all the nodes that share the given chain in the test
	backend := simulated.NewBackend(genesis, simulated.WithBlockGasLimit(50000000))
	backend.Commit() // ts will be now.
	return EVMChain{
		Backend:     backend,
		DeployerKey: owner,
		Users:       users,
	}
}
