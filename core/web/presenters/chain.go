package presenters

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	commonTypes "github.com/smartcontractkit/chainlink/v2/common/types"
)

type ChainResource struct {
	JAID
	Network string `json:"network"`
	Enabled bool   `json:"enabled"`
	Config  string `json:"config"` // TOML
}

// GetName implements the api2go EntityNamer interface
func (r ChainResource) GetName() string {
	return "chain"
}

// NewChainResource returns a new ChainResource for chain.
func NewChainResource(chain commonTypes.ChainStatusWithID) ChainResource {
	return ChainResource{
		JAID:    NewJAID(chain.RelayID.ChainID),
		Network: chain.RelayID.Network,
		Config:  chain.Config,
		Enabled: chain.Enabled,
	}
}

type NodeResource struct {
	JAID
	ChainID string `json:"chainID"`
	Name    string `json:"name"`
	Config  string `json:"config"` // TOML
	State   string `json:"state"`
}

// NewNodeResource returns a new NodeResource for node.
func NewNodeResource(node types.NodeStatus) NodeResource {
	return NodeResource{
		JAID:    NewPrefixedJAID(node.Name, node.ChainID),
		ChainID: node.ChainID,
		Name:    node.Name,
		State:   node.State,
		Config:  node.Config,
	}
}

// GetName implements the api2go EntityNamer interface
func (r NodeResource) GetName() string {
	return "node"
}
