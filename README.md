# 🛡️ PromptLeakFence (PLF)

[![Go Version](https://img.shields.io/github/go-mod/go-version/lukas-sgx/PromptLeakFence)](https://go.dev/)
[![License](https://img.shields.io/github/license/lukas-sgx/PromptLeakFence)](LICENSE)

**The first bi-directional LLM prompt firewall.**

PromptLeakFence is a lightweight, transparent security proxy designed to prevent prompt injection attacks and unauthorized sensitive information disclosure (PII, credentials, etc.) in AI systems. It sits between your application and your LLM provider, scanning every message for potential leaks.

## ✨ Key Features

- 🔍 **Real-time Interception**: Acts as a transparent proxy for LLM APIs.
- 🔒 **Bi-directional Filtering**: Scans both incoming prompts and outgoing assistant responses.
- 🛡️ **Policy-based Redaction**: Automatically redacts sensitive patterns (passwords, API keys, tokens) using customizable rules.
- 🚀 **Multi-Provider Support**: Pre-configured for Ollama, llama.cpp, LMStudio, Gemini, Claude, and more.
- 📊 **Audit Dashboard**: Specialized dashboard to visualize and analyze blocked attempts.
- 🛠️ **Lightweight & Fast**: Built with Go for high performance and low latency.

## 📦 Installation

### Prerequisites
- Go 1.25+
- Make

### Build
```bash
make build
```
This will generate the `plf` binary in the `bin/` directory.

## 🚀 Usage
### Starting the Security Proxy
The proxy requires a target LLM service to forward clean traffic to.

```bash
# Example: Running PLF in front of Ollama
sudo ./bin/plf proxy --target ollama --listen 8080 --verbose
```

**Supported Targets:**
- `ollama` (default port: 11434)
- `llama.cpp` (8080)
- `lmstudio` (1234)
- `oobabooga` (7860)
- `openwebui` (3000)
- `copilot` (5000)
- `gemini` (8080)
- `claude` (8080)

### Audit Dashboard
Visualize blocked leak attempts:

```bash
./bin/plf audit --port 9090
```
Then visit `http://127.0.0.1:9090`.

## ⚙️ Configuration

Security rules are defined in `configs/policy.yaml`. You can customize the `exclude` list to add patterns that should be redacted.

```yaml
policy:
  exclude:
    - "token"
    - "password"
    - "api_key"
    # ... add your custom sensitive keywords here
```

Every match found in a prompt or response will be replaced by `[INTERNAL_PROMPT_REDACTED]`.

## 📂 Project Structure

- `main.go`: Entry point for the CLI.
- `cmd/`: Command implementations (`proxy`, `audit`, `root`).
- `cmd/utils/`: Core logic for network redirection, policy parsing and launch helpers.
- `configs/`: Default security policies.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License
See LICENSE file for details.