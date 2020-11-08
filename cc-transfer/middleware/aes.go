package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
//key不能泄露
var PwdKey = []byte("DIS**#KKKDJJSKDI")

//AES加密（将data用key加密）
func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	//用PKCS7模式来填充欲加密data
	data = PKCS7Padding(data, blockSize)
	//预先定义结果变量
	res := make([]byte, len(data))
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blocMode.CryptBlocks(res, data)

	return res, nil
}

//AES解密（将sign用key解密）
func AesDecrypt(token []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()

	res := make([]byte, len(token))
	blocMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	blocMode.CryptBlocks(res, token)

	return PKCS7UnPadding(res)
}

//PKCS7 填充模式
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	//算出目标加密串长度离blockSize的整数倍还差多少
	padNum := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padNum)}复制padNum个，然后合并成新的字节切片返回，目的是将目标加密串填充为长度是blockSize的整数倍
	padText := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(ciphertext, padText...)
}

//PKCS7 反填充模式，也即删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		padNum := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - padNum)], nil
	}
}

//加密base64
func Encrypt(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, PwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

//解密
func Decrypt(pwd string) ([]byte, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	return AesDecrypt(pwdByte, PwdKey)
}