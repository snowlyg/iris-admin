package auth2

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
	"github.com/snowlyg/helper/dir"
)

var (
	sep    = []byte(".")
	pad    = []byte("=")
	padStr = string(pad)
)

// GetToken 雪花算法,支持分布式集群方式
func GetToken() (string, error) {
	v4 := uuid.NewV4()
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", fmt.Errorf("mutil: create token %w", err)
	}

	// 混入两个时间，防止并发token重复
	nodeBytes, _ := dir.Md5Byte(Base64Encode(node.Generate().Bytes()))
	uuidBytes, _ := dir.Md5Byte(Base64Encode(joinParts(Base64Encode(v4.Bytes()), []byte(nodeBytes))))
	token := joinParts(Base64Encode([]byte(uuidBytes)), Base64Encode([]byte(nodeBytes)))
	return string(Base64Encode([]byte(token))), nil
}

// joinParts
func joinParts(parts ...[]byte) []byte {
	return bytes.Join(parts, sep)
}

// Base64Encode
func Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.URLEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(buf, src)

	return bytes.TrimRight(buf, padStr) // JWT: no trailing '='.
}

// Base64Decode decodes "src" to jwt base64 url format.
// We could use the base64.RawURLEncoding but the below is a bit faster.
func Base64Decode(src []byte) ([]byte, error) {
	if n := len(src) % 4; n > 0 {
		// JWT: Because of no trailing '=' let's suffix it
		// with the correct number of those '=' before decoding.
		src = append(src, bytes.Repeat(pad, 4-n)...)
	}

	buf := make([]byte, base64.URLEncoding.DecodedLen(len(src)))
	n, err := base64.URLEncoding.Decode(buf, src)
	return buf[:n], err
}
