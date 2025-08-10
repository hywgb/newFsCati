# 架构说明（architecture.md）

## 1. 总览
- 拓扑：终端（坐席/WebRTC）→ SBC（OpenSIPS/Kamailio）→ FreeSWITCH 集群 → 运营商 SIP Trunk
- 业务面：Go/C++ 微服务（CTI、Dialer、Survey、Sample&Quota、QA、Analytics、Config、ASR Gateway）
- 数据面：PostgreSQL（OLTP）、Redis（缓存/锁）、ClickHouse（OLAP）、对象存储（录音）
- AI：本地 FunASR 流式识别（GPU），ASR Gateway 解耦

## 2. 组件职责
- SBC：注册/路由/限速/安全/拓扑隐藏；黏性路由到媒体节点
- FreeSWITCH：媒体锚定、录音、IVR、事件；`mod_audio_fork` 提供音频分流
- CTI Controller：ESL 连接池，外呼编排与状态同步
- Dialer Core（C++）：预测拨号与 CPS/并发/策略控制
- Survey Service：问卷/脚本/逻辑引擎
- Sample & Quota：样本池与配额控制
- ASR Gateway：WS 音频接入，FunASR 客户端，分类与判定
- QA Engine：录音质检、关键词/静默检测
- Analytics API：实时/历史指标；数据取自 ClickHouse
- Config Service：集中化目录与拨号方案（mod_xml_curl）

## 3. 核心数据流
1) 外呼与识别
- originate → progress-media → audio_fork → ASR → 判定 → uuid_kill 或交付坐席
2) 坐席执行
- 弹屏问卷 → 提交 → 结果码/预约 → 数据落库与事件
3) 录音与质检
- 接通录音 → 异步归档到对象存储 → 质检任务入队 → 评分与关键词

## 4. 部署与高可用
- 多可用区：SBC/VIP 高可用，FreeSWITCH 无共享会话水平扩展
- 数据层：Postgres（Patroni）、Redis（Sentinel/Cluster）、ClickHouse 多副本
- FunASR：多 GPU 实例，负载均衡与健康检查，亲和策略
- 配置：Helm/K8s；媒体节点资源隔离（CPU 亲和，IRQ 亲和）

## 5. 性能与可靠性
- RTP 端口固定；绕转码；Opus/WebRTC 与 G.711/PSTN 策略
- ASR 目标：<1.5s 判定；连接池与 backpressure；早停策略
- 录音异步上传与重试；S3 分区/冷热分层
- 预测拨号：Drop Rate 守护，CPS 自适应，SBC 限速兜底

## 6. 安全与合规
- TLS/SRTP，WSS（WebRTC），IP ACL；最小权限
- 审计日志、字段脱敏、隐私数据加密；外呼时间窗、双录提示

## 7. 可观测性
- 指标：FS/ASR/CTI/业务 KPI；Prometheus + Grafana
- 日志：结构化日志到 ELK；HOMER 采集 SIP Trace（HEP）
- 报警：Trunk 可用性、ASR 命中/延迟异常、录音积压、CPS 超阈值