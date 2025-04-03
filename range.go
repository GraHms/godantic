package godantic

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func (g *Validate) checkMinMax(f reflect.StructField, v reflect.Value, tree string) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	minTag := f.Tag.Get("min")
	maxTag := f.Tag.Get("max")

	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		length := v.Len()
		if minTag != "" {
			min, err := strconv.Atoi(minTag)
			if err == nil && length < min {
				return &Error{
					ErrType: "MIN_LENGTH_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must have at least %d items, but has %d", fieldName(f, tree), min, length),
				}
			}
		}
		if maxTag != "" {
			max, err := strconv.Atoi(maxTag)
			if err == nil && length > max {
				return &Error{
					ErrType: "MAX_LENGTH_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must have at most %d items, but has %d", fieldName(f, tree), max, length),
				}
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := v.Int()
		if minTag != "" {
			min, err := strconv.ParseInt(minTag, 10, 64)
			if err == nil && val < min {
				return &Error{
					ErrType: "MIN_VALUE_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must be at least %d, but was %d", fieldName(f, tree), min, val),
				}
			}
		}
		if maxTag != "" {
			max, err := strconv.ParseInt(maxTag, 10, 64)
			if err == nil && val > max {
				return &Error{
					ErrType: "MAX_VALUE_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must be at most %d, but was %d", fieldName(f, tree), max, val),
				}
			}
		}
	case reflect.Float32, reflect.Float64:
		val := v.Float()
		if minTag != "" {
			min, err := strconv.ParseFloat(minTag, 64)
			if err == nil && val < min {
				return &Error{
					ErrType: "MIN_VALUE_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must be at least %.2f, but was %.2f", fieldName(f, tree), min, val),
				}
			}
		}
		if maxTag != "" {
			max, err := strconv.ParseFloat(maxTag, 64)
			if err == nil && val > max {
				return &Error{
					ErrType: "MAX_VALUE_ERR",
					Path:    fieldName(f, tree),
					Message: fmt.Sprintf("The field <%s> must be at most %.2f, but was %.2f", fieldName(f, tree), max, val),
				}
			}
		}
	}
	return nil
}

func (g *Validate) checkNumericConstraints(f reflect.StructField, v reflect.Value, tree string) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	// Skip if not numeric
	if !(v.Kind() >= reflect.Int && v.Kind() <= reflect.Float64) {
		return nil
	}

	// Get tag values
	gtTag := f.Tag.Get("gt")
	ltTag := f.Tag.Get("lt")
	geTag := f.Tag.Get("ge")
	leTag := f.Tag.Get("le")
	multipleOfTag := f.Tag.Get("multiple_of")
	allowInfNaNTag := f.Tag.Get("allow_inf_nan")

	var value float64
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		value = v.Float()
		if allowInfNaNTag != "true" && (math.IsNaN(value) || math.IsInf(value, 0)) {
			return &Error{
				ErrType: "INVALID_FLOAT_ERR",
				Path:    fieldName(f, tree),
				Message: "The field <" + fieldName(f, tree) + "> cannot be NaN or infinite",
			}
		}
	default:
		return nil
	}

	// Constraint checks
	if gtTag != "" {
		if threshold, err := strconv.ParseFloat(gtTag, 64); err == nil && !(value > threshold) {
			return &Error{
				ErrType: "GREATER_THAN_ERR",
				Path:    fieldName(f, tree),
				Message: fmt.Sprintf("The field <%s> must be greater than %v", fieldName(f, tree), threshold),
			}
		}
	}
	if geTag != "" {
		if threshold, err := strconv.ParseFloat(geTag, 64); err == nil && !(value >= threshold) {
			return &Error{
				ErrType: "GREATER_EQUAL_ERR",
				Path:    fieldName(f, tree),
				Message: fmt.Sprintf("The field <%s> must be greater than or equal to %v", fieldName(f, tree), threshold),
			}
		}
	}
	if ltTag != "" {
		if threshold, err := strconv.ParseFloat(ltTag, 64); err == nil && !(value < threshold) {
			return &Error{
				ErrType: "LESS_THAN_ERR",
				Path:    fieldName(f, tree),
				Message: fmt.Sprintf("The field <%s> must be less than %v", fieldName(f, tree), threshold),
			}
		}
	}
	if leTag != "" {
		if threshold, err := strconv.ParseFloat(leTag, 64); err == nil && !(value <= threshold) {
			return &Error{
				ErrType: "LESS_EQUAL_ERR",
				Path:    fieldName(f, tree),
				Message: fmt.Sprintf("The field <%s> must be less than or equal to %v", fieldName(f, tree), threshold),
			}
		}
	}
	if multipleOfTag != "" {
		if base, err := strconv.ParseFloat(multipleOfTag, 64); err == nil && base != 0 && math.Mod(value, base) != 0 {
			return &Error{
				ErrType: "NOT_MULTIPLE_ERR",
				Path:    fieldName(f, tree),
				Message: fmt.Sprintf("The field <%s> must be a multiple of %v", fieldName(f, tree), base),
			}
		}
	}

	return nil
}
