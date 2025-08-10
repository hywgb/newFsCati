### 4.1 基础设施与平台
- [x] 建立仓库结构（backend-go, backend-cpp, frontend, deploy, docs）
- [x] Helm/K8s 清单：asr-gateway/cti-controller 基本模板（Deployment/Service）
- [ ] Helm/K8s 清单：FreeSWITCH、SBC、PG/Redis/CH、Prometheus/Grafana、ELK、HOMER
- [ ] CI/CD：构建、测试、镜像、Helm 部署流水线
### 4.2 CTI 与语音层
- [x] ESL 适配层（Go）：连接池、事件订阅、命令幂等（初版）
- [x] 音频分流：`uuid_audio_fork` 生命周期管理（在 PROGRESS_MEDIA 触发）
### 4.3 FunASR 集成
- [x] ASR Gateway：WS 接入、短语匹配、CTI 回调与指标
- [x] FunASR 客户端骨架（WS，占位 shadow 模式）
- [ ] 接入 FunASR 真流与转写解析，影子模式→强制模式
### 4.6 数据与合规
- [x] mod_xml_curl 规范文档 `docs/mod_xml_curl_spec.md`
## 5. 执行记录（持续更新）
- 2025-08-10（续）
  - [x] Config Service 增加 `/dialplan`，形成目录/拨号方案双接口
  - [x] ASR Gateway 指标增强：首判延迟/决策计数
  - [x] CTI 指标增强：PROGRESS_MEDIA 事件计数
  - [x] Helm 模板：asr-gateway 与 cti-controller 的 Deployment/Service
  - [x] FunASR 客户端骨架（WS）
  - [ ] 下一步：FreeSWITCH 配置模板与 Helm、FunASR 真流接入与影子/强制切换、SIPp 集成测试