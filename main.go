// Package main
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/chain5j/chain5j-keytool-shamir/shamir"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "secretBox",
	}

	cmdEncrypt = &cobra.Command{
		Use:   "encrypt to splice",
		Short: "encrypts to splice",
		Run:   runEncrypt,
	}

	cmdDecrypt = &cobra.Command{
		Use:   "decrypt to combine",
		Short: "decrypts to combine",
		Run:   runDecrypt,
	}

	// encrypt
	parts     int    // 密码分成多少片
	threshold int    // 门限
	secret    string // 原始密码
	// decrypt
	shares string // 密钥片

	// common
	inputFile   string // 输入文件
	outputFile  string // 输出文件
	encoding    string // 加密类型
	pwd         string // 加密密码
	sssType     string // 密钥共享的类型
	isSecretHex bool   // 原始密码是否是16进制
)

func init() {
	// encrypt
	cmdEncrypt.Flags().IntVarP(&parts, "parts", "n", 5, "total parts to split key into")
	cmdEncrypt.Flags().IntVarP(&threshold, "threshold", "m", 3, "minimum parts needed to decrypt")
	cmdEncrypt.Flags().StringVarP(&secret, "secret", "k", "", "secret need to shamir")

	// decrypt
	cmdDecrypt.Flags().StringVarP(&shares, "shares", "s", "", "Shares need to combine")

	addFlags := func(cmd *cobra.Command, op string) {
		cmd.Flags().StringVarP(&inputFile, "inputFile", "i", "", fmt.Sprintf("file to %s", op))
		cmd.Flags().StringVarP(&outputFile, "outputFile", "o", "", fmt.Sprintf("destination for %sed file", op))
		cmd.Flags().StringVarP(&encoding, "encoding", "e", shamir.EncodingHEX, "key encoding to use (hex, base64, aes)")
		cmd.Flags().StringVar(&pwd, "pwd", "", "the password of encrypt")
		cmd.Flags().StringVar(&sssType, "sssType", "galois", "the sss type to use (galois,gf256)")
		cmd.Flags().BoolVar(&isSecretHex, "isSecretHex", true, "where the secret is hex")
	}

	addFlags(cmdEncrypt, "encrypt")
	addFlags(cmdDecrypt, "decrypt")

	cmd.AddCommand(cmdEncrypt, cmdDecrypt)
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func runEncrypt(cmd *cobra.Command, args []string) {
	if secret == "" {
		if !fileExists(inputFile) {
			panic(fmt.Sprintf("inputFile file '%s' must exist", inputFile))
		}
		file, err := ioutil.ReadFile(inputFile)
		if err != nil {
			panic(fmt.Sprintf("Read inputFile err: %s", err))
		}
		secret = string(file)
	}

	b, err := shamir.Encrypt(parts, threshold, sssType, secret, isSecretHex, encoding, pwd)
	if err != nil {
		panic(err)
	}

	if outputFile != "" {
		fmt.Printf("Encrypting to '%s'\n", outputFile)
		parent := path.Dir(outputFile)
		if !fileExists(parent) {
			panic(fmt.Sprintf("outputFile directory '%s' does not exist", parent))
		}
		err := ioutil.WriteFile(outputFile, b, os.FileMode(400))
		if err != nil {
			panic(fmt.Sprintf("WriteFile err: %s", err))
		}
	} else {
		fmt.Println("=============================")
		fmt.Println(string(b))
	}

	fmt.Println("Encrypt Success!")
}

func runDecrypt(cmd *cobra.Command, args []string) {
	if shares == "" {
		if !fileExists(inputFile) {
			panic(fmt.Sprintf("inputFile file '%s' must exist", inputFile))
		}
		file, err := ioutil.ReadFile(inputFile)
		if err != nil {
			panic(fmt.Sprintf("Read inputFile err: %s", err))
		}
		shares = string(file)
		shares = strings.ReplaceAll(shares, "\n", ",")
	}

	resultBytes, err := shamir.Decrypt(sssType, shares, encoding, pwd)
	if err != nil {
		panic(err)
	}
	if outputFile != "" {
		parent := path.Dir(outputFile)
		os.MkdirAll(parent, os.ModePerm)
		if !fileExists(parent) {
			panic(fmt.Sprintf("outputFile directory '%s' does not exist", parent))
		}
		err = ioutil.WriteFile(outputFile, resultBytes, os.FileMode(400))
		if err != nil {
			panic(fmt.Sprintf("Failed to write outputFile file: %s", err))
		}
	} else {
		fmt.Println("=============================")
		if isSecretHex {
			fmt.Println(hexutil.Bytes2Hex(resultBytes))
		} else {
			fmt.Println(string(resultBytes))
		}
	}

	fmt.Println("Decrypt Success!")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
