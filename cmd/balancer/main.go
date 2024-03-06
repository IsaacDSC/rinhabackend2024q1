package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type LoadBalancer interface {
	Proxy(method string, path *url.URL, body *io.ReadCloser) (*http.Response, error)
}

func NewLoadBalance(lb LoadBalancer) LoadBalancer {
	return lb
}

type KindRoundRobin struct {
	listeners []string
	counter   int
}

func NewKindRoundRobin(listeners []string) *KindRoundRobin {
	return &KindRoundRobin{listeners, 0}
}

func (lb *KindRoundRobin) Proxy(method string, path *url.URL, body *io.ReadCloser) (*http.Response, error) {
	lb.counter += 1
	server := lb.listeners[lb.counter%len(lb.listeners)]
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", server, path), *body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json") // todo: adicionar headers e passsar a diante
	res, err := c.Do(req)
	if err != nil {
		return res, err
	}
	return res, err
}

type KindHash struct {
	listeners []string
}

func NewKindHash(listeners []string) *KindHash {
	return &KindHash{listeners}
}

func (lb *KindHash) Proxy(method string, path *url.URL, body *io.ReadCloser) (*http.Response, error) {
	h := sha256.New()
	h.Write([]byte(path.String()))
	md := h.Sum(nil)
	signature := new(big.Int).SetBytes(md).Uint64()
	server := lb.listeners[signature%uint64(len(lb.listeners))]

	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", server, path), *body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json") // todo: adicionar headers e passsar a diante
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}

func main() {
	//msg := "asd"
	//h := sha256.New()
	//h.Write([]byte(msg))
	////signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//md := h.Sum(nil)
	//signature := new(big.Int).SetBytes(md)
	//fmt.Println("Result: ", signature)
	//[]string{"http://localhost:3001", "http://localhost:3002", "http://localhost:3003"}
	entry := os.Getenv("ENTRIES")
	entries := strings.Split(entry, ",")
	if len(entries) == 0 {
		panic("Requirement entries")
		return
	}
	lb := NewLoadBalance(NewKindRoundRobin(entries))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL
		res, err := lb.Proxy(request.Method, path, &request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		writer.WriteHeader(res.StatusCode)
		writer.Write(b)
	})

	fmt.Println("Started Load Balancer")
	log.Fatal(http.ListenAndServe(":9999", mux))
}
