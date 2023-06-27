## DDNS 工具

一个用于动态更新 `IP` 到域名配置的工具，支持 `dnspod`、`f3322`、`godaddy`、`namesilo`、`he.net` 平台.

## 编译

安装 `golang/git` 环境后, 在项目目录执行以下命令编译

```
go build
```

即可完成编译，编译生成 `ddns` 可执行程序.

## 使用方法

`godaddy` 使用示例(请注意，`token` 为 `"API_KEY:API_SECRET"` 组合的字符串，域名为：`test.example.com`)

```bash
./ddns --godaddy --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```

`dnspod` 使用示例

```bash
./ddns --dnspod --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```

`namesilo` 使用示例

```bash
./ddns --namesilo --domain example.com --subdomain test --dnstype AAAA --token "1111111123123123"
```

`f3322` 使用示例

```bash
./ddns --f3322 -f3322user root -f3322passwd xxxxxxxx --domain example.f3322.net
```

`he.net` 使用示例

```bash
./ddns --henet --domain example.com --subdomain test --dnstype AAAA --token "A6z56I89bUghPk8h"
```

三方工具获取公网 `ip`(默认请求 `ipify.org` 获取 `ip`)使用示例

```bash
./ddns --dnspod --domain superpool.io --subdomain test --dnstype A --token "1111111:123123123" --command "curl https://ipv4.seeip.org"
```
