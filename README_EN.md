# DDNS Tool

[![actions workflow](https://github.com/Jackarain/ddns/actions/workflows/go.yml/badge.svg)](https://github.com/Jackarain/ddns/actions)
\
\
A tool for dynamically updating `IP` to domain name configurations, supporting platforms like `dnspod`, `f3322`, `godaddy`, `namesilo`, and `he.net`.

## Environment Preparation

Before starting the compilation, it is necessary to install `golang`/`git` environment and use the `git` command to clone the project locally.

```bash
git clone https://github.com/Jackarain/ddns.git
```

## Compilation Instructions

In the project directory, execute the following command:

```bash
go build
```

Upon successful compilation, an executable program named `ddns` will be generated.

## Usage

The `ddns` program can be scheduled to run in `crontab`, or it can be executed on a schedule with `systemd`. Below is a `crontab` usage example.

```bash
# Execute every 5 minutes
*/5 * * * * /path/to/ddns --dnspod --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```

Here is a systemd usage example.

```bash
# Edit /etc/systemd/system/ddns.service
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
# Edit /etc/systemd/system/ddns.timer
[Unit]
Description=DDNS Timer

[Timer]
OnBootSec=5min
OnUnitActiveSec=5min

[Install]
WantedBy=timers.target
```

```bash
# Start the timer service
systemctl start ddns.timer

# Set to start automatically at boot
systemctl enable ddns.timer
```

In the above examples, the `ddns` program will run 5 minutes after booting and then every 5 minutes. If you need to change the execution time, you can modify the `OnBootSec` and `OnUnitActiveSec` parameters in the `ddns.timer` file. For specific usage, please refer to the `systemd.timer` documentation.

`ddns` can run on devices like `routers` or `NAS`, thereby achieving dynamic updating of `IP` to domain name configurations on the `router` or `NAS`.

## Parameter Explanation and Usage Examples

Here's how you might use the `godaddy` command:

```bash
/path/to/ddns --godaddy --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```

In this example, token is a string comprised of `"API_KEY:API_SECRET"`. The domain here would be `test.example.com`.

Here's an example of how to use `dnspod`:

```bash
/path/to/ddns --dnspod --domain example.com --subdomain test --dnstype AAAA --token "1111111:123123123"
```

Here's how you might use `namesilo`:

```bash
/path/to/ddns --namesilo --domain example.com --subdomain test --dnstype AAAA --token "1111111123123123"
```

Here's an example of how to use `f3322` and `oray`:

```bash
/path/to/ddns --f3322 -user root -passwd xxxxxxxx --domain example.f3322.net

/path/to/ddns --oray -user root -passwd xxxxxxxx --domain example.vicp.net
```

Here's how you might use `he.net`:

```bash
/path/to/ddns --henet --domain example.com --subdomain test --dnstype AAAA --token "A6z56I89bUghPk8h"
```

Here's an example of how to obtain the public `ip` by sending a `curl` request to `ipv4.seeip.org`:

```bash
/path/to/ddns --dnspod --domain example.com --subdomain test --dnstype A --token "1111111:123123123" --command "curl https://ipv4.seeip.org"
```

By default, `ddns` queries `ipify.org` to obtain the public `ip`.

## Support and Feedback

If you encounter any issues during use, or have any suggestions and feedback, feel free to submit an [Issue](https://github.com/Jackarain/ddns/issues) or [Pull Request](https://github.com/Jackarain/ddns/pulls) via the Github page of this project.
