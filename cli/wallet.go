package cli

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/cli/helper"
	"github.com/filecoin-project/venus-wallet/core"
	"github.com/filecoin-project/venus-wallet/errcode"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var walletSetPassword = &cli.Command{
	Name:    "set-password",
	Aliases: []string{"setpwd"},
	Usage:   "Store a credential for a keystore file",
	Action: func(cctx *cli.Context) error {
		pw, err := gopass.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		pw2, err := gopass.GetPasswdPrompt("Enter Password again:", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		if !bytes.Equal(pw, pw2) {
			return errors.New("the input passwords are inconsistent")
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		err = api.SetPassword(ctx, string(pw2))
		if err != nil {
			return err
		}
		fmt.Println("Password set successfully")
		return nil
	},
}

var walletUnlock = &cli.Command{
	Name:  "unlock",
	Usage: "unlock the wallet and release private key",
	Action: func(cctx *cli.Context) error {
		pw, err := gopass.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		err = api.Unlock(ctx, string(pw))
		if err != nil {
			return err
		}
		fmt.Println("wallet unlock successfully")
		return nil
	},
}

var walletLock = &cli.Command{
	Name:  "lock",
	Usage: "Restrict the use of secret keys after locking wallet",
	Action: func(cctx *cli.Context) error {
		pw, err := gopass.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		err = api.Lock(ctx, string(pw))
		if err != nil {
			return err
		}
		fmt.Println("wallet lock successfully")
		return nil
	},
}

var walletLockState = &cli.Command{
	Name:  "lock-state",
	Usage: "unlock the wallet and release private key",
	Action: func(cctx *cli.Context) error {
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		locked := api.LockState(ctx)
		state := "unlocked"
		if locked {
			state = "locked"
		}
		fmt.Printf("wallet state: %s\n", state)
		return nil
	},
}

var walletNew = &cli.Command{
	Name:      "new",
	Usage:     "Generate a new key of the given type",
	ArgsUsage: "[bls|secp256k1 (default secp256k1)]",
	Action: func(cctx *cli.Context) error {
		t := core.KeyType(cctx.Args().First())
		if t == "" {
			t = core.KTSecp256k1
		}
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		ctx := helper.ReqContext(cctx)
		defer closer()
		nk, err := api.WalletNew(ctx, t)
		if err != nil {
			return err
		}
		fmt.Println(nk.String())
		return nil
	},
}

var walletList = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List wallet address",
	Action: func(cctx *cli.Context) error {
		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		addrs, err := api.WalletList(ctx)
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			fmt.Println(addr.String())
		}
		return nil
	},
}

var walletExport = &cli.Command{
	Name:      "export",
	Usage:     "export keys",
	ArgsUsage: "[address]",
	Action: func(cctx *cli.Context) error {
		if !cctx.Args().Present() {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		addr, err := address.NewFromString(cctx.Args().First())
		if err != nil {
			return err
		}

		pw, err := gopass.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}

		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		if err := api.VerifyPassword(ctx, string(pw)); err != nil {
			return err
		}
		ki, err := api.WalletExport(ctx, addr)
		if err != nil {
			return err
		}
		b, err := json.Marshal(ki)
		if err != nil {
			return err
		}

		fmt.Println(hex.EncodeToString(b))
		return nil
	},
}

var walletImport = &cli.Command{
	Name:      "import",
	Usage:     "import keys",
	ArgsUsage: "[<path> (optional, will read from stdin if omitted)]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "format",
			Usage: "specify input format for key",
			Value: "hex-venus",
		},
	},
	Action: func(cctx *cli.Context) error {
		var inpdata []byte
		if !cctx.Args().Present() || cctx.Args().First() == "-" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter private key: ")
			indata, err := reader.ReadBytes('\n')
			if err != nil {
				return err
			}
			inpdata = indata

		} else {
			fdata, err := ioutil.ReadFile(cctx.Args().First())
			if err != nil {
				return err
			}
			inpdata = fdata
		}

		var ki core.KeyInfo
		switch cctx.String("format") {
		case "hex-venus":
			data, err := hex.DecodeString(strings.TrimSpace(string(inpdata)))
			if err != nil {
				return err
			}

			if err := json.Unmarshal(data, &ki); err != nil {
				return err
			}
		case "json-venus":
			if err := json.Unmarshal(inpdata, &ki); err != nil {
				return err
			}
		case "gfc-json":
			var f struct {
				KeyInfo []struct {
					PrivateKey []byte
					SigType    int
				}
			}
			if err := json.Unmarshal(inpdata, &f); err != nil {
				return xerrors.Errorf("failed to parse go-filecoin key: %s", err)
			}

			gk := f.KeyInfo[0]
			ki.PrivateKey = gk.PrivateKey
			switch gk.SigType {
			case 1:
				ki.Type = core.KTSecp256k1
			case 2:
				ki.Type = core.KTBLS
			default:
				return fmt.Errorf("unrecognized key type: %d", gk.SigType)
			}
		default:
			return fmt.Errorf("unrecognized format: %s", cctx.String("format"))
		}

		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		addr, err := api.WalletImport(ctx, &ki)
		if err != nil {
			return err
		}

		fmt.Printf("imported key %s successfully!\n", addr)
		return nil
	},
}

var walletSign = &cli.Command{
	Name:      "sign",
	Usage:     "sign a message",
	ArgsUsage: "<signing address> <hexMessage>",
	Action: func(cctx *cli.Context) error {
		if !cctx.Args().Present() || cctx.NArg() != 2 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		addr, err := address.NewFromString(cctx.Args().First())

		if err != nil {
			return err
		}

		msg, err := hex.DecodeString(cctx.Args().Get(1))

		if err != nil {
			return err
		}

		pw, err := gopass.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}

		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		if err := api.VerifyPassword(ctx, string(pw)); err != nil {
			return err
		}
		sig, err := api.WalletSign(ctx, addr, msg, core.MsgMeta{})
		if err != nil {
			return err
		}
		sigBytes := append([]byte{byte(sig.Type)}, sig.Data...)
		fmt.Println(hex.EncodeToString(sigBytes))
		return nil
	},
}

var walletDel = &cli.Command{
	Name:      "del",
	Usage:     "del a wallet and message",
	ArgsUsage: "<address>",
	Action: func(cctx *cli.Context) error {
		if !cctx.Args().Present() || cctx.NArg() != 1 {
			return helper.ShowHelp(cctx, errcode.ErrParameterMismatch)
		}
		addr, err := address.NewFromString(cctx.Args().First())
		if err != nil {
			return err
		}

		pw, err := gopass.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
		if err != nil {
			return err
		}

		api, closer, err := helper.GetFullAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()
		ctx := helper.ReqContext(cctx)
		if err := api.VerifyPassword(ctx, string(pw)); err != nil {
			return err
		}
		if err = api.WalletDelete(ctx, addr); err != nil {
			return err
		}

		fmt.Println("success")
		return nil
	},
}
