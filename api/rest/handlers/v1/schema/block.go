package schema

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/core/types"
)

// BlockResponse defines types.Block HTTP response.
type BlockResponse map[string]interface{}

// Render implements render.Renderer interface.
func (e BlockResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NewBlockResponse builds a new serializable BlockResponse object.
func NewBlockResponse(block *types.Block) (BlockResponse, error) {
	resp, err := RPCMarshalBlock(block, true, true)
	if err != nil {
		return nil, fmt.Errorf("RPCMarshalBlock: %w", err)
	}

	return resp, nil
}
