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
- 2025-08-10（续3）
  - [x] FreeSWITCH ConfigMap：`modules.conf.xml` 启用所需模块、`internal.xml` 配置 TLS/SRTP/NAT、`xml_curl.conf.xml` 绑定 config-service
  - [x] 部署挂载上述 ConfigMap 到容器路径
  - [ ] 下一步：完善 external profile、SBC前置对接配置、证书与密钥挂载、安全加固、性能参数（RTP队列/内核调优）