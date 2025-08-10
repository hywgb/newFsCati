# CTI Controller 设计文档

## 1. 目标
- 管理与 FreeSWITCH 的 ESL 通信，编排外呼、控制媒体分流、接收 ASR 判定并进行快速决策。

## 2. ESL 连接
- 连接池：多路连接，心跳与自动重连
- 订阅事件：`CHANNEL_CREATE`, `CHANNEL_PROGRESS`, `CHANNEL_PROGRESS_MEDIA`, `CHANNEL_ANSWER`, `CHANNEL_HANGUP_COMPLETE`, `DTMF`, `RECORD_START/STOP`

## 3. 命令集
- 外呼：`bgapi originate {...}sofia/external/... &bridge(sofia/internal/agent@domain)`
- 分流：`bgapi uuid_audio_fork <uuid> start wss://asr-gw/stream?...`
- 停止：`bgapi uuid_audio_fork <uuid> stop all`
- 录音：`record_session=/recordings/${uuid}.wav`
- 结束：`api uuid_kill <uuid> NORMAL_CLEARING`

## 4. 状态机
- new → progress → progress_media（启动分流）→ answer（继续识别或转坐席）→ hangup（落库）
- 任意时刻：收到高置信 ASR → 结束并归档

## 5. 幂等与一致性
- 以 `call_attempt_id`/`uuid` 作为幂等键；重复命令去重
- 数据同步：先写 PG 再发命令或反之需补偿；使用 outbox/event 槽

## 6. 失败与重试
- originate 超时/失败：按策略回拨或失败落库
- audio_fork 启动失败：记录并回退到信令策略
- 结果回写失败：重试队列与告警

## 7. 指标
- `fs_active_calls`, `originate_latency_ms`, `audio_fork_attach_ratio`, `early_media_detect_ratio`, `hangup_by_asr_ratio`

## 8. 安全
- ESL 密码与网段 ACL；命令白名单；速率限制

## 9. 配置
- 租户/任务级参数：`ignore_early_media`, `call_timeout`, `asr_mode`, `asr_threshold`, `proof_seconds`