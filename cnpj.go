package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"unicode"

	"github.com/gorilla/mux"
)

// Info Empresa ..
type Info struct {
	DadosEmpresa Empresa `json:"empresa,omitempty"`
}

// Empresa dados ..
type Empresa struct {
	Cnpj              string               `json:"cnpj"`
	UltimaAtualizacao string               `json:"ultima_atualizacao"`
	Abertura          string               `json:"abertura"`
	Nome              string               `json:"nome"`
	Fantasia          string               `json:"fantasia"`
	Status            string               `json:"status"`
	Tipo              string               `json:"tipo"`
	Situacao          string               `json:"situacao"`
	CapitalSocial     string               `json:"capital_social"`
	DadosEndereco     Endereco             `json:"endereco,omitempty"`
	DadosContato      Contato              `json:"contato,omitempty"`
	DadosAtividade    []AtividadePrincipal `json:"atividade_principal"`
}

type Endereco struct {
	Bairro      string `json:"bairro"`
	Logradouro  string `json:"logradouro"`
	Numero      string `json:"numero"`
	Cep         string `json:"cep"`
	Municipio   string `json:"municipio"`
	Uf          string `json:"uf"`
	Complemento string `json:"complemento"`
}

// Contato Empresa
type Contato struct {
	Telefone string `json:"telefone"`
	Email    string `json:"email"`
}

// AtividadePrincipal Empresa
type AtividadePrincipal struct {
	Text string `json:"text"`
	Code string `json:"code"`
}

// verificar se um número é inteiro
func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func GetCnpjEndpoint(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)

	if isInt(params["id"]) {

		response, err := http.Get("https://www.receitaws.com.br/v1/cnpj/" + params["id"])

		dataMensagem, _ := ioutil.ReadAll(response.Body)

		payload := make(map[string]interface{})
		json.Unmarshal(dataMensagem, &payload)

		//verificando se o CNPJ é valido
		if payload["status"] == "ERROR" {
			dataMensagem, _ := json.Marshal(payload)
			w.Header().Set("Content-Type", "application/json")
			w.Write(dataMensagem)
			return
		} else if err != nil {
			fmt.Printf("A solicitação HTTP falhou com erro %s\n", err)
		} else {

			var info Info

			//populando os dados do JSON na struct EMPRESA
			err = json.Unmarshal(dataMensagem, &info.DadosEmpresa)

			//populando os dados do JSON na struct EMPRESA->ENDERECO
			err = json.Unmarshal(dataMensagem, &info.DadosEmpresa.DadosEndereco)

			//populando os dados do JSON na struct EMPRESA->CONTATO
			err = json.Unmarshal(dataMensagem, &info.DadosEmpresa.DadosContato)

			if err != nil {
				fmt.Println(err)
				return
			}
			dataR, _ := json.Marshal(info)

			fmt.Println(string(dataR))
			w.Header().Set("Content-Type", "application/json")
			w.Write(dataR)
		}

	} else {

		info := `{erro:'Informe corretamente o seu CNPJ'}`

		mensagemErro, _ := json.Marshal(info)

		w.Header().Set("Content-Type", "application/json")
		w.Write(mensagemErro)
	}

}
