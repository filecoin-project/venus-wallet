package helper

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/venus-wallet/api"
	"github.com/filecoin-project/venus-wallet/api/remotecli"
	"github.com/filecoin-project/venus-wallet/api/remotecli/httpparse"
	"github.com/filecoin-project/venus-wallet/common"
	"github.com/filecoin-project/venus-wallet/filemgr"
	"github.com/howeyc/gopass"
	logging "github.com/ipfs/go-log/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/urfave/cli/v2"
)

var log = logging.Logger("cli")

const (
	metadataTraceContext = "traceContext"
)

type ctxKey string

const (
	ctxPWD ctxKey = "pwd"
)

// custom CLI error

type ErrCmdFailed struct {
	msg string
}

func (e *ErrCmdFailed) Error() string {
	return e.msg
}

func GetAPIInfo(ctx *cli.Context) (httpparse.APIInfo, error) {
	p, err := homedir.Expand(ctx.String("repo"))
	if err != nil {
		return httpparse.APIInfo{}, fmt.Errorf("could not expand home dir (repo): %w", err)
	}
	r, err := filemgr.NewFS(p, nil)
	if err != nil {
		return httpparse.APIInfo{}, fmt.Errorf("could not open repo at path: %s; %w", p, err)
	}

	ma, err := r.APIEndpoint()
	if err != nil {
		return httpparse.APIInfo{}, fmt.Errorf("could not get api endpoint: %w", err)
	}

	token, err := r.APIToken()
	if err != nil {
		log.Warnf("Couldn't load CLI token, capabilities may be limited: %v", err)
	}

	return httpparse.APIInfo{
		Addr:  ma,
		Token: token,
	}, nil
}

func GetRawAPI(ctx *cli.Context) (string, http.Header, error) {
	ainfo, err := GetAPIInfo(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("could not get API info: %w", err)
	}

	if err := dial(ainfo.Addr); err != nil {
		return "", nil, err
	}

	addr, err := ainfo.DialArgs()
	if err != nil {
		return "", nil, fmt.Errorf("could not get DialArgs: %w", err)
	}

	return addr, ainfo.AuthHeader(), nil
}

func dial(addr string) error {
	ma, err := multiaddr.NewMultiaddr(addr)
	if err == nil {
		_, addr, err := manet.DialArgs(ma)
		if err != nil {
			return err
		}
		dialer := net.Dialer{
			Timeout: time.Second * 2,
		}
		_, err = dialer.Dial("tcp", addr)
		return err
	}

	return nil
}

func GetAPI(ctx *cli.Context) (common.ICommon, jsonrpc.ClientCloser, error) {
	addr, headers, err := GetRawAPI(ctx)
	if err != nil {
		return nil, nil, err
	}

	return remotecli.NewCommonRPC(ctx.Context, addr, headers)
}
func GetFullAPI(ctx *cli.Context) (api.IFullAPI, jsonrpc.ClientCloser, error) {
	addr, headers, err := GetRawAPI(ctx)
	if err != nil {
		return nil, nil, err
	}
	return remotecli.NewFullNodeRPC(ctx.Context, addr, headers)
}

func GetFullAPIWithPWD(ctx *cli.Context) (api.IFullAPI, jsonrpc.ClientCloser, error) {
	addr, headers, err := GetRawAPI(ctx)
	if err != nil {
		return nil, nil, err
	}
	err = withPWD(ctx)
	if err != nil {
		return nil, nil, err
	}
	return remotecli.NewFullNodeRPC(ctx.Context, addr, headers)
}

func DaemonContext(cctx *cli.Context) context.Context {
	if mtCtx, ok := cctx.App.Metadata[metadataTraceContext]; ok {
		return mtCtx.(context.Context)
	}
	return context.Background()
}

// ReqContext returns context for cli execution. Calling it for the first time
// installs SIGTERM handler that will close returned context.
// Not safe for concurrent execution.
func ReqContext(cctx *cli.Context) context.Context {
	tCtx := DaemonContext(cctx)

	ctx, done := context.WithCancel(tCtx)
	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		done()
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	return ctx
}

func withPWD(cctx *cli.Context) error {
	pwd, err := gopass.GetPasswdPrompt("Password:", true, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	cctx.Context = context.WithValue(cctx.Context, ctxPWD, pwd)
	return nil
}

// nolint
func withCategory(cat string, cmd *cli.Command) *cli.Command {
	cmd.Category = cat
	return cmd
}

func ShowHelp(cctx *cli.Context, err error) error {
	return &PrintHelpErr{Err: err, Ctx: cctx}
}

type PrintHelpErr struct {
	Err error
	Ctx *cli.Context
}

func (e *PrintHelpErr) Error() string {
	return e.Err.Error()
}

func (e *PrintHelpErr) Unwrap() error {
	return e.Err
}

func (e *PrintHelpErr) Is(o error) bool {
	_, ok := o.(*PrintHelpErr)
	return ok
}
