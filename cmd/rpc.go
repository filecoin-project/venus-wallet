package cmd

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/ipfs-force-community/venus-wallet/api"
	"github.com/ipfs-force-community/venus-wallet/api/permission"
	"github.com/ipfs-force-community/venus-wallet/api/remotecli/httpparse"
	"github.com/ipfs-force-community/venus-wallet/build"
	"github.com/ipfs-force-community/venus-wallet/core"
	"golang.org/x/xerrors"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	logging "github.com/ipfs/go-log/v2"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

var log = logging.Logger("main")

// httpparse cors setting
func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
		w.Header().Set("Content-Type", "application/json")
		if strings.ToLower(r.Method) == "options" {
			_, _ = fmt.Fprintf(w, "")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Start the interface service and bind the address
func ServeRPC(a api.IFullAPI, stop build.StopFunc, addr multiaddr.Multiaddr) error {
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register("Filecoin", api.PermissionedFullAPI(a))
	ah := &Handler{
		Verify: a.AuthVerify,
		Next:   rpcServer.ServeHTTP,
	}
	http.Handle("/rpc/v0", CorsMiddleWare(ah))
	lst, err := manet.Listen(addr)
	if err != nil {
		return xerrors.Errorf("could not listen: %w", err)
	}
	srv := &http.Server{Handler: http.DefaultServeMux}
	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		if err := srv.Shutdown(context.TODO()); err != nil {
			log.Errorf("shutting down RPC server failed: %s", err)
		}
		if err := stop(context.TODO()); err != nil {
			log.Errorf("graceful shutting down failed: %s", err)
		}
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	log.Infof("start rpc server at [%s] ...", addr)
	return srv.Serve(manet.NetListener(lst))
}

// JWT verify
type Handler struct {
	Verify func(ctx context.Context, token string) ([]permission.Permission, error)
	Next   http.HandlerFunc
}

// JWT verify
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.RemoteAddr[:len("127.0.0.1")] == "127.0.0.1" {
		ctx = permission.WithIPPerm(ctx)
	}
	token := r.Header.Get(httpparse.ServiceToken)
	if token == "" {
		token = r.FormValue("token")
		if token != "" {
			token = "Bearer " + token
		}
	}
	if token != "" {
		if !strings.HasPrefix(token, "Bearer ") {
			log.Warn("missing Bearer prefix in auth header")
			w.WriteHeader(401)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		allow, err := h.Verify(ctx, token)
		if err != nil {
			log.Warnf("JWT Verification failed: %s", err)
			w.WriteHeader(401)
			return
		}

		ctx = permission.WithPerm(ctx, allow)
	}
	strategyToken := r.Header.Get(httpparse.WalletStrategyToken)
	ctx = context.WithValue(ctx, core.CtxKeyStrategy, strategyToken)
	h.Next(w, r.WithContext(ctx))
}
