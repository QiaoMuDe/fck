# curl 响应高亮设计方案

## 一、目标

借鉴 cat 命令的 chroma 高亮实现，对 curl 的响应进行智能高亮：
1. 响应头：键值对格式高亮
2. 响应体：根据 Content-Type 自动选择语言高亮

## 二、响应头高亮方案

### 格式分析
```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 256
Set-Cookie: name=value; Path=/
```

### 高亮规则
| 部分 | 颜色 | 说明 |
|------|------|------|
| 状态行 HTTP 版本 | 灰色 | `HTTP/1.1` |
| 状态码 200 | 绿色 | 2xx 成功 |
| 状态码 301/302 | 黄色 | 3xx 重定向 |
| 状态码 400/500 | 红色 | 4xx/5xx 错误 |
| 状态描述 | 默认 | `OK` |
| Header Key | 蓝色 | `Content-Type` |
| 冒号 | 灰色 | `:` |
| Header Value | 默认 | `application/json` |

## 三、响应体高亮方案

### Content-Type 映射表

| Content-Type | 高亮语言 | 说明 |
|--------------|---------|------|
| `application/json` | json | JSON 数据 |
| `application/xml` | xml | XML 数据 |
| `text/html` | html | HTML 页面 |
| `text/css` | css | CSS 样式 |
| `text/javascript` | javascript | JavaScript |
| `application/javascript` | javascript | JavaScript |
| `text/plain` | 无 | 纯文本，不高亮 |
| `application/octet-stream` | 无 | 二进制，不高亮 |

### 实现逻辑
1. 从响应头获取 `Content-Type`
2. 提取 MIME 类型（去掉 charset 等后缀）
3. 查表获取对应语言
4. 使用 chroma 进行高亮

## 四、代码结构

```
internal/commands/curl/
├── cmd_curl.go      # 主入口
├── types.go         # 类型定义
├── formatter.go     # 输出格式化（改造）
└── highlight.go     # 高亮处理（新增）
```

## 五、高亮模块设计

### 1. 响应头高亮函数
```go
// highlightHeaders 高亮响应头
func highlightHeaders(headers http.Header, status string, color bool) string
```

### 2. 响应体高亮函数
```go
// highlightBody 根据 Content-Type 高亮响应体
func highlightBody(body []byte, contentType string, color bool) string
```

### 3. Content-Type 解析函数
```go
// detectLanguage 根据 Content-Type 检测语言
func detectLanguage(contentType string) string
```

## 六、集成方案

### formatter.go 改造
1. 添加 `highlight` 字段控制是否高亮
2. 改造 `PrintHeaders` 支持高亮
3. 改造 `PrintBody` 支持根据 Content-Type 高亮
4. 添加 `PrintHighlighted` 统一高亮输出

### 使用方式
```bash
# 自动根据 Content-Type 高亮
fck curl --color https://api.example.com/users

# 强制高亮（即使 --color 未指定，检测到终端也自动启用）
fck curl https://api.example.com/users
```

## 七、边缘案例

1. **无 Content-Type**：默认使用 text/plain，不高亮
2. **未知 Content-Type**：尝试从扩展名检测，失败则不高亮
3. **二进制内容**：检测为二进制时不高亮
4. **大响应体**：超过阈值时截断或不高亮
5. **JSON 解析失败**：即使 Content-Type 是 json，解析失败也原样输出

## 八、实现优先级

1. **P0** - 响应头高亮（键值对格式）
2. **P0** - JSON 响应体高亮
3. **P1** - HTML/XML 响应体高亮
4. **P2** - 其他语言高亮
5. **P2** - 自动检测终端支持颜色

## 九、测试用例

1. JSON API 响应高亮
2. HTML 页面响应高亮
3. XML 响应高亮
4. 纯文本响应（不高亮）
5. 响应头键值对高亮
6. 不同状态码颜色（200/301/404/500）
