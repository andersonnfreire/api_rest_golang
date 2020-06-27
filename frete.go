package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Recebendo os dados
type Dados struct {
	Token                string       `json:"token,omitempty"`
	CodigoPlataforma     string       `json:"codigo_plataforma,omitempty"`
	DadosRemetente       Remetente    `json:"remetente,omitempty"`
	DadosDestinario      Destinatario `json:"destinatario,omitempty"`
	DadosVolumes         []Volumes    `json:"volumes,omitempty"`
	Filtro               int          `json:"filtro,omitempty"`
	Canal                string       `json:"canal,omitempty"`
	Limite               int          `json:"limite,omitempty"`
	CotacaoPlafaforma    int          `json:"cotacao_plataforma,omitempty"`
	RetornarConsolidacao bool         `json:"retornar_consolidacao,omitempty"`
}

// Recebendo os volumes do frete
type Volumes struct {
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
	InfoTransportadora []Transportadora `json:"transportadoras"`
}

// resposta com as transportadoras disponiveis
type Transportadora struct {
	Nome         string  `json:"nome"`
	Servico      string  `json:"servico"`
	PrazoEntrega int     `json:"prazo_entrega"`
	PrecoFrete   float64 `json:"preco_frete"`
}

// recebendo as mensagens de erros de conversao de tipo
type conversaoErro struct {
	Value string
	Field string
}

// resposta formatada da mensagem de erro
type mensagemErro struct {
	Status     string `json:"status"`
	Informacao string `json:"mensagem"`
	Codigo     int    `json:"codigo,omitempty"`
}

// imprimindo as mensagens de erros
func formatandoMensagemErro(errors conversaoErro, err error, w http.ResponseWriter) bool {

	var info mensagemErro

	if err != nil { // caso aconteça algum erro na conversão dos tipos
		dataErro, _ := json.Marshal(err)

		//armazenando erro de conversao do tipo
		json.Unmarshal(dataErro, &errors)

		// setando o status e codigo da mensagem de erro para o conversao de tipo
		info.Status = "Erro: " + http.StatusText(http.StatusBadRequest)
		info.Codigo = http.StatusBadRequest
		//verificando se caso o usuario informou o tipo errado da variavel
		if errors.Field != "" && errors.Value != "" {
			info.Informacao = "O campo " + errors.Field + " não está correto pois foi digitado um(a) " + errors.Value + ""
		} else { // caso o usuario não informou nada
			info.Informacao = "Preencha os campos da forma correta"
		}

		dataConversao, _ := json.Marshal(info)

		w.Header().Set("Content-Type", "application/json")
		w.Write(dataConversao)
		return true
	} else {
		return false
	}
}
func respostaHttp(dataResposta []byte, statusCode int, infoErro mensagemErro, errors conversaoErro, err error, w http.ResponseWriter) bool {

	// verificando se a resposta Http foi se um erro
	if statusCode >= 400 && statusCode < 512 {
		infoErro.Status = "Erro: " + http.StatusText(statusCode)
	} else if statusCode >= 300 && statusCode < 309 { // verificando se a resposta Http foi se um redirecionamento
		infoErro.Status = "URI foi mudada"
	} else if statusCode == 200 { // verificando se a resposta Http foi um sucesso
		return false
	}

	infoErro.Informacao = string(dataResposta)
	infoErro.Codigo = statusCode
	corpoResultado, err := json.Marshal(infoErro)

	// verificando se aconteceu algum erro na formatacao da mensagem de erro
	if formatandoMensagemErro(errors, err, w) {
		return true
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(corpoResultado)
	return true

}

func verificarConexao(infoErro mensagemErro, errors conversaoErro, err error, w http.ResponseWriter) bool {
	//caso não exista conexão com a internet
	if err != nil {
		infoErro.Status = "Erro: " + http.StatusText(http.StatusServiceUnavailable)
		infoErro.Informacao = "Verifique sua conexão de rede."
		infoErro.Codigo = http.StatusServiceUnavailable
		mensagemErro, err := json.Marshal(infoErro)

		// verificando se aconteceu algum erro na formatacao da mensagem de erro
		if formatandoMensagemErro(errors, err, w) {
			return true
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(mensagemErro)
		return true
	}
	return false
}

func PostFreteEndpoint(w http.ResponseWriter, req *http.Request) {

	//requisitando a seguinte url para verificar se possui acesso a internet
	response, err := http.Get("http://clients3.google.com/generate_204")
	var infoErro mensagemErro
	var errors conversaoErro

	// verificando se possui acesso a internet
	if verificarConexao(infoErro, errors, err, w) {
		return
	} else if response.StatusCode == 204 { // caso obtenha conexão com a internet

		data, err := ioutil.ReadAll(req.Body)

		// verificando se possui algum erro na leitura do corpo da requisição
		if err != nil {

			infoErro.Status = "Erro"
			infoErro.Informacao = "O servidor não entendeu a requisição pois está com uma sintaxe inválida."
			infoErro.Codigo = http.StatusBadRequest

			mensagemErro, err := json.Marshal(infoErro)

			// verificando se aconteceu algum erro na formatacao da mensagem de erro
			if formatandoMensagemErro(errors, err, w) {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(mensagemErro)
			return
		}

		var dado Dados
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

		//armazenando o corpo da requisição na struct Dados
		err = json.Unmarshal(data, &dado)

		if formatandoMensagemErro(errors, err, w) {
			fmt.Println("OK")
			return
		}

		dataR, err := json.Marshal(dado)

		// verificando se aconteceu algum erro na formatacao da mensagem de erro
		if formatandoMensagemErro(errors, err, w) {
			return
		}

		// fazendo a requisição para calcular o frete
		response, err := http.Post("https://freterapido.com/api/external/embarcador/v1/quote-simulator", "application/json", bytes.NewBuffer(dataR))

		// verificando se aconteceu algum erro
		if formatandoMensagemErro(errors, err, w) {
			return
		}

		dataResult, err := ioutil.ReadAll(response.Body)

		// verificando se aconteceu erro na leitura da requisicao
		if formatandoMensagemErro(errors, err, w) {
			return
		} else if respostaHttp(dataResult, response.StatusCode, infoErro, errors, err, w) { // se caso a resposta http for != 200 e encerrado
			return
		} else {
			// verificando se aconteceu erro na leitura da requisicao
			var resultado Resposta

			//armazenando os dados recebidos na struct
			err = json.Unmarshal(dataResult, &resultado)

			// caso aconteça algum erro na atribuição dos dados na struct
			if formatandoMensagemErro(errors, err, w) {
				return
			}

			corpoResultado, err := json.Marshal(resultado)

			// verificando se aconteceu erro na leitura da requisicao
			if formatandoMensagemErro(errors, err, w) {
				return
			}
			//exibindo a resposta da cotação em json
			w.Header().Set("Content-Type", "application/json")
			w.Write(corpoResultado)
			return
		}

	}

}
