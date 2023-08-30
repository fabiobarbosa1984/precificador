package calculadora

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func IsFeriado(date time.Time) bool {
	return feriados[dateToString(date)]
}

func stringToDate(dateString string) (time.Time, error) {
	if dateString != "" {
		parts := strings.Split(dateString, "/")
		day, err := strconv.Atoi(parts[0])
		if err != nil {
			return time.Time{}, err
		}
		month, err := strconv.Atoi(parts[1])
		if err != nil {
			return time.Time{}, err
		}
		year, err := strconv.Atoi(parts[2])
		if err != nil {
			return time.Time{}, err
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, nil
}

func diaTrabalhoTotal(dataInicial, dataFinal time.Time) (int, error) {
	diasUteis := 0
	dataFinalTratada := dataFinal.AddDate(0, 0, -1)

	for date := dataInicial; !date.After(dataFinalTratada); date = date.AddDate(0, 0, 1) {
		if !IsFeriado(date) && date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			diasUteis++
		}
	}

	log.Println("Data inicial: ", dataInicial)
	log.Println("Data final considerada: ", dataFinalTratada)

	log.Println("Total de DU: ", diasUteis)
	return diasUteis, nil
}

func dateToString(date time.Time) string {
	return date.Format("02/01/2006")
}

func corrigirVencimento(vencimento string) (string, error) {
	date, err := stringToDate(vencimento)
	if err != nil {
		return "", err
	}

	datefim := date.AddDate(0, 0, 6)

	for ; !date.After(datefim); date = date.AddDate(0, 0, 1) {
		if !IsFeriado(date) && date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			return dateToString(date), nil
		}
	}
	return "", nil
}
