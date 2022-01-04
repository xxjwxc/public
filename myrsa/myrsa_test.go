package myrsa

import (
	"fmt"
	"testing"
)

var Pubkey = `-----BEGIN 公钥-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwz/s5ovLPtAcBVMx8jnQ
GQhrPC40DaGdVegU99KuK6Kx6vkMhQU6BxxJts3EYGzuxahFovjuI6MTlYMp2wak
HZwMMnW7Yc3bgPBREtIr8GQCZwdhW71desmC5BoWpzLSg+e+UEvOzwrE8tJECkXW
QdZkaXKS3vnrP+ULbQoJf/dNv2+/YVR1o5ZDmkwhAO8I7yM6Qa+65a9nrStXvnHm
RgMbMIOl5BM6/pBJ4fLzJ4hobW81q/laPmPbm0ow1cN59iTursd4DJm38YSaVYqb
GcA7FT5TV73PtrA1FB3XlKXtgIKHXsR5XzP+NEdO+pyHboioL4Pr7pFEJ6v19d06
IQIDAQAB
-----END 公钥-----
`

var Pirvatekey = `-----BEGIN 私钥-----
MIIEpgIBAAKCAQEAwz/s5ovLPtAcBVMx8jnQGQhrPC40DaGdVegU99KuK6Kx6vkM
hQU6BxxJts3EYGzuxahFovjuI6MTlYMp2wakHZwMMnW7Yc3bgPBREtIr8GQCZwdh
W71desmC5BoWpzLSg+e+UEvOzwrE8tJECkXWQdZkaXKS3vnrP+ULbQoJf/dNv2+/
YVR1o5ZDmkwhAO8I7yM6Qa+65a9nrStXvnHmRgMbMIOl5BM6/pBJ4fLzJ4hobW81
q/laPmPbm0ow1cN59iTursd4DJm38YSaVYqbGcA7FT5TV73PtrA1FB3XlKXtgIKH
XsR5XzP+NEdO+pyHboioL4Pr7pFEJ6v19d06IQIDAQABAoIBAQC/hEi/q4flSQTz
RDPNwV+Z7mQhV8C/TjOiPE+09vbY3nFeZoQdRo8wwlKb+SIS40ciongL79jHJALl
uQ6pRM5eLN7Z8BmpSd9xjkg4CegHmFGy9c5NymWTN24oiF6ICpXrxLks0e89rvaY
qB8NZItRcRZ1SGlARiy3M9gNULcYyKi7Vfx3ewukxbsY9mX6jB5+1qRLtjsOfy3f
HK7OYOs6d7ymPLI37jqCqpKmsif8MTV2mIAIboVN+6mTPv6X3okF8zOrhDWmscgU
5wmZJFt5DozZU4x6gXqFWUJMrmHfWniJnaYYHnULXtcRTUVjycm8le9ositsfBBc
ziAKb+gxAoGBAOXABoAf16CTGDazFjkwh4c0Frva8KKhEqo7hGTTjt3sS6V8u7Mt
GvZmf902WLNhfZAPoQ+O5Ngx3FSbZvESWTGUk9FMUABE92WCxWyE8DNKKeb9Kyx6
+tuZFQC/6tlo4/WRm80La75OJhyUqmOQQSJpM0idw1+sp5T97RC5+FUzAoGBANmO
y8VLjpCqKDEhsnHzBrtmIAbLcQe27ZRuymZgfnr68PvqMLKMiEZwMFNemNYvTGtp
I5133fWI0m3BhaxHQfI3sBMt1jHi6EXCE5gogl2gXVuQvN8XtUg6IvSWhrD0YtK4
1Bd5IPUjOuNLpfDX16tGjKXnGFmLPMzN1u3t20tbAoGBALBlxRfuWtIw3fBxg+iY
+BW4ypOlQAi9fuUxGS+ItzJw0IvYvwyM3xy8CgRAS84+Vfeb6F9XqSDM94wGXP1O
xyioGO4jl3D9gq1vwEDXuMzIbm+phdJ7AcxFNrkCoUAXpp7PEz5VPH465kwfYMtc
4IWZHATvDCiTGX/tjmy/PIm9AoGBANNaenPtd07rP6ibh/RTmRKtoCd5tRE9kYlG
KLNUwwtOhpb1aOHMzQdBLnGP0QMjaCZhOgxcyvEiPuwJuYcootRhbVj0isZkHirG
5KpJkHzMsmWmMxa4vZCxigv7wFZg1TDKBqHXN0FvPGJct5VG22q1WyZBX9J+Bk8h
GdCD5ytJAoGBALot/lOQZiczJ2O4TYLN6R4tnDGLfT9lmyBu8BvpVTxb3NjYPMZa
fsw3L3IyedE/omwm+OgBeqTR8Xk0ZgD7B29hoTY5cyBNUWmKNsowcBcDMgSCaTr4
Vp20lFpW8zKyHwj1DeXabvT3dUCQC5KIqE6Bv6GYjLa1CHTJSz3dEVVZ
-----END 私钥-----
`

func TestMain(t *testing.T) {
	// 公钥加密私钥解密
	src, err := PublicEncrypt("你好，hello world", Pubkey)
	fmt.Println(src, err)
	src, err = PriKeyDecrypt(src, Pirvatekey)
	fmt.Println(src, err)

	// 公钥解密私钥加密
	src, err = PriKeyEncrypt("你好，hello world", Pirvatekey)
	fmt.Println(src, err)
	src, err = PublicDecrypt(src, Pubkey)
	fmt.Println(src, err)
}
