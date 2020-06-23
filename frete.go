package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Recebendo os dados
type Dados struct {
	Token                string       `json:"token"`
	CodigoPlataforma     string       `json:"codigo_plataforma"`
	DadosRemetente       Remetente    `json:"remetente"`
	DadosDestinario      Destinatario `json:"destinatario"`
	DadosVolumes         Volumes      `json:"volumes"`
	Filtro               int          `json:"filtro,omitempty"`
	Canal                string       `json:"canal,omitempty"`
	Limite               int          `json:"limite,omitempty"`
	CotacaoPlafaforma    int          `json:"cotacao_plataforma,omitempty"`
	RetornarConsolidacao bool         `json:"retornar_consolidacao,omitempty"`
}

// Recebendo os volumes do frete
type Volumes []struct {
	Tipo           int     `json:"tipo"`
	Sku            string  `json:"sku"`
	Descricao      string  `json:"descricao"`
	Quantidade     int64   `json:"quantidade"`
	Altura         float64 `json:"altura"`
	Largura        float64 `json:"largura"`
	Comprimento    float64 `json:"comprimento"`
	Peso           float64 `json:"peso"`
	Valor          float64 `json:"valor"`
	VolumesProduto int64   `json:"volumes_produto,omitempty"`
	Consolidar     bool    `json:"consolidar,omitempty"`
	Sobreposto     bool    `json:"sobreposto,omitempty"`
	Tombar         bool    `json:"tombar,omitempty"`
}

// recebendo o cnpj do remetente
type Remetente struct {
	Cnpj string `json:"cnpj"`
}

// destinatario do calculo de frete
type Destinatario struct {
	TipoPessoa        int64                `json:"tipo_pessoa"`
	CnpjCpf           string               `json:"cnpj_cpf"`
	InscricaoEstadual string               `json:"inscricao_estadual,omitempty"`
	DadosEndereco     EnderecoDestinatario `json:"endereco"`
}

// endereco do destinario
type EnderecoDestinatario struct {
	Cep string `json:"cep"`
}

// Enviando a resposta
type Resposta struct {
	InfoTransportadora Transportadora `json:"transportadoras"`
}

// resposta com as transportadoras disponiveis
type Transportadora []struct {
	Nome         string  `json:"nome,omitempty"`
	Servico      string  `json:"servico,omitempty"`
	PrazoEntrega int     `json:"prazo_entrega,omitempty"`
	PrecoFrete   float64 `json:"preco_frete,omitempty"`
}

// recebendo as mensagens de erros
type Erros struct {
	Mensagem string
	Value    string
	Field    string
}

// imprimindo as mensagens de erros
func MensagemErro(errors Erros, err error, w http.ResponseWriter) bool {
	if err != nil { // caso aconteça algum erro na conversão dos tipos
		dataErro, _ := json.Marshal(err)
		json.Unmarshal(dataErro, &errors)
		errors.Mensagem = "Erro: Você atribuiu um tipo de valor errado para o campo"
		imprimirErrors, _ := json.Marshal(errors)

		w.Header().Set("Content-Type", "application/json")
		w.Write(imprimirErrors)
		return true
	} else {
		return false
	}
}

func PostFreteEndpoint(w http.ResponseWriter, req *http.Request) {
	//verificando se possui acesso a internet
	if !connected() {
		// caso o usuario não informou corretamente o cnpj
		info := `{erro:'SEM INTERNET'}`

		mensagemErro, _ := json.Marshal(info)

		w.Header().Set("Content-Type", "application/json")
		w.Write(mensagemErro)
		return
	}

	data, _ := ioutil.ReadAll(req.Body)

	var dado Dados
	var errors Erros

	//armazenando o corpo da requisição na struct dados
	err := json.Unmarshal(data, &dado)

	if MensagemErro(errors, err, w) {
		return
	} else {
		// adicionando os seguintes valores para os campo de token, codigoPlataforma ...para fazer a requisição
		dado.Token = "c8359377969ded682c3dba5cb967c07b"
		dado.CodigoPlataforma = "588604ab3"
		dado.DadosRemetente.Cnpj = "17184406000174"

		// caso o usuario digite o tipo da pessoa == 1, os seguintes campos serão nulos, pois só
		// é necessario informa o cnpj_cpf e inscrição estadual caso o tipo da pessoa == 2
		if dado.DadosDestinario.TipoPessoa == 1 {
			dado.DadosDestinario.CnpjCpf = ""
			dado.DadosDestinario.InscricaoEstadual = " "
		}

		dataR, _ := json.Marshal(dado)

		response, _ := http.Post("https://freterapido.com/api/external/embarcador/v1/quote-simulator", "application/json", bytes.NewBuffer(dataR))

		// caso aconteça algum erro na resposta da cotação
		if response.StatusCode == 400 {
			dataErr, _ := ioutil.ReadAll(response.Body)
			w.Header().Set("Content-Type", "application/json")
			w.Write(dataErr)
			return
		} else if response.StatusCode == 200 { // caso der tudo certo na resposta da cotação
			dataResult, _ := ioutil.ReadAll(response.Body)

			var resultado Resposta
			json.Unmarshal(dataResult, &resultado)

			//armazenando os dados recebidos na struct
			err := json.Unmarshal(dataResult, &resultado)

			// caso aconteça algum erro na atribuição dos dados na struct
			if MensagemErro(errors, err, w) {
				return
			} else {

				corpoResultado, _ := json.Marshal(resultado)
				//exibindo a resposta da cotação em json
				w.Header().Set("Content-Type", "application/json")
				w.Write(corpoResultado)
				return
			}

		}

	}

}
