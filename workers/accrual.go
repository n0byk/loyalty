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

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

var iteration = 0

func AccrualAskWorker() {

	c := time.Tick(time.Second)
	for range c {
		go AccrualAskWorkerRunner()
	}
}

func AccrualAskWorkerRunner() {

	order, _, err := config.App.Storage.GetNewOrder(context.Background())
	if err != nil {
		config.App.Logger.Error("AccrualAskWorkerRunner", zap.Error(err))
	}

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

			var result accrualResponse
			if err := json.Unmarshal(body, &result); err != nil {
				config.App.Logger.Error("AccrualAskWorkerRunner", zap.Error(err))
			}
			fmt.Printf("%+v \n", result)

			if result.Status == "PROCESSED" {
				config.App.Storage.SetOrderStatus(context.Background(), order.OrderID, "PROCESSED")
				config.App.Storage.PostAccrue(context.Background(), result.Order, float32(result.Accrual))

				config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("message", order.OrderID))
			}
			if result.Status == "INVALID" {
				config.App.Storage.SetOrderStatus(context.Background(), order.OrderID, "INVALID")
				config.App.Storage.PostAccrue(context.Background(), result.Order, float32(result.Accrual))
				config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("message", order.OrderID))
			}
			config.App.Logger.Info("AccrualAskWorkerRunner", zap.String("iteration", strconv.Itoa(iteration)))
			iteration++
		}
	}
}
