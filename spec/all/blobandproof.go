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

package all

import (
	"encoding/json"
	"fmt"

	"github.com/pk910/dynamic-ssz/sszutils"

	"github.com/ethpandaops/go-eth-engine-client/spec"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
	"github.com/ethpandaops/go-eth-engine-client/spec/osaka"
	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// BlobAndProof is a fork-agnostic blob-and-proof. Cancun V1 (reused by
// prague) carries a single KZG proof; osaka V2 (reused by amsterdam)
// carries a list of cell proofs. The union holds both shapes — Proof for
// V1, Proofs for V2 — and the active Version selects which is meaningful.
type BlobAndProof struct {
	Version version.DataVersion

	Blob   cancun.Blob
	Proof  cancun.KZGProof   // cancun/prague (V1)
	Proofs []cancun.KZGProof // osaka/amsterdam (V2)
}

func (b *BlobAndProof) viewType() (any, error) {
	switch b.Version {
	case version.DataVersionCancun, version.DataVersionPrague:
		return (*cancun.BlobAndProof)(nil), nil
	case version.DataVersionOsaka,
		version.DataVersionAmsterdam,
		version.DataVersionBogota:
		return (*osaka.BlobAndProof)(nil), nil
	default:
		return nil, fmt.Errorf("BlobAndProof: unsupported version %d", b.Version)
	}
}

func (b *BlobAndProof) populateVersion(v version.DataVersion) {
	b.Version = v
}

// ToView returns a fresh fork-specific BlobAndProof populated from b.
func (b *BlobAndProof) ToView() (any, error) {
	return toViewByCopy(b)
}

// FromView populates b from a fork-specific BlobAndProof.
func (b *BlobAndProof) FromView(view any) error {
	if b.Version == version.DataVersionUnknown {
		switch view.(type) {
		case *cancun.BlobAndProof:
			b.Version = version.DataVersionCancun
		case *osaka.BlobAndProof:
			b.Version = version.DataVersionOsaka
		default:
			return fmt.Errorf("BlobAndProof: unsupported view type %T", view)
		}
	}

	return copyByName(view, b)
}

// ToVersioned converts b into a *spec.VersionedBlobAndProof.
func (b *BlobAndProof) ToVersioned() (*spec.VersionedBlobAndProof, error) {
	out := &spec.VersionedBlobAndProof{}
	if err := toVersioned(b.Version, b, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates b from src.
func (b *BlobAndProof) FromVersioned(src *spec.VersionedBlobAndProof) error {
	return fromVersioned(b, src)
}

// MarshalJSON implements json.Marshaler.
func (b *BlobAndProof) MarshalJSON() ([]byte, error) {
	view, err := b.ToView()
	if err != nil {
		return nil, err
	}

	return json.Marshal(view)
}

// UnmarshalJSON implements json.Unmarshaler. Version must be set first.
func (b *BlobAndProof) UnmarshalJSON(input []byte) error {
	inst, err := newViewInstance(b)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(input, inst); err != nil {
		return err
	}

	return copyByName(inst, b)
}

// MarshalSSZ marshals the blob-and-proof under the active fork view.
func (b *BlobAndProof) MarshalSSZ() ([]byte, error) {
	ds := globalDynSSZ()

	return marshalSSZDyn(b, ds, make([]byte, 0, sizeSSZDyn(b, ds)))
}

// MarshalSSZTo marshals the blob-and-proof under the active fork view, appending to buf.
func (b *BlobAndProof) MarshalSSZTo(buf []byte) ([]byte, error) {
	return marshalSSZDyn(b, globalDynSSZ(), buf)
}

// SizeSSZ returns the SSZ size of the blob-and-proof under the active fork view.
func (b *BlobAndProof) SizeSSZ() int {
	return sizeSSZDyn(b, globalDynSSZ())
}

// UnmarshalSSZ decodes the blob-and-proof under the active fork view.
// Version must be set first.
func (b *BlobAndProof) UnmarshalSSZ(buf []byte) error {
	return unmarshalSSZDyn(b, b, globalDynSSZ(), buf)
}

// HashTreeRoot computes the hash-tree-root under the active fork view.
func (b *BlobAndProof) HashTreeRoot() ([32]byte, error) {
	return globalDynSSZ().HashTreeRoot(b)
}

// HashTreeRootWith computes the hash-tree-root into hh under the active fork view.
func (b *BlobAndProof) HashTreeRootWith(hh sszutils.HashWalker) error {
	return hashTreeRootWithDyn(b, globalDynSSZ(), hh)
}

// String returns the JSON representation.
func (b *BlobAndProof) String() string {
	out, err := b.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(out)
}
