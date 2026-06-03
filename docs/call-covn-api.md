# `/api/file/*` 调用与签名说明

本文档面向外部调用方，说明如何调用文件接口以及如何生成请求签名。

## 接口范围

以下接口需要携带文件签名请求头：

| Method | Path | 说明 |
| --- | --- | --- |
| GET | `/api/file/list` | 获取上传目录下的全部文件 |
| GET | `/api/file/listByDir` | 获取指定目录下的文件 |
| POST | `/api/file/upload` | 上传文件 |

其他接口，例如 `/api/meta`、`/api/app/*`、`/public/*`，不使用这套文件签名认证。

## 请求头

每次请求都需要携带以下请求头：

| Header | 必填 | 说明 |
| --- | --- | --- |
| `X-Timestamp` | 是 | Unix 秒级时间戳，例如 `1762012800` |
| `X-Sign` | 是 | 小写 MD5 hex 签名 |

服务端默认允许时间戳与当前时间相差不超过 300 秒。调用方应保证机器时间准确，建议开启 NTP。

## 签名规则

签名原文固定为 5 行文本：

```text
METHOD
PATH
CANONICAL_QUERY
TIMESTAMP
SECRET
```

也就是：

```text
METHOD + "\n" + PATH + "\n" + CANONICAL_QUERY + "\n" + TIMESTAMP + "\n" + SECRET
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `METHOD` | HTTP 方法，大写，例如 `GET`、`POST` |
| `PATH` | URL path，不包含 scheme、host 和 query，例如 `/api/file/listByDir` |
| `CANONICAL_QUERY` | 标准化后的 query string；没有 query 时为空字符串 |
| `TIMESTAMP` | 与 `X-Timestamp` 完全一致的 Unix 秒级时间戳字符串 |
| `SECRET` | 双方约定的共享密钥 |

对签名原文计算 MD5，并转为小写十六进制字符串，作为 `X-Sign`。

请求体不参与签名。上传文件时只签 method、path、query、timestamp 和 secret。

## Query 标准化

`CANONICAL_QUERY` 的规则与 Go 标准库 `url.Values.Encode()` 一致：

- 按参数名升序排序。
- 参数值使用 URL 编码。
- 多个参数使用 `&` 连接。
- 没有 query 参数时使用空字符串。
- `X-Timestamp` 和 `X-Sign` 是请求头，不放入 query，也不参与 `CANONICAL_QUERY`。

示例：

```text
原始 query: z=last&dirName=avatars&a=first
标准 query: a=first&dirName=avatars&z=last
```

## 签名示例

假设：

```text
METHOD: GET
PATH: /api/file/listByDir
CANONICAL_QUERY: dirName=avatars
TIMESTAMP: 1762012800
SECRET: test-secret
```

签名原文为：

```text
GET
/api/file/listByDir
dirName=avatars
1762012800
test-secret
```

计算：

```bash
printf 'GET\n/api/file/listByDir\ndirName=avatars\n1762012800\ntest-secret' | md5
```

得到的小写 MD5 hex 即为 `X-Sign`。

## cURL 示例

以下示例需要先安装可用的 `md5` 命令。macOS 使用 `md5`，Linux 通常使用 `md5sum`，二者输出格式不同，请只取 hash 字符串。

### 查询全部文件

```bash
BASE_URL='http://localhost:6301'
SECRET='test-secret'
METHOD='GET'
PATH='/api/file/list'
QUERY=''
TIMESTAMP="$(date +%s)"
SIGN="$(printf '%s\n%s\n%s\n%s\n%s' "$METHOD" "$PATH" "$QUERY" "$TIMESTAMP" "$SECRET" | md5 | awk '{print $NF}')"

curl -sS \
  -H "X-Timestamp: $TIMESTAMP" \
  -H "X-Sign: $SIGN" \
  "$BASE_URL$PATH"
```

### 查询指定目录

```bash
BASE_URL='http://localhost:6301'
SECRET='test-secret'
METHOD='GET'
PATH='/api/file/listByDir'
QUERY='dirName=avatars'
TIMESTAMP="$(date +%s)"
SIGN="$(printf '%s\n%s\n%s\n%s\n%s' "$METHOD" "$PATH" "$QUERY" "$TIMESTAMP" "$SECRET" | md5 | awk '{print $NF}')"

curl -sS \
  -H "X-Timestamp: $TIMESTAMP" \
  -H "X-Sign: $SIGN" \
  "$BASE_URL$PATH?$QUERY"
```

### 上传文件

`dirName` 是 multipart 表单字段，不是 query 参数，因此不参与签名。

```bash
BASE_URL='http://localhost:6301'
SECRET='test-secret'
METHOD='POST'
PATH='/api/file/upload'
QUERY=''
TIMESTAMP="$(date +%s)"
SIGN="$(printf '%s\n%s\n%s\n%s\n%s' "$METHOD" "$PATH" "$QUERY" "$TIMESTAMP" "$SECRET" | md5 | awk '{print $NF}')"

curl -sS \
  -X POST \
  -H "X-Timestamp: $TIMESTAMP" \
  -H "X-Sign: $SIGN" \
  -F "file=@./example.txt" \
  -F "dirName=avatars" \
  "$BASE_URL$PATH"
```

## JavaScript 示例

Node.js 调用方可以按下面方式生成签名：

```js
import crypto from "node:crypto";

function buildFileSign({ method, path, query = "", timestamp, secret }) {
  const raw = [
    method.toUpperCase(),
    path,
    query,
    String(timestamp),
    secret,
  ].join("\n");

  return crypto.createHash("md5").update(raw).digest("hex");
}

const timestamp = Math.floor(Date.now() / 1000);
const sign = buildFileSign({
  method: "GET",
  path: "/api/file/listByDir",
  query: new URLSearchParams({ dirName: "avatars" }).toString(),
  timestamp,
  secret: "test-secret",
});

const response = await fetch("http://localhost:6301/api/file/listByDir?dirName=avatars", {
  headers: {
    "X-Timestamp": String(timestamp),
    "X-Sign": sign,
  },
});
```

## 响应格式

接口统一返回 HTTP 200，业务状态通过 `code` 判断：

```json
{
  "code": 1000,
  "message": "success",
  "data": {}
}
```

常见认证失败：

| code | 说明 | 常见原因 |
| --- | --- | --- |
| `1007` | 无效的 token | 缺少请求头、时间戳格式错误、时间戳过期 |
| `1010` | 权限不足 | 签名错误，或生产环境未配置共享密钥 |

## 配置项

服务端通过以下配置启用文件签名：

```yaml
auth:
  file_sign:
    secret: "your-shared-secret"
    expire_seconds: 300
```

也可以通过环境变量注入：

```bash
export COVN_AUTH_FILE_SIGN_SECRET='your-shared-secret'
export COVN_AUTH_FILE_SIGN_EXPIRE_SECONDS='300'
```

生产环境必须配置非空 `secret`。开发环境如果 `secret` 为空，服务端会跳过文件签名校验。
