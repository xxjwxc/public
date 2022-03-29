package myrsa

import (
	"github.com/wenzhenxi/gorsa"
)

// PublicEncrypt 公钥加密
func PublicEncrypt(src string, pubkey string) (string, error) {
	return gorsa.PublicEncrypt(src, pubkey)
}

// PriKeyDecrypt 私钥解密
func PriKeyDecrypt(src string, Pirvatekey string) (string, error) {
	return gorsa.PriKeyDecrypt(src, Pirvatekey)
}

// PriKeyEncrypt 私钥加密
func PriKeyEncrypt(src string, pirvatekey string) (string, error) {
	return gorsa.PriKeyEncrypt(src, pirvatekey)
}

// PublicDecrypt 公钥解密
func PublicDecrypt(src string, pubkey string) (string, error) {
	return gorsa.PublicDecrypt(src, pubkey)
}
