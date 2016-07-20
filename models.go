package main

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

type TransactionType int
type OrderStatus int

const (
	DRAFT OrderStatus = iota
	ENTERED
	CANCELED
	PAID
	APPROVED
	REJECTED
	RE_ENTERED
	CLOSED
)

const (
	PAYMENT TransactionType = iota
	CANCEL
)

func (o *Order) Save() error {
	o.Id = uuid.NewV4().String()

	if err := session.Query("INSERT INTO \"order\" (id,reference) VALUES (?,?)", o.Id, o.Reference).Exec(); err != nil {
		iris.Logger.Dangerf("%s", err)
		return err
	}
	return nil
}

func (o *Order) Get(id string) {
	if err := session.Query("SELECT id, reference FROM \"neurose_order\".\"orders\" WHERE id = ? ", id).Scan(&o.Id, &o.Reference); err != nil {
		iris.Logger.Dangerf("%s", err)
	}
}

func (i *OrderItem) Save() error {
	i.Id = uuid.NewV4().String()

	if err := session.Query("INSERT INTO \"order_item\" (id,price,quantity,order_id) VALUES (?,?,?,?)", i.Id, i.Price, i.Quantity, i.Order.Id).Exec(); err != nil {
		iris.Logger.Dangerf("%s", err)
		return err
	}
	return nil
}
