package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"github.com/prometheus/common/log"
	"strings"
)

type RsaKey struct {
	PublicKey  rsa.PublicKey
	PrivateKey rsa.PrivateKey
}

func GenPublicPrivateKey() (*RsaKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	publicKey := privateKey.PublicKey

	rsaKey := RsaKey{PublicKey: publicKey, PrivateKey: *privateKey}
	return &rsaKey, nil

	/*cipherText, e := rsa.EncryptPKCS1v15(rand.Reader, rsaKey.publicKey, []byte(msg))
	if e != nil {
		println("error:" , e)
	}

	fmt.Println(string(cipherText))

	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, rsaKey.privateKey, cipherText)
	fmt.Println(string(plainText))*/
}

func AesEncrypt(orig string, password string) string {
	key := GetSuitablePassword(password)
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string, password string) string {
	key := GetSuitablePassword(password)
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

var size1 = 16
var size2 = 24
var size3 = 32

func GetSuitablePassword(password string) string {
	passwdLen := len(password)
	if passwdLen <= size1 {
		return password + strings.Repeat("1", size1 - passwdLen)
	} else if passwdLen <= size2 && passwdLen > size1{
		return password + strings.Repeat("1", size2 - passwdLen)
	} else if passwdLen <= size3 && passwdLen > size2{
		return password + strings.Repeat("1", size3 - passwdLen)
	} else {
		return password[:size3]
	}
}
