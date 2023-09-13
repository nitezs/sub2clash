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

### `/clash`,`/meta`

获取 Clash/Clash.Meta 配置链接

| Query 参数     | 类型     | 是否必须              | 说明                                                                                                                                                                 |
|--------------|--------|-------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| sub          | string | sub/proxy 至少有一项存在 | 订阅链接（可以输入多个，用 `,` 分隔）                                                                                                                                              |
| proxy        | string | sub/proxy 至少有一项存在 | 节点分享链接（可以输入多个，用 `,` 分隔）                                                                                                                                            |
| refresh      | bool   | 否（默认 `false`）     | 强制刷新配置（默认缓存 5 分钟）                                                                                                                                                  |
| template     | string | 否                 | 外部模板                                                                                                                                                               |
| ruleProvider | string | 否                 | 格式 `[Behavior,Url,Group,Prepend],[Behavior,Url,Group,Prepend],...`，其中 `Group` 是该规则集所走的策略组名，`Prepend` 为 bool 类型，如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部（会调整到MATCH规则之前） | 
| rule         | string | 否                 | 格式 `[Rule,Prepend],[Rule,Prepend]...`，其中 `Prepend` 为 bool 类型，如果为 `true` 规则将被添加到规则列表顶部，否则添加到规则列表底部（会调整到MATCH规则之前）                                                   | 
| autoTest     | bool   | 否（默认 `false`）     | 指定国家策略组是否自动测速                                                                                                                                                      |
| lazy         | bool   | 否（默认 `false`）     | 自动测速是否启用 lazy                                                                                                                                                      |

## 默认模板

- [Clash](./templates/template_clash.yaml)
- [Clash.Meta](./templates/template_meta.yaml)

## 已知问题

[代理链接解析](./parser)还没有经过严格测试，可能会出现解析错误的情况，如果出现问题请提交 issue

## TODO

