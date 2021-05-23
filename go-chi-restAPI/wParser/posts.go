package wParser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Product struct {
	Id    string   `json:"Pid"`
	Name  string   `json:"name"`
	Price int      `json:"price"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Transaction struct {
	Tid        string   `json:"Tid,omitempty"`
	Ip         string   `json:"Ip"`
	BuyerId    string   `json:"Cid"`
	Device     string   `json:"Device"`
	ProductIds []string `json:"ProductIds"`
	DType      []string `json:"dgraph.type,omitempty"`
}

type Products []Product

type Clients []Client

type Transactions []Transaction

func GetBody(date, url string) []byte {
	now := time.Now()
	sec := now.Unix()
	if date == "" {
		date = strconv.FormatInt(sec, 10)
	}
	str := []string{url, "?date=", date}
	dataurl := strings.Join(str, "")
	resp, err := http.Get(dataurl)

	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func ReadClients(date string) Clients {
	var clientList Clients
	body := GetBody(date, "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/buyers")
	err := json.Unmarshal(body, &clientList)
	if err != nil {
		panic(err)
	}
	return clientList
}

func ReadProducts(date string) Products {
	var productsList Products
	body := GetBody(date, "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/products")
	responseString := string(body)
	pList := strings.Split(responseString, "\n")
	fmt.Println(len(pList))

	for i := 0; i < len(pList)-1; i++ {
		p := strings.Split(pList[i], "'")
		size := len(p)
		price, err := strconv.Atoi(p[size-1])
		if err != nil {
			panic(err)
		}
		prod := Product{
			Id:    p[0],
			Name:  strings.Join(p[1:size-1], " "),
			Price: price,
			DType: []string{"Product"},
		}
		productsList = append(productsList, prod)
	}
	return productsList

}

func ReadTransactions(date string) Transactions {
	var transactionList Transactions
	body := GetBody(date, "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com/transactions")
	responseString := string(body)
	tmplist := strings.Split(responseString, "#")

	for i := 1; i < len(tmplist); i++ {
		t := strings.Split(tmplist[i], "\x00")
		s := string(t[4])
		res := strings.Split(s[1:len(s)-1], ",")
		Tr := Transaction{
			Tid:        t[0],
			BuyerId:    t[1],
			Ip:         t[2],
			Device:     t[3],
			ProductIds: res,
			DType:      []string{"Transaction"},
		}
		transactionList = append(transactionList, Tr)
	}
	return transactionList
}
