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
	DadosAtividade    []AtividadePrincipal `json:"atividade_principal,omitempty"`
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
	//requisitando a seguinte url para verificar se possui acesso a internet

	response, err := http.Get("http://clients3.google.com/generate_204")
	var infoErro mensagemErro
	var errors conversaoErro
	var info Info

	// verificando se possui acesso a internet
	if verificarConexao(infoErro, errors, err, w) {
		return
	} else if response.StatusCode == 204 { // caso tenha acesso a internet
		params := mux.Vars(req)

		//verificando se o usuario o digitou uma string no cnpj
		if isInt(params["id"]) {

			// fazendo a requisição para obter os dados da empresa de acordo com o cnpj
			response, err := http.Get("https://www.receitaws.com.br/v1/cnpj/" + params["id"])

			// verificando se possui erros
			if formatandoMensagemErro(errors, err, w) {
				return
			}

			dataMensagem, err := ioutil.ReadAll(response.Body)

			// verificando se possui erros ao ler o escopo da resposta
			if formatandoMensagemErro(errors, err, w) {
				return
			} else if respostaHttp(dataMensagem, response.StatusCode, infoErro, errors, err, w) { // se caso a resposta http for != 200 e encerrado
				return
			} else {

				// map para armazenar possiveis erros de CNPJ invalido
				statusMensagem := make(map[string]string)

				err = json.Unmarshal(dataMensagem, &statusMensagem)

				//verificando se o CNPJ é valido e verificando se possui algum erro
				if statusMensagem["status"] == "ERROR" {
					infoErro.Status = "Erro: " + http.StatusText(http.StatusBadGateway)
					infoErro.Informacao = statusMensagem["message"]
					infoErro.Codigo = http.StatusBadGateway

					dataMensagem, err := json.Marshal(infoErro)

					// verificando se aconteceu algum erro na formatacao da mensagem de erro
					if formatandoMensagemErro(errors, err, w) {
						return
					}

					w.Header().Set("Content-Type", "application/json")
					w.Write(dataMensagem)
					return
				} else {

					//populando os dados do JSON na struct EMPRESA
					err = json.Unmarshal(dataMensagem, &info.DadosEmpresa)

					//verifica se os dados da empresa foram inseridos
					if formatandoMensagemErro(errors, err, w) {
						return
					} else {

						//populando os dados do JSON na struct EMPRESA->ENDERECO
						err = json.Unmarshal(dataMensagem, &info.DadosEmpresa.DadosEndereco)

						//verifica se o Endereco da empresa foi inserido
						if formatandoMensagemErro(errors, err, w) {
							return
						} else {

							//populando os dados do JSON na struct EMPRESA->CONTATO
							err = json.Unmarshal(dataMensagem, &info.DadosEmpresa.DadosContato)

							//verifica se o Contato da empresa foi inserido
							if formatandoMensagemErro(errors, err, w) {
								return
							}

							dataR, err := json.Marshal(info)

							//verifica se os dados foram inseridos com sucesso
							if formatandoMensagemErro(errors, err, w) {
								return
							}

							w.Header().Set("Content-Type", "application/json")
							w.Write(dataR)
							return
						}
					}
				}
			}

		} else {
			// caso o usuario não informou corretamente o cnpj
			infoErro.Status = "Erro: " + http.StatusText(http.StatusBadRequest)
			infoErro.Codigo = http.StatusBadRequest
			infoErro.Informacao = "Verifique se o CNPJ foi informado corretamente"
			mensagemErro, _ := json.Marshal(info)

			w.Header().Set("Content-Type", "application/json")
			w.Write(mensagemErro)
			return
		}
	}
}
