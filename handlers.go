package main

import (
	"fmt"
	"github.com/kataras/iris"
)

type OrdersAPI struct {
	*iris.Context
}
type OrderAPI struct {
	*iris.Context
}
type OrderItemsAPI struct {
	*iris.Context
}
type OrderItemAPI struct {
	*iris.Context
}
type TransactionsAPI struct {
	*iris.Context
}
type TransactionAPI struct {
	*iris.Context
}

func init() {
	iris.API("/orders", OrdersAPI{})
	iris.API("/orders/:orderId", OrderAPI{})
	iris.API("/orders/:orderId/items", OrderItemsAPI{})
	iris.API("/orders/:orderId/items/:orderItemId", OrderItemAPI{})
	iris.API("/orders/:orderId/transactions", TransactionsAPI{})
	iris.API("/orders/:orderId/transactions/:transactionId", TransactionAPI{})
}

func (api OrdersAPI) Post() {
	order := Order{}

	err := api.ReadJSON(&order)

	if err != nil {
		iris.Logger.Dangerf("%s", err)
		//api.SetStatusCode(iris.StatusBadRequest)
		api.JSON(iris.StatusBadRequest, err)
		return
	}

	err = order.Save()
	if err != nil {
		api.SetStatusCode(iris.StatusInternalServerError)
	} else {
		api.SetStatusCode(iris.StatusCreated)
		api.SetHeader("location", fmt.Sprintf("/orders/%s/%s", config.App.Apiversion, order.Id))
	}

}

func (api OrderAPI) Get() {
	order := Order{}
	order.Get(api.Param("orderId"))
	api.JSON(iris.StatusOK, order)
}

func (api OrderItemsAPI) Post() {
	item := OrderItem{}

	err := api.ReadJSON(&item)

	if err != nil {
		iris.Logger.Dangerf("%s", err)
		api.SetStatusCode(iris.StatusBadRequest)
		return
	}

	order := Order{Id: api.Param("orderId")}
	item.Order = &order

	err = item.Save()
	if err != nil {
		api.SetStatusCode(iris.StatusInternalServerError)
	} else {
		api.SetStatusCode(iris.StatusCreated)
		api.SetHeader("location", fmt.Sprintf("/orders/%s/items/%s", order.Id, item.Id))
	}

}
