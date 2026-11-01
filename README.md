# 学习项目：阿根廷战术回旋

![img.png](img.png)

---

这个项目引用最小依赖实现 Agent Loop，以达到学习目的
> 即使用 `net/http` 而不是 `openai/openai-go`

#### 支持的工具：
- `read_file`
- `write_file`
- `run_shell`

## 如何使用

#### 手动编译

```bash
git clone https://github.com/unicornfairy864/unf-harness.git
cd unf-harness

# 手动编译
go build -o uh.exe

# 复制配置文件
cp config.yaml.example config.yaml
```

#### 对话示例

```bash     
# 可使用行末带空格的 " \" 换行

> .\uh.exe

User > Good Noon!!
[System]Success.

● Thinking
The user just said "Good Noon!!" - a greeting. I should respond warmly. There's no task here. I have no tools to use for this. Just respond.

● Assistant
Good Noon to you too! ☀️

How can I help you today? Whether it's reading a file, answering a question, or something else — I'm here for it. 😊

User > 

```

#### 工具调用示例
```shell

```

#### 可选参数

```shell
> .\uh.exe -h
Usage of uh:
  -M int
        最大输出token数 (default 4096)
  -S string
        system_prompt (default "You are a helpful assistant.")
  -T float
        temperature (default 0.2)
  -k string
        api_key (default "sk-7d3fabacf3e2463bac3b9516032e8c9a")
  -m string
        选择使用的模型 (default "deepseek-flash")
  -r string
        reasoning_effort (default "high")
  -s    是否启用流式输出
  -t duration
        最大输出时间 (default 5m0s)
  -thinking
        是否开启思考（关闭请用 -thinking=false） (default true)

```