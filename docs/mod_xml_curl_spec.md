# mod_xml_curl 接口规范

## 目录（/directory）
- 方法：GET
- 参数：
  - domain: string（默认 default）
  - key/section 兼容：按需扩展
- 返回：`document[type=freeswitch/xml]>section[name=directory]>result{status=success}>data>directory>domain`
  - `domain[name]`
  - `params>param[name,value]`
  - `groups>group>users>user[id]>params>param`

## 拨号方案（/dialplan）
- 方法：GET
- 参数：
  - context: string（默认 public）
  - destination_number: string（可用于动态路由）
- 返回：`document>section[name=dialplan]>result>data>dialplan>context`
  - `context[name]>extension[name]>condition[field,expression]>action[application,data]`

## 鉴权与安全
- 支持 HMAC-SHA256 签名（Header：X-Signature）或 Basic Auth
- IP ACL：仅允许媒体网段访问

## 错误处理
- 返回 `result status=not found` 或 HTTP 404

## 性能
- 可选缓存（ETag/短时缓存）；数据库查询加索引；结构化日志