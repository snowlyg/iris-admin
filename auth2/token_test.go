package auth2

import (
	"bytes"
	"testing"

	"github.com/snowlyg/helper/arr"
)

func TestGetToken(t *testing.T) {
	token, err := getToken()
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Generate token is fail.")
	}
	if token1, err := getToken(); err != nil {
		t.Error(err)
	} else if token == "" {
		t.Error("Generate token is fail.")
	} else if token == token1 {
		t.Errorf("token[%s] token1[%s] is repeat", token, token1)
	}
}

func TestJoinParts(t *testing.T) {
	afterJoin := joinParts([]byte("header"), []byte("footer"))
	want := []byte("header.footer")
	if bytes.Compare(afterJoin, want) > 0 {
		t.Errorf("Join parts want %s but get %s", string(want), string(afterJoin))
	}
}

func TestBase64Encode(t *testing.T) {
	want := []byte("header")
	baseEncode := Base64Encode(want)
	afterDecode, err := Base64Decode(baseEncode)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(afterDecode, want) > 0 {
		t.Errorf("Base64Encode and Base64Decode not effect")
	}
}

func BenchmarkGetToken(b *testing.B) {
	b.Run("Benchmark test get token", func(b *testing.B) {
		tokens := Token{CheckArrayType: *arr.NewCheckArrayType(b.N)}
		for i := 0; i < b.N; i++ {
			token, err := getToken()
			if err != nil {
				b.Error(err)
			}
			if token == "" {
				b.Error("Generate token is fail.")
			}
			if tokens.Check(token) {
				b.Fatalf("token is repeat")
			}
			tokens.Add(token)
		}
	})
}

type Token struct {
	arr.CheckArrayType
}
