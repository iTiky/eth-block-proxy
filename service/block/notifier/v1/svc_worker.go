package v1

import (
	"bytes"
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// worker is the main service work loop.
func (svc *BlockNotifierSvc) worker() {
	pollTimer := time.NewTimer(1 * time.Second)
	pollTimer.Stop()

	for {
		// Wait for Poll / Stop event
		pollTimer.Reset(svc.pollDur)
		select {
		case <-pollTimer.C:
		case <-svc.stopCh:
			return
		}

		svc.workerStep()
	}
}

// workerStep defines the main service work loop step (iteration).
func (svc *BlockNotifierSvc) workerStep() {
	// Get the latest block number (not the whole block to reduce network load)
	latestBlockNumber, err := svc.reader.GetLatestBlockNumber(context.Background())
	if err != nil {
		svc.Logger(context.TODO()).Error().
			Err(err).
			Msg("Polling the latest block number")
		return
	}

	// Check if a new block was minted and get the new block data
	if !svc.checkNewBlockReceived(latestBlockNumber) {
		return
	}

	newBlock, err := svc.reader.GetBlock(context.Background(), latestBlockNumber)
	if err != nil {
		svc.Logger(context.TODO()).Error().
			Uint64("block", latestBlockNumber).
			Err(err).
			Msg("Fetching the new block data")
		return
	}
	if newBlock == nil {
		return
	}

	// Handle chain fork event
	if svc.checkChainForked(newBlock) {
		svc.Logger(context.TODO()).Info().
			Uint64("block", latestBlockNumber).
			Msg("Chain fork event detected")
		// Reset state
		svc.latestBlock = nil
		// Emit event
		if svc.chainForkHandler != nil {
			go svc.chainForkHandler()
		}
		return
	}

	// Handle new block event
	svc.Logger(context.TODO()).Info().
		Uint64("block", latestBlockNumber).
		Msg("New block event detected")
	// Update state
	svc.adjustPollDur(newBlock)
	svc.latestBlock = newBlock
	// Emit event
	if svc.newBlockHandler != nil {
		go svc.newBlockHandler(*newBlock)
	}
}

// checkNewBlockReceived checks if received "latest block number" is a new block.
func (svc *BlockNotifierSvc) checkNewBlockReceived(blockNumber uint64) bool {
	if svc.latestBlock == nil {
		return true
	}

	if svc.latestBlock.NumberU64() != blockNumber {
		return true
	}

	return false
}

// checkChainForked checks if the Ethereum chain was reordered.
func (svc *BlockNotifierSvc) checkChainForked(newBlock *types.Block) bool {
	if svc.latestBlock == nil {
		return false
	}

	// Check if the new block has a lower index
	if newBlock.NumberU64() <= svc.latestBlock.NumberU64() {
		return true
	}

	// Re-request the previous "latest block" (chain is forked if block wasn't found)
	latestBlockNew, err := svc.reader.GetBlock(context.Background(), svc.latestBlock.NumberU64())
	if err != nil {
		svc.Logger(context.TODO()).Error().
			Uint64("block", svc.latestBlock.NumberU64()).
			Err(err).
			Msg("Re-requesting the latest block")
		return false
	}
	if latestBlockNew == nil {
		return true
	}

	// Check if hashes are equal
	if !bytes.Equal(svc.latestBlock.Hash().Bytes(), latestBlockNew.Hash().Bytes()) {
		return true
	}

	return false
}

// adjustPollDur adjusts the polling timeout to lower the number of GetLatestBlockNumber() requests.
func (svc *BlockNotifierSvc) adjustPollDur(newBlock *types.Block) {
	if svc.latestBlock == nil {
		return
	}

	if newBlock.Time() <= svc.latestBlock.Time() {
		return
	}

	if svc.latestBlock.NumberU64()+1 == newBlock.NumberU64() {
		newPollDur := time.Duration(newBlock.Time()-svc.latestBlock.Time()+1) * time.Second
		svc.Logger(context.TODO()).Debug().
			Dur("duration", newPollDur).
			Msg("Polling duration adjusted")
		svc.pollDur = newPollDur
	}
}
