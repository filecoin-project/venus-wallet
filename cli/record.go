package cli

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-state-types/builtin/v8/market"
	"github.com/filecoin-project/go-state-types/builtin/v8/paych"
	"github.com/filecoin-project/venus-wallet/cli/helper"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/urfave/cli/v2"

	types2 "github.com/filecoin-project/venus/venus-shared/types/wallet"
)

var recordCmd = &cli.Command{
	Name:  "record",
	Usage: "manipulate sign record",
	Subcommands: []*cli.Command{
		recordQuery,
	},
}

var recordQuery = &cli.Command{
	Name:  "query",
	Usage: "query sign record",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "address to query",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "sign type to query",
		},
		&cli.TimestampFlag{
			Name:     "from",
			Aliases:  []string{"after", "f"},
			Usage:    "from time to query",
			Timezone: time.Local,
			Layout:   "2006-1-2-15:04:05",
		},
		&cli.TimestampFlag{
			Name:     "to",
			Aliases:  []string{"before"},
			Timezone: time.Local,
			Usage:    "to time to query",
			Layout:   "2006-1-2-15:04:05",
		},
		&cli.IntFlag{
			Name:  "limit",
			Usage: "limit to query",
		},
		&cli.IntFlag{
			Name:    "offset",
			Aliases: []string{"skip"},
			Usage:   "offset to query",
		},
		&cli.BoolFlag{
			Name:  "error",
			Usage: "query error record",
		},
		&cli.StringFlag{
			Name:  "id",
			Usage: "query record by id",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "verbose output",
			Aliases: []string{"v"},
		},
	},
	Action: func(cctx *cli.Context) error {
		api, closer, err := helper.GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		ctx := helper.ReqContext(cctx)

		QueryParams := types.QuerySignRecordParams{}

		if cctx.IsSet("address") {
			addrStr := cctx.String("address")
			addr, err := address.NewFromString(addrStr)
			if err != nil {
				return fmt.Errorf("parse address %s : %w", addrStr, err)
			}
			QueryParams.Signer = addr
		}

		if cctx.IsSet("type") {
			t := types.MsgType(cctx.String("type"))
			_, ok := wallet.SupportedMsgTypes[t]
			if !ok {

				fmt.Println("supported types:")
				for k := range wallet.SupportedMsgTypes {
					fmt.Println(k)
				}
				return fmt.Errorf("unsupported type %s", t)
			}
			QueryParams.Type = t
		}
		if cctx.IsSet("from") {
			from := cctx.Timestamp("from")
			QueryParams.After = *from
		}
		if cctx.IsSet("to") {
			to := cctx.Timestamp("to")
			QueryParams.Before = *to
		}
		if cctx.IsSet("limit") {
			limit := cctx.Int("limit")
			QueryParams.Limit = limit
		}
		if cctx.IsSet("offset") {
			offset := cctx.Int("offset")
			QueryParams.Skip = offset
		}
		if cctx.IsSet("error") {
			QueryParams.IsError = cctx.Bool("error")
		}
		if cctx.IsSet("id") {
			QueryParams.ID = cctx.String("id")
		}

		records, err := api.QuerySignRecord(ctx, &QueryParams)
		if err != nil {
			return fmt.Errorf("query sign record: %w", err)
		}
		// output in table format
		w := helper.NewTabWriter(cctx.App.Writer)
		if cctx.Bool("verbose") {
			fmt.Fprintln(w, "ID\tSIGNER\tTYPE\tTIME\tDETAIL\tERROR")
			for _, r := range records {
				errStr := "no error"
				if r.Err != nil {
					errStr = r.Err.Error()
				}
				detail, err := getDetail(&r)
				if err != nil {
					return fmt.Errorf("get detail: %w", err)
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", r.ID, r.Signer, r.Type, r.CreateAt, detail, errStr)
			}
		} else {
			fmt.Fprintln(w, "SIGNER\tTYPE\tTIME\tERROR")
			for _, r := range records {
				errStr := "no error"
				if r.Err != nil {
					errStr = r.Err.Error()
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Signer, r.Type, r.CreateAt, errStr)
			}
		}
		w.Flush()

		return nil
	},
}

func getDetail(r *types.SignRecord) (string, error) {
	var ret string
	t, ok := wallet.SupportedMsgTypes[r.Type]
	if !ok {
		return "", fmt.Errorf("unsupported type %s", r.Type)
	}

	wrap := func(err error) error {
		return fmt.Errorf("get detail: %w", err)
	}

	if r.Msg == nil {
		return "", wrap(fmt.Errorf("msg is nil"))
	}

	if r.Type == types.MTVerifyAddress || r.Type == types.MTUnknown {
		// encode into hex string
		hs := hex.EncodeToString(r.Msg)
		return fmt.Sprintf("Hex:%s.", hs), nil
	}

	signObj := reflect.New(t.Type).Interface()
	if err := wallet.CborDecodeInto(r.Msg, signObj); err != nil {
		return "", fmt.Errorf("decode msg:%w", err)
	}
	switch r.Type {
	case types.MTDealProposal:
		deal := signObj.(*market.DealProposal)
		cid, err := deal.Cid()
		if err != nil {
			return "", wrap(err)
		}
		ret = fmt.Sprintf("DealProposal:%s; Client:%s; Provider:%s.", cid.String(), deal.Client.String(), deal.Provider.String())
	case types.MTClientDeal:
		deal := signObj.(*market.ClientDealProposal)
		cid, err := deal.Proposal.Cid()
		if err != nil {
			return "", wrap(err)
		}
		ret = fmt.Sprintf("ClientDeal:%s; Client:%s; Provider:%s.", cid.String(), deal.Proposal.Client.String(), deal.Proposal.Provider.String())
	case types.MTDrawRandomParam:
		param := signObj.(*types2.DrawRandomParams)
		ret = fmt.Sprintf("Pers:%d ; Round:%s; Entropy :%s.", param.Pers, param.Round, param.Entropy)
	case types.MTSignedVoucher:
		voucher := signObj.(*paych.SignedVoucher)
		ret = fmt.Sprintf("Channel:%s; Amount:%s; Lane:%d .", voucher.ChannelAddr.String(), voucher.Amount.String(), voucher.Lane)
	case types.MTStorageAsk:
		ask := signObj.(*storagemarket.StorageAsk)
		ret = fmt.Sprintf("Miner:%s; Price:%s; VerifiedPrice:%s.", ask.Miner.String(), ask.Price.String(), ask.VerifiedPrice.String())
	case types.MTAskResponse:
		resp := signObj.(*network.AskResponse)
		ret = fmt.Sprintf("Miner:%s; Price:%s; VerifiedPrice:%s.", resp.Ask.Ask.Miner.String(), resp.Ask.Ask.Price.String(), resp.Ask.Ask.VerifiedPrice.String())
		return ret, nil
	case types.MTNetWorkResponse:
		resp := signObj.(*network.Response)
		if resp.State != storagemarket.StorageDealUnknown {
			resp := signObj.(network.Response)
			ret = fmt.Sprintf("State:%s; ProposalCid:%s.", storagemarket.DealStates[resp.State], resp.Proposal)
		}
	case types.MTBlock:
		block := signObj.(*types.BlockHeader)
		ret = fmt.Sprintf("Height:%d ; Miner:%s.", block.Height, block.Miner.String())
	case types.MTChainMsg:
		msg := signObj.(*types.Message)
		ret = fmt.Sprintf("To:%s; Value:%s; Method:%s.", msg.To.String(), msg.Value.String(), msg.Method)
	case types.MTProviderDealState:
		deal := signObj.(*storagemarket.ProviderDealState)
		ret = fmt.Sprintf("ProposalCid:%s; State:%s.", deal.ProposalCid.String(), storagemarket.DealStates[deal.State])
	default:
		ret = "unknown message type"
	}

	return ret, nil
}
