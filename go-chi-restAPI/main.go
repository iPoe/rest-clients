package main

import (
	"api/user/reg/wParser"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

type Person struct {
	Uid   string   `json:"uid,omitempty"`
	Cid   string   `json:"Cid,omitempty"`
	Name  string   `json:"name,omitempty"`
	Age   int      `json:"age,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

var m = make(map[string]wParser.Product)

//End-point that allows user send a date and then update the db data by that date
func FetchDatabyDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	date := r.URL.Query().Get("date")
	d := ""
	if len(date) != 0 {
		t, err := time.Parse("01/02/2006", date)
		if err != nil {
			log.Fatal(err)
		}
		d = strconv.FormatInt(t.Unix(), 10)
	}

	UpdateDb(d)
}

//Simple End-point that updates the db to the current moment
func UpDb(w http.ResponseWriter, r *http.Request) {
	UpdateDb("")
}

//Given a  date it performs and update of data in the current db
func UpdateDb(d string) {
	fmt.Println(d)
	if len(d) == 0 {
		d = ""
	}
	DropData()
	inserClients(d)
	listaProductos := wParser.ReadProducts(d)
	AddProducts(listaProductos)
	transsactList := wParser.ReadTransactions(d)
	AddTransactions(transsactList)
}

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

var dg = newClient()
var ctx = context.Background()

func AddProducts(productlist wParser.Products) {
	// var m = make(map[string]interface{})

	for i := range productlist {
		m[productlist[i].Id] = productlist[i]
	}
	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(productlist)
	mu.SetJson = pb
	response, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		fmt.Println(response.Json)
		log.Fatal(err)

	}
	//fmt.Println(err)

}

func inserClients(time string) {
	// fmt.Println("====Connecting to dgrap====")
	// fmt.Println("====Parsing clients data====")
	Clist := wParser.ReadClients("")
	fmt.Println(len(Clist))
	op := &api.Operation{} //Falta añadir los esquemas de producto y transaccion
	op.Schema = `	
		name: string @index(exact) . 
		age: int .
		price: int .
		Cid: string @index(exact)  .
		Pid: string @index(exact)  .
		Tid: string @index(exact)  .
		Ip: string .
		Device: string .
		ProductIds: [string].

		
		type Person {
			name: string
			age: int
			Cid: string
		}
		type Product{
			name: string
			price: int
			Pid: string
		}
		type Transaction{
			Tid: string
			Ip: string
			Cid: string
			Device: string
			ProductIds: [string]
		}
		`
	var plist []Person
	for i := range Clist {
		p := Person{
			Cid:   Clist[i].Id,
			Name:  Clist[i].Name,
			Age:   Clist[i].Age,
			DType: []string{"Person"},
		}
		plist = append(plist, p)
	}

	if err := dg.Alter(ctx, op); err != nil {
		log.Fatal(err)
	}
	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(plist)
	mu.SetJson = pb
	response, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		fmt.Println(response.Json)
		log.Fatal(err)

	}

}

//Given a list of transactions it populates the db with it
func AddTransactions(TransacList wParser.Transactions) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	fmt.Println(len(TransacList))
	pb, err := json.Marshal(TransacList)
	if err != nil {
		log.Fatal(err)
	}
	mu.SetJson = pb
	response, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println(response.Json)

}

//Erases all existing data in the moment
func DropData() {
	op := &api.Operation{
		DropAll: true,
	}
	if err := dg.Alter(ctx, op); err != nil {
		log.Fatal(err)
	}

}

func main() {
	port := "9000"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	log.Printf("Starting up on http://localhost:%s", port)
	r := chi.NewRouter()
	r.Get("/", UpDb)                                //Fetch data from today
	r.Get("/data/", FetchDatabyDate)                //Fetch data from en exact date
	r.Mount("/clients", clientsResource{}.Routes()) //Use a client resource to manage its requests
	// getSimilarBuyers()

	log.Fatal(http.ListenAndServe(":"+port, r))

}
