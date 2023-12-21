package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ResponseViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ResponseBrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func FindViaCep(cep string) (*ResponseViaCep, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var responseViaCep ResponseViaCep
	err = json.Unmarshal(body, &responseViaCep)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &responseViaCep, nil
}

func FindBrasilApi(cep string) (*ResponseBrasilApi, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var responseBrasilApi ResponseBrasilApi
	err = json.Unmarshal(body, &responseBrasilApi)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return &responseBrasilApi, nil
}

func main() {

	cep := "65909001"
	chViaCep := make(chan *ResponseViaCep)
	chBrasilApi := make(chan *ResponseBrasilApi)

	go func() {
		responseViaCep, err := FindViaCep(cep)
		if err != nil {
			fmt.Println(err)
		}
		chViaCep <- responseViaCep
	}()

	go func() {
		responseBrasilApi, err := FindBrasilApi(cep)
		if err != nil {
			fmt.Println(err)
		}
		chBrasilApi <- responseBrasilApi
	}()

	select {
	case responseViaCep := <-chViaCep:
		fmt.Println("API VIA CEP")
		fmt.Println("CEP: ", responseViaCep.Cep)
		fmt.Println("Logradouro: ", responseViaCep.Logradouro)
		fmt.Println("Complemento: ", responseViaCep.Complemento)
		fmt.Println("Bairro: ", responseViaCep.Bairro)
		fmt.Println("Localidade: ", responseViaCep.Localidade)
		fmt.Println("UF: ", responseViaCep.Uf)
		fmt.Println("IBGE: ", responseViaCep.Ibge)
		fmt.Println("DDD: ", responseViaCep.Ddd)
		fmt.Println("SIAFI: ", responseViaCep.Siafi)
	case responseBrasilApi := <-chBrasilApi:
		fmt.Println("API BRASIL API")
		fmt.Println("CEP: ", responseBrasilApi.Cep)
		fmt.Println("Logradouro: ", responseBrasilApi.Street)
		fmt.Println("Bairro: ", responseBrasilApi.Neighborhood)
		fmt.Println("Localidade: ", responseBrasilApi.City)
		fmt.Println("UF: ", responseBrasilApi.State)
		fmt.Println("Serviço: ", responseBrasilApi.Service)

	case <-time.After(time.Second * 1):
		fmt.Println("timeout error")
	}

}
