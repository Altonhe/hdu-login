# ⚡ hdu-login ![Go](https://github.com/altonhe/hdu-login/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/altonhe/hdu-login)](https://goreportcard.com/report/github.com/altonhe/hdu-login) [![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?logo=sourcegraph)](https://sourcegraph.com/github.com/altonhe/hdu-login)

杭州电子科技大学校园网 Wi-Fi 登录 / 深澜（srun）校园网模拟登录

2021 年暑期学校网络大改造，更换了新的深澜上网认证系统。Vidar-Team 信息安全实验室需要适配学校新的认证系统进行自动登录接入网络，因此有了本项目。

本项目在原有的基础上，同时兼容旧版的网络登录。
## 开始使用

```bash
# 克隆项目
git clone git@github.com:vidar-team/srun-login.git

# 编译项目
cd srun-login/cmd && go build .

# 模拟登录
./srun-login --username=<REDACTED> --password=<REDACTED>
```

## License

MIT License
