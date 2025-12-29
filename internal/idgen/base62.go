package idgen

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = int64(len(alphabet))

func EncodeBase62(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	var buf [11]byte // enough for 64-bit base62
	i := len(buf)

	for n > 0 {
		i--
		buf[i] = alphabet[n%base]
		n /= base
	}

	return string(buf[i:])
}

