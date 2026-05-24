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

package paris

// Byte-length constants for fixed-size Engine API types.
const (
	// Hash32Length is the number of bytes in a Hash32.
	Hash32Length = 32
	// AddressLength is the number of bytes in an Address.
	AddressLength = 20
	// BloomLength is the number of bytes in a Bloom filter.
	BloomLength = 256
	// PayloadIDLength is the number of bytes in a PayloadID.
	PayloadIDLength = 8
)

// SSZ list-size limits used in Engine API containers introduced in Paris.
const (
	// MaxBytesPerTransaction is the maximum number of bytes in a single
	// transaction (2**30, EIP-4844).
	MaxBytesPerTransaction = 1 << 30
	// MaxTransactionsPerPayload is the maximum number of transactions in
	// an ExecutionPayload (2**20, Bellatrix).
	MaxTransactionsPerPayload = 1 << 20
	// MaxExtraDataBytes is the maximum length of the extraData field
	// (2**5, Bellatrix).
	MaxExtraDataBytes = 32
	// MaxErrorMessageLength is the maximum length of a PayloadStatus
	// validation error string when encoded as SSZ.
	MaxErrorMessageLength = 1024
)
