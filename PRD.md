# 迷你虾 CLI（minixia）— 产品需求文档

> **版本**：2.0 · **日期**：2026-03-31
> **产品名称**：minixia CLI
> **定位**：跨平台命令行工具，让开发者在 Windows / macOS / Ubuntu 下一键安装并使用迷你虾全部能力

---

## 0. 术语表

| 术语 | 说明 |
|------|------|
| **apikey** | 用户的迷你虾 API 密钥，用于验证身份与路由消息 |
| **role** | 机器人实例标识，同一 apikey 下可有多个 role |
| **Base URL** | 迷你虾服务端根地址，默认 `https://api.minixia.app` |
| **profile** | 本地配置档，存储 apikey / role / base_url 等，避免每次手动输入 |

---

## 1. 产品目标

1. 提供 **一行命令即可安装** 的体验（`curl ... | sh`），覆盖 Windows（WSL / Git Bash / PowerShell）、macOS、Ubuntu/Debian。
2. 封装迷你虾全部 HTTP API，让用户无需手写 curl 即可完成：发送消息、轮询指令、Webhook 配置、MQTT 订阅。
3. 支持 **配置持久化**，一次 `minixia init` 后，后续命令无需重复输入 apikey。
4. 提供丰富的快捷用法与管道支持，便于集成进 CI/CD、cron、脚本。

---

## 2. 安装与卸载

### 2.1 一键安装（install.sh）

用户在终端执行一行命令即可完成安装：

```bash
curl -fsSL https://cli.minixia.app/install.sh | sh
```

**install.sh 职责：**

1. **检测操作系统与 CPU 架构**
   - `uname -s` → `Linux` / `Darwin` / `MINGW*|MSYS*|CYGWIN*`（Windows Git Bash / WSL）
   - `uname -m` → `x86_64` / `arm64` / `aarch64`
2. **下载对应的预编译二进制**
   - 下载地址格式：`https://cli.minixia.app/releases/latest/minixia-{os}-{arch}{ext}`
   - 例：`minixia-darwin-arm64`、`minixia-linux-x86_64`、`minixia-windows-x86_64.exe`
3. **校验完整性**（SHA256 校验和）
4. **安装到系统 PATH**
   - Linux / macOS：默认移动到 `/usr/local/bin/minixia`（无权限时回退到 `~/.local/bin/`）
   - Windows（WSL / Git Bash）：安装到 `~/.local/bin/` 并提示添加 PATH
5. **赋予执行权限**（`chmod +x`）
6. **验证安装成功**（执行 `minixia --version`）
7. **打印后续指引**（提示运行 `minixia init`）

#### 2.1.1 Windows PowerShell 安装

针对原生 Windows（非 WSL）提供 PowerShell 安装方式：

```powershell
irm https://cli.minixia.app/install.ps1 | iex
```

**install.ps1 职责：**

1. 检测架构（`$env:PROCESSOR_ARCHITECTURE`）
2. 下载 `.exe` 二进制到 `$env:LOCALAPPDATA\minixia\`
3. 自动将该目录添加到用户 PATH（`[Environment]::SetEnvironmentVariable`）
4. 验证安装成功

#### 2.1.2 包管理器安装（未来规划）

```bash
# macOS
brew install minixia

# Ubuntu/Debian
sudo apt install minixia

# Windows
scoop install minixia
# 或
winget install minixia
```

### 2.2 卸载

```bash
minixia uninstall          # 删除二进制与配置目录（交互确认）
minixia uninstall --force  # 跳过确认直接删除
```

### 2.3 自更新

```bash
minixia upgrade            # 检查并升级到最新版本
minixia upgrade --check    # 仅检查，不执行升级
```

---

## 3. 配置管理

### 3.1 交互式初始化

```bash
minixia init
```

交互流程：

```
🦐 欢迎使用迷你虾 CLI！

? 请输入你的 API Key：ak_xxxxxxxxxxxx
? 请输入角色名称（默认 bot）：my-robot
? 请输入服务器地址（默认 https://api.minixia.app）：
? 是否设为默认 profile？(Y/n)：Y

✅ 配置已保存至 ~/.minixia/config.toml
🚀 快速开始：minixia send "Hello, World!"
```

### 3.2 配置文件结构

配置目录：`~/.minixia/`

```
~/.minixia/
├── config.toml      # 全局配置（当前 profile、通用设置）
└── profiles/
    ├── default.toml  # 默认 profile
    └── work.toml     # 自定义 profile
```

**config.toml 示例：**

```toml
current_profile = "default"
auto_upgrade_check = true
```

**profiles/default.toml 示例：**

```toml
apikey = "ak_xxxxxxxxxxxx"
role = "bot"
base_url = "https://api.minixia.app"
```

### 3.3 配置命令

```bash
minixia config list                       # 列出所有 profile
minixia config show                       # 显示当前 profile 详情
minixia config set <key> <value>          # 修改当前 profile 的配置项
minixia config use <profile_name>         # 切换 profile
minixia config create <profile_name>      # 交互式创建新 profile
minixia config delete <profile_name>      # 删除 profile
```

### 3.4 配置优先级

**命令行参数 > 环境变量 > profile 配置文件 > 默认值**

| 来源 | 示例 |
|------|------|
| 命令行 | `--apikey ak_xxx` |
| 环境变量 | `MINIXIA_APIKEY=ak_xxx` |
| profile | `~/.minixia/profiles/default.toml` |
| 默认值 | role=`bot`、base_url=`https://api.minixia.app` |

支持的环境变量：

| 环境变量 | 对应配置 |
|----------|----------|
| `MINIXIA_APIKEY` | apikey |
| `MINIXIA_ROLE` | role |
| `MINIXIA_BASE_URL` | base_url |
| `MINIXIA_PROFILE` | 指定使用的 profile |

---

## 4. 核心命令

### 4.1 发送消息：`minixia send`

```bash
minixia send [OPTIONS] <content>
```

**参数与选项：**

| 选项 | 短写 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `<content>` | — | 是 | — | 消息内容（文本 / 文件路径） |
| `--apikey` | `-k` | 否* | profile | API 密钥 |
| `--role` | `-r` | 否 | profile 或 `bot` | 角色名称 |
| `--type` | `-t` | 否 | `text` | 消息类型：`text` / `image` / `voice` / `voice_call` |
| `--level` | `-l` | 否 | `info` | 消息级别：`info` / `warning` / `error` |
| `--group` | `-g` | 否 | `""` | 消息分组 |
| `--url` | `-u` | 否 | `""` | 关联外部链接 |
| `--message-id` | `-m` | 否 | 自动 UUID | 幂等键 |
| `--silent` | `-s` | 否 | false | 静默模式，仅输出 JSON 结果 |
| `--dry-run` | | 否 | false | 仅打印将要发送的请求，不实际发送 |

\* 若已通过 `minixia init` 配置，则无需每次传递。

#### 4.1.1 发送文本

```bash
# 最简用法（已配置 profile 时）
minixia send "Hello, World!"

# 完整参数
minixia send --apikey ak_xxx --role "bot ceo" --level warning --group ops "服务器 CPU 超过 90%"

# 从 stdin 读取（管道支持）
echo "部署完成" | minixia send --role deployer -

# 从文件读取文本
minixia send --role reporter --type text @report.txt
```

#### 4.1.2 发送图片

```bash
# 本地图片文件（自动 base64 编码）
minixia send --type image ./screenshot.png

# 指定角色和分组
minixia send -t image -r monitor -g dashboard ./chart.png
```

**处理逻辑：**
- 检测文件是否存在且为合法图片格式（png / jpg / jpeg / gif / webp）
- 读取文件并进行 base64 编码
- 校验编码后大小是否超出套餐限制，超出则提示

#### 4.1.3 发送语音

```bash
# 文本转语音（TTS）后发送
minixia send --type voice "欢迎使用迷你虾"

# 直接发送音频文件
minixia send --type voice ./notification.mp3
```

**处理逻辑：**
- 若 content 为纯文本：对文本进行 URL 编码（如查询参数或约定字段中传输）后提交 **server**；Azure 语音合成（TTS）由 **mini-claw-server** 调用 Azure 完成，CLI 不持有 TTS 密钥
- 若 content 为音频文件路径（.mp3 / .wav / .ogg）：读取文件并 base64 编码 → 提交 **server**
- Azure TTS 相关配置（`AZURE_TTS_KEY`、`AZURE_TTS_REGION`、`AZURE_TTS_VOICE` 等）在 **mini-claw-server** 侧通过环境变量或配置文件维护，详见服务端部署说明

#### 4.1.4 发送语音通话邀请

```bash
# 文本转语音后作为来电内容
minixia send --type voice_call "你有一条新的告警通知"

# 使用音频文件
minixia send --type voice_call ./alert.mp3
```

**处理逻辑：**
- 与 voice 类似，但最终 content 格式为 `data:audio/<subtype>;base64,...`
- App 收到后直达语音通话页面

### 4.2 轮询指令：`minixia query`

```bash
minixia query [OPTIONS]
```

**选项：**

| 选项 | 短写 | 默认值 | 说明 |
|------|------|--------|------|
| `--apikey` | `-k` | profile | API 密钥 |
| `--role` | `-r` | profile 或 `bot` | 角色名称 |
| `--limit` | `-n` | `20` | 单次拉取条数 |
| `--watch` | `-w` | false | 持续轮询模式 |
| `--interval` | `-i` | `5` | watch 模式下轮询间隔（秒） |
| `--auto-ack` | | false | 拉取后自动确认 |
| `--exec` | `-e` | — | 对每条指令执行的 shell 命令（`$CONTENT`、`$TYPE`、`$CMD_ID` 可用） |
| `--output` | `-o` | `table` | 输出格式：`table` / `json` / `raw` |

#### 4.2.1 单次拉取

```bash
# 拉取一次指令
minixia query

# 指定拉取条数
minixia query --limit 5

# JSON 格式输出（适合脚本处理）
minixia query -o json
```

输出示例：

```
┌──────────────────┬──────────┬──────┬──────────────────────────┬─────────────────────┐
│ CMD_ID           │ ROLE     │ TYPE │ CONTENT                  │ TIME                │
├──────────────────┼──────────┼──────┼──────────────────────────┼─────────────────────┤
│ cmd-uuid-001     │ bot      │ text │ 帮我总结今日报表          │ 2026-03-31 14:30:00 │
│ cmd-uuid-002     │ bot      │ text │ 重启服务                  │ 2026-03-31 14:32:00 │
└──────────────────┴──────────┴──────┴──────────────────────────┴─────────────────────┘
📋 共 2 条指令。使用 minixia ack <CMD_ID> 确认已处理。
```

#### 4.2.2 持续监听模式

```bash
# 每 5 秒轮询一次，有新指令则打印
minixia query --watch

# 自定义间隔
minixia query --watch --interval 10

# 监听并自动执行脚本
minixia query --watch --exec 'echo "收到指令: $CONTENT" >> /var/log/minixia.log'

# 监听 + 自动确认 + 执行
minixia query --watch --auto-ack --exec './handle_command.sh "$CONTENT" "$TYPE"'
```

#### 4.2.3 确认指令

```bash
# 确认单条
minixia ack <cmd_id>

# 确认多条
minixia ack <cmd_id_1> <cmd_id_2> <cmd_id_3>

# 确认全部已拉取的指令
minixia ack --all
```

### 4.3 Webhook 配置：`minixia webhook`

```bash
minixia webhook [OPTIONS]
```

**选项：**

| 选项 | 短写 | 说明 |
|------|------|------|
| `--apikey` | `-k` | API 密钥 |
| `--role` | `-r` | 角色名称 |
| `--url` | `-u` | Webhook 回调地址（HTTPS） |
| `--remove` | | 移除已配置的 Webhook |
| `--test` | | 发送测试回调验证连通性 |

```bash
# 设置 Webhook
minixia webhook --url https://my-server.com/webhook

# 测试 Webhook
minixia webhook --test

# 移除 Webhook
minixia webhook --remove
```

### 4.4 MQTT 连接：`minixia mqtt`

```bash
minixia mqtt [OPTIONS]
```

**选项：**

| 选项 | 短写 | 说明 |
|------|------|------|
| `--apikey` | `-k` | API 密钥 |
| `--role` | `-r` | 角色名称 |
| `--broker` | `-b` | MQTT Broker 地址 |
| `--username` | | MQTT 用户名 |
| `--password` | | MQTT 密码 |
| `--exec` | `-e` | 对每条指令执行的 shell 命令 |

```bash
# 连接并订阅指令
minixia mqtt --broker mqtt://broker.minixia.app:1883

# 连接并自动处理
minixia mqtt --broker mqtt://broker.minixia.app:1883 --exec './handle.sh "$CONTENT"'
```

运行后 CLI 将作为前台进程持续运行，订阅 `cmd/{apikey}/{role}` topic，收到消息后输出到 stdout 或执行 `--exec` 指定的命令。

### 4.5 健康检查：`minixia status`

```bash
minixia status
```

输出示例：

```
🦐 迷你虾服务状态

  服务器：https://api.minixia.app
  状态  ：✅ 在线（延迟 42ms）
  版本  ：CLI v1.2.0
  Profile：default
  API Key：ak_xxxx****xxxx
  角色  ：bot
```

### 4.6 帮助与版本

```bash
minixia --help              # 全局帮助
minixia <command> --help    # 子命令帮助
minixia --version           # 版本号
```

---

## 5. 跨平台支持矩阵

### 5.1 安装方式

| 平台 | 安装方式 | Shell 环境 |
|------|----------|------------|
| **macOS** (Intel/Apple Silicon) | `curl ... \| sh` / `brew` | zsh / bash |
| **Ubuntu / Debian** | `curl ... \| sh` / `apt` | bash |
| **其他 Linux**（CentOS/Fedora/Arch） | `curl ... \| sh` | bash / zsh |
| **Windows + WSL** | `curl ... \| sh`（WSL 内） | bash |
| **Windows + Git Bash** | `curl ... \| sh` | bash |
| **Windows PowerShell** | `irm ... \| iex` | PowerShell 5.1+ / pwsh 7+ |

### 5.2 预编译目标

| 操作系统 | 架构 | 文件名 |
|----------|------|--------|
| Linux | x86_64 | `minixia-linux-x86_64` |
| Linux | arm64 | `minixia-linux-arm64` |
| macOS | x86_64 | `minixia-darwin-x86_64` |
| macOS | arm64 | `minixia-darwin-arm64` |
| Windows | x86_64 | `minixia-windows-x86_64.exe` |
| Windows | arm64 | `minixia-windows-arm64.exe` |

### 5.3 平台差异处理

| 功能 | Linux / macOS | Windows |
|------|---------------|---------|
| 配置目录 | `~/.minixia/` | `%LOCALAPPDATA%\minixia\` |
| PATH 安装 | `/usr/local/bin` 或 `~/.local/bin` | `%LOCALAPPDATA%\minixia\` |
| TTS 依赖 | 无需额外依赖 | 无需额外依赖 |
| MQTT 客户端 | 内嵌 | 内嵌 |
| 管道输入 | 完整支持 | PowerShell 部分支持 |

---

## 6. install.sh 详细设计

### 6.1 脚本流程

```
开始
  ├─ 检测是否有 curl 或 wget
  ├─ 检测 OS (uname -s)
  │   ├─ Darwin  → macos
  │   ├─ Linux   → linux
  │   └─ MINGW/MSYS/CYGWIN → windows (Git Bash/WSL)
  ├─ 检测架构 (uname -m)
  │   ├─ x86_64/amd64 → x86_64
  │   └─ arm64/aarch64 → arm64
  ├─ 构造下载 URL
  ├─ 下载二进制 + SHA256 校验文件
  ├─ 校验完整性
  │   ├─ 通过 → 继续
  │   └─ 失败 → 报错退出
  ├─ 确定安装路径
  │   ├─ 有 sudo 权限 → /usr/local/bin/
  │   └─ 无 sudo 权限 → ~/.local/bin/ (并提示 PATH)
  ├─ 移动二进制 + chmod +x
  ├─ 验证安装 (minixia --version)
  └─ 打印成功信息与后续步骤
```

### 6.2 关键特性

- **幂等安装**：重复运行不会出错，已安装时提示升级
- **离线友好**：支持 `MINIXIA_BINARY_URL` 环境变量指定本地或内网下载源
- **代理支持**：自动使用 `HTTP_PROXY` / `HTTPS_PROXY` 环境变量
- **无侵入**：不修改 shell 配置文件（.bashrc / .zshrc），仅在安装到 `~/.local/bin` 时提示手动添加 PATH
- **可配参数**：
  - `MINIXIA_INSTALL_DIR`：自定义安装目录
  - `MINIXIA_VERSION`：指定版本号（默认 latest）
  - `MINIXIA_BINARY_URL`：自定义下载源

```bash
# 指定版本安装
curl -fsSL https://cli.minixia.app/install.sh | MINIXIA_VERSION=1.2.0 sh

# 指定安装目录
curl -fsSL https://cli.minixia.app/install.sh | MINIXIA_INSTALL_DIR=/opt/bin sh
```

### 6.3 install.ps1 流程（Windows PowerShell）

```
开始
  ├─ 检测 PowerShell 版本（≥ 5.1）
  ├─ 检测架构 ($env:PROCESSOR_ARCHITECTURE)
  ├─ 构造下载 URL (.exe)
  ├─ 下载到临时目录 (Invoke-WebRequest)
  ├─ 校验 SHA256 (Get-FileHash)
  ├─ 创建安装目录 ($env:LOCALAPPDATA\minixia\)
  ├─ 复制二进制
  ├─ 检查并添加到用户 PATH
  ├─ 验证安装 (minixia --version)
  └─ 打印成功信息
```

---

## 7. 输出格式与体验

### 7.1 输出规范

- **人类友好模式**（默认）：彩色输出、表格排版、emoji 状态指示
- **机器友好模式**（`--silent` / `--output json`）：纯 JSON 输出，零多余字符，方便管道和 jq 处理
- **错误输出**：错误信息输出到 stderr，正常结果输出到 stdout

```bash
# 人类友好
minixia send "测试"
# ✅ 消息已发送 (message_id: 550e8400-...)

# 机器友好
minixia send --silent "测试"
# {"code":200,"message":"success","data":{"message_id":"550e8400-...","status":"ok"}}
```

### 7.2 错误提示

错误提示应包含：原因、建议操作、参考文档链接。

```
❌ 发送失败：API Key 无效 (ERR_INVALID_APIKEY)

   请检查你的 API Key 是否正确：
   → 运行 minixia config show 查看当前配置
   → 运行 minixia init 重新配置
   → 文档：https://docs.minixia.app/cli/troubleshooting
```

### 7.3 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 成功 |
| 1 | 通用错误 |
| 2 | 参数错误 / 用法错误 |
| 3 | 认证失败（apikey 无效） |
| 4 | 权限不足（套餐限制） |
| 5 | 网络错误 / 服务不可用 |
| 6 | 配额超限 |
| 7 | 限流 |

---

## 8. 实用场景示例

### 8.1 CI/CD 集成

```bash
# GitHub Actions 中使用
- name: Notify via Minixia
  run: |
    curl -fsSL https://cli.minixia.app/install.sh | sh
    minixia send --apikey ${{ secrets.MINIXIA_APIKEY }} --role ci --level info \
      "✅ Build #${{ github.run_number }} 部署成功"
```

### 8.2 服务器监控 cron

```bash
# crontab: 每 5 分钟检查磁盘
*/5 * * * * disk_usage=$(df -h / | tail -1 | awk '{print $5}'); \
  if [ "${disk_usage%\%}" -gt 90 ]; then \
    minixia send -t text -l error "磁盘使用率 $disk_usage，请及时清理"; \
  fi
```

### 8.3 作为指令执行代理

```bash
# 持续监听并执行 App 下发的指令
minixia query --watch --interval 5 --auto-ack \
  --exec 'bash -c "$CONTENT" 2>&1 | minixia send -r bot -'
```

### 8.4 日志实时推送

```bash
# 将日志尾部实时推送到手机
tail -f /var/log/app.log | while read line; do
  minixia send "$line"
done
```

---

## 9. 技术选型

| 维度 | 方案 |
|------|------|
| **语言** | Go（编译为静态二进制，无运行时依赖，跨平台交叉编译成熟） |
| **CLI 框架** | [cobra](https://github.com/spf13/cobra) + [viper](https://github.com/spf13/viper) |
| **HTTP 客户端** | Go 标准库 `net/http` |
| **MQTT 客户端** | [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang)（内嵌，无需用户安装） |
| **TTS** | 不在 CLI 内合成。`send --type voice` / `voice_call`：**文本内容** URL 编码后提交 **server**；**音频文件**读取后 base64 编码再提交 **server**。Azure 合成由 **mini-claw-server** 在文本路径上调用 |
| **配置格式** | TOML（viper 原生支持） |
| **终端 UI** | [lipgloss](https://github.com/charmbracelet/lipgloss) + [tablewriter](https://github.com/olekukonez/tablewriter) |
| **构建发布** | GoReleaser + GitHub Actions（自动交叉编译 6 个目标平台） |

---

## 10. 发布物清单

| 文件 | 说明 |
|------|------|
| `install.sh` | Unix 一键安装脚本 |
| `install.ps1` | Windows PowerShell 一键安装脚本 |
| `minixia-{os}-{arch}` × 6 | 预编译二进制 |
| `checksums.sha256` | 所有二进制的 SHA256 校验和 |
| `CHANGELOG.md` | 版本更新日志 |

---

## 11. 里程碑

| 阶段 | 范围 | 目标 |
|------|------|------|
| **M1：基础可用** | `install.sh` + `init` + `send`（text）+ `status` + `config` | 可一键安装，可发文本消息 |
| **M2：指令闭环** | `query` + `ack` + `--watch` + `--exec` | 可拉取和处理指令 |
| **M3：富媒体** | `send`（image / voice / voice_call）；纯文本语音由服务端 Azure TTS | 支持图片和语音发送 |
| **M4：高级连接** | `webhook` + `mqtt` | 支持全部三种指令接收方式 |
| **M5：生态完善** | `upgrade` + `uninstall` + brew/apt/scoop + shell 补全 | 完善安装体验与生态 |

---

## 12. 非功能性需求

| 需求 | 标准 |
|------|------|
| **启动速度** | 冷启动 < 100ms |
| **二进制大小** | 单平台 < 15MB |
| **零依赖** | 用户机器上无需预装 Node.js / Python / Java 等运行时 |
| **安全** | apikey 配置文件权限 0600；不在 shell history 提示明文 key |
| **可观测** | `--verbose` / `--debug` 输出详细请求日志便于排查 |
| **国际化** | CLI 界面默认中文，通过 `LANG` 环境变量自动切换中英文 |


