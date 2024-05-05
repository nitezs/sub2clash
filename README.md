# sub2clash

> Sing-box 用户？看看另一个项目 [sub2sing-box](https://github.com/nitezs/sub2sing-box)

将订阅链接转换为 Clash、Clash.Meta 配置  
[预览](https://www.nite07.com/sub)

## 特性

- 开箱即用的规则、策略组配置
- 自动根据节点名称按国家划分策略组
- 多订阅合并
- 自定义 Rule Provider、Rule
- 支持多种协议
  - Shadowsocks
  - ShadowsocksR
  - Vmess
  - Vless （Clash.Meta）
  - Trojan
  - Hysteria （Clash.Meta）
  - Hysteria2 （Clash.Meta）

## 使用

### 部署

- [docker compose](./docker-compose.yml)
- 运行[二进制文件](https://github.com/nitezs/sub2clash/releases/latest)

### 配置

可以通过编辑 .env 文件来修改默认配置，docker 直接添加环境变量

| 变量名                | 说明                                           | 默认值                |
| --------------------- | ---------------------------------------------- | --------------------- |
| PORT                  | 端口                                           | `8011`                |
| META_TEMPLATE         | 默认 meta 模板文件名                           | `template_meta.yaml`  |
| CLASH_TEMPLATE        | 默认 clash 模板文件名                          | `template_clash.yaml` |
| REQUEST_RETRY_TIMES   | Get 请求重试次数                               | `3`                   |
| REQUEST_MAX_FILE_SIZE | Get 请求订阅文件最大大小（byte）               | `1048576`             |
| CACHE_EXPIRE          | 订阅缓存时间（秒）                             | `300`                 |
| LOG_LEVEL             | 日志等级，可选值 `debug`,`info`,`warn`,`error` | `info`                |
| SHORT_LINK_LENGTH     | 短链长度                                       | `6`                   |

### API

[API 文档](./API.md)

### 模板

可以通过变量自定义模板中的策略组代理节点  
具体参考下方默认模板

- `<all>` 为添加所有节点
- `<countries>` 为添加所有国家策略组
- `<地区二位字母代码>` 为添加指定地区所有节点，例如 `<hk>` 将添加所有香港节点

#### 默认模板

- [Clash](./templates/template_clash.yaml)
- [Clash.Meta](./templates/template_meta.yaml)
