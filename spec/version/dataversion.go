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

// Package version defines the DataVersion enum used by the fork-agnostic
// Versioned* wrappers in spec/. Names match the execution-layer fork in
// which the corresponding Engine API structures were introduced.
package version

import (
	"fmt"
	"strings"
)

// DataVersion identifies the execution-layer fork whose Engine API spec a
// container conforms to.
type DataVersion uint64

const (
	// DataVersionUnknown is an unknown / unset data version.
	DataVersionUnknown DataVersion = iota
	// DataVersionParis is the initial Engine API release (the Merge).
	DataVersionParis
	// DataVersionShanghai adds withdrawals (ExecutionPayloadV2,
	// PayloadAttributesV2, ExecutionPayloadBodyV1).
	DataVersionShanghai
	// DataVersionCancun adds blob transactions and the parent beacon block
	// root (ExecutionPayloadV3, PayloadAttributesV3, BlobsBundleV1).
	DataVersionCancun
	// DataVersionPrague adds execution-layer triggered requests
	// (GetPayloadResponseV4 with execution_requests).
	DataVersionPrague
	// DataVersionOsaka switches to cell proofs (BlobsBundleV2,
	// BlobAndProofV2, GetPayloadResponseV5).
	DataVersionOsaka
	// DataVersionAmsterdam adds the block access list and slot number
	// (ExecutionPayloadV4, PayloadAttributesV4, ExecutionPayloadBodyV2,
	// GetPayloadResponseV6, BlobCellsAndProofsV1).
	DataVersionAmsterdam
	// DataVersionBogota adds inclusion lists (EIP-7805): PayloadAttributesV5
	// with inclusionListTransactions, engine_newPayloadV6 with a 5th
	// inclusionListTransactions parameter, and engine_getInclusionListV1.
	DataVersionBogota
)

var dataVersionStrings = [...]string{
	"unknown",
	"paris",
	"shanghai",
	"cancun",
	"prague",
	"osaka",
	"amsterdam",
	"bogota",
}

var dataVersionMap = map[string]DataVersion{
	`"paris"`:     DataVersionParis,
	`"shanghai"`:  DataVersionShanghai,
	`"cancun"`:    DataVersionCancun,
	`"prague"`:    DataVersionPrague,
	`"osaka"`:     DataVersionOsaka,
	`"amsterdam"`: DataVersionAmsterdam,
	`"bogota"`:    DataVersionBogota,
}

// String returns a string representation of the data version.
func (d DataVersion) String() string {
	if int(d) >= len(dataVersionStrings) {
		return "unknown"
	}

	return dataVersionStrings[d]
}

// MarshalJSON implements json.Marshaler.
func (d DataVersion) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", d.String()), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DataVersion) UnmarshalJSON(input []byte) error {
	lower := strings.ToLower(string(input))

	v, ok := dataVersionMap[lower]
	if !ok {
		return fmt.Errorf("unrecognised data version %s", string(input))
	}

	*d = v

	return nil
}

// DataVersionFromString turns a fork name into a DataVersion. It returns an
// error if the fork is not recognised.
func DataVersionFromString(fork string) (DataVersion, error) {
	var v DataVersion

	return v, v.UnmarshalJSON(fmt.Appendf(nil, "%q", fork))
}
