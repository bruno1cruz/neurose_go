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

func init() {
	iris.API("/orders", OrdersAPI{})
	iris.API("/orders/:orderId", OrderAPI{})
	iris.API("/orders/:orderId/items", OrderItemsAPI{})
	iris.API("/orders/:orderId/items/:orderItemId", OrderItemAPI{})
	iris.API("/orders/:orderId/transactions", TransactionsAPI{})
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
		api.SetHeader("location", fmt.Sprintf("/orders/%s", order.Id))
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
		api.SetStatusCode(iris.StatusInternalServerError)
		return
	}

	order := Order{Id: api.Param("orderId")}

	err = order.Get(order.Id)
	if err != nil {
		iris.Logger.Dangerf("%s", err)
		api.SetStatusCode(iris.StatusInternalServerError)
		return
	}

	err = order.AddItem(item)
	if err != nil {
		iris.Logger.Dangerf("%s", err)
		api.SetStatusCode(iris.StatusInternalServerError)
		return
	}

	item.Order = &order

	err = item.Save()
	if err != nil {
		api.SetStatusCode(iris.StatusInternalServerError)
	} else {
		api.SetStatusCode(iris.StatusCreated)
		api.SetHeader("location", fmt.Sprintf("/orders/%s/items/%s", order.Id, item.Id))
	}

}

func (api TransactionsAPI) Post() {

	transaction := Transaction{}

	err := api.ReadJSON(&transaction)
	if err != nil {
		iris.Logger.Dangerf("%s", err)
		api.SetStatusCode(iris.StatusInternalServerError)
		return
	}

	order := Order{Id: api.Param("orderId")}

	err = order.Get(order.Id)
	if err != nil {
		iris.Logger.Dangerf("%s", err)
		api.SetStatusCode(iris.StatusInternalServerError)
		return
	}

	err = order.AddTransaction(transaction)
	if err != nil {
		iris.Logger.Dangerf("%s", err)
		api.SetStatusCode(iris.StatusInternalServerError)
		return
	}

	transaction.Order = &order

	err = transaction.Save()
	if err != nil {
		api.SetStatusCode(iris.StatusInternalServerError)
	} else {
		api.SetStatusCode(iris.StatusCreated)
	}

}

func (api OrderItemAPI) Delete() {

	item := OrderItem{}
	item.Get(api.Param("orderItemId"))

	err := item.Delete()
	if err != nil {
		iris.Logger.Dangerf("%s", err)
		api.SetStatusCode(iris.StatusInternalServerError)
	}

}
