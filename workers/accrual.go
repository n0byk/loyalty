package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/n0byk/loyalty/config"
	"go.uber.org/zap"
)

const (
	OrderStatusRegistered = "REGISTERED"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusInvalid    = "INVALID"
)

type Order struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

var iteration = 0

func AccrualAskWorker() {
	c := time.Tick(time.Second)
	for range c {
		for i := 0; i < config.AppConfig.AccrualAskWorkerCount; i++ {
			go AccrualAskWorkerRunner()
		}
	}
}

func AccrualAskWorkerRunner() {
	config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("message", "No response from request"))

	orders, _, err := config.App.Storage.GetNewOrder(context.Background())
	if err != nil {
		config.App.Logger.Error("AccrualAskWorkerRunner", zap.Error(err))
	}

	if len(orders) > 0 {
		for _, order := range orders {
			if order.OrderID != "" && len(order.OrderNumber) > 0 {
				resp, err := http.Get(config.AppConfig.AccrualSystemAddress + "/api/orders/" + order.OrderNumber)
				if err != nil {
					config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("message", "No response from request"))
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						config.App.Logger.Error("AccrualAskWorkerRunner", zap.Error(err))
					}

					var result Order
					if err := json.Unmarshal(body, &result); err != nil {
						config.App.Logger.Error("AccrualAskWorkerRunner", zap.Error(err))
					}
					fmt.Printf("%+v \n", result)

					if result.Status == OrderStatusProcessed {
						config.App.Storage.SetOrderStatus(context.Background(), order.OrderID, OrderStatusProcessed)
						config.App.Storage.PostAccrue(context.Background(), result.Order, float32(result.Accrual))

						config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("message", order.OrderID))
					}
					if result.Status == OrderStatusInvalid {
						config.App.Storage.SetOrderStatus(context.Background(), order.OrderID, OrderStatusProcessed)
						config.App.Storage.PostAccrue(context.Background(), result.Order, float32(result.Accrual))
						config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("message", order.OrderID))
					}
					config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("iteration", strconv.Itoa(iteration)))
					iteration++
				}
			}
		}
	}

}
