# ITACATI 级别 CATI 系统落地过程文档（process.md）

本文件是项目的单一事实来源（SSOT），记录从需求分析到设计、任务拆分、执行、测试、问题与解决、复盘总结的全流程。目标：100% 覆盖 ITACATI（电访专家）行业功能，并结合行业最佳实践与高性能架构要求，构建生产级系统。

---

## 1. 需求分析（Industry-grade，无任何简化）

### 1.1 业务目标
- 覆盖 CATI 全流程：问卷脚本 → 样本与配额 → 外呼（预览/预测/手动）→ 坐席调查 → 录音质检 → 报表分析 → 合规审计。
- 与 FreeSWITCH 深度集成，支持运营商 SIP Trunk、SBC 前置、WebRTC/软电话；结合本地 FunASR 对早期媒体与接通后的自动播报进行识别，提前判定无效呼叫状态（关机/停机/空号/不在服务区等），快速释放资源。
- 多租户、可横向扩展、观察性完善、合规与安全可控。

### 1.2 必须覆盖的 ITACATI 行业功能清单（功能映射）
- 问卷脚本与逻辑
  - 题型库：单选、多选、矩阵、开放题、量表、排序、填空、录音型问题等
  - 逻辑：分支/跳转、随机化、管道、显示逻辑、必答校验、互斥/配额联动
  - 版本管理、草稿/发布、预览与A/B版本
- 样本与配额
  - 样本导入（CSV/XLSX/接口）、清洗去重、号码验证、DNC/黑名单、时区/时段策略
  - 配额维度（地区/性别/年龄/渠道等）与动态配额分配；优先级与抽样策略
- 外呼引擎
  - 模式：手动、预览、预测（含 Drop Rate 约束）
  - 呼叫节奏：CPS、并发阈值、失败回拨策略（Busy/No Answer/Reject/Operator）
  - 早期媒体识别（FunASR）：关机/停机/空号/不在服务区/忙/彩铃/信箱/智能助理
  - 号码策略：主叫号池、区号策略、重试窗口、号码置换
- 坐席工作台
  - 软电话控制条（接/挂/保持/转接/外拨）、快捷键、状态机（Ready/Not ready/After-call）
  - 问卷执行器（高可用表单、校验/自动保存/断点续）
  - 弹屏：样本画像、历史沟通记录、配额命中提示
  - 结果码（Disposition）与预约回拨、备注/标签
- 质检与合规
  - 录音全量/抽检、质检评分卡、关键词与静默检测、双录提示校验
  - 权限/审计：字段级脱敏、操作审计、数据保留策略、导出加密
- 报表与分析
  - 实时看板：ASR/ACD、并发、CPS、命中/接通、配额进度
  - 成果报表：题项统计/交叉分析、配额完成、坐席效率、拨号效率
  - 明细导出（CSV/Excel/SPSS）、API 对接（DWH/BI）
- 系统与配置
  - 多租户、角色权限、字典/短语管理（ASR 词典）、黑名单/DNC 管理
  - FreeSWITCH/Trunk/SBC/安全策略配置；配置中心（mod_xml_curl）

### 1.3 非功能性需求
- 高可用：SBC 主备、FreeSWITCH 多实例、FunASR 多实例、数据库高可用
- 高性能：200–1000 并发扩展能力；早期媒体识别 < 1.5s；录音异步归档
- 安全与合规：TLS/SRTP，最小权限，隐私保护，审计留痕
- 可观测性：指标、日志、追踪、SIP Trace（HEP/HOMER）

---

## 2. 架构与框架设计

### 2.1 技术栈
- 后端：Go（API、CTI 控制、Agent Gateway、Survey/Quota、Analytics）、C++（预测拨号器、质检引擎、Paradata 处理）
- 通信：REST/gRPC、WebSocket、Kafka/NATS 事件总线
- 语音：FreeSWITCH（mod_sofia/mod_audio_fork/mod_event_socket）、FunASR 本地流式模型
- 数据：PostgreSQL、Redis、ClickHouse、对象存储（S3/OBS）
- 前端：React + TypeScript + Ant Design + ECharts + SIP.js（或 Verto）

### 2.2 关键子系统
- CTI Controller（Go）：ESL 连接池、外呼编排、事件处理、与 ASR Gateway 协作
- ASR Gateway（Go）：接收音频（WS），对接 FunASR，文本聚合与分类器，回传判定
- Dialer Core（C++）：预测/预览/手动调度、CPS/并发控制、Drop Rate 约束
- Survey Service（Go）：问卷引擎与脚本
- Sample & Quota（Go）：样本/配额、抽样/锁、DNC/黑名单
- QA Engine（C++/Go）：录音质检、关键词/静默检测、评分卡
- Analytics API（Go）：实时/历史报表与数据服务
- Config Service（Go）：mod_xml_curl 配置下发、系统参数

### 2.3 关键数据流
- 外呼→FS→CHANNEL_PROGRESS_MEDIA→audio_fork→ASR→决策→ESL uuid_kill/交付坐席
- 坐席→问卷提交→结果码/预约→写库与事件推送→配额更新
- 录音归档→对象存储→质检队列→评分结果入库

---

## 3. 完整建议（Industry best practices）
- 前置 SBC（OpenSIPS/Kamailio）实施拓扑隐藏、限速、号码规范化、DoS 防护、黏性路由
- FreeSWITCH 只做媒体与业务逻辑，配置集中化（mod_xml_curl），绕转码优先
- FunASR 采用 GPU 流式模型，词典热更新；ASR Gateway 解耦会话与后端模型
- 预测拨号采用自适应控制（PID/ErlangA + 多臂老虎机），约束放空率
- 数据分层：OLTP（PG）与 OLAP（CH）分离；录音走对象存储生命周期
- DevSecOps：基础设施即代码（Helm/K8s）、密钥管理（KMS）、审计
- 可观测：Prometheus 指标+ELK+HOMER，逐条路由可追踪

---

## 4. 开发任务与 TODO（分阶段、完备）

### 4.1 基础设施与平台
- [x] 建立仓库结构（backend-go, backend-cpp, frontend, deploy, docs）
- [ ] Helm/K8s 清单：SBC、FreeSWITCH、ASR Gateway、FunASR、PG/Redis/CH、Prometheus/Grafana、ELK、HOMER
- [ ] CI/CD：构建、测试、镜像、Helm 部署流水线

### 4.2 CTI 与语音层
- [ ] FreeSWITCH 配置模板化（profiles、dialplan、modules）与 mod_xml_curl 服务
- [ ] ESL 适配层（Go）：连接池、事件订阅、命令幂等
- [x] 音频分流：`uuid_audio_fork` 生命周期管理（设计文档完成）
- [ ] 号码策略与主叫号池管理
- [ ] Trunk/SBC/编解码策略与 NAT 处理

### 4.3 FunASR 集成
- [ ] 部署流式模型与 VAD；GPU/CPU 混部与健康检查
- [x] ASR Gateway：WS 接入、FunASR 客户端、文本聚合、分类器、阈值治理（设计文档完成）
- [x] 词典管理：短语库、热更新 API、租户级覆盖（初版 `config/phrases.yml`）
- [ ] 决策回写：REST/gRPC 到 CTI Controller，幂等落库

### 4.4 业务服务
- [ ] Survey Service：问卷模型、逻辑引擎、版本管理、预览
- [ ] Sample & Quota：导入/去重/抽样/配额/时间窗/DNC
- [ ] Dialer Core：预测/预览/手动、CPS/并发、失败回拨、Drop Rate 约束
- [ ] Agent Gateway：WS 推送、坐席状态机、事件分发
- [ ] QA Engine：录音质检、关键词/静默检测、评分卡
- [ ] Analytics API：实时看板与历史报表
- [ ] Config Service：XML 渲染、系统参数、租户配置

### 4.5 前端
- [ ] 坐席工作台：问卷执行器、呼叫控制、结果码、预约回拨
- [ ] 任务与配额管理：Campaign 配置、CPS/并发、失败策略、配额看板
- [ ] 问卷设计器：题型库、逻辑、随机化/管道
- [ ] 样本管理：导入/清洗、分层/优先级、黑名单/DNC
- [ ] 看板与报表：实时指标、交叉分析、明细导出
- [ ] 质检台：录音回放、打分卡、关键词/静默提示
- [ ] 系统设置：SIP/Trunk、ASR 词典、权限、合规策略

### 4.6 数据与合规
- [x] 数据库建模与迁移：PG/CH 表结构与索引（初版 DDL 文档完成）
- [ ] 对象存储与生命周期：录音归档、证据片段
- [x] 权限与审计：RBAC、字段脱敏、审计日志（规范文档完成）
- [x] 合规模块：时间窗、双录提示、隐私与留痕（规范文档完成）

### 4.7 可观测与测试
- [ ] 指标导出：FS、ASR、CTI、业务 KPI
- [ ] 日志与追踪：结构化日志、链路ID、SIP Trace
- [x] 压测与验收测试计划文档
- [ ] 验收：覆盖率、端到端脚本、回滚策略

---

## 5. 执行记录（持续更新）
- 2025-08-10
  - [x] 初始化文档目录与 process.md
  - [x] 创建功能映射文档 `docs/itacati-feature-mapping.md`
  - [x] 输出 `architecture.md` 与 `test-plan.md`
  - [x] 生成 `config/phrases.yml` 初版
  - [x] 创建项目骨架目录结构（backend-go/backend-cpp/frontend/deploy/config）
  - [x] 新增设计文档：`asr-gateway.md`, `cti-controller.md`, `dialer-core.md`, `frontend-spec.md`, `security-compliance.md`, `database-schema.md`
  - [x] 生成 Helm Chart 骨架（asr-gateway/cti-controller/freeswitch）
  - [ ] 下一步：实现 CTI Controller ESL 适配与 ASR 回调 API，编写单元与集成测试

---

## 6. 测试计划（总览）
- 单元测试：问卷逻辑、抽样/配额、ESL 命令生成、分类器规则
- 集成测试：FS 事件→audio_fork→ASR→决策链；录音归档一致性
- 性能测试：并发 200/500/1000 路媒体，ASR 延迟与命中率；Dialer 稳定性与 Drop Rate
- 端到端：任务全生命周期；前端交互稳定性；数据一致性与报表校验
- 安全与合规：TLS/SRTP、最小权限、审计、数据保留

---

## 7. 问题与解决（将随项目推进更新）
- 待记录

---

## 8. 总结（阶段性占位）
- 待记录