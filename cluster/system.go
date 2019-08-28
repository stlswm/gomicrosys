package cluster

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang_microservice_assistant/io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var clusterIP []string
var clusterKey string
var system map[string]string

func init() {
	system = make(map[string]string)
}

// 设置系统秘钥
func SetClusterKey(key string) error {
	if len(key) != 32 {
		return errors.New("集群秘钥必须是32位字符串")
	}
	clusterKey = key
	return nil
}

// 添加集群服务器
func AddClusterMemberServer(ip string) error {
	if ip == "" {
		return errors.New("请设置正确的ip地址")
	}
	clusterIP = append(clusterIP, ip)
	return nil
}

// 是否是集群成员服务器
func IsClusterMemberServer(ip string) bool {
	if ip == "" {
		return false
	}
	for _, ipAddr := range clusterIP {
		if ip == ipAddr {
			return true
		}
	}
	return false
}

// 验证请求是否来原与内部系统
func IsInnerReq(authStr string, random string) bool {
	return authStr == GeneratorAuthKey(random)
}

// 添加内部系统
func AddServer(alias string, domain string) error {
	if alias == "" {
		return errors.New("请设置正确的系统别名")
	}
	if domain == "" {
		return errors.New("请设置正确的请求地址")
	}
	system[alias] = domain
	return nil
}

// 获取内部系统域名
func GetSystemDomain(alias string) (error, string) {
	domain, ok := system[alias]
	if ok {
		return nil, domain
	}
	return errors.New("系统不存在"), ""
}

// 随机字符串
func GetRandomString(length int) string {
	var result []byte
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bs := []byte(str)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bs[r.Intn(len(bs))])
	}
	return string(result)
}

// 生成授权秘钥
func GeneratorAuthKey(randomStr string) string {
	ctx := md5.New()
	ctx.Write([]byte(clusterKey + randomStr))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 内部系统Json请求
func InnerJsonReq(alias string, router string, data interface{}) (error, *io.Package) {
	err, domain := GetSystemDomain(alias)
	if err != nil {
		return err, nil
	}
	router = "/" + strings.TrimLeft(router, "/")
	dataByte, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", domain+router, bytes.NewBuffer(dataByte))
	if err != nil {
		return err, nil
	}
	random := GetRandomString(32)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cluster-Random", random)
	req.Header.Set("Cluster-Auth", GeneratorAuthKey(random))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := &io.Package{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err, nil
	}
	return nil, res
}
