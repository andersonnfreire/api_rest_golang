package main

import (
	"encoding/json"
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
	Cnpj              string             `json:"cnpj"`
	UltimaAtualizacao string             `json:"ultima_atualizacao"`
	Abertura          string             `json:"abertura"`
	Nome              string             `json:"nome"`
	Fantasia          string             `json:"fantasia"`
	Status            string             `json:"status"`
	Tipo              string             `json:"tipo"`
	Situacao          string             `json:"situacao"`
	CapitalSocial     string             `json:"capital_social"`
	DadosEndereco     Endereco           `json:"endereco,omitempty"`
	DadosContato      Contato            `json:"contato,omitempty"`
	DadosAtividade    AtividadePrincipal `json:"atividade_principal"`
}

// Endereço da empresa
type Endereco struct {
	Bairro      string `json:"bairro"`
	Logradouro  string `json:"logradouro"`
	Numero      string `json:"numero"`
	Cep         string `json:"cep"`
	Municipio   string `json:"municipio"`
	Uf          string `json:"uf"`
	Complemento string `json:"complemento"`
}

// Contato da Empresa
type Contato struct {
	Telefone string `json:"telefone"`
	Email    string `json:"email"`
}

// AtividadePrincipal da Empresa
type AtividadePrincipal []struct {
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
func connected() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}
func GetCnpjEndpoint(w http.ResponseWriter, req *http.Request) {
	//verificando se possui acesso a internet
	if !connected() {
		// caso o usuario não informou corretamente o cnpj
		info := `{erro:'SEM INTERNET'}`

		mensagemErro, _ := json.Marshal(info)

		w.Header().Set("Content-Type", "application/json")
		w.Write(mensagemErro)
		return
	}
	params := mux.Vars(req)

	//verificando se o usuario o digitou uma string no cnpj
	if isInt(params["id"]) {

		response, _ := http.Get("https://www.receitaws.com.br/v1/cnpj/" + params["id"])

		// verificando caso aconteça algum erro na resposta da requisição
		if response.StatusCode == 400 {
			dataErr, _ := ioutil.ReadAll(response.Body)
			w.Header().Set("Content-Type", "application/json")
			w.Write(dataErr)
			return
		} else if response.StatusCode == 200 {
			dataMensagem, _ := ioutil.ReadAll(response.Body)

			// caso tenha erros ao ler o response.Body, esse map será responsável por armazenar dessa leitura
			statusMensagem := make(map[string]interface{})

			json.Unmarshal(dataMensagem, &statusMensagem)

			//verificando se o CNPJ é valido e verificando se possui algum erro
			if statusMensagem["status"] == "ERROR" {
				dataMensagem, _ := json.Marshal(statusMensagem)
				w.Header().Set("Content-Type", "application/json")
				w.Write(dataMensagem)
				return
			} else {

				var info Info
				var errors Erros
				//populando os dados do JSON na struct EMPRESA
				err := json.Unmarshal(dataMensagem, &info.DadosEmpresa)

				//verifica se os dados da empresa foram inseridos
				if MensagemErro(errors, err, w) {
					return
				} else {

					//populando os dados do JSON na struct EMPRESA->ENDERECO
					err := json.Unmarshal(dataMensagem, &info.DadosEmpresa.DadosEndereco)

					//verifica se o Endereco da empresa foi inserido
					if MensagemErro(errors, err, w) {
						return
					} else {

						//populando os dados do JSON na struct EMPRESA->CONTATO
						err := json.Unmarshal(dataMensagem, &info.DadosEmpresa.DadosContato)

						//verifica se o Contato da empresa foi inserido
						if MensagemErro(errors, err, w) {
							return
						}

						dataR, _ := json.Marshal(info)

						w.Header().Set("Content-Type", "application/json")
						w.Write(dataR)
						return
					}
				}
			}
		}

	} else {
		// caso o usuario não informou corretamente o cnpj
		info := `{erro:'Informe corretamente o seu CNPJ'}`

		mensagemErro, _ := json.Marshal(info)

		w.Header().Set("Content-Type", "application/json")
		w.Write(mensagemErro)
		return
	}

}
