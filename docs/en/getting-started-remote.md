# Venus wallet

1. venus-wallet is a remote wallet that provides policy for Filecoin and supports JsonRPC2.0 call. It can dynamically configure whether various data types to be signed are signed or not.
2. The project is decoupled from Lotus and Venus independently, and can be called by different implementations of Filecoin.

## quickstart

### 1-Downloadcode

```
git clone https://github.com/filecoin-project/venus-wallet.git
```

### 2-Compile

- go version ^1.15

```shell script
# Setting BLS compilation environment variables
export CGO_CFLAGS_ALLOW="-D__BLST_PORTABLE__"
export CGO_CFLAGS="-D__BLST_PORTABLE__"

# Compile the current platform executable
make

# If you need to cross compile Linux versions on MAC
# You need to install GCC related files (you can also download files to local via GitHub and install them locally by brew)
brew install FiloSottile/musl-cross/musl-cross
make linux
```

### 3-Startserviceprocess

```shell script
# It starts on the Mainnetwork by default(--network=main)
# The address begins with f
$ ./venus-wallet run

# Start in test network
# The address begins with t
$ ./venus-wallet run  --network=test
```

### 4-Configurationintroduction

- Default file location `~/.venus_wallet/config.toml`

```toml
[API]
  # The HTTP listening address of the local process
  ListenAddress = "/ip4/0.0.0.0/tcp/5678/http"

[DB]
  # Data files that embedded store the database  by default
  Conn = "~/.venus_wallet/keystore.sqlit"
  Type = "sqlite"
  DebugMode = true

[JWT]
  # JWT token hex，If it is not configured, it will be generated randomly
  Token = "65794a68624763694f694a49557a49314e694973496e523563434936496b705856434a392e65794a42624778766479493657794a795a57466b4969776964334a70644755694c434a7a615764754969776959575274615734695858302e7133787a356f75634f6f543378774d5463743870574d42727668695f67697a4f7a365142674b2d6e4f7763"
  # JWT secret hex，If it is not configured, it will be generated randomly
  Secret = "7c40ce66a492e35ac828e8333a5703e38b23add87f29bd8fc7343989e08b3458"

[Factor]
  # keystore private key symmetric encryption variable
  ScryptN = 262144
  ScryptP = 1

[Strategy]
  # Strategy level，0：Don't turn on strategy 1：Verify only the data type to be signed 2：Verify the data type to be signed, and verify the message type with the method policy configured
  Level = 2
  NodeURL = "/ip4/127.0.0.1/tcp/2345/http"

[APIRegisterHub]
  # The URL of the gateway. If not configured, the gateway will not be connected
  RegisterAPI = ["/ip4/127.0.0.1/tcp/45132"]
  # The token of the gateway
  Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdG1pbmVyIiwicGVybSI6ImFkbWluIiwiZXh0IjoiIn0.oakIfSg1Iiv1T2F1BtH1bsb_1GeXWuirdPSjvE5wQLs"
  SupportAccounts = ["testminer"]
```

# CLIoperationguide

## ViewHelp

```shell script
$ ./venus-wallet -h


NAME:
   venus-wallet - A new cli application

USAGE:
   venus-wallet [global options] command [command options] [arguments...]

VERSION:
   1.0.0'+gitc04f451.dirty'

COMMANDS:
   run                   Start a venus wallet process
   auth                  Manage RPC permissions
   log                   Manage logging
   strategy, st          Manage logging
   new                   Generate a new key of the given type
   list, ls              List wallet address
   export                export keys
   import                import keys
   sign                  sign a message
   del                   del a wallet and message
   set-password, setpwd  Store a credential for a keystore file
   unlock                unlock the wallet and release private key
   lock                  Restrict the use of secret keys after locking wallet
   lockState, lockstate  unlock the wallet and release private key
   help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

- The operation here is flat and single-layer. Different from the `./venus wallet list` operation of Venus or Lotus, only `./venus-wallet list` is needed in venus-wallet.
- Some commands are XX processed, such as `strategy`, which can be directly replaced by `st`.

## BasicoperationofVenusWallet

### Thestateofthewallet

1. Set the key of private key symmetric encryption

```shell script
# ./venus-wallet setpwd (aliase)
$ ./venus-wallet set-password
Password:******
Enter Password again:******

# res
Password set successfully
```

> Note: this password is only stored in memory for symmetric encryption of the private key. Once the service process exits in any form, it cannot be restored. Therefore, the private key managed by this program needs to be backed up by itself or directly.

- After setting the password, the default state of the wallet is unlock

2. Lock Wallet
   > After the wallet is locked, the functions of signing, generating new address, importing and exporting private key will be disabled, which will affect the remote call chain, so please use it with caution.

```shell script
$ ./venus-wallet lock
Password:******

# res
wallet lock successfully
```

3. unlock wallet
   > After unlocking, all functions of the wallet will be released.

```shell script
$ ./venus-wallet unlock
Password:******

# res
wallet unlock successfully
```

4. View the wallet status

```shell script
$ ./venus-wallet lockstate

#res
wallet state: unlocked
```

### Privatekeymanagement

1. Generate new random private key
   > venus-wallet new [command options] [bls|secp256k1 (default secp256k1)]

```shell script
$ ./venus-wallet new

#res
t12mchblwgi243re5i2pg2harmnqvm6q3rwb2cnpy
```

- The default type is secp256k1. You can also use `./venus-wallet new bls` to generate BLS private key

2. Import the private key
   > venus-wallet import [command options] [\<path\> (optional, will read from stdin if omitted)]

```shell script
$ ./venus-wallet import
Enter private key:7b2254797065223a22736563703235366b31222c22507269766174654b6579223a22626e765665386d53587171346173384633654c647a7438794a6d68764e434c377132795a6c6657784341303d227d

#res
imported key t12mchblwgi243re5i2pg2harmnqvm6q3rwb2cnpy successfully!
```

3. Export the private key
   > venus-wallet export [command options] [address]

```shell script
$ ./venus-wallet export t12mchblwgi243re5i2pg2harmnqvm6q3rwb2cnpy

# res
7b2254797065223a22736563703235366b31222c22507269766174654b6579223a22626e765665386d53587171346173384633654c647a7438794a6d68764e434c377132795a6c6657784341303d227d
```

4. View address list

```shell script
$ ./venus-wallet list

t3uktqgxtagiyk5cxrjn5h4wq4v247saxtfukfi6zsvt4sek2q2ufkg27biasg7247zhdpm2kpotukwsapr7pa
t3rcgmzisnusxvwrwvi7l5hcuissvmluvkrzfuehjdfawba75qlv3mxl6rtnxitt33z5fuwds76rbcyafhxrua
t12mchblwgi243re5i2pg2harmnqvm6q3rwb2cnpy
```

> Show all private key corresponding address, there are spec and bls two kinds of address 5. Delete the specified private key
> venus-wallet del [command options] \<address\>

```shell script
$ ./venus-wallet del t12mchblwgi243re5i2pg2harmnqvm6q3rwb2cnpy

#res
success
```

### JWTauthoritymanagement

For remote access interface authorization

1. Gets the remote connection string
   > venus-wallet auth api-info [command options] [arguments...]

```shell script
$ ./venus-wallet auth api-info --perm admin

#res
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.q3xz5oucOoT3xwMTct8pWMBrvhi_gizOz6QBgK-nOwc:/ip4/0.0.0.0/tcp/5678/http
```

- perm has four kinds of permissions: read, write, sign and Admin. They are generated by the corresponding `JWT` configuration in the configuration file and will not change dynamically.

### Config in venus

format: `token:muitiaddr`.

```json
{
        "walletModule": {
                "defaultAddress": "f3ueri27yppflsxodo66r2u4jajw5d4lhrzlcv4ncx7efrrxyivnrsufi7wuvdjmpbepwb2npvj7wglla6gtcq",
                "passphraseConfig": {
                        "scryptN": 2097152,
                        "scryptP": 1
                },
                "remoteEnable": true,
                "remoteBackend": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIl19.gCLPHlI5r9lyxfbPoeU8nSGQI9CpUBaBGA54EzgZ9vE_e78f9e6c-9033-4144-8992-a1890ad76ead:/ip4/192.168.5.64/tcp/5678/http"
        }
}
```

### Config in lotus

format: `token:muitiaddr` .the same reason as before

```toml
[Wallet]
   RemoteBackend = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIl19.gCLPHlI5r9lyxfbPoeU8nSGQI9CpUBaBGA54EzgZ9vE_e78f9e6c-9033-4144-8992-a1890ad76ead:/ip4/192.168.5.64/tcp/5678/http"
  #EnableLedger = false
  #DisableLocal = false
```
