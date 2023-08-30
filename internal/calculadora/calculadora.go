package calculadora

import (
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
}

type Calculo struct {
	Titulo         Titulo
	DataLiquidacao time.Time
	DU             int
	Preco          float64
	Taxa           float64
}

type fluxoPagamento struct {
	dataPagamento time.Time
	du            int
	pagamento     float64
	valorPresente float64
}

func NovoCalculo(tipo TipoTitulo, vencimento time.Time, liquidacao time.Time) Calculo {
	nc := Calculo{
		Titulo:         Titulo{Tipo: tipo, Vencimento: vencimento},
		DataLiquidacao: liquidacao,
		Taxa:           0.100696,
	}

	//define o valor do cupom com base no tipo de título
	nc.definirCupom()

	//pre calcula a quantidade de dias úteis
	var err error
	nc.DU, err = diaTrabalhoTotal(nc.DataLiquidacao, nc.Titulo.Vencimento)
	if err != nil {
		panic("Erro ao calcular a quantidade de dias úteis")
	}

	return nc
}

func (calc *Calculo) PrecificarLTN() {
	calc.Preco = truncar(1000/math.Pow(1+calc.Taxa, truncar(float64(calc.DU)/252, 14)), 6)
}

// define o percentual do cupom ajustado para pagamentos semestrais, de acordo com o tipo de título
func (calc *Calculo) definirCupom() {
	ajusteSemestral := 0.5
	if calc.Titulo.Tipo == NTN_F {
		calc.Titulo.Cupom = math.Pow(1.10, ajusteSemestral) - 1
	} else if calc.Titulo.Tipo == NTN_B {
		calc.Titulo.Cupom = math.Pow(1.06, ajusteSemestral) - 1
	}
}

func truncar(value float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Trunc(value*shift) / shift
}
