package main

import (
	// "encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/kataras/iris"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const DEFAULT_CASSANDRA_HOST = "127.0.0.1"

const (
	Version = "   Kraken Release: 0.0.1"

	banner = `          
            /\          
           /  \ 
          /"\/"\ 
         / \  / \
         |(@)(@)|
         )  __  (
        //'))(( \\
       (( ((  )) ))        
        \\ ))(( //` + Version + ` `
)

var session *gocql.Session

var config = getConfiguration()

type Config struct {
	App struct {
		Address    string
		Port       string
		Apiversion string
	}
	Cassandra struct {
		ContactPoints string
		Port          string
		ProtoVersion  int
		KeySpace      string
	}
}

func getConfiguration() Config {
	data, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	t := Config{}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- Configuration:\n%v\n\n", t)

	return t
}

func main() {

	var err error
	iris.Config.DisableBanner = true

	cluster := gocql.NewCluster("172.17.0.2")
	cluster.ProtoVersion = 4
	cluster.Keyspace = "neurose_order"
	// cluster := gocql.NewCluster(config.Cassandra.ContactPoints)
	// cluster.ProtoVersion = config.Cassandra.ProtoVersion
	// cluster.Keyspace = config.Cassandra.KeySpace

	session, err = cluster.CreateSession()

	iris.Logger.PrintBanner(banner, "Iniciado em: "+config.App.Address+":"+config.App.Port)

	if err != nil {
		iris.Logger.Dangerf("nao foi possivel conectar ao cassandra", err)
		return
	}

	defer session.Close()

	iris.Listen(config.App.Address + ":" + config.App.Port)
}
