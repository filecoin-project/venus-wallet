package cli

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/filemgr"
	"github.com/mitchellh/go-homedir"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/prometheus/common/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	metadataTraceContext = "traceContext"
)

// custom CLI error

type ErrCmdFailed struct {
	msg string
}

func (e *ErrCmdFailed) Error() string {
	return e.msg
}

type APIInfo struct {
	Addr  multiaddr.Multiaddr
	Token []byte
}

func (a APIInfo) DialArgs() (string, error) {
	_, addr, err := manet.DialArgs(a.Addr)
	if strings.HasPrefix(addr, "0.0.0.0:") {
		addr = "127.0.0.1:" + addr[8:]
	}
	return "ws://" + addr + "/rpc/v0", err
}

func (a APIInfo) AuthHeader() http.Header {
	if len(a.Token) != 0 {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer "+string(a.Token))
		return headers
	}
	return nil
}

func GetAPIInfo(ctx *cli.Context) (APIInfo, error) {
	p, err := homedir.Expand(ctx.String("repo"))
	if err != nil {
		return APIInfo{}, xerrors.Errorf("cound not expand home dir (repo): %w", err)
	}
	r, err := filemgr.NewFS(p, nil)
	if err != nil {
		return APIInfo{}, xerrors.Errorf("could not open repo at path: %s; %w", p, err)
	}

	ma, err := r.APIEndpoint()
	if err != nil {
		return APIInfo{}, xerrors.Errorf("could not get api endpoint: %w", err)
	}

	token, err := r.APIToken()
	if err != nil {
		log.Warnf("Couldn't load CLI token, capabilities may be limited: %v", err)
	}

	return APIInfo{
		Addr:  ma,
		Token: token,
	}, nil
}

func GetRawAPI(ctx *cli.Context) (string, http.Header, error) {
	ainfo, err := GetAPIInfo(ctx)
	if err != nil {
		return "", nil, xerrors.Errorf("could not get API info: %w", err)
	}

	addr, err := ainfo.DialArgs()
	if err != nil {
		return "", nil, xerrors.Errorf("could not get DialArgs: %w", err)
	}

	return addr, ainfo.AuthHeader(), nil
}

func GetAPI(ctx *cli.Context) (api.ICommon, jsonrpc.ClientCloser, error) {
	addr, headers, err := GetRawAPI(ctx)
	if err != nil {
		return nil, nil, err
	}

	return api.NewCommonRPC(ctx.Context, addr, headers)
}
func GetFullNodeAPI(ctx *cli.Context) (api.IFullAPI, jsonrpc.ClientCloser, error) {
	addr, headers, err := GetRawAPI(ctx)
	if err != nil {
		return nil, nil, err
	}
	return api.NewFullNodeRPC(ctx.Context, addr, headers)
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

var Commands = []*cli.Command{
	authCmd,
	walletCmd,
	logCmd,
}

//nolint
func withCategory(cat string, cmd *cli.Command) *cli.Command {
	cmd.Category = cat
	return cmd
}
