package cryptography

import (
	"bytes"
	"math"
)

// pad.go implements pkcs7 padding functions

func pad(in []byte, blocksize int) []byte {
	if blocksize <= 1 || blocksize >= 256 {
		return in
	}

	m := math.Max(float64(blocksize-len(in)), 0)
	return append(in, bytes.Repeat([]byte{byte(m)}, int(m))...)
}

func safeUnpad(in []byte, blocksize int) []byte {
	l := len(in)
	if l != blocksize || l%blocksize != 0 {
		return in
	}

	lastByteInteger := int(in[l-1]) // aka padding size
	padding := bytes.Repeat([]byte{byte(lastByteInteger)}, lastByteInteger)

	if lastByteInteger > blocksize || !bytes.HasSuffix(in, padding) || lastByteInteger == 0 {
		return in
	}

	return in[:blocksize-lastByteInteger]
}
