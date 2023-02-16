package cryptography

import (
	"bytes"
	"math"
)

// pad.go implements pkcs7 padding functions

func pad(in []byte, blocksize int) []byte {
	m := math.Max(float64(blocksize-len(in)), 0)
	return append(in, bytes.Repeat([]byte{byte(m)}, int(m))...)
}

func safeUnpad(in []byte, blocksize int) []byte {
	if len(in) != blocksize {
		return in
	}

	lastByteInteger := int(in[len(in)-1])
	padding := bytes.Repeat([]byte{byte(lastByteInteger)}, lastByteInteger)

	if lastByteInteger > len(in) {
		return in
	} else if !bytes.Contains(in, padding) {
		return in
	}

	return in[:blocksize-lastByteInteger]
}
