### 4.1 基础设施与平台
- [x] Helm/K8s 清单：asr-gateway/cti-controller 基本模板（Deployment/Service）
- [x] Helm/K8s 清单：FreeSWITCH Deployment/Service/ConfigMap（xml_curl 绑定 config-service）
- [ ] Helm/K8s 清单：SBC、PG/Redis/CH、Prometheus/Grafana、ELK、HOMER
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
- 2025-08-10（续2）
  - [x] FreeSWITCH Helm 模板（Deployment/Service/ConfigMap）并挂载 xml_curl
  - [x] asr-gateway/cti-controller 增加 Prometheus 抓取注解
  - [x] 创建 SIPp 测试脚本目录 `tests/sipp`
  - [x] ASR Gateway 接入 FunASR 转写回调（影子/强制模式），完善指标
  - [ ] 下一步：生成 Grafana 仪表盘与 Prometheus 抓取配置、SIPp 场景脚本与CI流水线、FS镜像参数化及示例配置