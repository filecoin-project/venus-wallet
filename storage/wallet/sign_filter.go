package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/filecoin-project/venus-wallet/config"

	"github.com/filecoin-project/venus/venus-shared/types"
)

type ISignMsgFilter interface {
	CheckSignMsg(ctx context.Context, signMsg SignMsg) error
}

type SignMsg struct {
	SignType types.MsgType
	Data     interface{}
}

type SignFilter struct {
	cfg *config.SignFilter
}

func NewSignFilter(cfg *config.SignFilter) *SignFilter {
	return &SignFilter{cfg}
}

func (filter *SignFilter) CheckSignMsg(ctx context.Context, signMsg SignMsg) error {
	if len(filter.cfg.Expr) == 0 {
		return nil
	}

	j, err := json.MarshalIndent(signMsg, "", "  ")
	if err != nil {
		return err
	}

	var out bytes.Buffer

	c := exec.Command("sh", "-c", filter.cfg.Expr)
	c.Stdin = bytes.NewReader(j)
	c.Stdout = &out
	c.Stderr = &out

	switch err := c.Run().(type) {
	case nil:
		return nil
	case *exec.ExitError:
		return fmt.Errorf("run customer sign check fail (%s)", out.String())
	default:
		return fmt.Errorf("filter cmd run error %w", err)
	}
}
