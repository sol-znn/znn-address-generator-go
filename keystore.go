package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"fmt"

	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

// BIP-44 derivation path: m/44'/73404'/{account}'
const znnCoinType = 73404
const hardenedOffset uint32 = 0x80000000

// slip10Key holds an ed25519 SLIP-0010 derived key and its chain code.
type slip10Key struct {
	key       [32]byte
	chainCode [32]byte
}

// slip10MasterKey derives the master key from a BIP-39 seed, per SLIP-0010.
func slip10MasterKey(seed []byte) slip10Key {
	mac := hmac.New(sha512.New, []byte("ed25519 seed"))
	mac.Write(seed)
	sum := mac.Sum(nil)

	var k slip10Key
	copy(k.key[:], sum[:32])
	copy(k.chainCode[:], sum[32:])
	return k
}

// slip10DeriveChild derives a hardened child key, per SLIP-0010.
func slip10DeriveChild(parent slip10Key, index uint32) slip10Key {
	data := make([]byte, 37)
	data[0] = 0x00
	copy(data[1:33], parent.key[:])
	binary.BigEndian.PutUint32(data[33:], hardenedOffset+index)

	mac := hmac.New(sha512.New, parent.chainCode[:])
	mac.Write(data)
	sum := mac.Sum(nil)

	var k slip10Key
	copy(k.key[:], sum[:32])
	copy(k.chainCode[:], sum[32:])
	return k
}

// deriveAccountKey derives the ed25519 seed for m/44'/73404'/{account}'.
func deriveAccountKey(seed []byte, account int) [32]byte {
	k := slip10MasterKey(seed)
	k = slip10DeriveChild(k, 44)
	k = slip10DeriveChild(k, znnCoinType)
	k = slip10DeriveChild(k, uint32(account))
	return k.key
}

// addressFromPublicKey builds a Zenon "z1q..." address from an ed25519
// public key, matching Address.fromPublicKey in znn_sdk_dart.
func addressFromPublicKey(publicKey ed25519.PublicKey) string {
	digest := sha3.Sum256(publicKey)
	core := make([]byte, 20)
	core[0] = 0x00 // Address.userByte
	copy(core[1:], digest[:19])
	return bech32Encode("z", core)
}

// keyStore mirrors the subset of znn_sdk_dart's KeyStore used by this tool.
type keyStore struct {
	mnemonic string
	seed     []byte
}

// newRandomKeyStore generates a random 256-bit entropy mnemonic and its seed.
func newRandomKeyStore() (*keyStore, error) {
	entropy := make([]byte, 32)
	if _, err := rand.Read(entropy); err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("generating mnemonic: %w", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	return &keyStore{mnemonic: mnemonic, seed: seed}, nil
}

// deriveAddressesByRange derives addresses [left, right) as in
// KeyStore.deriveAddressesByRange.
func (ks *keyStore) deriveAddressesByRange(left, right int) []string {
	addresses := make([]string, 0, right-left)
	for i := left; i < right; i++ {
		accountSeed := deriveAccountKey(ks.seed, i)
		priv := ed25519.NewKeyFromSeed(accountSeed[:])
		pub := priv.Public().(ed25519.PublicKey)
		addresses = append(addresses, addressFromPublicKey(pub))
	}
	return addresses
}
