// Package keytool
package shamir

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func encode(part []byte, enc string, pwd string) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("加密出错%s\n", err)
		}
	}()
	switch enc {
	case EncodingBASE64:
		return base64.StdEncoding.EncodeToString(part)
	case EncodingAES:
		encrypted := AesEncryptECB(part, []byte(pwd))
		return hex.EncodeToString(encrypted)
	default:
		return hex.EncodeToString(part)
	}
}

func decode(part, enc string, pwd string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("解密出错%s\n", err)
		}
	}()
	switch enc {
	case EncodingBASE64:
		return base64.StdEncoding.DecodeString(part)
	case EncodingAES:
		bytes, err := hex.DecodeString(part)
		if err != nil {
			return nil, err
		}
		decrypted := AesDecryptECB(bytes, []byte(pwd))
		return decrypted, nil
	default:
		return hex.DecodeString(part)
	}
}
