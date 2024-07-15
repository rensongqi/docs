
# Getting Started

使用`github.com/go-acme/lego/v4`生成指定域名ssl证书

```bash
go get github.com/go-acme/lego/v4
```

# Generate private key

```bash
# 使用openssl生成私钥
openssl ecparam -genkey -name prime256v1 -out private-key.pem

# 查看生成的私钥
openssl ec -in private-key.pem -text -noout

# 导出私钥为BEGIN EC PRIVATE KEY格式
openssl ec -in private-key.pem -outform PEM -out private-key.ec.pem

```

# Code
```go
package main
import (
   "crypto"
   "crypto/ecdsa"
   "crypto/x509"
   "encoding/pem"
   "fmt"
   "github.com/gin-gonic/gin"
   "github.com/go-acme/lego/v4/certcrypto"
   "github.com/go-acme/lego/v4/certificate"
   "github.com/go-acme/lego/v4/lego"
   "github.com/go-acme/lego/v4/providers/dns/alidns"
   "github.com/go-acme/lego/v4/registration"
   "github.com/sirupsen/logrus"
   "io"
   "net/http"
   "os"
   "path"
   "runtime"
   "strings"
)
// 阿里云token
const (
   AliyunRegionID        = "cn-shanghai"
   AliyunAccessKeyId     = "xxx"
   AliyunAccessKeySecret = "xxx"
)
type MyUser struct {
   Email        string
   Registration *registration.Resource
   key          crypto.PrivateKey
}
func (u *MyUser) GetEmail() string {
   return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
   return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
   return u.key
}
// loadPrivateKey 读取并解析私钥文件
func loadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
   bytes, err := os.ReadFile(path)
   if err != nil {
      return nil, fmt.Errorf("could not read private key file: %w", err)
   }
   block, _ := pem.Decode(bytes)
   if block == nil || block.Type != "EC PRIVATE KEY" {
      return nil, fmt.Errorf("failed to decode PEM block containing private key")
   }
   privateKey, err := x509.ParseECPrivateKey(block.Bytes)
   if err != nil {
      return nil, fmt.Errorf("could not parse private key: %w", err)
   }
   return privateKey, nil
}
func generateCertificate(domain string) (err error, pemBody, keyBody []byte) {
   var sslFilename string
   if strings.HasPrefix(domain, "*.") {
      domainSplit := strings.Split(domain, "*.")
      sslFilename = domainSplit[1]
   } else {
      sslFilename = domain
   }
   keyFileName := fmt.Sprintf("data/%s.key", sslFilename)
   pemFileName := fmt.Sprintf("data/%s.pem", sslFilename)
   privateKey, err := loadPrivateKey("ssl_ecdsa_private_key.pem")
   if err != nil {
      logrus.Error(err)
      return
   }
   myUser := MyUser{
      Email: "cowaadmin@cowarobot.com",
      key:   privateKey,
   }
   config := lego.NewConfig(&myUser)
   // 配置密钥的类型和密钥申请的地址
   config.CADirURL = lego.LEDirectoryProduction
   config.Certificate.KeyType = certcrypto.RSA2048
   // 创建一个client与CA服务器通信
   client, err := lego.NewClient(config)
   if err != nil {
      logrus.Error(err)
      return
   }
   configC := alidns.NewDefaultConfig()
   configC.RegionID = AliyunRegionID
   configC.APIKey = AliyunAccessKeyId
   configC.SecretKey = AliyunAccessKeySecret
   provider, err := alidns.NewDNSProviderConfig(configC)
   if err != nil {
      logrus.Error(err)
      return
   }
   if err = client.Challenge.SetDNS01Provider(provider); err != nil {
      logrus.Error("Could not create DNS provider: %v", err)
   }
   reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
   if err != nil {
      logrus.Error(err)
      return
   }
   myUser.Registration = reg
   request := certificate.ObtainRequest{
      Domains: []string{domain}, // 支持多个域名
      Bundle:  true,             // 如果是true，将把颁发者证书一起返回: certificates.IssuerCertificate
   }
   certs, err := client.Certificate.Obtain(request)
   if err != nil {
      logrus.Error(err)
      return
   }
   pemBody = certs.Certificate
   keyBody = certs.PrivateKey
   err = os.WriteFile(pemFileName, certs.Certificate, 0600)
   if err != nil {
      logrus.Error(err)
      return
   }
   err = os.WriteFile(keyFileName, certs.PrivateKey, 0600)
   if err != nil {
      logrus.Error(err)
      return
   }
   return
}
type SSLRecord struct {
   Domain string `json:"domain" binding:"required"`
}
func generateSSLCertificate(c *gin.Context) {
   var sslRecord SSLRecord
   if err := c.ShouldBindJSON(&sslRecord); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
         "error": fmt.Sprintf("%v", err),
      })
      return
   }
   err, pemBody, keyBody := generateCertificate(sslRecord.Domain)
   if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
         "error": fmt.Sprintf("%v", err),
      })
      return
   }
   c.JSON(http.StatusOK, gin.H{
      "ssl_pem": pemBody,
      "ssl_key": keyBody,
   })
}
func init() {
   // 设置json文本输出
   logFormat := &logrus.JSONFormatter{
      TimestampFormat: "2006-01-02 15:04:05.000",
      PrettyPrint:     false,
      CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
         fileName := fmt.Sprintf("%v:%v", path.Base(frame.File), frame.Line)
         return frame.Function, fileName
      },
   }
   logrus.SetReportCaller(true)
   logrus.SetFormatter(logFormat)
   //同时写文件和屏幕
   stdOutWrite := io.Writer(os.Stdout)
   //fileAndStdoutWriter := io.MultiWriter(os.Stdout, logrus)
   logrus.SetOutput(stdOutWrite)
   //设置最低loglevel
   logrus.SetLevel(logrus.InfoLevel)
   logrus.Info("日志模块初始化完成")
}
func main() {
   r := gin.Default()
   {
      r.POST("/api/v1/ssl", generateSSLCertificate) // 生成ssl证书
   }
   _ = r.Run(":6666")
}
```

# Usage

```bash
curl -X POST -d '{"domain": "*.rsq.com"}' -H "Content-Type: application/json" http://127.0.0.1:6666/api/v1/ssl
```