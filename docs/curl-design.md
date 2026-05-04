# curl 子命令设计方案

## 一、命令定位

一个简洁的 HTTP 客户端工具，类似 curl，但输出更友好。

## 二、命令格式

```
fck curl [options] <url>
```

## 三、核心参数

| 参数 | 简写 | 说明 |
|------|------|------|
| `--request` | `-X` | HTTP 方法 (GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS)，默认 GET |
| `--data` | `-d` | 请求体数据 |
| `--header` | `-H` | 请求头，多个用逗号分隔 |
| `--output` | `-o` | 输出到文件 |
| `--include` | `-i` | 显示响应头 |
| `--silent` | `-s` | 静默模式，只输出响应体 |
| `--verbose` | `-v` | 显示完整请求/响应详情 |
| `--form` | `-F` | multipart/form-data 格式发送数据 |
| `--user` | `-u` | 用户名密码认证 (user:password) |
| `--location` | `-L` | 跟随重定向 |
| `--max-time` | `-m` | 最大执行时间 |
| `--retry` |  | 失败重试次数 |
| `--pretty` | `-p` | 格式化 JSON 输出 |
| `--color` |  | 启用彩色输出 |

## 四、使用示例

```bash
# 简单 GET
fck curl https://api.example.com/users

# POST JSON
fck curl -X POST https://api.example.com/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","age":25}'

# POST 表单
fck curl -X POST https://api.example.com/login \
  -d 'username=admin&password=123456'

# 文件上传
fck curl -F "file=@/path/to/image.jpg" https://api.example.com/upload

# 下载并保存
fck curl -o result.json https://api.example.com/users

# 显示响应头
fck curl -i https://api.example.com/users

# 完整详情
fck curl -v https://api.example.com/users

# Basic 认证
fck curl -u admin:password https://api.example.com/admin

# 跟随重定向
fck curl -L https://bit.ly/xxx

# 静默模式（只输出响应体）
fck curl -s https://api.example.com/users | jq '.data[0].name'
```

## 五、输出格式

### 普通模式

```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 256

{
  "code": 0,
  "data": [
    {"id": 1, "name": "张三"},
    {"id": 2, "name": "李四"}
  ]
}
```

### Verbose 模式

```
> GET /users HTTP/1.1
> Host: api.example.com
> User-Agent: fck/1.8
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Content-Length: 256
<
{
  "code": 0,
  "data": [...]
}

Time: 45ms
Size: 256 bytes
```

## 六、技术实现

### 文件结构

```
internal/
├── cli/
│   └── curl.go                 # CLI 参数定义
├── commands/
│   └── curl/
│       ├── cmd_curl.go         # 主入口
│       ├── request.go          # HTTP 请求构建
│       ├── response.go         # 响应处理
│       └── formatter.go        # 输出格式化
```

### 依赖

- 标准库 `net/http` - HTTP 客户端
- `gitee.com/MM-Q/color` - 彩色输出

## 七、实现优先级

### P0 - 核心功能
- [ ] 基础 GET/POST/PUT/DELETE 请求
- [ ] Header 和 Body 支持
- [ ] JSON 格式化输出

### P1 - 增强功能
- [ ] 文件上传 (multipart/form-data)
- [ ] 文件下载 (-o)
- [ ] Verbose 模式 (-v)
- [ ] 跟随重定向 (-L)
- [ ] 认证支持 (-u)

### P2 - 高级功能
- [ ] 失败重试
- [ ] 超时控制
- [ ] 代理支持

## 八、边缘案例

1. **URL 解析**：处理相对路径和绝对路径
2. **编码问题**：自动检测响应编码
3. **大文件**：流式处理下载/上传
4. **重定向循环**：限制最大重定向次数
5. **超时控制**：连接超时和读取超时分离
6. **证书验证**：HTTPS 自签名证书处理

## 九、测试用例

1. 基础 GET 请求
2. POST JSON 数据
3. POST 表单数据
4. 自定义 Header
5. 文件上传
6. 文件下载
7. 跟随重定向
8. Verbose 输出格式
9. 静默模式
10. 彩色输出控制
