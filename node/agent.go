package node

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xyths/hs"
	"github.com/xyths/hs/convert"
	"github.com/xyths/hs/exchange"
	"github.com/xyths/hs/exchange/gateio"
	"go.uber.org/zap"
	"strings"
)

// Agent is for v4 API, Node is for v2 API.
type Agent struct {
	Config Config
	Sugar  *zap.SugaredLogger
	Gate   *gateio.SpotV4
}

func NewAgent(config Config) (*Agent, error) {
	return &Agent{Config: config}, nil
}

func (a *Agent) Init(ctx context.Context) error {
	l, err := hs.NewZapLogger(a.Config.Log)
	if err != nil {
		return err
	}
	a.Sugar = l.Sugar()
	a.Sugar.Info("Logger initialized")

	a.Gate = gateio.NewSpotV4(a.Config.Exchange.Key, a.Config.Exchange.Secret, a.Config.Exchange.Host, a.Sugar)
	a.Sugar.Infof("Exchange initialized")
	return nil
}

// release all resources
func (a *Agent) Close(ctx context.Context) error {
	a.Sugar.Info("Logs are synchronized")
	_ = a.Sugar.Sync()
	return nil
}

func (a *Agent) Balance(ctx context.Context) error {
	balances, err := a.Gate.Balance(ctx)
	if err != nil {
		a.Sugar.Errorf("get spot balances error: %s", err)
		return err
	}
	value := decimal.Zero
	for _, b := range balances {
		currency := b.Currency
		price := decimal.NewFromInt(1)
		if strings.ToLower(currency) != "usdt" {
			symbol := fmt.Sprintf("%s_usdt", currency)
			if p, err1 := a.Gate.LastPrice(ctx, symbol); err1 == nil {
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
		a.Sugar.Info(str)
	}
	str := fmt.Sprintf("---------------------------------\nAll: %s", value)
	fmt.Println(str)
	a.Sugar.Info(str)
	return nil
}

func (a *Agent) ListOrders(ctx context.Context, symbol string) error {
	orders, err := a.Gate.ListOpenOrders(ctx, symbol)
	if err != nil {
		a.Sugar.Errorf("check %s open order error: %s", symbol, err)
		return err
	}
	for _, o := range orders {
		b, err1 := json.MarshalIndent(o, "", "  ")
		if err1 != nil {
			a.Sugar.Errorf("marshal order error: %s", err1)
			continue
		}
		fmt.Println(string(b))
		a.Sugar.Info(string(b))
	}
	return nil
}

func (a *Agent) PlaceOrder(ctx context.Context, symbol, clientId, side, orderType, price, amount, total string) (exchange.Order, error) {
	a.Sugar.Infof("place order, symbol %s, clientId %s, side %s, type %s, price %s, amount %s, total %s", symbol, clientId, side, orderType, price, amount, total)

	switch orderType {
	case "limit":
		price2 := decimal.RequireFromString(price)
		amount2 := decimal.RequireFromString(amount)
		switch side {
		case "buy":
			order, err := a.Gate.BuyLimit(ctx, symbol, clientId, price2, amount2)
			if err == nil {
				a.Sugar.Infof("order id is %d", order)
			} else {
				a.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		case "sell":
			order, err := a.Gate.SellLimit(ctx, symbol, clientId, price2, amount2)
			if err == nil {
				a.Sugar.Infof("order id is %d", order)
			} else {
				a.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		default:
			a.Sugar.Errorf("bad order side: %s", side)
			return exchange.Order{}, nil
		}
	case "market":
		symbol2, err := a.Gate.GetSymbol(ctx, symbol)
		if err != nil {
			return exchange.Order{}, err
		}
		switch side {
		case "buy":
			total2 := decimal.RequireFromString(total)
			order, err := a.Gate.BuyMarket(ctx, symbol2, clientId, total2)
			if err == nil {
				a.Sugar.Infof("order id is %d", order)
			} else {
				a.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		case "sell":
			amount2 := decimal.RequireFromString(amount)
			order, err := a.Gate.SellMarket(ctx, symbol2, clientId, amount2)
			if err == nil {
				a.Sugar.Infof("order id is %d", order)
			} else {
				a.Sugar.Errorf("place order error: %s", err)
			}
			return order, err
		default:
			a.Sugar.Errorf("bad order side: %s", side)
			return exchange.Order{}, nil
		}
	default:
		a.Sugar.Errorf("bad order type: %s", orderType)
		return exchange.Order{}, nil
	}
}

func (a *Agent) CancelOrder(ctx context.Context, symbol, orderId string) (exchange.Order, error) {
	a.Sugar.Infof("will cancel order %s", orderId)
	order, err := a.Gate.CancelOrder(ctx, symbol, convert.StrToUint64(orderId))
	if err != nil {
		a.Sugar.Errorf("cancel order error: %s", err)
	} else {
		a.Sugar.Infof("order %s is cancelled", orderId)
		fmt.Printf("order %s is cancelled\n", orderId)
	}
	return order, nil
}

func (a *Agent) TxHistory(ctx context.Context, symbol, order string) error {
	a.Sugar.Infof("get %s trade history, orderId %s", symbol, order)
	trades, err := a.Gate.MyTrades(ctx, symbol, order)
	if err != nil {
		return err
	}
	for _, t := range trades {
		b, err1 := json.MarshalIndent(t, "", "  ")
		if err1 != nil {
			a.Sugar.Errorf("marshal trade error: %s", err1)
			continue
		}
		a.Sugar.Info(string(b))
		fmt.Println(string(b))
	}
	return nil
}
