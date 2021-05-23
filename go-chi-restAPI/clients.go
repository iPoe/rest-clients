package main

//Emigrar en este paquete los endpoints encargados de responder con la
//info de los clientes

import (
	"bytes"
	"context"
	"encoding/json"

	"fmt"
	"io"
	"log"
	"net/http"

	//"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/go-chi/chi"
)

type Ips struct {
	Ip string
}
type Iplist []Ips

type OwnerName struct {
	Owner []struct {
		Name string `json:"name"`
	} `json:"owner"`
}
type IpArray struct {
	Tran []struct {
		IP string `json:"Ip"`
	} `json:"tran"`
}
type MyData struct {
	Tran []struct {
		Cid string `json:"Cid"`
		IP  string `json:"Ip"`
	} `json:"tran"`
	Tran2 []struct {
		Name string `json:"name"`
	} `json:"tran2"`
}
type PidArray struct {
	Array []struct {
		ProductIds []string `json:"ProductIds"`
		Tid        string   `json:"Tid"`
		Price        int   `json:"price"`


	} `json:"owner"`
}

type clientsResource struct{}

func (rs clientsResource) Routes() chi.Router {

	r := chi.NewRouter()
	r.Get("/", rs.List) // Get a list of all the clients in the system

	r.Route("/{id}", func(r chi.Router) {
		r.Use(PostCtx)
		r.Get("/", rs.Get) // GET /client/{id} Read a client by its ID and return its transactions

	})

	return r
}

func (rs clientsResource) List(w http.ResponseWriter, r *http.Request) {
	q := `{
		datos(func:type(Person)){
			Cid
			name
			age
		}		
	}`
	resp, err := dg.NewTxn().Query(ctx, q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respByte := bytes.NewReader(resp.Json)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, respByte); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs clientsResource) Get(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	fmt.Println(id)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	data := getClientOrders(id)
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Fatal(err)
	}
	m["simBuyers"] = getSimilarBuyers(id)
	// fmt.Println(m["simBuyers"])
	newData, err := json.Marshal(m)
	respByte := bytes.NewReader(newData)

	if _, err := io.Copy(w, respByte); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func AuxOwnerName(id string) string {
	list := make(map[string]string)
	list["$a"] = id
	q := `query tran($a:string){
		owner(func:eq(Cid,$a)){
			name
	   }
	}`
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, list)
	if err != nil {
		log.Fatal(err)
	}
	var r OwnerName
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(r.Owner))
	ans := r.Owner[0].Name
	return ans
}

func RemoveDuplicateValues(StringSlice []string, id string) []string {
	keys := make(map[string]bool)
	list := []string{}
	owner := AuxOwnerName(id)
	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range StringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			if entry != owner {
				list = append(list, entry)
			}
		}
	}

	return list
}

func getSimilarBuyers(id string) []string {
	list := make(map[string]string)
	list["$a"] = id
	q := `query tran($a:string){
		Tran(func:eq(Cid,$a)){
			Ip

	   }	 
	}`
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, list)

	var r IpArray
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Fatal(err)
	}
	ips := r.Tran

	listaIps := make(map[string]string)
	var namesxTransaccion []string
	for i := range ips {
		namesxTransaccion = append(namesxTransaccion, getNamesxTransact(ips[i].IP, listaIps)...)

	}
	namesNotDup := RemoveDuplicateValues(namesxTransaccion, id)
	return namesNotDup

}

//Returns the names of users that used the same Ip address
func getNamesxTransact(ipaddress string, list map[string]string) []string {
	// Given an Ip address it returns a list
	// with all the names that used that ip
	list["$a"] = ipaddress

	q := `query tran($a:string){
		tran(func:type(Transaction))@filter(eq(Ip,$a)){
			CC as Cid
			Ip
	   } 
		tran2(func:eq(Cid,val(CC))){
			name
		}
		 
	
	}`
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, list)
	// fmt.Printf("%s\n", resp.Json)
	if err != nil {
		log.Fatal(err)
	}

	var r2 MyData
	err = json.Unmarshal(resp.Json, &r2)
	if err != nil {
		log.Fatal(err)
	}
	n := r2.Tran2
	// fmt.Println(n)
	var res []string
	for i := range n {
		res = append(res, n[i].Name)
	}
	return res

}

func getClientOrders(id string) []byte {
	list := make(map[string]string)
	list["$a"] = id
	q := `query tran($a:string){
		owner(func:eq(Cid,$a)){
			ProductIds
			Tid

	   }	 
	}`
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, list)
	if err != nil {
		log.Fatal(err)
	}
	var r PidArray
	err = json.Unmarshal(resp.Json, &r)
	if err != nil{
		log.Fatal(err)
	}
	//fmt.Println(r)
	//orders := r.Array
	for j := range r.Array{
		r.Array[j].Price = getOrderPrice(r.Array[j].ProductIds)
	}
	ans,err := json.Marshal(r)
	if err != nil{
		log.Fatal(err)
	}
	return ans

}

func getOrderPrice(plist []string) int{
	ans := 0
	for i := range plist {
		ans = ans + m[plist[i]]
	}
	return ans
}
