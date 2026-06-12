// Copyright © 2026 ethPandaOps.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package engine defines the fork-agnostic interfaces for an Ethereum Engine
// API client. The concrete JSON-RPC implementation lives in the http
// subpackage.
//
// Every method provider exposes two entry points:
//
//   - a default method that operates on the memory-efficient spec.Versioned*
//     wrappers (only the active fork's struct is allocated), and
//   - an *Agnostic counterpart that operates on the flat spec/all union
//     types, which are more convenient but allocate every fork field.
//
// Prefer the default (Versioned) methods; reach for the Agnostic variants
// only when a flat union value is already in hand.
package engine

import (
	"context"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/all"
	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/identification"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// Service is the minimal interface implemented by every Engine API client.
type Service interface {
	// Address returns the address provided to the client.
	Address() string
}

// NewPayloadProvider is the interface for submitting execution payloads
// (engine_newPayloadV1..V5, selected by the request's Version).
type NewPayloadProvider interface {
	// NewPayload submits an execution payload for validation.
	NewPayload(
		ctx context.Context,
		request *spec.VersionedNewPayloadRequest,
	) (*paris.PayloadStatus, error)

	// NewPayloadAgnostic is the fork-agnostic-typed variant of NewPayload.
	NewPayloadAgnostic(
		ctx context.Context,
		request *all.NewPayloadRequest,
	) (*paris.PayloadStatus, error)
}

// ForkchoiceUpdatedProvider is the interface for updating fork choice
// (engine_forkchoiceUpdatedV1..V4, selected by the request's Version).
type ForkchoiceUpdatedProvider interface {
	// ForkchoiceUpdated updates the fork choice and optionally starts a
	// payload build.
	ForkchoiceUpdated(
		ctx context.Context,
		request *spec.VersionedForkchoiceUpdatedRequest,
	) (*paris.ForkchoiceUpdatedResponse, error)

	// ForkchoiceUpdatedAgnostic is the fork-agnostic-typed variant of
	// ForkchoiceUpdated.
	ForkchoiceUpdatedAgnostic(
		ctx context.Context,
		request *all.ForkchoiceUpdatedRequest,
	) (*paris.ForkchoiceUpdatedResponse, error)
}

// GetPayloadProvider is the interface for retrieving built payloads
// (engine_getPayloadV2..V6, selected by the supplied version). Paris's
// engine_getPayloadV1 returns a bare execution payload with no response
// container and is not covered here.
type GetPayloadProvider interface {
	// GetPayload retrieves a previously-initiated payload build.
	GetPayload(
		ctx context.Context,
		dataVersion version.DataVersion,
		payloadID paris.PayloadID,
	) (*spec.VersionedGetPayloadResponse, error)

	// GetPayloadAgnostic is the fork-agnostic-typed variant of GetPayload.
	GetPayloadAgnostic(
		ctx context.Context,
		dataVersion version.DataVersion,
		payloadID paris.PayloadID,
	) (*all.GetPayloadResponse, error)
}

// GetBlobsProvider is the interface for fetching blobs from the EL blob
// pool (engine_getBlobsV1..V4, selected by the supplied version).
type GetBlobsProvider interface {
	// GetBlobs fetches blobs all-or-nothing (engine_getBlobsV1 / V2).
	GetBlobs(
		ctx context.Context,
		request *spec.VersionedGetBlobsRequest,
	) (*spec.VersionedGetBlobsResponse, error)

	// GetBlobsAgnostic is the fork-agnostic-typed variant of GetBlobs.
	GetBlobsAgnostic(
		ctx context.Context,
		request *all.GetBlobsRequest,
	) (*all.GetBlobsResponse, error)

	// GetBlobsPartial fetches blobs with per-blob partial responses
	// (engine_getBlobsV3, osaka only): missing blobs are nil entries rather
	// than failing the whole call. Single-fork, so no Versioned/Agnostic
	// split.
	GetBlobsPartial(
		ctx context.Context,
		request *osaka.GetBlobsPartialRequest,
	) (*osaka.GetBlobsPartialResponse, error)

	// GetBlobsCells fetches partial blob cell matrices (engine_getBlobsV4,
	// amsterdam only). It has no Versioned/Agnostic split because the method
	// exists in a single fork. A nil response slice indicates the EL is
	// syncing or unable to serve blob-pool data.
	GetBlobsCells(
		ctx context.Context,
		request *amsterdam.GetBlobsCellsRequest,
	) (*amsterdam.GetBlobsCellsResponse, error)
}

// GetPayloadBodiesProvider is the interface for retrieving execution
// payload bodies (engine_getPayloadBodiesByHashV1/V2 and ByRangeV1/V2,
// selected by the supplied version). Result entries are nil for blocks the
// EL does not have.
type GetPayloadBodiesProvider interface {
	// GetPayloadBodiesByHash retrieves bodies for the given block hashes.
	GetPayloadBodiesByHash(
		ctx context.Context,
		dataVersion version.DataVersion,
		blockHashes []paris.Hash32,
	) ([]*spec.VersionedExecutionPayloadBody, error)

	// GetPayloadBodiesByHashAgnostic is the fork-agnostic-typed variant.
	GetPayloadBodiesByHashAgnostic(
		ctx context.Context,
		dataVersion version.DataVersion,
		blockHashes []paris.Hash32,
	) ([]*all.ExecutionPayloadBody, error)

	// GetPayloadBodiesByRange retrieves bodies for a block-number range.
	GetPayloadBodiesByRange(
		ctx context.Context,
		dataVersion version.DataVersion,
		start uint64,
		count uint64,
	) ([]*spec.VersionedExecutionPayloadBody, error)

	// GetPayloadBodiesByRangeAgnostic is the fork-agnostic-typed variant.
	GetPayloadBodiesByRangeAgnostic(
		ctx context.Context,
		dataVersion version.DataVersion,
		start uint64,
		count uint64,
	) ([]*all.ExecutionPayloadBody, error)
}

// CapabilitiesProvider is the interface for exchanging supported methods.
type CapabilitiesProvider interface {
	// ExchangeCapabilities advertises the consensus client's supported
	// methods and returns the execution client's supported methods
	// (engine_exchangeCapabilities).
	ExchangeCapabilities(ctx context.Context, supported []string) ([]string, error)
}

// ClientVersionProvider is the interface for exchanging client identity
// (engine_getClientVersionV1).
type ClientVersionProvider interface {
	// ClientVersion sends the consensus client's version and returns the
	// execution client's version information.
	ClientVersion(
		ctx context.Context,
		clientVersion *identification.ClientVersion,
	) ([]*identification.ClientVersion, error)
}
