# go-eth-engine-client

Go client types and codecs for the Ethereum Engine API — the consensus-layer
to execution-layer interface defined in
[ethereum/execution-apis](https://github.com/ethereum/execution-apis).

The library exposes one Go package per Engine-API fork (`paris`, `shanghai`,
`cancun`, `prague`, `osaka`, `amsterdam`) under `spec/`, plus a top-level
`spec/` package with fork-agnostic `Versioned*` wrappers. Type names follow
the execution-layer fork rather than the Engine API's `V1`/`V2`/`V3` suffix
— the parameter type to `engine_newPayloadV1` is `paris.NewPayloadRequest`,
the parameter to `engine_newPayloadV3` is `cancun.NewPayloadRequest`, and so
on.

Every spec type provides:

  * **JSON** marshaling in Engine API wire format (camelCase field names,
    hex-encoded byte strings, hex-encoded `QUANTITY` integers).
  * **SSZ** marshaling and hash-tree-root via
    [`pk910/dynamic-ssz`](https://github.com/pk910/dynamic-ssz), matching the
    SSZ transport defined in [execution-apis PR #764](https://github.com/ethereum/execution-apis/pull/764).

## Fork packages

| Go package   | Engine API spec |
| ------------ | --------------- |
| `paris`      | [paris.md](https://github.com/ethereum/execution-apis/blob/main/src/engine/paris.md)         |
| `shanghai`   | [shanghai.md](https://github.com/ethereum/execution-apis/blob/main/src/engine/shanghai.md)   |
| `cancun`     | [cancun.md](https://github.com/ethereum/execution-apis/blob/main/src/engine/cancun.md)       |
| `prague`     | [prague.md](https://github.com/ethereum/execution-apis/blob/main/src/engine/prague.md)       |
| `osaka`      | [osaka.md](https://github.com/ethereum/execution-apis/blob/main/src/engine/osaka.md)         |
| `amsterdam`  | [amsterdam.md](https://github.com/ethereum/execution-apis/blob/main/src/engine/amsterdam.md) |

## Status

This is the initial spec-types implementation. A higher-level client
(HTTP/JSON-RPC + SSZ REST transport with capability negotiation) will follow
in subsequent releases.

## Acknowledgements

The package layout and `Versioned*` wrapper pattern are inspired by
[Attestant's `go-eth2-client`](https://github.com/attestantio/go-eth2-client).
See [NOTICE](./NOTICE) for full attribution.

## License

Apache-2.0. See [LICENSE](./LICENSE).
