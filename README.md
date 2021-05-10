### Venus-wallet
[![Go Report Card](https://goreportcard.com/badge/github.com/filecoin-project/venus-wallet)](https://goreportcard.com/report/github.com/filecoin-project/venus-wallet)
![Go](https://github.com/filecoin-project/venus-wallet/workflows/Go/badge.svg)

- a private key management tool
- an independent and brief business structure
- provide services such as private key management 
and data signing for local and remote calls via RPC or CLI.

---
### Get Started
#### 1. Build
```
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"

make 
```
- If the test or target application crashes with an "illegal instruction" exception [after copying to an older system], rebuild with `CGO_CFLAGS` environment variable set to `-O -D__BLST_PORTABLE__`. Don't forget `-O`!

#### 2. Setup 
```
# start daemon
$ ./venus-wallet run

# set password to protect wallet security (Used for AES encryption, private key, root token seed)
$ ./venus-wallet set-password
Password:******
Enter Password again:******
```

#### 3. Get remote connect string
> JWT Token restricts access to RPC interface calls
```
# --perm 
# "read","write","sign","admin" 
./venus-wallet auth api-info --perm admin
```

#### 4. Get strategy token
- Strategy token restricts the authority of business execution
- How to generate strategy token for remote service [Venus wallet cli](https://venus.filecoin.io/Venus%20wallet.html#basic-operation-of-venus-wallet)
- URL append strategy token `<JWT token>:/ip4/0.0.0.0/tcp/5678/http:<Strategy token>`



> Once we have a connection string, we can connect to the remote wallet through it.

---

### [How to access remote wallet](./example)
---
### Config
```
[API]
ListenAddress = "/ip4/0.0.0.0/tcp/5678/http"

[DB]
Conn = "[homePath]/keystore.sqlit"
Type = "sqlite"
DebugMode = true

[Factor]
# aes variable
ScryptN = 262144
ScryptP = 1

[JWT]
#  hex JWT token, generate by secret
Token = "" 
# hex JWT secret, randam generate first init
Secret = ""

[Strategy]
# strategy level
# 0：nouse  1：only check struct  2：check struct and msg.method
Level = 2
# need config when Level = 2 and get the actor for msg.to
NodeURL = "/ip4/127.0.0.1/tcp/2345/http"

```
```
---
### Package concept
```
+-- api // RPC service interface permission setting
|
+-- build // dependency injection
|
+-- cli  // shell cmd
|
+-- cmd  // service startup entry
|
+-- config // config provider
|
+-- core // constant 
|
+-- crypto // private key 
|
+-- filemgr // local file manager, Ps:config,database
|
+-- log // log set
|
+-- middeleware // middleware such as link tracking, data reporting
|
+-- signature // signature verification
|
+-- sotrage // the wallet keystore implementation
|
+-- version // git version by ldflags

```


