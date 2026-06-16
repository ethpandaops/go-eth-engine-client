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

package spec

import (
	"errors"

	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/bogota"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/paris"
	"github.com/ethpandaops/go-eth-engine-client/spec/shanghai"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// VersionedPayloadAttributes wraps the per-fork PayloadAttributes types.
// Prague and Osaka reuse cancun's PayloadAttributesV3. Bogota's V5 extends
// amsterdam's V4 with an `inclusionListTransactions` field.
type VersionedPayloadAttributes struct {
	Version version.DataVersion

	Paris     *paris.PayloadAttributes
	Shanghai  *shanghai.PayloadAttributes
	Cancun    *cancun.PayloadAttributes
	Prague    *cancun.PayloadAttributes
	Osaka     *cancun.PayloadAttributes
	Amsterdam *amsterdam.PayloadAttributes
	Bogota    *bogota.PayloadAttributes
}

// IsEmpty returns true if no attributes are set for the current version.
func (v *VersionedPayloadAttributes) IsEmpty() bool {
	switch v.Version {
	case version.DataVersionParis:
		return v.Paris == nil
	case version.DataVersionShanghai:
		return v.Shanghai == nil
	case version.DataVersionCancun:
		return v.Cancun == nil
	case version.DataVersionPrague:
		return v.Prague == nil
	case version.DataVersionOsaka:
		return v.Osaka == nil
	case version.DataVersionAmsterdam:
		return v.Amsterdam == nil
	case version.DataVersionBogota:
		return v.Bogota == nil
	default:
		return true
	}
}

// Timestamp returns the timestamp of the attributes.
func (v *VersionedPayloadAttributes) Timestamp() (uint64, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return 0, errors.New("no paris attributes")
		}

		return v.Paris.Timestamp, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return 0, errors.New("no shanghai attributes")
		}

		return v.Shanghai.Timestamp, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return 0, errors.New("no cancun attributes")
		}

		return v.Cancun.Timestamp, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return 0, errors.New("no prague attributes")
		}

		return v.Prague.Timestamp, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return 0, errors.New("no osaka attributes")
		}

		return v.Osaka.Timestamp, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return 0, errors.New("no amsterdam attributes")
		}

		return v.Amsterdam.Timestamp, nil
	case version.DataVersionBogota:
		if v.Bogota == nil {
			return 0, errors.New("no bogota attributes")
		}

		return v.Bogota.Timestamp, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// SuggestedFeeRecipient returns the suggested fee recipient.
func (v *VersionedPayloadAttributes) SuggestedFeeRecipient() (paris.Address, error) {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return paris.Address{}, errors.New("no paris attributes")
		}

		return v.Paris.SuggestedFeeRecipient, nil
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return paris.Address{}, errors.New("no shanghai attributes")
		}

		return v.Shanghai.SuggestedFeeRecipient, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return paris.Address{}, errors.New("no cancun attributes")
		}

		return v.Cancun.SuggestedFeeRecipient, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return paris.Address{}, errors.New("no prague attributes")
		}

		return v.Prague.SuggestedFeeRecipient, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return paris.Address{}, errors.New("no osaka attributes")
		}

		return v.Osaka.SuggestedFeeRecipient, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return paris.Address{}, errors.New("no amsterdam attributes")
		}

		return v.Amsterdam.SuggestedFeeRecipient, nil
	case version.DataVersionBogota:
		if v.Bogota == nil {
			return paris.Address{}, errors.New("no bogota attributes")
		}

		return v.Bogota.SuggestedFeeRecipient, nil
	default:
		return paris.Address{}, errors.New("unknown version")
	}
}

// Withdrawals returns the withdrawals list (shanghai+).
func (v *VersionedPayloadAttributes) Withdrawals() ([]*shanghai.Withdrawal, error) {
	switch v.Version {
	case version.DataVersionParis:
		return nil, errors.New("no withdrawals in paris")
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return nil, errors.New("no shanghai attributes")
		}

		return v.Shanghai.Withdrawals, nil
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return nil, errors.New("no cancun attributes")
		}

		return v.Cancun.Withdrawals, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return nil, errors.New("no prague attributes")
		}

		return v.Prague.Withdrawals, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return nil, errors.New("no osaka attributes")
		}

		return v.Osaka.Withdrawals, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return nil, errors.New("no amsterdam attributes")
		}

		return v.Amsterdam.Withdrawals, nil
	case version.DataVersionBogota:
		if v.Bogota == nil {
			return nil, errors.New("no bogota attributes")
		}

		return v.Bogota.Withdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ParentBeaconBlockRoot returns the parent beacon block root (cancun+).
func (v *VersionedPayloadAttributes) ParentBeaconBlockRoot() (paris.Hash32, error) {
	switch v.Version {
	case version.DataVersionParis, version.DataVersionShanghai:
		return paris.Hash32{}, errors.New("no parent beacon block root before cancun")
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return paris.Hash32{}, errors.New("no cancun attributes")
		}

		return v.Cancun.ParentBeaconBlockRoot, nil
	case version.DataVersionPrague:
		if v.Prague == nil {
			return paris.Hash32{}, errors.New("no prague attributes")
		}

		return v.Prague.ParentBeaconBlockRoot, nil
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return paris.Hash32{}, errors.New("no osaka attributes")
		}

		return v.Osaka.ParentBeaconBlockRoot, nil
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return paris.Hash32{}, errors.New("no amsterdam attributes")
		}

		return v.Amsterdam.ParentBeaconBlockRoot, nil
	case version.DataVersionBogota:
		if v.Bogota == nil {
			return paris.Hash32{}, errors.New("no bogota attributes")
		}

		return v.Bogota.ParentBeaconBlockRoot, nil
	default:
		return paris.Hash32{}, errors.New("unknown version")
	}
}

// String returns a JSON representation of the active attributes.
func (v *VersionedPayloadAttributes) String() string {
	switch v.Version {
	case version.DataVersionParis:
		if v.Paris == nil {
			return ""
		}

		return v.Paris.String()
	case version.DataVersionShanghai:
		if v.Shanghai == nil {
			return ""
		}

		return v.Shanghai.String()
	case version.DataVersionCancun:
		if v.Cancun == nil {
			return ""
		}

		return v.Cancun.String()
	case version.DataVersionPrague:
		if v.Prague == nil {
			return ""
		}

		return v.Prague.String()
	case version.DataVersionOsaka:
		if v.Osaka == nil {
			return ""
		}

		return v.Osaka.String()
	case version.DataVersionAmsterdam:
		if v.Amsterdam == nil {
			return ""
		}

		return v.Amsterdam.String()
	case version.DataVersionBogota:
		if v.Bogota == nil {
			return ""
		}

		return v.Bogota.String()
	default:
		return "unknown version"
	}
}
