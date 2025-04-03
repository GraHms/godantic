package godantic

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (g *Validate) checkDecimalConstraints(f reflect.StructField, v reflect.Value, tree string) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if !(v.Kind() == reflect.Float32 || v.Kind() == reflect.Float64) {
		return nil
	}

	value := v.Float()

	str := strconv.FormatFloat(value, 'f', -1, 64) // string sem notação científica
	str = strings.TrimRight(str, "0")              // remove zeros à direita
	str = strings.TrimSuffix(str, ".")             // remove ponto solto

	parts := strings.Split(str, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	// Tags
	maxDigitsTag := f.Tag.Get("max_digits")
	decimalPlacesTag := f.Tag.Get("decimal_places")

	if maxDigitsTag != "" {
		if maxDigits, err := strconv.Atoi(maxDigitsTag); err == nil {
			totalDigits := len(strings.TrimLeft(intPart, "-")) + len(decPart)
			if totalDigits > maxDigits {
				return &Error{
					ErrType: "MAX_DIGITS_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must have at most %d total digits (got %d)", fieldName(f, tree), maxDigits, totalDigits),
				}
			}
		}
	}

	if decimalPlacesTag != "" {
		if decPlaces, err := strconv.Atoi(decimalPlacesTag); err == nil {
			if len(decPart) > decPlaces {
				return &Error{
					ErrType: "DECIMAL_PLACES_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must have at most %d decimal places (got %d)", fieldName(f, tree), decPlaces, len(decPart)),
				}
			}
		}
	}

	return nil
}
