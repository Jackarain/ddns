# DDNS 工具

一个用于动态更新 `IP` 到域名配置的工具，支持 `dnspod`、`f3322`、`godaddy`、`namesilo`、`he.net` 平台.

## 环境准备

在开始编译前，需要安装 `golang`/`git` 环境，并使用 `git` 命令将项目克隆到本地

```bash
git clone https://github.com/Jackarain/ddns.git
```

## 编译方法

在项目目录下执行以下命令

```bash
go build
```

编译完成后，会生成名为 `ddns` 的可执行程序。


## 使用方法

通常可以将 ddns 程序放在 `crontab` 中定时执行，也可以使用 `systemd` 定时执行，以下是 `crontab` 的使用示例

```bash
# 每 5 分钟执行一次
*/5 * * * * /path/to/ddns --dnspod --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```

以下是 `systemd` 的使用示例

```bash
# 编辑 /etc/systemd/system/ddns.service
[Unit]
Description=DDNS Service
After=network.target

[Service]
WorkingDirectory=/tmp/
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
Restart=no
ExecStart=/path/to/ddns --dnspod --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"

[Install]
WantedBy=multi-user.target
```

```bash
# 编辑 /etc/systemd/system/ddns.timer
[Unit]
Description=DDNS Timer

[Timer]
OnBootSec=5min
OnUnitActiveSec=5min

[Install]
WantedBy=timers.target
```

```bash
# 启动定时服务
systemctl start ddns.timer

# 设置开机自启
systemctl enable ddns.timer
```
以上示例中，`ddns` 程序会在开机后 5 分钟执行一次，之后每 5 分钟执行一次，如果需要修改执行时间，可以修改 `ddns.timer` 文件中的 `OnBootSec` 和 `OnUnitActiveSec` 参数，具体使用方法可以参考 `systemd.timer` 的文档。

`ddns` 可以运行在路由器或 `NAS` 等设备上，这样就可以实现在路由器或 `NAS` 上实现动态更新 `IP` 到域名配置的功能

## 参数说明及使用示例

`godaddy` 使用示例

```bash
/path/to/ddns --godaddy --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```
在这个示例中，`token` 是由 `"API_KEY:API_SECRET"` 组成的字符串，域名为：`test.example.com`。

`dnspod` 使用示例

```bash
/path/to/ddns --dnspod --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```

`namesilo` 使用示例

```bash
/path/to/ddns --namesilo --domain example.com --subdomain test --dnstype AAAA --token "1111111123123123"
```

`f3322` 使用示例

```bash
/path/to/ddns --f3322 -f3322user root -f3322passwd xxxxxxxx --domain example.f3322.net
```

`he.net` 使用示例

```bash
/path/to/ddns --henet --domain example.com --subdomain test --dnstype AAAA --token "A6z56I89bUghPk8h"
```

通过 `curl` 请求 `ipv4.seeip.org` 获取公网 `ip` 使用示例

```bash
/path/to/ddns --dnspod --domain example.com --subdomain test --dnstype A --token "1111111:123123123" --command "curl https://ipv4.seeip.org"
```
默认情况下，`ddns` 请求 `ipify.org` 以获取公网 `ip`


## 支持和反馈

如果您在使用过程中遇到任何问题，或有任何建议和反馈，欢迎通过本项目的 `Github` 页面提交 [Issue](https://github.com/Jackarain/ddns/issues) 或 [Pull Request](https://github.com/Jackarain/ddns/pulls)
