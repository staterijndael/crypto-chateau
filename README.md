# Chateau

Chateau RPC protocol library and generator

## Installation

```bash
go install github.com/oringik/crypto-chateau/cmd/chateau-gen@latest
```

## Example generation

```bash
chateau-gen -language=go -chateau_file=examples/reverse/contract/reverse.chateau -codegen_output=examples/reverse/codegen
```