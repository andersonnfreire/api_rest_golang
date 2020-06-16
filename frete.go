package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Dados struct {
	Token            string `json:"token"`
	CodigoPlataforma string `json:"codigo_plataforma"`
	Remetente        struct {
		Cnpj string `json:"cnpj"`
	} `json:"remetente"`
	Destinatario struct {
		TipoPessoa        int64  `json:"tipo_pessoa"`
		CnpjCpf           string `json:"cnpj_cpf"`
		InscricaoEstadual string `json:"inscricao_estadual,omitempty"`
		Endereco          struct {
			Cep string `json:"cep"`
		} `json:"endereco"`
	} `json:"destinatario"`
	Volumes []struct {
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
	} `json:"volumes"`
	Filtro               int    `json:"filtro,omitempty"`
	Canal                string `json:"canal,omitempty"`
	Limite               int    `json:"limite,omitempty"`
	CotacaoPlafaforma    int    `json:"cotacao_plataforma,omitempty"`
	RetornarConsolidacao bool   `json:"retornar_consolidacao,omitempty"`
}
type Resposta struct {
	Transportadoras []struct {
		Nome         string  `json:"nome,omitempty"`
		Servico      string  `json:"servico,omitempty"`
		PrazoEntrega int64   `json:"prazo_entrega,omitempty"`
		PrecoFrete   float64 `json:"preco_frete,omitempty"`
	} `json:"transportadoras"`
}

type Erros struct {
	Mensagem string
	Value    string
	Field    string
}

func PostFreteEndpoint(w http.ResponseWriter, req *http.Request) {

	data, _ := ioutil.ReadAll(req.Body)

	var dado Dados
	var errors Erros
	//armazenando os dados recebidos na struct
	err := json.Unmarshal(data, &dado)

	if err != nil {

		dataErro, _ := json.Marshal(err)
		json.Unmarshal(dataErro, &errors)
		errors.Mensagem = "Erro"
		imprimirErrors, _ := json.Marshal(errors)

		w.Header().Set("Content-Type", "application/json")
		w.Write(imprimirErrors)
	}

	dado.Token = "c8359377969ded682c3dba5cb967c07b"
	dado.CodigoPlataforma = "588604ab3"
	dado.Remetente.Cnpj = "17184406000174"

	if dado.Destinatario.TipoPessoa == 1 {
		dado.Destinatario.CnpjCpf = ""
		dado.Destinatario.InscricaoEstadual = " "
	}

	dataR, _ := json.Marshal(dado)

	response, err := http.Post("https://freterapido.com/api/external/embarcador/v1/quote-simulator", "application/json", bytes.NewBuffer(dataR))

	if response.StatusCode == 400 {
		dataErr, _ := ioutil.ReadAll(response.Body)
		fmt.Println(dataErr)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataErr)
	}

	if err != nil {
		fmt.Printf("A solicitação HTTP falhou com erro %s\n", err)
	} else {
		dataResult, _ := ioutil.ReadAll(response.Body)

		var resultado Resposta

		//armazenando os dados recebidos na struct
		err := json.Unmarshal(dataResult, &resultado)

		if err != nil {
			fmt.Println(err)
			return
		}

		//exibindo somente as transportadoras de EXPRESSO FR e Correios
		for i := 0; i < len(resultado.Transportadoras); i++ {
			if strings.Compare(resultado.Transportadoras[i].Nome, "EXPRESSO FR") != 0 && strings.Compare(resultado.Transportadoras[i].Nome, "CORREIOS") != 0 {
				resultado.Transportadoras = append(resultado.Transportadoras[:i], resultado.Transportadoras[i+1:]...)
				i--
			}

		}
		corpoResultado, _ := json.Marshal(resultado)

		w.Header().Set("Content-Type", "application/json")
		w.Write(corpoResultado)

	}

}
