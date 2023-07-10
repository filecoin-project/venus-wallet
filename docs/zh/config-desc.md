# 配置文件解析

- 默认文件位置 “~/.venus_wallet/config.toml” 钱包配置文件务必备份好
```toml
[API]
  # 本地进程http监听地址
  ListenAddress = "/ip4/0.0.0.0/tcp/5678/http"

[DB]
  # 默认内嵌存储数据库数据文件
  Conn = "~/.venus_wallet/keystore.sqlit"
  Type = "sqlite"
  DebugMode = true

# 用于远程主动请求wallet进行签名时，进行合法性验证
[JWT]
  # JWT token hex，未配置情况下会随机生成
  Token = "65794a68624763694f694a49557a49314e694973496e523563434936496b705856434a392e65794a42624778766479493657794a795a57466b4969776964334a70644755694c434a7a615764754969776959575274615734695858302e7133787a356f75634f6f543378774d5463743870574d42727668695f67697a4f7a365142674b2d6e4f7763"
  # JWT secret hex，未配置情况下会随机生成
  Secret = "7c40ce66a492e35ac828e8333a5703e38b23add87f29bd8fc7343989e08b3458"

[Factor]
  # keystore私钥对称加密变量, 强烈建议不要随意修改这两个参数，后期考虑会从配置中删除
  ScryptN = 262144
  ScryptP = 1

[SignFilter]
  Expr = "" # 不填写任何东西，表示不开启过滤功能；具体使用方式请参考 [订单过滤器]文档

[APIRegisterHub]
  # gateway的URL，不配置则不连接gateway
  RegisterAPI = ["/ip4/127.0.0.1/tcp/45132"]
  # 用于访问gateway的token；其实是auth服务产生的token
  Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdG1pbmVyIiwicGVybSI6ImFkbWluIiwiZXh0IjoiIn0.oakIfSg1Iiv1T2F1BtH1bsb_1GeXWuirdPSjvE5wQLs"
  SupportAccounts = ["authTestUser"] # 建议不再使用此字段，后面会删除
```