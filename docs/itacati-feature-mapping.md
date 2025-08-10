# ITACATI 功能覆盖与实现映射（Feature Mapping）

目标：确保系统实现 100% 覆盖 ITACATI（电访专家）现有功能，并结合行业最佳实践扩展，提供一一对应的实现方案与责任归属。

## 1. 核心域功能

### 1.1 问卷与访谈脚本
- 覆盖项：题型库、分支/跳转、随机化、管道、必答、显示逻辑、版本管理、预览、A/B
- 实现：
  - 服务：`survey-service`（Go）
  - 数据：`surveys`, `survey_questions`, `survey_versions`
  - 前端：问卷设计器（React+TS），执行器在坐席工作台渲染
  - 额外：Form autosave、断点续、前后端一致性校验

### 1.2 样本与配额
- 覆盖项：批量导入、清洗去重、号码校验、DNC/黑名单、时区、时间窗策略、配额维度与动态分配
- 实现：
  - 服务：`sample-quota`（Go）
  - 数据：`samples`, `campaign_samples`, `quotas`
  - 规则：租户级 DNC，号码归属/时区库（本地 Geo DB）

### 1.3 外呼引擎
- 覆盖项：手动、预览、预测；CPS 与并发控制；失败回拨策略；主叫号池；号码置换
- 实现：
  - 服务：`dialer-core`（C++）调度；`cti-controller`（Go）执行
  - 数据：`campaigns`, `call_attempts`, `callerid_pools`
  - 策略：PID + ErlangA + 多臂老虎机；运营商限速自适应

### 1.4 号码状态识别（FunASR）
- 覆盖项：关机/停机/空号/不在服务区/忙/彩铃/信箱/智能助理
- 实现：
  - 服务：`asr-gateway`（Go） + 本地 FunASR（GPU）
  - FS：`mod_audio_fork` 在 progress-media 即分流
  - 数据：`call_attempts` 扩展 asr 字段；`asr_events`

### 1.5 坐席工作台
- 覆盖项：软电话控制、脚本执行、弹屏、结果码、预约回拨、快捷键、状态机
- 实现：
  - 前端：`agent-app`（React+TS+SIP.js）
  - 服务：`agent-gateway`（WS）、`cti-controller`
  - 数据：`agents`, `dispositions`, `callbacks`

### 1.6 质检与合规
- 覆盖项：录音全量/抽检、评分卡、关键词与静默检测、双录提示、审计日志、脱敏、留存
- 实现：
  - 服务：`qa-engine`（C++/Go）、`audit-service`（Go）
  - 数据：`qa_scores`, `audit_logs`
  - 存储：对象存储分桶与生命周期、KMS 加密

### 1.7 报表与分析
- 覆盖项：实时看板、效率与配额进度、题项统计与交叉、导出（CSV/Excel/SPSS）、API
- 实现：
  - 服务：`analytics-api`（Go） + ClickHouse
  - 前端：`admin-app` 看板 + 报表

### 1.8 系统与配置
- 覆盖项：多租户、RBAC、字典/短语管理、黑名单/DNC 管理、Trunk/SBC/安全
- 实现：
  - 服务：`identity`（Go）、`config-service`（Go）
  - FS：`mod_xml_curl` 拉取目录/拨号方案

## 2. 最佳实践增强项
- 早期媒体快速识别与提前挂断（节省时长）
- 影子模式（observe-only）上线，A/B 比对误判率
- 证据片段留存（音频/转写），支持申诉与训练
- 可观测性：全链路指标、SIP Trace、事件审计
- DevSecOps：CI/CD、IaC、密钥管理、审计留痕

## 3. 合规映射
- 外呼时间窗、双录提示、DNC、隐私保护、数据保留与销毁

## 4. 验收标准
- 覆盖度：上述功能点全部可用并通过测试用例
- 性能：并发与延迟指标达标；预测拨号 Drop Rate 符合法规
- 准确性：ASR 分类综合准确率 ≥ 97%，强规则 ≥ 99%