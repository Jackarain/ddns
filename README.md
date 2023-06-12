## DDNS 工具

一个用于动态更新IP到指定域名的工具.

## 编译

安装golang/git环境后, 在项目目录执行以下命令编译
```
go build
```

即可完成编译，编译生成ddns可执行程序.


## 使用方法

godaddy使用示例(请注意，token为 `"API_KEY:API_SECRET"` 组合的字符串)
```
./ddns --godaddy --domain superpool.io --subdomain test --dnstype AAAA --token "1111111:123123123"
```

dnspod使用示例
```
./ddns --dnspod --domain superpool.io --subdomain test --dnstype AAAA --token "1111111:123123123"
```

namesilo使用示例
```
./ddns --namesilo --domain superpool.io --subdomain test --dnstype AAAA --token "1111111123123123"
```

f3322使用示例
```
./ddns --f3322 -f3322user root -f3322passwd xxxxxxxx --domain test.f3322.net
```

三方工具获取公网ip(默认请求 ipify.org 获取ip)使用示例
```
./ddns --dnspod --domain superpool.io --subdomain test --dnstype A --token "1111111:123123123" --command "curl https://ipv4.seeip.org"
```
