package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

var MyKey string = "WKqyCgplpuRgqm38"

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func Encrypt(origData []byte) ([]byte, error) {
	return AesEncrypt(origData, []byte(MyKey))
}

func Decrypt(origData []byte, clientos string) ([]byte, error) {
	return AesDecrypt(origData, []byte(MyKey))
}


func HexEncode(oridata []byte) (encodestr string){
	dstlen := hex.EncodedLen(len(oridata))
	dst := make([]byte, dstlen)

	hex.Encode(dst, oridata)
	encodestr = string(dst)
	return

}

func HexDecode(oridata []byte) (decodestr []byte){
	hex2strlen := hex.DecodedLen(len(oridata))
	str := make([]byte, hex2strlen)
	hex.Decode(str, oridata)
	decodestr = str
	return
}