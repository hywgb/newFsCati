<<<<<<< Current (Your changes)
# newFsCati
base sip freeswtich cati 
=======
# CATI 生产系统（FreeSWITCH + FunASR + ITACATI 功能覆盖）

本项目旨在实现 100% 覆盖 ITACATI 行业功能的 CATI 生产系统，结合 FreeSWITCH 与本地 FunASR，实现高性能外呼与早期媒体识别、预测拨号、坐席工作台、录音质检与报表分析。

- 设计与过程：见 `docs/process.md`
- 功能映射：见 `docs/itacati-feature-mapping.md`
- 架构说明：见 `docs/architecture.md`
- 测试计划：见 `docs/test-plan.md`
- 配置词典：`config/phrases.yml`

目录结构（示意）：
- backend-go：Go 微服务（CTI、Survey、Quota、Analytics、ASR Gateway 等）
- backend-cpp：C++ 组件（Dialer Core、质检引擎）
- frontend：坐席与管理前端
- deploy/helm：K8s 部署清单（Helm Charts）
- docs：项目文档
- config：ASR 词典等配置 
>>>>>>> Incoming (Background Agent changes)
