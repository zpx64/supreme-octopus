package conversions

import (
	"github.com/zpx64/supreme-octopus/internal/db"
	"github.com/zpx64/supreme-octopus/pkg/nationalize"
)

func CountriesToNationalization(countries []nationalize.Country) []db.Nationalization {
	out := make([]db.Nationalization, len(countries))
	for i := range countries {
		out[i] = db.Nationalization{
			CountryCode: countries[i].Id,
			Probability: countries[i].Probability,
		}
	}
	return out
}

func CountryToIds(countries []nationalize.Country) []string {
	out := make([]string, 0, len(countries))
	for _, e := range countries {
		out = append(out, e.Id)
	}
	return out
}

func CountryToProbalities(countries []nationalize.Country) []float64 {
	out := make([]float64, 0, len(countries))
	for _, e := range countries {
		out = append(out, e.Probability)
	}
	return out
}
