package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rinhabackend/shared/dto"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

func TestCreateTransactions(t *testing.T) {
	timeTesting := time.Now()
	var (
		wg sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		internalTimer := time.Now()
		input, err := json.Marshal(map[string]any{
			"valor":     1000,
			"tipo":      "c",
			"descricao": "descricao",
		})
		assert.NoError(t, err)
		payload := bytes.NewBuffer(input)
		res, err := http.Post("http://localhost:3000/clientes/2/transacoes", "application/json", payload)
		assert.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, 200, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		fmt.Println("Body:", string(body))
		fmt.Println("GoRoutine01:", time.Now().Sub(internalTimer))
		var data Response
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
	}()

	go func() {
		internalTimer := time.Now()
		defer wg.Done()
		input, err := json.Marshal(map[string]any{
			"valor":     1000,
			"tipo":      "c",
			"descricao": "descricao",
		})
		assert.NoError(t, err)
		payload := bytes.NewBuffer(input)
		res, err := http.Post("http://localhost:3000/clientes/2/transacoes", "application/json", payload)
		assert.NoError(t, err)
		defer res.Body.Close()
		assert.Equal(t, 200, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		fmt.Println("Body:", string(body))
		fmt.Println("GoRoutine02:", time.Now().Sub(internalTimer))
		var data Response
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
	}()

	wg.Wait()

	res, err := http.Get("http://localhost:3000/clientes/2/extrato")
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
	var account dto.GetBalanceOutputDTO
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &account)
	assert.NoError(t, err)
	assert.Equal(t, int(2000), account.Ballance.Total)
	assert.Equal(t, 2, len(account.Transactions))
	fmt.Println("")
	fmt.Println("extrato: ", string(body))
	fmt.Println("")
	fmt.Println("TimeEnd Test:", time.Now().Sub(timeTesting))
}
