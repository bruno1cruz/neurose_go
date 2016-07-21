package main

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

type TransactionType int
type OrderStatus int

const (
	ORDER_DRAFT OrderStatus = iota
	ORDER_ENTERED
	ORDER_CANCELED
	ORDER_PAID
	ORDER_APPROVED
	ORDER_REJECTED
	ORDER_RE_ENTERED
	ORDER_CLOSED
)

const (
	TRANSACTION_PAYMENT TransactionType = iota
	TRANSACTION_CANCEL
)

type Order struct {
	Id        string      `json:"id"`
	Reference string      `json:"reference"`
	Number    string      `json:number`
	Status    OrderStatus `json:status`
	CreatedAt time.Time   `json:createdAt`
	UpdatedAt time.Time   `json:updatedAt`
	Notes     string      `json:updatedAt`
	Price     int         `json:price`
	Payed     int         `json:-`
}

type OrderItem struct {
	Id       string `json:"sku"`
	Price    int    `json:"unit_price"`
	Quantity int    `json:"quantity"`
	Order    *Order `json:"-"`
}

type Transaction struct {
	Id                string          `json:"id"`
	ExternalId        string          `json:"external_id"`
	Amount            int             `json:amount`
	Type              TransactionType `json:type`
	AuthorizationCode string          `json:authorization_code`
	CardBrand         string          `json:card_brand`
	CardBin           string          `json:card_bin`
	CardLast          string          `json:card_last`
	Order             *Order          `json:"-"`
}

func (o *Order) Save() error {
	o.Id = uuid.NewV4().String()

	return session.Query("INSERT INTO \"order\" (id,reference) VALUES (?,?)", o.Id, o.Reference).Exec()
}

func (o *Order) Get(id string) error {
	return session.Query("SELECT id, reference, price, payed FROM \"order\" WHERE id = ? ", id).Scan(&o.Id, &o.Reference, &o.Price, &o.Payed)
}

func (o *Order) AddItem(i OrderItem) error {
	o.Price = o.Price + i.Price*i.Quantity
	return o.updatePrice()
}

func (o *Order) RemoveItem(i OrderItem) error {
	o.Price = o.Price - i.Price*i.Quantity
	return o.updatePrice()
}

func (o *Order) updatePrice() error {
	sql := "UPDATE \"order\" SET price=? WHERE id=? "
	return session.Query(sql, o.Price, o.Id).Exec()
}

func (o *Order) AddTransaction(t Transaction) error {

	if TRANSACTION_PAYMENT == t.Type {
		o.Payed = o.Payed + t.Amount
	} else if TRANSACTION_CANCEL == t.Type {
		o.Payed = o.Payed - t.Amount
	} else {
		errors.New(fmt.Sprintf("transaction type[%d] invalid", t.Type))
	}

	sql := "UPDATE \"order\" SET payed=? WHERE id=? "
	return session.Query(sql, o.Payed, o.Id).Exec()
}

func (i *OrderItem) Get(id string) error {
	o := Order{}
	err := session.Query("SELECT id, price, quantity,order_id FROM \"order_item\" WHERE id = ? ", id).Scan(&i.Id, &i.Price, &i.Quantity, &o.Id)

	if err == nil {
		i.Order = &o
	}
	return err
}

func (i *OrderItem) Save() error {
	i.Id = uuid.NewV4().String()

	return session.Query("INSERT INTO \"order_item\" (id,price,quantity,order_id) VALUES (?,?,?,?)", i.Id, i.Price, i.Quantity, i.Order.Id).Exec()
}

func (i *OrderItem) Delete() error {

	i.Order.Get(i.Order.Id)

	err := i.Order.RemoveItem(*i)
	if err == nil {
		return session.Query("UPDATE \"order_item\" SET deleted = true WHERE id = ?", i.Id).Exec()
	}
	return err

}

func (t *Transaction) Save() error {
	t.Id = uuid.NewV4().String()

	return session.Query("INSERT INTO \"transaction\" (id,amount,type,order_id) VALUES (?,?,?,?)", t.Id, t.Amount, t.Type, t.Order.Id).Exec()
}
