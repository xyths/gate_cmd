package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xyths/hs"
	"github.com/xyths/hs/convert"
	"github.com/xyths/hs/exchange/gateio"
	"go.uber.org/zap"
	"strings"
)

type Node struct {
	Config Config
	Sugar  *zap.SugaredLogger
	Gate   *gateio.V2
}

func NewNode(config Config) (*Node, error) {
	return &Node{Config: config}, nil
}

func (n *Node) Init(ctx context.Context) error {
	l, err := hs.NewZapLogger(n.Config.Log)
	if err != nil {
		return err
	}
	n.Sugar = l.Sugar()
	n.Sugar.Info("Logger initialized")

	n.Gate = gateio.NewV2(n.Config.Exchange.Key, n.Config.Exchange.Secret, n.Config.Exchange.Host, n.Sugar)
	n.Sugar.Infof("Exchange initialized")
	return nil
}

// release all resources
func (n *Node) Close(ctx context.Context) error {
	n.Sugar.Info("Logs are synchronized")
	_ = n.Sugar.Sync()
	return nil
}

func (n *Node) Balance(ctx context.Context) error {
	balance, err := n.Gate.SpotBalanceDetail()
	if err != nil {
		n.Sugar.Errorf("get spot balance error: %s", err)
		return err
	}
	value := decimal.Zero
	for currency, b := range balance {
		price := decimal.NewFromInt(1)
		if strings.ToLower(currency) != "usdt" {
			symbol := fmt.Sprintf("%s_usdt", currency)
			if p, err1 := n.Gate.LastPrice(symbol); err1 == nil {
				price = p
			}
		}
		aa := b.Available // available amount
		av := price.Mul(aa)
		la := b.Locked
		lv := price.Mul(la)
		sa := aa.Add(la)
		sv := av.Add(lv)
		value = value.Add(sv)
		str := fmt.Sprintf(`%s		Amount	Value(USDT)
	Available	%s	%s
	Locked   	%s	%s
	Sum      	%s	%s`, currency, aa, av, la, lv, sa, sv)
		fmt.Println(str)
		n.Sugar.Info(str)
	}
	str := fmt.Sprintf("---------------------------------\nAll: %s", value)
	fmt.Println(str)
	n.Sugar.Info(str)
	return nil
}

func (n *Node) ListOrders(ctx context.Context) error {
	orders, err := n.Gate.OpenOrders()
	if err != nil {
		n.Sugar.Errorf("check open order error: %s", err)
		return err
	}
	for _, o := range orders {
		b, err1 := json.MarshalIndent(o, "", "  ")
		if err1 != nil {
			n.Sugar.Errorf("marshal order error: %s", err1)
			continue
		}
		fmt.Println(string(b))
		n.Sugar.Info(string(b))
	}
	return nil
}

func (n *Node) PlaceOrder(ctx context.Context, symbol, clientId, side, orderType, price, amount, total string) (uint64, error) {
	n.Sugar.Infof("place order, symbol %s, clientId %s, side %s, type %s, price %s, amount %s, total %s", symbol, clientId, side, orderType, price, amount, total)

	switch orderType {
	case "limit":
		price2 := decimal.RequireFromString(price)
		amount2 := decimal.RequireFromString(amount)
		switch side {
		case "buy":
			order, err := n.Gate.BuyLimit(symbol, clientId, price2, amount2)
			if err == nil {
				n.Sugar.Infof("order id is %d", order)
			} else {
				n.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		case "sell":
			order, err := n.Gate.SellLimit(symbol, clientId, price2, amount2)
			if err == nil {
				n.Sugar.Infof("order id is %d", order)
			} else {
				n.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		default:
			n.Sugar.Errorf("bad order side: %s", side)
			return 0, nil
		}
	case "market":
		symbol2, err := n.Gate.GetSymbol(ctx, symbol)
		if err != nil {
			return 0, err
		}
		switch side {
		case "buy":
			total2 := decimal.RequireFromString(total)
			order, err := n.Gate.BuyMarket(symbol2, clientId, total2)
			if err == nil {
				n.Sugar.Infof("order id is %d", order)
			} else {
				n.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		case "sell":
			amount2 := decimal.RequireFromString(amount)
			order, err := n.Gate.SellMarket(symbol2, clientId, amount2)
			if err == nil {
				n.Sugar.Infof("order id is %d", order)
			} else {
				n.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		default:
			n.Sugar.Errorf("bad order side: %s", side)
			return 0, nil
		}
	default:
		n.Sugar.Errorf("bad order type: %s", orderType)
		return 0, nil
	}
}

func (n *Node) CancelOrder(ctx context.Context, symbol, order string) error {
	n.Sugar.Infof("will cancel order %s", order)
	err := n.Gate.CancelOrder(symbol, convert.StrToUint64(order))
	if err != nil {
		n.Sugar.Errorf("cancel order error: %s", err)
	} else {
		n.Sugar.Infof("order %s is cancelled", order)
		fmt.Printf("order %s is cancelled\n", order)
	}
	return nil
}

func (n *Node) TxHistory(ctx context.Context, symbol, order string) error {
	n.Sugar.Infof("get %s trade history, orderId %s", symbol, order)
	trades, err := n.Gate.MyTradeHistory(symbol)
	if err != nil {
		return err
	}
	for _, t := range trades {
		b, err1 := json.MarshalIndent(t, "", "  ")
		if err1 != nil {
			n.Sugar.Errorf("marshal trade error: %s", err1)
			continue
		}
		n.Sugar.Info(string(b))
		fmt.Println(string(b))
	}
	return nil
}
