# Dialer Core（预测拨号器）设计

## 1. 目标
- 最大化坐席占用，最小化放空率（Drop Rate），满足法规阈值（如 ≤ 3%）。

## 2. 输入信号
- 在线坐席数、平均通话时长（ACD）、接通率（ASR，Answer-Seizure Ratio）、队列长度、CPS 限制、时段因子、ASR 号码识别命中率与时延

## 3. 策略
- 控制器：PID 控制拨出速率，目标为目标占用率（Occupancy）
- 排队模型：Erlang A/S 估计接通与等待；
- 探索/利用：多臂老虎机，根据人群/号段/时段调整拨出比例
- 约束：Drop Rate 实时估计与硬阈值；CPS 上限；坐席可用性

## 4. 输出
- 时间片（1s）内拨出批次规模与目标 CPS；失败回拨计划

## 5. 接口（gRPC）
- `PredictNextBatch(online_agents, asr, acd, cps, backlog) -> {batch_size, target_cps}`
- `Feedback(attempt_id, result_code, talk_time)`

## 6. 容错
- 信号缺失时回退到保守策略；异常波动进入降级模式

## 7. 指标
- `target_cps`, `actual_cps`, `drop_rate`, `agent_occupancy`, `abandon_rate`