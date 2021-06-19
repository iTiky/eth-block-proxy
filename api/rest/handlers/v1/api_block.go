package v1

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/itiky/eth-block-proxy/api/rest/handlers/v1/schema"
)

// BlockGet is a handler for: GET /block/{blockNumber}.
// nolint: errcheck
func (a ApiHandlers) BlockGet(w http.ResponseWriter, r *http.Request) {
	block := r.Context().Value(CtxKeyBlock).(*types.Block)
	blockBz, err := schema.NewBlockResponse(block)
	if err != nil {
		render.Render(w, r, NewErrResponseRender(err))
		return
	}

	if err := render.Render(w, r, blockBz); err != nil {
		render.Render(w, r, NewErrResponseRender(err))
		return
	}
}

// BlockGetTx is a handler for: GET /block/{blockNumber}/txs/{txHash}.
// nolint: errcheck
func (a ApiHandlers) BlockGetTx(w http.ResponseWriter, r *http.Request) {
	block := r.Context().Value(CtxKeyBlock).(*types.Block)
	txHashParam := strings.TrimPrefix(chi.URLParam(r, UrlParamTxHash), "0x")

	txHashRaw, err := hex.DecodeString(txHashParam)
	if err != nil {
		render.Render(w, r, NewErrResponseInvalidRequest(fmt.Errorf("{txHash} HEX string parsing: %w", err)))
		return
	}

	txBz := schema.NewTxResponse(block, common.BytesToHash(txHashRaw))
	if txBz == nil {
		render.Render(w, r, NewErrResponseNotFound())
		return
	}

	if err := render.Render(w, r, txBz); err != nil {
		render.Render(w, r, NewErrResponseRender(err))
		return
	}
}

// BlockCtx middleware loads a types.Block object by blockNumber from the URL parameters and store it to the request context.
// nolint: errcheck
func (a ApiHandlers) BlockCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		blockParam := chi.URLParam(r, UrlParamBlock)

		blockIdx := uint64(0)
		if blockParam != "latest" {
			v, err := strconv.ParseUint(blockParam, 10, 64)
			if err != nil {
				render.Render(w, r, NewErrResponseInvalidRequest(err))
				return
			}
			blockIdx = v
		}

		block, err := a.blockCacheSvc.GetBlock(ctx, blockIdx)
		if err != nil {
			render.Render(w, r, NewErrResponseInternal(err))
			return
		}
		if block == nil {
			render.Render(w, r, NewErrResponseNotFound())
			return
		}

		ctx = context.WithValue(ctx, CtxKeyBlock, block)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
