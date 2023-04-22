// Package keytool
package shamir

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/chain5j/chain5j-keytool-shamir/galois"
	"github.com/chain5j/chain5j-keytool-shamir/gf256"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

const (
	// EncodingHEX ...
	EncodingHEX = "hex"
	// EncodingBASE64 ...
	EncodingBASE64 = "base64"
	// EncodingAES ...
	EncodingAES = "aes"

	// SssTypeGalois ...
	SssTypeGalois = "galois"
	// SssTypeGf256 ...
	SssTypeGf256 = "gf256"
)

// Encrypt 分片
//
// parts：私钥总分片数量
// threshold：门限数
// sssType：shamir的方式(galois,gf256)
// secret：需要分片的私钥
// isSecretHex：私钥是否时16进制
// encoding：Encrypt时输出的结果类型(hex,base64,aes)
// pwd：当encoding=aes时,使用pwd进行加密返回结果
func Encrypt(parts, threshold int, sssType string, secret string, isSecretHex bool, encoding, pwd string) ([]byte, error) {
	if secret == "" {
		return nil, errors.New("plain is empty")
	}
	err := checkEncoding(encoding)
	if err != nil {
		return nil, err
	}
	err = checkSssType(sssType)
	if err != nil {
		return nil, err
	}
	var paramsBytes []byte
	if isSecretHex {
		if !hexutil.IsHex(secret) {
			return nil, errors.New("secret is not hex")
		}
		paramsBytes = hexutil.Hex2Bytes(secret)
	} else {
		paramsBytes = []byte(secret)
	}
	var b bytes.Buffer // 直接定义一个 Buffer 变量，而不用初始化
	if sssType == SssTypeGalois {
		split, err := galois.Split(paramsBytes, parts, threshold)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Shamir split err: %s", err))
		}
		for _, part := range split {
			s := encode(part, encoding, pwd)
			b.WriteString(s)
			b.WriteString("\n")
		}
	} else {
		split, err := gf256.Split(byte(parts), byte(threshold), paramsBytes)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Shamir split err: %s", err))
		}
		for i, part := range split {
			part = append([]byte{i}, part...)
			s := encode(part, encoding, pwd)
			// // byte其实是uint8的别名，byte 和 uint8 之间可以直接进行互转。目前，只能将0 ~ 255范围的int转成byte
			b.WriteString(s)
			b.WriteString("\n")
		}
	}
	return b.Bytes(), nil
}

// Decrypt 恢复私钥
// sssType：shamir的方式(galois,gf256)
// shares：需要恢复私钥的密码片集，密钥片之间用”,“分割
// encoding：Encrypt时输出的结果类型(hex,base64,aes)
// pwd：当encoding=aes时,使用pwd进行加密返回结果
func Decrypt(sssType string, shares string, encoding, pwd string) ([]byte, error) {
	var (
		resultBytes []byte
		err         error
	)

	if shares == "" {
		return nil, errors.New("shares is empty")
	}
	err = checkEncoding(encoding)
	if err != nil {
		return nil, err
	}
	err = checkSssType(sssType)
	if err != nil {
		return nil, err
	}
	shamirPart := strings.Split(shares, ",")
	if len(shamirPart) == 0 {
		return nil, errors.New("shares is empty")
	}

	len := 0
	for _, p := range shamirPart {
		if p == "" {
			continue
		}
		len++
	}

	if sssType == SssTypeGalois {
		partsArray := make([][]byte, len)
		for i, p := range shamirPart {
			bytes, err := decode(p, encoding, pwd)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Decode plain err: %s", err))
			}
			partsArray[i] = bytes
		}

		resultBytes, err = galois.Combine(partsArray)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Shamir combine err: %s", err))
		}
	} else {
		partsArray := make(map[byte][]byte, len)
		for _, p := range shamirPart {
			bytes, err := decode(p, encoding, pwd)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Decode plain err: %s", err))
			}
			partsArray[bytes[0]] = bytes[1:]
		}
		resultBytes = gf256.Combine(partsArray)
	}
	return resultBytes, err
}

// 确信输出格式是否支持
func checkEncoding(encoding string) error {
	switch encoding {
	case EncodingHEX, EncodingBASE64:
		return nil
	case EncodingAES:
		return nil
	default:
		return errors.New(fmt.Sprintf("unknown encoding '%s'", encoding))
	}
}

// checkSssType ...
// @sssType: string
// returns:
// #1: error
func checkSssType(sssType string) error {
	switch sssType {
	case SssTypeGalois, SssTypeGf256:
		return nil
	default:
		return errors.New(fmt.Sprintf("unknown sss type '%s'", sssType))
	}
}
