# znn-address-generator-go

Concurrently generate Zenon (NoM) mnemonics and addresses, then save them to disk.
Optionally search for addresses with a specific prefix and/or suffix.

A Go port of [znn-address-generator](https://github.com/Sol-Sanctum/znn-address-generator), with the
same BIP-39 mnemonics and `m/44'/73404'/{account}'` derivation path as `znn_sdk_dart`.

## Build

```bash
git clone https://github.com/sol-znn/znn-address-generator-go.git && cd znn-address-generator-go
go build -o znn-address-generator-go .
./znn-address-generator-go -n 1000
```

## Releases

Download compiled binaries for Windows/Linux/macOS from the [latest release](https://github.com/sol-znn/znn-address-generator-go/releases/latest).

## Options

```text
Zenon Address Generator
Usage: znn-address-generator-go [OPTIONS]

Options
  -t, --threads        number of cpu threads (default 4)
  -n, --mnemonics       number of mnemonics to generate (default 100)
  -a, --addresses       number of addresses per mnemonic (default 5)
  -o, --output          output file name (default ./results/results-<time>.txt)
  -p, --prefix          save addresses that match a specific prefix
  -s, --suffix          save addresses that match a specific suffix
  -i, --infinite        generate infinite mnemonics
      --all             save all generated mnemonics, even if they don't match
  -v, --verbose         display mnemonics as they're being generated
  -h, --help            Displays help information
      --version         Displays version information
```

## Usage

You can mix and match whatever settings you want.

### Generate 1000 mnemonics with 8 threads

```bash
znn-address-generator-go -n 1000 -t 8
```

### Find an address with suffix 321

```bash
znn-address-generator-go -s 321
```

### Find an address with prefix zzz and suffix 321

```bash
znn-address-generator-go -p zzz -s 321
```

### Find an address with suffix 321 and save all generated mnemonics to a separate file

```bash
znn-address-generator-go -s 321 --all
```

### Keep searching for addresses with suffix 321 until stopped

```bash
znn-address-generator-go -s 321 --infinite
# Ctrl-C to end execution
```
