// stm: #unit
package common

import (
	"context"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/filecoin-project/venus-wallet/filemgr"
	"github.com/filecoin-project/venus-wallet/storage"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestCommon_AuthVerify(t *testing.T) {
	// stm: @VENUSWALLET_API_COMMON_AUTH_VERIFY_001, @VENUSWALLET_API_COMMON_AUTH_NEW_001
	t.Parallel()
	var c Common

	cng, err := filemgr.RandJWTConfig()
	require.NoError(t, err)
	sec, err := hex.DecodeString(cng.Secret)
	require.NoError(t, err)

	app := fxtest.New(t,
		fx.Provide(func() *jwt.HMACSHA { return jwt.NewHS256(sec) }),
		fx.Provide(func() storage.IRecorder { return nil }),
		fx.Populate(&c),
	)
	defer app.RequireStart().RequireStop()

	type args struct {
		token string
	}

	type testCase struct {
		args    args
		want    []auth.Permission
		wantErr bool
	}

	tests := map[string]*testCase{
		"invalid-token-verify": {args: args{"invalid-token"}, want: nil, wantErr: true},
	}

	ctx := context.Background()

	validTokenCase := &testCase{want: []auth.Permission{"admin", "sign", "write"}, wantErr: false}
	token, err := c.AuthNew(ctx, validTokenCase.want)
	require.NoError(t, err)
	validTokenCase.args.token = string(token)

	tests["valid-token-verify"] = validTokenCase

	for tName, tt := range tests {
		t.Run(tName, func(t *testing.T) {
			got, err := c.AuthVerify(ctx, tt.args.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("AuthVerify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthVerify() got = %v, want %v", got, tt.want)
			}
		})
	}
}
