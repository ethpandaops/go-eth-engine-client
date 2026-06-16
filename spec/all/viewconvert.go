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

// Package all provides fork-agnostic union types for the Engine API spec
// structures. Each union type carries a Version plus the union of fields
// across every fork that defines the structure, and converts to/from the
// fork-specific "view" types in the per-fork spec packages and the
// spec.Versioned* wrappers.
//
// Unlike go-eth2-client's all package — which uses dynssz "view"
// codegen to marshal a single struct under multiple SSZ schemas — these
// union types delegate SSZ and JSON to the concrete fork view produced by
// ToView/FromView. The public API (ToView/FromView/ToVersioned/
// FromVersioned + Marshal/UnmarshalSSZ/JSON) matches; only the marshaling
// internals differ (a reflective field-copy rather than generated view
// marshalers).
package all

import (
	"errors"
	"fmt"
	"reflect"

	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"

	"github.com/ethpandaops/go-eth-engine-client/spec/version"
)

// viewProvider maps a union type's active Version to the fork-specific view
// type (a typed nil pointer used purely for its dynamic type).
type viewProvider interface {
	viewType() (any, error)
}

// viewer is implemented by a union type's ToView method.
type viewer interface {
	ToView() (any, error)
}

// fromViewer is implemented by a union type's FromView method.
type fromViewer interface {
	FromView(view any) error
}

// versionSetter lets fromVersioned pin the authoritative Version on the
// destination before FromView runs, so a view type shared by multiple
// versions (e.g. cancun's ExecutionPayload, reused by prague + osaka) does
// not get its version downgraded by FromView's type inference.
type versionSetter interface {
	populateVersion(v version.DataVersion)
}

// versionFieldName maps a DataVersion to the per-fork field name on a
// spec.Versioned* struct.
func versionFieldName(v version.DataVersion) (string, error) {
	switch v {
	case version.DataVersionParis:
		return "Paris", nil
	case version.DataVersionShanghai:
		return "Shanghai", nil
	case version.DataVersionCancun:
		return "Cancun", nil
	case version.DataVersionPrague:
		return "Prague", nil
	case version.DataVersionOsaka:
		return "Osaka", nil
	case version.DataVersionAmsterdam:
		return "Amsterdam", nil
	case version.DataVersionBogota:
		return "Bogota", nil
	default:
		return "", fmt.Errorf("unsupported version %d", v)
	}
}

// newViewInstance allocates a fresh instance of the fork-specific view type
// indicated by p.viewType().
func newViewInstance(p viewProvider) (any, error) {
	view, err := p.viewType()
	if err != nil {
		return nil, err
	}

	return reflect.New(reflect.TypeOf(view).Elem()).Interface(), nil
}

// toVersioned populates a spec.Versioned* struct from a union source. It
// calls src.ToView to produce the fork-specific view and assigns it to the
// dst field whose name matches srcVersion, also setting dst.Version.
func toVersioned(srcVersion version.DataVersion, src viewer, dst any) error {
	view, err := src.ToView()
	if err != nil {
		return err
	}

	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Pointer || dv.IsNil() {
		return errors.New("toVersioned: dst must be a non-nil pointer")
	}

	dv = dv.Elem()

	versionField := dv.FieldByName("Version")
	if !versionField.IsValid() {
		return fmt.Errorf("toVersioned: %T has no Version field", dst)
	}

	versionField.Set(reflect.ValueOf(srcVersion))

	fieldName, err := versionFieldName(srcVersion)
	if err != nil {
		return fmt.Errorf("toVersioned: %w", err)
	}

	f := dv.FieldByName(fieldName)
	if !f.IsValid() {
		return fmt.Errorf("toVersioned: %T has no %s field", dst, fieldName)
	}

	rv := reflect.ValueOf(view)
	if !rv.Type().AssignableTo(f.Type()) {
		return fmt.Errorf("toVersioned: view type %T not assignable to %s field of type %s",
			view, fieldName, f.Type())
	}

	f.Set(rv)

	return nil
}

// fromVersioned populates a union destination from a spec.Versioned* struct
// by extracting the field matching src.Version and feeding it to
// dst.FromView.
func fromVersioned(dst fromViewer, src any) error {
	sv := reflect.ValueOf(src)
	if sv.Kind() != reflect.Pointer || sv.IsNil() {
		return errors.New("fromVersioned: src must be a non-nil pointer")
	}

	sv = sv.Elem()

	versionField := sv.FieldByName("Version")
	if !versionField.IsValid() {
		return fmt.Errorf("fromVersioned: %T has no Version field", src)
	}

	v, ok := versionField.Interface().(version.DataVersion)
	if !ok {
		return fmt.Errorf("fromVersioned: Version field on %T is %T not version.DataVersion",
			src, versionField.Interface())
	}

	fieldName, err := versionFieldName(v)
	if err != nil {
		return fmt.Errorf("fromVersioned: %w", err)
	}

	f := sv.FieldByName(fieldName)
	if !f.IsValid() {
		return fmt.Errorf("fromVersioned: %T has no %s field", src, fieldName)
	}

	if f.Kind() == reflect.Pointer && f.IsNil() {
		return fmt.Errorf("fromVersioned: %T.%s is nil for Version=%s", src, fieldName, v)
	}

	if vs, ok := dst.(versionSetter); ok {
		vs.populateVersion(v)
	}

	return dst.FromView(f.Interface())
}

// copyByName copies exported fields from src to dst by matching field names.
// Both must be pointers to structs (or structs). Directly-assignable fields
// are assigned; pointer/slice fields with differing element types are copied
// recursively. Fields present on one side but not the other are skipped.
func copyByName(src, dst any) error {
	sv := reflect.ValueOf(src)
	dv := reflect.ValueOf(dst)

	if sv.Kind() == reflect.Pointer {
		if sv.IsNil() {
			return nil
		}

		sv = sv.Elem()
	}

	if dv.Kind() == reflect.Pointer {
		if dv.IsNil() {
			return errors.New("copyByName: destination is a nil pointer")
		}

		dv = dv.Elem()
	}

	if sv.Kind() != reflect.Struct || dv.Kind() != reflect.Struct {
		return fmt.Errorf("copyByName: src kind=%s, dst kind=%s; both must be structs", sv.Kind(), dv.Kind())
	}

	for i := range dv.NumField() {
		df := dv.Type().Field(i)
		if !df.IsExported() {
			continue
		}

		sf := sv.FieldByName(df.Name)
		if !sf.IsValid() {
			continue
		}

		if err := copyValue(sf, dv.Field(i)); err != nil {
			return fmt.Errorf("field %s: %w", df.Name, err)
		}
	}

	return nil
}

// copyValue copies a single field, recursing into pointers and slices when
// the source and destination types differ.
func copyValue(src, dst reflect.Value) error {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)

		return nil
	}

	switch src.Kind() {
	case reflect.Pointer:
		if dst.Kind() != reflect.Pointer {
			return fmt.Errorf("incompatible kinds %s -> %s", src.Kind(), dst.Kind())
		}

		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))

			return nil
		}

		newDst := reflect.New(dst.Type().Elem())
		if err := copyByName(src.Interface(), newDst.Interface()); err != nil {
			return err
		}

		dst.Set(newDst)

		return nil
	case reflect.Slice:
		if dst.Kind() != reflect.Slice {
			return fmt.Errorf("incompatible kinds %s -> %s", src.Kind(), dst.Kind())
		}

		if src.IsNil() {
			dst.Set(reflect.Zero(dst.Type()))

			return nil
		}

		n := src.Len()
		newDst := reflect.MakeSlice(dst.Type(), n, n)

		for i := range n {
			if err := copyValue(src.Index(i), newDst.Index(i)); err != nil {
				return err
			}
		}

		dst.Set(newDst)

		return nil
	default:
		return fmt.Errorf("cannot copy %s to %s", src.Type(), dst.Type())
	}
}

// toViewByCopy allocates a fresh fork-specific view for src's active Version
// and field-copies src into it.
func toViewByCopy(src viewProvider) (any, error) {
	inst, err := newViewInstance(src)
	if err != nil {
		return nil, err
	}

	if err := copyByName(src, inst); err != nil {
		return nil, err
	}

	return inst, nil
}

// The four helpers below mirror go-eth2-client's spec/all SSZ plumbing.
// They accept the statically-implemented viewProvider and assert to the
// generated dynamic-view interfaces (sszutils.DynamicView*) at runtime via
// any(p). The runtime assertion is what breaks the chicken-and-egg between
// these hand-written methods and the dynssz-gen output: the package
// compiles before generation (the assertion simply fails with a clear
// "generated SSZ code missing" error), and once the *_ssz.go files exist
// the assertions succeed and dispatch to the per-view fast paths.

// marshalSSZDyn marshals p under the SSZ schema of its active Version's view.
func marshalSSZDyn(p viewProvider, ds sszutils.DynamicSpecs, buf []byte) ([]byte, error) {
	view, err := p.viewType()
	if err != nil {
		return nil, err
	}

	m, ok := any(p).(sszutils.DynamicViewMarshaler)
	if !ok {
		return nil, fmt.Errorf("%T: generated SSZ code missing", p)
	}

	fn := m.MarshalSSZDynView(view)
	if fn == nil {
		return nil, fmt.Errorf("%T: no view marshaler for the active version", p)
	}

	return fn(ds, buf)
}

// sizeSSZDyn returns the SSZ size of p under its active Version's view.
func sizeSSZDyn(p viewProvider, ds sszutils.DynamicSpecs) int {
	view, err := p.viewType()
	if err != nil {
		return 0
	}

	s, ok := any(p).(sszutils.DynamicViewSizer)
	if !ok {
		return 0
	}

	fn := s.SizeSSZDynView(view)
	if fn == nil {
		return 0
	}

	return fn(ds)
}

// unmarshalSSZDyn decodes buf into p under its active Version's view and
// propagates the version to any nested union children.
func unmarshalSSZDyn(p viewProvider, vs versionSetter, ds sszutils.DynamicSpecs, buf []byte) error {
	view, err := p.viewType()
	if err != nil {
		return err
	}

	u, ok := any(p).(sszutils.DynamicViewUnmarshaler)
	if !ok {
		return fmt.Errorf("%T: generated SSZ code missing", p)
	}

	fn := u.UnmarshalSSZDynView(view)
	if fn == nil {
		return fmt.Errorf("%T: no view unmarshaler for the active version", p)
	}

	if err := fn(ds, buf); err != nil {
		return err
	}

	vs.populateVersion(currentVersion(p))

	return nil
}

// hashTreeRootWithDyn writes p's hash-tree-root into hh under its active
// Version's view.
func hashTreeRootWithDyn(p viewProvider, ds sszutils.DynamicSpecs, hh sszutils.HashWalker) error {
	view, err := p.viewType()
	if err != nil {
		return err
	}

	h, ok := any(p).(sszutils.DynamicViewHashRoot)
	if !ok {
		return fmt.Errorf("%T: generated SSZ code missing", p)
	}

	fn := h.HashTreeRootWithDynView(view)
	if fn == nil {
		return fmt.Errorf("%T: no view hasher for the active version", p)
	}

	return fn(ds, hh)
}

// currentVersion reads the Version field of a union via reflection. Used by
// unmarshalSSZDyn to re-propagate the version after the generated decoder
// has populated nested children.
func currentVersion(p any) version.DataVersion {
	rv := reflect.ValueOf(p)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	f := rv.FieldByName("Version")
	if !f.IsValid() {
		return version.DataVersionUnknown
	}

	v, _ := f.Interface().(version.DataVersion)

	return v
}

// globalDynSSZ returns the process-global dynamic-ssz instance used by the
// standard (non-Dyn) SSZ entrypoints.
func globalDynSSZ() *dynssz.DynSsz {
	return dynssz.GetGlobalDynSsz()
}
