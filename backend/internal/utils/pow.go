package utils

import (
	"github.com/zpx64/supreme-octopus/internal/vars"

	"github.com/zpx64/supreme-octopus/pkg/cryptograph"
)

// concat pow with string
func PowCat(str, pow string) string {
	if vars.PowRightCat {
		return str + pow
	}
	return pow + str
}

func HashPassWithPows(pass, localPow string) string {
	return cryptograph.HashPass(
		PowCat(
			PowCat(pass, localPow),
			vars.GlobalPow,
		),
	)
}
