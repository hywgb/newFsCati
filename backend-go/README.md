# Backend (Go)

包含服务：
- cti-controller：接收 ASR 判定回调（/asr/decision），未来对接 ESL 控制 FreeSWITCH
- asr-gateway：WS 接入（/stream），短语词典热加载（/config/reload），未来对接 FunASR

快速开始：
- `go mod tidy`
- 运行 CTI：`go run ./cmd/cti-controller`
- 运行 ASR GW：`go run ./cmd/asr-gateway`
- 测试：`go test ./...`