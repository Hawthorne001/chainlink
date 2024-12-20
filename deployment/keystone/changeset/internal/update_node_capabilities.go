package internal

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type UpdateNodeCapabilitiesImplRequest struct {
	Chain             deployment.Chain
	ContractSet       *kslib.ContractSet
	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability

	UseMCMS bool
}

func (req *UpdateNodeCapabilitiesImplRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	if req.ContractSet == nil {
		return fmt.Errorf("registry is nil")
	}

	return nil
}

func UpdateNodeCapabilitiesImpl(lggr logger.Logger, req *UpdateNodeCapabilitiesImplRequest) (*UpdateNodesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}
	// collect all the capabilities and add them to the registry
	var capabilities []kcr.CapabilitiesRegistryCapability
	for _, cap := range req.P2pToCapabilities {
		capabilities = append(capabilities, cap...)
	}
	op, err := kslib.AddCapabilities(lggr, req.ContractSet, req.Chain, capabilities, req.UseMCMS)
	if err != nil {
		return nil, fmt.Errorf("failed to add capabilities: %w", err)
	}

	p2pToUpdates := map[p2pkey.PeerID]NodeUpdate{}
	for id, caps := range req.P2pToCapabilities {
		p2pToUpdates[id] = NodeUpdate{Capabilities: caps}
	}

	updateNodesReq := &UpdateNodesRequest{
		Chain:        req.Chain,
		P2pToUpdates: p2pToUpdates,
		ContractSet:  req.ContractSet,
		Ops:          op,
		UseMCMS:      req.UseMCMS,
	}
	resp, err := UpdateNodes(lggr, updateNodesReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update nodes: %w", err)
	}
	return resp, nil
}
