# 迷你虾 CLI（minixia）

跨平台的 **[迷你虾](https://minixia.app)** 命令行工具：发送消息、拉取/确认指令、Webhook 与 MQTT 集成，以及配置与多环境（profile）管理。

本仓库：**[github.com/jinwoll/minixia-cli](https://github.com/jinwoll/minixia-cli)**  
模块路径：`github.com/jinwoll/minixia-cli`  
二进制名称：`minixia`（Windows 安装后为 `minixia.exe`）。

---

## 环境要求

- **从源码编译**：安装 [Go](https://go.dev/dl/)，版本不低于本仓库 `go.mod` 中的 `go` 行（当前为 **1.26.1**）。
- **一键安装脚本**：需要 `curl` 或 `wget`（Linux / macOS / Git Bash），或 **PowerShell 5.1+**（Windows）。
- **本地打全平台包**：安装 [GoReleaser](https://goreleaser.com/install/)（可选，用于交叉编译发布产物）。

---

## 安装

### 方式一：GitHub 一键安装（推荐）

安装脚本会从 **GitHub Releases 最新版** 的「附件」里下载二进制，因此你需要先至少发布过一次 Release（见下文「本地构建并发布」）。

**Linux / macOS / WSL / Git Bash**（默认分支名 `main`，若你用 `master` 请改 URL）：

```sh
curl -fsSL https://raw.githubusercontent.com/jinwoll/minixia-cli/main/install.sh | sh
```

**Windows PowerShell**：

```powershell
irm https://raw.githubusercontent.com/jinwoll/minixia-cli/main/install.ps1 | iex
```

安装脚本支持的环境变量（节选）：

| 变量 | 说明 |
|------|------|
| `MINIXIA_INSTALL_DIR` | 自定义安装目录 |
| `MINIXIA_BINARY_URL` | 二进制根地址（默认 `https://github.com/jinwoll/minixia-cli/releases/latest/download`） |
| `MINIXIA_VERSION` | 仅 shell 脚本保留；当前下载逻辑以 `latest/download` 为准 |

自用 fork 时，把脚本里的仓库Owner/仓库名换成你的，或导出 `MINIXIA_BINARY_URL=https://github.com/you/minixia-cli/releases/latest/download`。

### 方式二：`go install`

```sh
go install github.com/jinwoll/minixia-cli@latest
```

可执行文件在 `$(go env GOPATH)/bin`，请将该目录加入 `PATH`。

### 方式三：克隆后本地单平台编译

```sh
git clone https://github.com/jinwoll/minixia-cli.git
cd minixia-cli
go build -o minixia .
```

Windows：`go build -o minixia.exe .`

带版本信息：

```sh
VERSION=0.1.0
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -o minixia \
  -ldflags "-s -w -X github.com/jinwoll/minixia-cli/cmd.version=${VERSION} -X github.com/jinwoll/minixia-cli/cmd.commit=${COMMIT} -X github.com/jinwoll/minixia-cli/cmd.date=${DATE}" \
  .
```

---

## 本地构建所有平台的安装包

在同一台机器上交叉编译 **linux / darwin / windows** 的 **amd64 / arm64**，使用项目里的 `.goreleaser.yml`。

### 1. 仅本地打包（不上传）

适合自测 `dist/` 产物：

```sh
cd minixia-cli
goreleaser release --snapshot --clean
```

完成后 **`dist/`** 目录中会包含：

- `minixia-linux-x86_64`、`minixia-linux-arm64`
- `minixia-darwin-x86_64`、`minixia-darwin-arm64`
- `minixia-windows-x86_64.exe`、`minixia-windows-arm64.exe`
- `checksums.sha256`

命名与 `install.sh`、`install.ps1`、`minixia upgrade` 一致。

### 2. 发布到 GitHub Releases（让别人能一键安装）

前提：代码已推到 `https://github.com/jinwoll/minixia-cli.git`，且 **`install.sh` / `install.ps1` 已在默认分支上**（别人执行的是 raw GitHub 上的脚本）。

**步骤 A — 打 tag 并用 GoReleaser 自动发版（推荐）**

1. 提交并推送当前代码到 `main`。
2. 创建并推送版本 tag（与前端显示的版本号一致，建议带 `v`）：

   ```sh
   git tag v0.1.0
   git push origin v0.1.0
   ```

3. 在 [GitHub 新建 Personal Access Token](https://github.com/settings/tokens)，勾选 **`repo`**（对私有仓库）或对当前公有仓库具备 **Contents**、**Releases** 写权限。
4. 在终端设置环境变量后执行正式发行（**不要用** `--snapshot`）：

   ```sh
   export GITHUB_TOKEN=ghp_xxxx   # Windows PowerShell: $env:GITHUB_TOKEN="ghp_..."
   goreleaser release --clean
   ```

GoReleaser 会根据当前 tag 编译全部平台，并创建/更新 **GitHub Release**，上传上述二进制与 `checksums.sha256`。

**步骤 B — 手动上传（无 GoReleaser 或想完全手控）**

1. 执行 `goreleaser release --snapshot --clean` 得到 `dist/`。
2. 在 GitHub 网页 **Releases → Draft a new release**，Tag 填 `v0.1.0`，把 `dist/` 里 **全部** `minixia-*` 文件和 **`checksums.sha256`** 拖到附件区，再发布。

或用 [GitHub CLI](https://cli.github.com/)：

```sh
gh release create v0.1.0 dist/checksums.sha256 dist/minixia-* --repo jinwoll/minixia-cli --generate-notes
```

> **注意**：`releases/latest/download` 只会解析「标为 Latest 的那一个 Release」。新建版本后 GitHub 通常会把最新一条标成 Latest；若你手工取消 Latest，一键安装会仍指向 GitHub 认定的 latest。

---

## 快速开始

1. `minixia init` — 交互式保存 API Key、角色、服务器地址。  
2. `minixia send "Hello, World!"`  
3. `minixia version` 或 `minixia --version`

---

## 配置说明

### 配置目录

- **Windows**：`%LOCALAPPDATA%\minixia\`
- **Linux / macOS**：`~/.minixia/`

内含 `config.toml` 与 `profiles/<名>.toml`。

### 环境变量

| 变量 | 含义 |
|------|------|
| `MINIXIA_APIKEY` | API 密钥 |
| `MINIXIA_ROLE` | 角色名 |
| `MINIXIA_BASE_URL` | 服务器根地址 |
| `MINIXIA_PROFILE` | 使用的 profile 名 |
| `MINIXIA_GITHUB_OWNER` | （可选）`upgrade` 查询的仓库 Owner，默认 `jinwoll` |
| `MINIXIA_GITHUB_REPO` | （可选）仓库名，默认 `minixia-cli` |

合并顺序：**命令行 flag > 环境变量 > profile 文件 > 默认值**。

### 全局命令行参数

```text
  -k, --apikey string    API 密钥
  -r, --role string      角色名称
      --base-url string  服务器地址
      --profile string   指定 profile
      --verbose          详细日志
      --debug            调试信息
```

---

## 命令一览

| 命令 | 说明 |
|------|------|
| `minixia init` | 交互式创建默认 profile |
| `minixia send` | 发送消息 |
| `minixia query` | 轮询拉取指令 |
| `minixia ack` | 确认指令 |
| `minixia webhook` | Webhook 管理 |
| `minixia mqtt` | MQTT 订阅指令 |
| `minixia status` | 服务健康检查 |
| `minixia config` | profile 管理 |
| `minixia upgrade` | 从 **GitHub Releases** 检查并替换当前二进制（`--check` 仅检查） |
| `minixia uninstall` | 卸载 |
| `minixia version` | 版本信息 |

### 子命令示例

```sh
minixia send "纯文本"
minixia send --type image ./screenshot.png
echo "管道" | minixia send -
minixia query --watch --interval 5
minixia webhook --url https://example.com/hook
minixia mqtt --broker tcp://127.0.0.1:1883
```

更多见 `minixia <cmd> -h`。

---

## 许可证

发布到 GitHub 前建议在仓库根目录添加 `LICENSE`；若有该文件，以其中条款为准。

---

## 相关文档

- 同仓库 `PRD.md`（若存在）。
- 服务端行为以迷你虾后端与官方 API 文档为准。
