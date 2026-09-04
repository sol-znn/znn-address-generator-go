package main

// Standard Bech32 (BIP-173) encoding, matching the `bech32` Dart package
// used by znn_sdk_dart (final checksum XOR constant of 1, not bech32m).

const bech32Charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var bech32Generator = [5]uint32{
	0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3,
}

func bech32Polymod(values []byte) uint32 {
	chk := uint32(1)
	for _, v := range values {
		top := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ uint32(v)
		for i := 0; i < 5; i++ {
			if (top>>uint(i))&1 == 1 {
				chk ^= bech32Generator[i]
			}
		}
	}
	return chk
}

func bech32HrpExpand(hrp string) []byte {
	result := make([]byte, 0, len(hrp)*2+1)
	for _, c := range hrp {
		result = append(result, byte(c)>>5)
	}
	result = append(result, 0)
	for _, c := range hrp {
		result = append(result, byte(c)&31)
	}
	return result
}

func bech32CreateChecksum(hrp string, data []byte) []byte {
	values := append(bech32HrpExpand(hrp), data...)
	values = append(values, 0, 0, 0, 0, 0, 0)
	polymod := bech32Polymod(values) ^ 1

	checksum := make([]byte, 6)
	for i := range checksum {
		checksum[i] = byte((polymod >> uint(5*(5-i))) & 31)
	}
	return checksum
}

// convertBits regroups a slice of `from`-bit values into `to`-bit values,
// mirroring convertBech32Bits from znn_sdk_dart's bech32.dart.
func convertBits(data []byte, from, to uint, pad bool) []byte {
	acc := uint32(0)
	bits := uint(0)
	maxv := uint32(1<<to) - 1
	var result []byte

	for _, value := range data {
		v := uint32(value)
		acc = (acc << from) | v
		bits += from
		for bits >= to {
			bits -= to
			result = append(result, byte((acc>>bits)&maxv))
		}
	}

	if pad {
		if bits > 0 {
			result = append(result, byte((acc<<(to-bits))&maxv))
		}
	}

	return result
}

// bech32Encode encodes hrp + core bytes (8-bit) into a Bech32 address string.
func bech32Encode(hrp string, core []byte) string {
	data := convertBits(core, 8, 5, true)
	checksum := bech32CreateChecksum(hrp, data)
	combined := append(append([]byte{}, data...), checksum...)

	out := make([]byte, 0, len(hrp)+1+len(combined))
	out = append(out, hrp...)
	out = append(out, '1')
	for _, b := range combined {
		out = append(out, bech32Charset[b])
	}
	return string(out)
}
