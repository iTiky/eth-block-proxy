package schema

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxResponse defines types.Transaction HTTP response.
type TxResponse RPCTransaction

// Render implements render.Renderer interface.
func (e *TxResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewTxResponse builds a new serializable TxResponse object.
func NewTxResponse(block *types.Block, txHash common.Hash) *TxResponse {
	return (*TxResponse)(newRPCTransactionFromBlockHash(block, txHash))
}
