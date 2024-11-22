// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {USDCBridgeMigrator} from "../../../../pools/USDC/USDCBridgeMigrator.sol";
import {HybridLockReleaseUSDCTokenPoolSetup} from "../USDCTokenPoolSetup.t.sol";

contract USDCBridgeMigrator_excludeTokensFromBurn is HybridLockReleaseUSDCTokenPoolSetup {
  function test_excludeTokensWhenNoMigrationProposalPending_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(USDCBridgeMigrator.NoMigrationProposalPending.selector));

    vm.startPrank(OWNER);

    s_usdcTokenPool.excludeTokensFromBurn(SOURCE_CHAIN_SELECTOR, 1e6);
  }
}
