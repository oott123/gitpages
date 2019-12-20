package hash

import (
	"crypto/sha1"
	"fmt"
	"io"
)

func SHA1(in string) string {
	h := sha1.New()
	_, _ = io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}
