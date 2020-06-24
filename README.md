## DDNS 工具

一个用于动态更新IP到指定域名的工具.

## 编译

安装golang/git环境后, 在项目目录执行以下命令编译
```
go build
```

即可完成编译，编译生成ddns可执行程序.


## 使用方法

godaddy使用示例
```
./ddns --godaddy --domain superpool.io --subdomain test --dnstype AAAA --token "1111111:123123123"
```

dnspod使用示例
```
./ddns --dnspod --domain superpool.io --subdomain test --dnstype AAAA --token "1111111:123123123"
```

三方工具获取ip(默认请求ipify.org获取ip)使用示例
```
./ddns --dnspod --domain superpool.io --subdomain test --dnstype A --token "1111111:123123123" --command "curl" --args "http://ipecho.net/plain"
```
