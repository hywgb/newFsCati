# ASR Gateway 设计文档

## 1. 目标与职责
- 从 FreeSWITCH `mod_audio_fork` 接收流式音频（WS），对接本地 FunASR 流式识别服务。
- 增量转写、文本聚合、短语匹配、置信度融合，与信令原因码互补。
- 在置信度达标时快速返回判定给 CTI Controller；支持影子模式与证据留存。

## 2. 接口设计
### 2.1 来自 FreeSWITCH（WS）
- URL：`wss://asr-gw:10000/stream?uuid={uuid}&sr=8000&channels=1&ptime=20`
- 输入：二进制音频帧（PCM 8k，20ms），辅以 JSON metadata（可扩展）
- 返回：心跳/ack（可选）

### 2.2 FunASR 后端
- 连接：gRPC 或 WebSocket（根据部署而定），连接池与健康检查
- 功能：VAD、Streaming ASR、可选 PUNC

### 2.3 判定回传（到 CTI Controller）
- REST：`POST /asr/decision`
- 消息：
```json
{
  "uuid": "...",
  "result": "power_off|suspended|invalid_number|out_of_service|busy|voicemail|color_ring_back|operator_ivr|answered_human",
  "confidence": 0.0,
  "latency_ms": 0,
  "transcript": "",
  "mode": "early|post_answer",
  "audio_proof_uri": "s3://.../proof.wav",
  "fallback": false
}
```

## 3. 分类与融合
- 规则优先：短语词典（`config/phrases.yml`）命中即早停；支持正则/模糊匹配
- 统计兜底：累计置信度 ≥ θ（默认 0.75），窗口 1.2–1.5s
- 信令融合：若已收到明确 Q.850/SIP 失败码，直接映射；否则以 ASR 高置信覆盖

## 4. 会话与容错
- 会话：按 `uuid` 一致性路由，跟踪起止、速率、统计
- 超时：无音频/无结果超时 → fallback；记录原因
- 重试：FunASR 连接断开自动快速重连；实例健康探测与摘除

## 5. 证据留存
- 环形缓冲：最近 N 秒原始音频（例如 10s）
- 命中后：将证据片段（音频+转写）异步上传对象存储并返回 URI

## 6. 性能指标
- 端到端首判 ≤ 1.5s；强规则 ≤ 800ms
- 并发：A10 单卡 80–120 路；A100 200+；CPU 节点 20–40 路（仅回退）
- 背压：当 FunASR 队列过长时主动拒绝新会话（返回 503），CTI 退回信令策略

## 7. 配置与治理
- `phrases.yml` 热更新（租户级覆盖），灰度生效
- 影子模式：仅观察不挂断，记录预测与误判
- 参数：阈值、窗口、早停、证据秒数、最大会话

## 8. 指标与日志
- Prom：`asr_decision_latency_ms`, `asr_confidence`, `asr_hit_rate`, `active_sessions`, `backend_queue_len`
- 日志：结构化，含 `uuid`/tenant/campaign；错误带上下文

## 9. 安全
- WSS + mTLS（可选）；鉴权（token/签名）；IP ACL

## 10. 灰度与回滚
- 按租户/任务开关；A/B；一键禁用 audio_fork 路由