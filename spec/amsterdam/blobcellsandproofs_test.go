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

package amsterdam_test

import (
	"testing"

	"github.com/ethpandaops/go-eth-engine-client/spec/amsterdam"
	"github.com/ethpandaops/go-eth-engine-client/spec/cancun"
)

// TestBlobCellsAndProofsSSZNullable checks that the defined-by-us amsterdam
// getBlobsV4 container round-trips through SSZ, including null cell/proof
// entries encoded as List[T, 1] (the nullable convention).
func TestBlobCellsAndProofsSSZNullable(t *testing.T) {
	cell := &amsterdam.BlobCell{}
	cell[0] = 0xaa
	cell[amsterdam.CellLength-1] = 0xbb

	proof := &cancun.KZGProof{}
	proof[0] = 0xcc

	orig := &amsterdam.BlobCellsAndProofs{
		// [present, null, present]
		BlobCells: []*amsterdam.BlobCell{cell, nil, cell},
		Proofs:    []*cancun.KZGProof{proof, nil, proof},
	}

	encoded, err := orig.MarshalSSZ()
	if err != nil {
		t.Fatalf("MarshalSSZ: %v", err)
	}

	var decoded amsterdam.BlobCellsAndProofs
	if err := decoded.UnmarshalSSZ(encoded); err != nil {
		t.Fatalf("UnmarshalSSZ: %v", err)
	}

	if len(decoded.BlobCells) != 3 || len(decoded.Proofs) != 3 {
		t.Fatalf("length mismatch: cells=%d proofs=%d", len(decoded.BlobCells), len(decoded.Proofs))
	}

	if decoded.BlobCells[1] != nil {
		t.Fatalf("expected null cell at index 1")
	}

	if decoded.BlobCells[0] == nil || *decoded.BlobCells[0] != *cell {
		t.Fatalf("cell 0 roundtrip mismatch")
	}

	if decoded.Proofs[1] != nil {
		t.Fatalf("expected null proof at index 1")
	}

	if decoded.Proofs[0] == nil || *decoded.Proofs[0] != *proof {
		t.Fatalf("proof 0 roundtrip mismatch")
	}
}
