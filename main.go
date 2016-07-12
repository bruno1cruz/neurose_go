package main

import (
	// "fmt"
	// "log"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	// "github.com/valyala/fasthttp"
)

// var cluster *gocql.ClusterConfig
var session *gocql.Session

type Order struct {
	Id        string `json:"-"`
	Reference string `json:"reference"`
}

type OrderItemAPI struct {
	*iris.Context
}
type OrdersAPI struct {
	*iris.Context
}
type OrderAPI struct {
	*iris.Context
}

func main() {

	cluster := gocql.NewCluster("172.17.0.2")
	cluster.ProtoVersion = 4
	cluster.Keyspace = "neurose_order"

	session, _ = cluster.CreateSession()
	defer session.Close()

	iris.API("/orders", OrdersAPI{})
	iris.API("/orders/:orderId", OrderAPI{})
	iris.API("/orders/:orderId/items", OrderItemAPI{})

	iris.Listen(":3000")
}

func (api OrdersAPI) Post() {
	order := Order{}

	err := json.Unmarshal(api.PostBody(), &order)

	if err != nil {
		log.Error(err)
		api.SetStatusCode(iris.StatusBadRequest)
		return
	}

	order.Save()

	// api.SetStatusCode(iris.StatusCreated)
	// api.SetHeader("location", fmt.Sprintf("/orders/%s", order.Id))
}

func (api OrderAPI) Get() {
	order := Order{}
	order.Get(api.Param("orderId"))
	api.JSON(iris.StatusOK, order)
}

func (o *Order) Save() {
	// session, _ := cluster.CreateSession()
	// defer session.Close()

	o.Id = uuid.NewV4().String()

	if err := session.Query("INSERT INTO \"order\" (id,reference) VALUES (?,?)", o.Id, o.Reference).Exec(); err != nil {
		log.Error(err)
	}
}

func (o *Order) Get(id string) {
	log.Debugf("get order %s", id)

	// session, _ := cluster.CreateSession()
	// defer session.Close()

	if err := session.Query("SELECT id, reference FROM \"order\" WHERE id = ? ", id).Scan(&o.Id, &o.Reference); err != nil {
		log.Error(err)
	}
}
