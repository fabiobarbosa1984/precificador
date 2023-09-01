package calculadora

import (
	"fmt"
	"math"
	"time"
)

type TipoTitulo string

const (
	LTN   TipoTitulo = "LTN"
	NTN_F TipoTitulo = "NTN-F"
	NTN_B TipoTitulo = "NTN-B"
	LFT   TipoTitulo = "LFT"
)

type Titulo struct {
	Tipo       TipoTitulo
	Vencimento time.Time
	Cupom      float64
	valorFace  float64
}

type Calculo struct {
	Titulo         Titulo
	DataLiquidacao time.Time
	primeiroCupom  time.Time
	DU             int
	Preco          float64
	Taxa           float64
	fluxoPagamento []fluxoPagamento
}

type fluxoPagamento struct {
	dataPagamento time.Time
	du            int
	pagamento     float64
	valorPresente float64
}

func NovoCalculo(tipo TipoTitulo, vencimento time.Time, liquidacao time.Time, taxa float64) Calculo {
	nc := Calculo{
		Titulo:         Titulo{Tipo: tipo, Vencimento: vencimento},
		DataLiquidacao: liquidacao,
		Taxa:           taxa,
	}

	//define o valor do cupom com base no tipo de título
	nc.definirCupom()

	//pre calcula a quantidade de dias úteis
	var err error
	nc.DU, err = diaTrabalhoTotal(nc.DataLiquidacao, nc.Titulo.Vencimento)
	if err != nil {
		panic("Erro ao calcular a quantidade de dias úteis")
	}

	//calcula o fluxo de pagamentos
	if nc.Titulo.Tipo == NTN_F || nc.Titulo.Tipo == NTN_B {
		nc.calcularPrimeiroCupom()
		nc.calcularFluxoPagamento()
	}

	return nc
}

func (calc *Calculo) calcularFluxoPagamento() {
	calc.DataLiquidacao.Month()

	if calc.Titulo.Tipo == NTN_F {
		calc.calcularFluxoPagamentoNTN_F()
	}
}

func (calc *Calculo) calcularPrimeiroCupom() {
	var primeiroCupom time.Time

	if calc.Titulo.Tipo == NTN_F {
		cup07 := time.Date(calc.DataLiquidacao.Year(), time.July, 1, 0, 0, 0, 0, time.UTC)
		cup01 := time.Date(calc.DataLiquidacao.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)

		if calc.DataLiquidacao.After(cup07) {
			primeiroCupom = cup01
		}
	} else {
		cup02aa := time.Date(calc.DataLiquidacao.Year(), time.February, 1, 0, 0, 0, 0, time.UTC)
		cup02pa := time.Date(calc.DataLiquidacao.Year()+1, time.February, 1, 0, 0, 0, 0, time.UTC)
		cup08 := time.Date(calc.DataLiquidacao.Year(), time.August, 1, 0, 0, 0, 0, time.UTC)

		if calc.DataLiquidacao.After(cup08) {
			primeiroCupom = cup02pa
		} else if calc.DataLiquidacao.Before(cup02aa) {
			primeiroCupom = cup02aa
		} else {
			primeiroCupom = cup08
		}
	}
	calc.primeiroCupom = primeiroCupom
}

func (calc *Calculo) calcularFluxoPagamentoNTN_F() {
	cupomAtual := calc.primeiroCupom

	for cupomAtual.Before(calc.Titulo.Vencimento) {
		du, _ := diaTrabalhoTotal(calc.DataLiquidacao, cupomAtual)
		vf := arred(calc.Titulo.valorFace*calc.Titulo.Cupom, 5)
		vp := valorPresente(vf, calc.Taxa, du)

		fp := fluxoPagamento{
			dataPagamento: cupomAtual,
			du:            du,
			pagamento:     vf,
			valorPresente: vp,
		}

		calc.fluxoPagamento = append(calc.fluxoPagamento, fp)
		cupomAtual = cupomAtual.AddDate(0, 6, 0)
	}

	du, _ := diaTrabalhoTotal(calc.DataLiquidacao, cupomAtual)
	vf := arred(calc.Titulo.valorFace*(1+calc.Titulo.Cupom), 5)
	vp := valorPresente(vf, calc.Taxa, du)

	fp := fluxoPagamento{
		dataPagamento: cupomAtual,
		du:            du,
		pagamento:     vf,
		valorPresente: vp,
	}
	calc.fluxoPagamento = append(calc.fluxoPagamento, fp)

	for _, fp := range calc.fluxoPagamento {
		fmt.Println(fp.dataPagamento, "\t", fp.du, "\t", fp.pagamento, "\t", fp.valorPresente)
		calc.Preco = calc.Preco + fp.valorPresente
	}
	fmt.Println("PU: ", calc.Preco)
}

func (calc *Calculo) PrecificarLTN() {
	calc.Preco = valorPresente(1000.00, calc.Taxa, calc.DU)
}

func valorPresente(valorFace float64, taxa float64, du int) float64 {
	return truncar(valorFace/math.Pow(1+taxa, truncar(float64(du)/252, 14)), 6)
}

// define o percentual do cupom ajustado para pagamentos semestrais, de acordo com o tipo de título
func (calc *Calculo) definirCupom() {
	ajusteSemestral := 0.5
	if calc.Titulo.Tipo == NTN_F {
		calc.Titulo.Cupom = math.Pow(1.10, ajusteSemestral) - 1
		calc.Titulo.valorFace = 1000.00
	} else if calc.Titulo.Tipo == NTN_B {
		calc.Titulo.Cupom = math.Pow(1.06, ajusteSemestral) - 1
		calc.Titulo.valorFace = 100.00
	}

}

func truncar(value float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Trunc(value*shift) / shift
}

func arred(value float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*shift) / shift
}
