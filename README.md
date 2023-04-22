# 密盒

## 简介

`chain5j-keytool-shamir` 是一个根据私钥进行密码共享分片的工具包。

## 安装使用

```shell
go get -u github.com/chain5j/chain5j-keytool-shamir
```

## 使用说明

### 参数说明

| 参数            | 说明                                                              |
|---------------|-----------------------------------------------------------------|
| -m            | (Encrypt)--parts,私钥总分片数量                                        |
| -m            | (Encrypt)--threshold,门限数                                        |
| -k            | (Encrypt)--secret,需要分片的私钥                                       |
| -s            | (Decrypt)--shares,需要恢复私钥的密码片集，密钥片之间用”,“分割                       |
| -i            | --inputFile,若Encrypt时,secret为空则读文件中的内容;若Decrypt时,shares为空则读文件内容 |
| -o            | --outputFile,输出的文件.若无,直接打印                                      |
| -e            | --encoding,Encrypt时输出的结果类型(hex,base64,aes)                      |
| --pwd         | pwd,当encoding=aes时,使用pwd进行加密返回结果                                |
| --sssType     | sssType,shamir的方式(galois,gf256)(默认galois)                       |
| --isSecretHex | secret是否为16进制(默认true)                                           |

### 加密

```go
go run main.go encrypt -n 5 -m 2 --sssType "gf256" -k "0ddb327ad1059662da1f02f1b8521bf0f69cf5cecc09a4d8fc7f928fc9726818" --isSecretHex = true -e aes --pwd "1234567890123456" 
```

### 解密

```go
go run main.go decrypt --sssType "gf256" -s "7b035ad7368edb425fc2a31d81a7bf0f91ff9e22352c7283250cebaa236986c381a2c013a214e89ed7f47a40e509a60d,d35de6c3007cb36f5597dad5a54da3b1368290010ac742993db4ecf378baeb4f15018fbacdfe8836938558a472d3f078" --isSecretHex = true -e aes --pwd "1234567890123456" 
```

## 证书

`chain5j-keytool-shamir` 的源码允许用户在遵循 [Apache 2.0 开源证书](LICENSE) 规则的前提下使用。

## 版权

Copyright@2023 chain5j

![chain5j](./chain5j.png)