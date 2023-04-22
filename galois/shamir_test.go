// Package galois
package galois

import (
	"fmt"
	"testing"

	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

func TestSplit(t *testing.T) {
	split, err := Split([]byte("0ddb327ad1059662da1f02f1b8521bf0f69cf5cecc09a4d8fc7f928fc9726818"), 5, 2)
	if err != nil {
		panic(err)
	}
	for i, v := range split {
		fmt.Println(i, hexutil.Bytes2Hex(v))
	}
	bytes, err := Combine([][]byte{split[0], split[1]})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
