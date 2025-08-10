# SBC 对接指南（OpenSIPS/Kamailio）

## 目标
- 前置SBC承载注册/路由/限速/防护/拓扑隐藏；与 FreeSWITCH 内外部profile对接；上连运营商Trunk。

## 关键点
- Topology Hiding：开启 `topology_hiding()`，移除私网头字段
- 权限与限速：`permissions` 检查来源网段、`pike`/`htable` 限速CPS
- 负载均衡：`dispatcher` 将呼叫分配到多台 FreeSWITCH（黏性可基于坐席）
- 号码规范化：E164 化或本地规则
- 健康检查：`OPTIONS` 到 FreeSWITCH

## 对接流程
1. 运营商Trunk → SBC（外侧监听 5060/5061）
2. SBC → FreeSWITCH external profile（默认 5080）
3. 坐席/WebRTC → SBC（WSS/TLS）→ FreeSWITCH internal profile（5061）

## 示例（OpenSIPS 片段）
- 入口判定→分发到 `dispatcher`；回源时增加 Header 标记；失败回退下一节点
- 限速：基于来源IP与全局CPS；异常拉黑短期

## 监控
- Exporter (Prometheus)、HepAgent（HOMER），记录SIP信令用于排障