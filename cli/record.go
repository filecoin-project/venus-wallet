package cli

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus-wallet/cli/helper"
	"github.com/filecoin-project/venus-wallet/storage/wallet"
	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/urfave/cli/v2"
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
			Name:    "from",
			Aliases: []string{"after", "f"},
			Usage:   "from time to query",
			Layout:  "2006/1/2/15:04:05",
		},
		&cli.TimestampFlag{
			Name:    "to",
			Aliases: []string{"before"},
			Usage:   "to time to query",
			Layout:  "2006/1/2/15:04:05",
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
			Name:  "verbose",
			Usage: "verbose output",
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
			fmt.Fprintln(w, "ID\tSIGNER\tTYPE\tTIME\tERROR")
			for _, r := range records {
				errStr := "no error"
				if r.Err != nil {
					errStr = r.Err.Error()
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.ID, r.Signer, r.Type, r.CreateAt, errStr)
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
