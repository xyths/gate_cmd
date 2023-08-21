package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xyths/gate_cmd/node"
	"github.com/xyths/gate_cmd/utils"
	"github.com/xyths/hs"
)

var (
	balanceCmd = &cli.Command{
		Action: balanceAction,
		Name:   "balance",
		Usage:  "get balance",
	}
	orderCmd = &cli.Command{
		Action: orderAction,
		Name:   "order",
		Usage:  "place or check order",
		Subcommands: []*cli.Command{
			{
				Action: listAction,
				Name:   "list",
				Usage:  "list open orders",
				Flags: []cli.Flag{
					utils.SymbolFlag,
				},
			},
			{
				Action: placeAction,
				Name:   "place",
				Usage:  "place an order",
				Flags: []cli.Flag{
					utils.SymbolFlag,
					utils.ClientIdFlag,
					utils.SideFlag,
					utils.TypeFlag,
					utils.PriceFlag,
					utils.AmountFlag,
					utils.TotalFlag,
				},
			},
			{
				Action: cancelAction,
				Name:   "cancel",
				Usage:  "cancel an open order",
				Flags: []cli.Flag{
					utils.SymbolFlag,
					utils.OrderFlag,
				},
			},
		},
		Flags: []cli.Flag{
		},
	}
	txCmd = &cli.Command{
		Action: txAction,
		Name:   "tx",
		Usage:  "get transaction history",
		Flags: []cli.Flag{
			utils.SymbolFlag,
			utils.OrderFlag,
		},
	}
)

func balanceAction(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := node.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	n, _ := node.NewAgent(cfg)
	err := n.Init(ctx.Context)
	if err != nil {
		return err
	}
	defer n.Close(ctx.Context)

	if err = n.Balance(ctx.Context); err != nil {
		return err
	}
	return nil
}

func orderAction(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := node.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	n, _ := node.NewAgent(cfg)
	err := n.Init(ctx.Context)
	if err != nil {
		return err
	}
	defer n.Close(ctx.Context)

	return nil
}

func listAction(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := node.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	n, _ := node.NewAgent(cfg)
	err := n.Init(ctx.Context)
	if err != nil {
		return err
	}
	defer n.Close(ctx.Context)

	if err = n.ListOrders(ctx.Context, ctx.String(utils.SymbolFlag.Name)); err != nil {
		return err
	}

	return nil
}

func placeAction(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := node.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	n, _ := node.NewAgent(cfg)
	err := n.Init(ctx.Context)
	if err != nil {
		return err
	}
	defer n.Close(ctx.Context)

	if orderId, err := n.PlaceOrder(
		ctx.Context,
		ctx.String(utils.SymbolFlag.Name),
		ctx.String(utils.ClientIdFlag.Name),
		ctx.String(utils.SideFlag.Name),
		ctx.String(utils.TypeFlag.Name),
		ctx.String(utils.PriceFlag.Name),
		ctx.String(utils.AmountFlag.Name),
		ctx.String(utils.TotalFlag.Name),
	); err != nil {
		return err
	} else {
		fmt.Printf("order id is %d\n", orderId)
	}

	return nil
}

func cancelAction(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := node.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	n, _ := node.NewAgent(cfg)
	err := n.Init(ctx.Context)
	if err != nil {
		return err
	}
	defer n.Close(ctx.Context)

	if _, err = n.CancelOrder(ctx.Context, ctx.String(utils.SymbolFlag.Name), ctx.String(utils.OrderFlag.Name)); err != nil {
		return err
	}

	return nil
}

func txAction(ctx *cli.Context) error {
	configFile := ctx.String(utils.ConfigFlag.Name)
	cfg := node.Config{}
	if err := hs.ParseJsonConfig(configFile, &cfg); err != nil {
		return err
	}
	n, _ := node.NewAgent(cfg)
	err := n.Init(ctx.Context)
	if err != nil {
		return err
	}
	defer n.Close(ctx.Context)

	if err = n.TxHistory(ctx.Context, ctx.String(utils.SymbolFlag.Name), ctx.String(utils.OrderFlag.Name)); err != nil {
		return err
	}

	return nil
}
