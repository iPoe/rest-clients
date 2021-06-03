package main

import (
	"bytes"
	"context"
	"encoding/json"

	"fmt"
	"io"
	"log"
	"net/http"

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
		Cid        string   `json:"Cid"`
		IP         string   `json:"Ip"`
		ProductIds []string `json:"ProductIds"`
	} `json:"tran"`
	Tran2 []struct {
		Name string `json:"name"`
	} `json:"tran2"`
}
type PidArray struct {
	Array []struct {
		ProductIds []string `json:"ProductIds"`
		Tid        string   `json:"Tid"`
		Price      int      `json:"price"`
	} `json:"owner"`
}

type clientsResource struct{}

//This function is used to mount the new routes
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

//Given an user ID return its transactions, similarbuyers and recommended products
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
	simbuyers, favPlist := getSimilarBuyers(id)
	m["simBuyers"] = simbuyers
	m["favProducts"] = favPlist

	// fmt.Println(m["simBuyers"])
	newData, err := json.Marshal(m)
	respByte := bytes.NewReader(newData)

	if _, err := io.Copy(w, respByte); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

//Get the name of an user given its ID
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
	fmt.Println("FLAG", (r.Owner))
	ans := r.Owner[0].Name
	return ans
}

//Given a list of strings it removes the duplicate values
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

//Given an user id the func returns the list of names of users that share
//the same ip addres of that user and also return a list of recommended products
func getSimilarBuyers(id string) ([]string, []string) {
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
	var Products []string
	for i := range ips {
		tnames, pnames := getNamesxTransact(ips[i].IP, listaIps, id)
		namesxTransaccion = append(namesxTransaccion, tnames...)
		Products = append(Products, pnames...)

	}
	namesNotDup := RemoveDuplicateValues(namesxTransaccion, id)
	Productnames := preferedProducts(Products, len(Products))

	return namesNotDup, Productnames

}

//Returns the names of users that used the same Ip address
func getNamesxTransact(ipaddress string, list map[string]string, id string) ([]string, []string) {
	// Given an Ip address it returns a list
	// with all the names that used that ip
	list["$a"] = ipaddress //excluding the user identified with "id"
	list["$o"] = id

	q := `query tran($a:string,$o:string){
		tran(func:type(Transaction))@filter(NOT eq(Cid,$o) AND eq(Ip,$a)){
			CC as Cid
			Ip
			ProductIds
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
	transactlist := r2.Tran

	var res []string
	for i := range n {
		res = append(res, n[i].Name)
	}

	var ans []string
	for i := range transactlist {
		ans = append(ans, transactlist[i].ProductIds...)
	}
	// fmt.Println(len(Productnames))
	return res, ans

}

// Given and ID a query is executed on the db and returns all the clients transactions
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
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(r)
	//orders := r.Array
	for j := range r.Array {
		r.Array[j].Price = getOrderPrice(r.Array[j].ProductIds)
		temp := Idtoname(r.Array[j].ProductIds)
		r.Array[j].ProductIds = temp
	}
	ans, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return ans

}

func getOrderPrice(plist []string) int {
	ans := 0
	for i := range plist {
		ans = ans + m[plist[i]].Price
	}
	return ans
}

func Idtoname(plist []string) []string {
	var ans []string
	for i := range plist {
		ans = append(ans, m[plist[i]].Name)
	}
	return ans
}

func preferedProducts(plist []string, total int) []string {
	keys := make(map[string]int)
	list := []string{}
	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range plist {
		_, value := keys[entry]
		if value {
			keys[entry] += 1
		} else {
			keys[entry] = 1
		}
	}
	min := float64(keys[plist[0]]) / float64(total)
	max := float64(keys[plist[0]]) / float64(total)

	for _, element := range keys {
		if (float64(element) / float64(total)) > max {
			max = float64(element) / float64(total)
		}
		if (float64(element) / float64(total)) < min {
			min = float64(element) / float64(total)

		}
	}
	taux := max - min
	thresh := taux * float64(0.75)
	for key, element := range keys {
		p := float64(element) / float64(total)
		if p >= thresh {
			list = append(list, m[key].Name)
		}
	}

	return list
}
