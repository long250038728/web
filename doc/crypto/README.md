## 非对称加密

### 常见的非对称加密有
RSA,SM2,ECDSA

### 公钥/私钥
* 私钥：解密、签名 （自己独有）
* 公钥：加密，验签 （可分享至外部人员）
注： 
1. 公钥用于加密数据的，私钥用于解密（由于公钥是到处持有如果公钥能解密就没意义了）
2. 私钥签名，公钥验签 （公钥是判断这个签名是否有效而不能解密这个签名原本的数据，只能起到校验的作用所以叫做签名及验签而不是加密）

## X.509
* 是国际标准，定义数字证书应该包含什么信息，每个信息的格式及类型。定义出请求文件格式

### 文件格式
* PEM是一种文本编码格式（用于保存密钥/证书的格式）
* CSR是证书签名请求 （里面包含公钥信息，你的组织信息，还有签名） 通过该方式可以把公钥传递格式为.csr(保证公钥及公钥不被篡改)，可通过CA验证后生成.crt文件（即有公钥，公钥的颁发者信息及签名，CA的认证后的签名————保证了这个是一个通过CA认证过的可以传递公钥的一个证书）
* .key跟.pem为保存私钥的格式
* .pub跟.pem为保存公钥的格式


### 公钥私钥生成
```
//通过openssl生成一个2048的私钥。文件名为 my.key
openssl genrsa -out my.key 2048

//通过私钥生成公钥
openssl rsa -in my.key -pubout -out my.pub
```

### 自签名证书
```
//通过openssl生成一个2048的私钥。文件名为 my.key
openssl genrsa -out my.key 2048

//通过openssl指定私钥生成CSR证书签名请求
openssl req -new -key my.key -out my.csr

//通过openssl生成自签名证书(使用的是my.key这个私钥)
openssl x509 -req -days 365 -in my.csr --signkey my.key -out my.crt
```

### 代码演示
加密解密
1. 获取公钥私钥内容
2. 通过pem.Decode生成block验证是否有误
3. 使用 x509.ParsePKIXPublicKey 或 x509.ParsePKCS1PrivateKey 通过公钥私钥内容获取公钥私钥对象
4. 使用 rsa.EncryptPKCS1v15 或 rsa.DecryptPKCS1v15 进行加密解密
```
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// ===== 1. 读取公钥 =====
	pubKeyData, err := os.ReadFile("my.pub")
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		panic("failed to decode public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	pub := pubInterface.(*rsa.PublicKey)

	// ===== 2. 公钥加密 =====
	plaintext := []byte("Hello, 非对称加密 RSA！")
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, plaintext)
	if err != nil {
		panic(err)
	}
	fmt.Printf("加密后的数据（Base64前）：%x\n", ciphertext)

	// ===== 3. 读取私钥 =====
	privKeyData, err := os.ReadFile("my.key")
	if err != nil {
		panic(err)
	}
	block, _ = pem.Decode(privKeyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		panic("failed to decode private key PEM")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	// ===== 4. 私钥解密 =====
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	fmt.Println("解密后的明文：", string(decrypted))
}

```
签名验签
1. 获取公钥私钥内容
2. 通过pem.Decode生成block验证是否有误
3. 使用 x509.ParsePKIXPublicKey 或 x509.ParsePKCS1PrivateKey 通过公钥私钥内容获取公钥私钥对象
4. 使用 rsa.SignPKCS1v15 或 rsa.VerifyPKCS1v15 进行签名跟验签，使用hash算法crypto.SHA256（生成出来的用于是32字节，由于是签名所以不需要解密回来的使用hash算法长度为32位的是可以的）
    * 需要使用base64.StdEncoding.EncodeToString(signature)对[]byte进行转换为string用于显示及传输
```
package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	// ===== 1. 读取私钥，用于签名 =====
	privKeyData, err := os.ReadFile("my.key")
	if err != nil {
		panic(err)
	}
	privBlock, _ := pem.Decode(privKeyData)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		panic("failed to decode private key PEM")
	}
	priv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		panic(err)
	}

	// ===== 2. 对内容生成 SHA256 摘要 =====
	message := []byte("Hello RSA Sign & Verify!")
	hashed := sha256.Sum256(message)

	// ===== 3. 私钥签名 =====
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:]) //sha256.Sum256(message)返回的是一个固定长度数组 [32]byte，而 RSA 的函数需要 字节切片 []byte
	if err != nil {
		panic(err)
	}
	fmt.Println("签名(Base64)：", base64.StdEncoding.EncodeToString(signature))

	// ===== 4. 读取公钥，用于验签 =====
	pubKeyData, err := os.ReadFile("my.pub")
	if err != nil {
		panic(err)
	}
	pubBlock, _ := pem.Decode(pubKeyData)
	if pubBlock == nil || pubBlock.Type != "PUBLIC KEY" {
		panic("failed to decode public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		panic(err)
	}
	pub := pubInterface.(*rsa.PublicKey)

	// ===== 5. 公钥验签 =====
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Println("❌ 验签失败！")
	} else {
		fmt.Println("✅ 验签成功，内容未被篡改！")
	}
}

```

### 注意
1. 由于RSA加密明文长度的限制，不适合加密长文本，常用的做法是通过RSA加密一个对称密钥，然后通过这个对称密钥进行解密