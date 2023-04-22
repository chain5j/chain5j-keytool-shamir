// Package gf256
package gf256

import (
	"fmt"

	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

func Example() {
	secret := "0ddb327ad1059662da1f02f1b8521bf0f69cf5cecc09a4d8fc7f928fc9726818" // our secret
	n := byte(5)                                                                 // create 30 shares
	k := byte(2)                                                                 // require 2 of them to combine

	bytes := hexutil.Hex2Bytes(secret)
	// shares, err := Split(n, k, []byte(secret)) // split into 30 shares
	shares, err := Split(n, k, bytes) // split into 30 shares
	if err != nil {
		fmt.Println(err)
		return
	}

	// select a random subset of the total shares
	subset := make(map[byte][]byte, k)
	for x, y := range shares { // just iterate since maps are randomized
		// subset[x] = y
		fmt.Println(x, hexutil.Bytes2Hex(y))
		if len(subset) == int(k) {
			break
		}
	}
	// 1 ddb93db7d0c076f013afa80fd38305ef636085efb6243e9ce45b53261f2b321d
	// 2 b61f2cfbd3944d5d53644d166eeb27cec77f158c38538b50cc370bc67ec0dc12
	// 3 667d2336d251adcf9ad4e7e8053a39d1528365ad427e1114d413ca6fa8998617
	// 4 60480e63d53c3b1cd3e99c240f3b638c94412e4a3fbdfad39cefbb1dbc0d1b0c
	// 5 b02a01aed4f9db8e1a5936da64ea7d9301bd5e6b4590609784cb7ab46a544109

	subset[byte(1)] = hexutil.Hex2Bytes("03b1fb1fbc2f6f27e93b4bf513de0b655130a02e98718bedf762dd989b22552a")
	subset[byte(2)] = hexutil.Hex2Bytes("110fbbb00b517fe8bc5790f9f5513bc1a3df5f1564f9fab2ea450ca16dd2127c")

	combine := Combine(subset)
	// combine two shares and recover the secret
	// recovered := string(Combine(subset))
	fmt.Println(hexutil.Bytes2Hex(combine))

	// Output: well hello there!
}
