package myaes

import "crypto/cipher"

type Tobytes struct {
	Cip     cipher.Block
	Pdgtext []byte
}

// Encrypt 使用AES加密文本,加密的文本不能为空
func (a *Tobytes) Encrypt(src []byte) (dst []byte) {
	src = a.padding(src)
	dst = make([]byte, len(src))
	var index int = 0
	for len(src) > 0 {
		a.Cip.Encrypt(dst[index:index+a.Cip.BlockSize()], src)
		index += a.Cip.BlockSize()
		src = src[a.Cip.BlockSize():]
	}
	return dst
}

// Decrypt 使用AES解密文本
func (a *Tobytes) Decrypt(src []byte) (dst []byte) {
	if len(src)%a.Cip.BlockSize() != 0 {
		return src
	}
	dst = make([]byte, len(src))
	var index int = 0
	for len(src) > 0 {
		a.Cip.Decrypt(dst[index:index+a.Cip.BlockSize()], src)
		index += a.Cip.BlockSize()
		src = src[a.Cip.BlockSize():]
	}
	return a.unpadding(dst)
}

// 使用AES加密文本的时候文本必须定长,即必须是16,24,32的整数倍,
func (a *Tobytes) padding(src []byte) (dst []byte) {
	pdg := a.Cip.BlockSize() - len(src)%a.Cip.BlockSize()
	p := a.Pdgtext[:pdg]
	p[pdg-1] = byte(pdg)
	return append(src, p...)
}

// 使用AES解密文本,解密收删除padding的文本
func (a *Tobytes) unpadding(src []byte) (dst []byte) {
	length := len(src)
	if length <= 0 {
		return src
	}
	return src[:(length - int(src[length-1]))]
}
