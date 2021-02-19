package utils

import "github.com/urfave/cli/v2"

var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.json",
		Usage:   "load configuration from `file`",
	}
	SymbolFlag = &cli.StringFlag{
		Name:    "symbol",
		Aliases: []string{"s"},
		Usage:   "trade order symbol, eg. btc_usdt",
	}
	ClientIdFlag = &cli.StringFlag{
		Name:    "clientId",
		Aliases: []string{"i"},
		Usage:   "order's clientId",
	}
	SideFlag = &cli.StringFlag{
		Name:    "side",
		Aliases: []string{"d"},
		Usage:   "trade side, buy/sell",
	}
	TypeFlag = &cli.StringFlag{
		Name:    "type",
		Aliases: []string{"t"},
		Usage:   "trade order type, limit/market",
	}
	PriceFlag = &cli.StringFlag{
		Name:    "price",
		Aliases: []string{"p"},
		Usage:   "order's price, only used by limit order",
	}
	AmountFlag = &cli.StringFlag{
		Name:    "amount",
		Aliases: []string{"m"},
		Usage:   "order's amount, buy or sell amount",
	}
	TotalFlag = &cli.StringFlag{
		Name:    "total",
		Aliases: []string{"a"},
		Usage:   "order's total, only used by market order",
	}
	OrderFlag = &cli.StringFlag{
		Name:    "order",
		Aliases: []string{"o"},
		Usage:   "order id will be cancelled",
	}
)
