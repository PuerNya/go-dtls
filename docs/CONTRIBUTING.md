# Contributing

[简体中文](CONTRIBUTING.md) | [English](CONTRIBUTING.en.md) | [Русский](CONTRIBUTING.ru.md)

感谢参与 `go-dtls`。变更应保持 DTLS 1.3 报文语义、协议安全性和现有性能，不为没有明确需求的扩展预建框架。

## 开发环境

- Go 1.26 或更高版本。
- golangci-lint v2.12.2，与 CI 使用的版本一致。
- Windows race 测试需要 Zig 0.17 和可用的 CGO 工具链。
- 变更涉及 wolfSSL 已支持的功能时，wolfSSL 互通测试和专项 benchmark 是必需门禁；其他变更可不安装 wolfSSL。

## 必需检查

所有 Go 源码必须使用 `gofmt` 格式化。每个 commit 和 pull request 都必须通过 lint、test、vet 和 race；Linux、macOS 及 CI 使用：

```sh
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
go test ./... -count=1
go vet ./...
go test -race ./... -count=1
```

协议、并发或资源生命周期变更还应执行重复测试：

```sh
go test ./... -count=10 -timeout=10m
```

Windows 运行前三项检查，并用仓库脚本替代上面的 race 命令：

```powershell
.\tools\test-race.ps1
```

GitHub Actions 会在每次 push 和 pull request 上分别执行四个必需检查。任何一项失败都必须修复，不得通过放宽断言、跳过测试或扩大 lint 排除范围绕过。

## 协议变更

协议行为以 RFC 9147、RFC 9846 及适用的相关 RFC 为准：

- `MUST`、`MUST NOT`、`REQUIRED`、`SHALL` 和 `SHALL NOT` 必须在客户端与服务端的发送、接收方向严格执行。
- `SHOULD` 类要求由客户端主动实现；服务端可以兼容对端偏差，但不得削弱认证、机密性、反重放、放大限制或状态一致性。
- `MAY` 和 `OPTIONAL` 能力一旦协商，其条件性强制要求必须完整实现。
- 保持一条 `WriteDatagram` 对应一条 Application Data record；不得隐式引入应用数据分片、排序或重传。
- 优先复用现有 record、flight、ACK、transcript、key schedule、CID 和错误封装，不重复实现同类状态机。

协议变更至少需要覆盖：

- 正常路径，以及 malformed、truncated、duplicate、unknown 和越界输入。
- 客户端与服务端、发送与接收、完整握手与恢复路径中所有受影响分支。
- 丢包、延迟、乱序和重复；重传或状态机变更应使用高重复次数验证。
- anti-replay、放大限制、AEAD 使用上限、secret 清除和资源上限等相关安全边界。
- timeout、关闭、goroutine、内存和 Listener session 的资源生命周期。

重传、握手或协议状态机变更应使用 `-count=100` 运行相关弱网和端到端测试；资源生命周期测试应至少使用 `-count=10`。

可选扩展应以明确的部署需求和可测试的协议设计为前提，不提前加入占位 API 或空框架。

## 性能验证

涉及握手、record、解析、分配或并发热路径的变更必须同时检查微基准和完整连接。完整连接是主要性能门禁，不能只用局部基准推断整体性能。

```sh
go test -run '^$' -bench . -benchmem
go test -run '^$' -bench '^BenchmarkConnectionHandshakeLifecycle$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkMutualTLSHandshakeLifecycle/(Full|Resumed)$' -benchmem -count=10
go test -run '^$' -bench '^BenchmarkCertificateCompression' -benchmem
go test -run '^$' -bench '^BenchmarkProtectedRecord(Seal|RoundTripInPlace)$' -benchmem -count=5
```

比较修改前后性能时：

- 固定 Go 版本、CPU、`-cpu` 和 `-benchtime`。
- 预编译修改前后的测试二进制；完整连接使用 `-cpu=1`、每个样本 500 次迭代，至少交替运行 8 轮并比较中位数。
- 同时检查 `ns/op`、`B/op`、`allocs/op`、完整连接吞吐和资源生命周期。
- 分配或内存出现差异时使用 profile 定位到具体对象和调用路径。
- benchmark、profile 和临时二进制保存在仓库外，不得提交生成物。

README 只记录当前环境下的代表性能，不记录 A/B、冻结基线或修改前后叙事；这些内容应作为 pull request 的验证证据。

## 互通验证

对端支持相关能力时，应补充第三方双向端到端测试。wolfSSL 测试可使用：

```powershell
$env:WOLFSSL_ROOT = 'C:\path\to\wolfssl'
go test -run TestInteropWolfSSL -v -count=10
```

wolfSSL 支持目标功能时，还必须增加并运行真实 UDP 专项 benchmark，覆盖 go-dtls -> go-dtls、go-dtls -> wolfSSL、wolfSSL -> go-dtls 和 wolfSSL -> wolfSSL。将对端 commit/构建配置、workload、计时边界、五轮中位数及能力限制作为 pull request 验证证据；跨实现 benchmark 数据不写入 README。

第三方实现未启用或不支持目标扩展时，应明确互通边界，不得把跳过测试表述为通过。

## 文档

- 公共 API、配置、错误或协议行为变化必须同步更新 Go doc 和全部 README 语言版本。
- 中文、英文和俄文内容必须保持一致，语言切换链接必须有效。
- 各语言 README 只描述现有能力、RFC 完成度、当前限制、使用方式、代表性能和测试覆盖。
- 未实现能力只需简要列出；开发计划、优先级、调研过程和历史性能比较不写入各语言 README。

## Commit 和 pull request

每个 commit subject 必须遵循 Conventional Commits：

```text
type(scope): description
```

例如：

```text
feat(cid): implement RFC 9853 return routability checks
fix(record): reject invalid truncated sequence numbers
docs(api): clarify datagram truncation behavior
```

Pull request 应说明适用的 RFC 章节、行为变化、测试结果、race 结果、性能数据和互通范围。破坏性变更使用 `!`，并在正文中说明迁移方式。
