# sub2clash

将订阅链接转换为 Clash、Clash.Meta 配置

## 特性

- 开箱即用的规则、策略组配置
- 自动根据节点名称按国家划分策略组
- 支持多订阅合并
- 支持多种协议
    -   [x] Shadowsocks
    -   [x] ShadowsocksR
    -   [x] Vmess
    -   [x] Vless
    -   [x] Trojan
    -   [ ] Hysteria
    -   [ ] TUIC
    -   [ ] WireGuard

## API

### /clash

获取 Clash 配置链接

| Query 参数 | 类型     | 说明                      |
|----------|--------|-------------------------|
| sub      | string | 订阅链接（可以输入多个订阅，用 `,` 分隔） |
| refresh  | bool   | 强制刷新配置（默认缓存 5 分钟）       |

### /meta

获取 Meta 配置链接

| Query 参数 | 类型     | 说明                      |
|----------|--------|-------------------------|
| sub      | string | 订阅链接（可以输入多个订阅，用 `,` 分隔） |
| refresh  | bool   | 强制刷新配置（默认缓存 5 分钟）       |

## TODO

- [ ] 完善日志功能
- [ ] 支持自动测速分组
- [ ] 完善配置模板
