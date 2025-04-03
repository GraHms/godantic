// Package godantic is written by Ismael GraHms
package godantic

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func getValueOf(val interface{}) reflect.Value {
	return reflect.ValueOf(val)
}
func isPtr(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr
}
func isStruct(value reflect.Value) bool {
	kind := value.Kind()
	return kind == reflect.Struct || (kind == reflect.Ptr && value.Elem().Kind() == reflect.Struct)
}

func isString(value reflect.Value) bool {
	return value.Kind() == reflect.String
}
func isTime(value reflect.Value) bool {
	return value.Type().ConvertibleTo(TimeType)
}
func isList(value reflect.Value) bool {
	return value.Kind() == reflect.Slice || value.Kind() == reflect.Array
}

var TimeType = reflect.TypeOf(time.Time{})

func (g *Validate) InspectStruct(val interface{}) error {
	enumMap := extractEnumValues(getValueOf(val), "")
	return g.inspect(val, "", 0, reflect.StructField{}, enumMap)
}

func (g *Validate) inspect(val interface{}, tree string, i int, f reflect.StructField, enumMap map[string]string) error {

	v := getValueOf(val)
	if _, ok := v.Interface().(*time.Time); ok {
		return nil
	}
	switch {
	case isPtr(v):
		return g.inspect(v.Elem().Interface(), tree, i, f, enumMap)
	case isStruct(v):
		return g.checkStruct(val, v, tree, enumMap)
	case isString(v):
		return g.checkString(v, tree, i, f)
	case isTime(v):
		return g.checkTime(v, tree)
	case isList(v):
		return g.checkList(v, tree, enumMap)
	default:
		return nil

	}
}

func getFormatRegex(formatTag string) string {
	switch formatTag {
	case "email":
		return `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	case "url":
		return `^(http|https)://[a-zA-Z0-9./?=_%-&]+$`
	case "date":
		return `^\d{4}-\d{2}-\d{2}$`
	case "time":
		return `^([01]\d|2[0-3]):([0-5]\d):([0-5]\d)$`
	case "uuid":
		return `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	case "ip":
		return `^(\d{1,3}\.){3}\d{1,3}$`
	case "credit_card":
		return `^\d{4}-\d{4}-\d{4}-\d{4}$`
	case "postal_code":
		return `^[a-zA-Z0-9]+$`
	case "phone":
		return `^\+[1-9]\d{1,14}$`
	case "ssn":
		return `^\d{3}-\d{2}-\d{4}$`
	case "credit_card_expiry":
		return `^(0[1-9]|1[0-2])\/(20\d{2}|2[1-9]\d{1})$`
	case "latitude":
		return `^(-?([0-8]?[0-9]\.\d+|90(\.0+)?))$`
	case "longitude":
		return `^(-?((1?[0-7]?|[0-9]?)[0-9]\.\d+|180(\.0+)?))$`
	case "hex_color":
		return `^#?([a-fA-F0-9]{6}|[a-fA-F0-9]{3})$`
	case "mac_address":
		return `^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`
	case "mz-msisdn":
		return `^258\d{9}$`
	case "mz-nuit":
		return `^\d{9}$`
	case "mz-bi":
		return `^\d{12}[A-Z]$`
	default:
		return ""
	}
}

func (g *Validate) formatValidation(f reflect.StructField, v reflect.Value, tree string) error {
	tag := f.Tag.Get("format")
	if tag == "" {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	fieldValue := v.String()
	formatRegex := getFormatRegex(tag)
	if matchRegexPattern(formatRegex, fieldValue, f, tree) != nil {
		return &Error{
			ErrType: fmt.Sprintf("INVALID_%s_ERR", strings.ToUpper(tag)),
			Path:    fieldName(f, tree),
			Message: fmt.Sprintf("error on field <%s>. the given value '%s' is not a valid %s", fieldName(f, tree), fieldValue, tag),
		}
	}

	return matchRegexPattern(formatRegex, fieldValue, f, tree)
}

func matchRegexPattern(regexpPattern, fieldValue string, f reflect.StructField, tree string) error {
	// Compile the regular expression pattern
	regexpPatternCompiled, err := regexp.Compile(regexpPattern)
	if err != nil {
		return err // Handle error
	}
	// Check if the field's value matches the regular expression pattern
	if !regexpPatternCompiled.MatchString(fieldValue) {
		return &Error{
			ErrType: "INVALID_PATTERN_ERR",
			Path:    fieldName(f, tree),
			Message: fmt.Sprintf("The field <%s> value '%s' does not match the required pattern: %s", fieldName(f, tree), fieldValue, regexpPattern),
		}
	}
	return nil
}

func (g *Validate) regexPattern(f reflect.StructField, v reflect.Value, tree string) error {
	regexpPattern := f.Tag.Get("regex")
	if regexpPattern == "" {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	fieldValue := v.String()

	// Compile the regular expression pattern
	return matchRegexPattern(regexpPattern, fieldValue, f, tree)

}

func (g *Validate) checkTime(v reflect.Value, tree string) error {
	timeValue := v.Interface().(time.Time)
	if timeValue.IsZero() {
		return &Error{
			ErrType: "INVALID_TIME_ERR",
			Path:    tree,
			Message: "The field <" + tree + "> cannot have an invalid time value",
		}
	}
	return nil
}

func (g *Validate) checkString(v reflect.Value, tree string, _ int, f reflect.StructField) error {
	passTag := f.Tag.Get("pass-empty")
	if passTag == "true" {
		return nil
	}
	s := strings.TrimSpace(v.String())
	if len(s) < 1 {
		return &Error{
			ErrType: "EMPTY_STRING_ERR",
			Path:    tree,
			Message: "The field <" + tree + "> cannot be an empty string",
		}
	}
	return nil
}

func (g *Validate) checkList(v reflect.Value, tree string, enumMap map[string]string) error {
	min := 1
	if g.IgnoreMinLen == true {
		min = 0
	}
	isLesserThanMinLength := v.Len() < min
	if isLesserThanMinLength {

		return &Error{
			ErrType: "EMPTY_LIST_ERR",
			Path:    tree,
			Message: fmt.Sprintf("Field <%s> must be with at least one value.", tree),
		}
	}
	for i := 0; i < v.Len(); i++ {
		err := g.inspect(v.Index(i).Interface(), tree, i, reflect.StructField{}, enumMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Validate) checkStruct(val interface{}, v reflect.Value, tree string, enumMao map[string]string) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if isTime(v.Field(i)) {
			// ignore time.Time fields, they are already checked in bindJSON
			continue
		}

		err := g.checkField(val, v, t, tree, i, enumMao)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Validate) checkField(val interface{}, v reflect.Value, t reflect.Type, tree string, i int, enumMap map[string]string) error {

	f := t.Field(i)
	if f.PkgPath != "" {
		// Field is unexported, handle it gracefully
		return nil
	}

	enums := f.Tag.Get("enum")
	if len(enums) == 0 {
		enums = f.Tag.Get("enums")
	}

	valField := v.Field(i)

	if tag := f.Tag.Get("binding"); tag == "ignore" && !reflect.DeepEqual(valField.Interface(), reflect.Zero(f.Type).Interface()) {
		return &Error{
			ErrType: "INVALID_FIELD_ERR",
			Path:    fieldName(f, tree),
			Message: fmt.Sprintf("Invalid field <%s>", fieldName(f, tree)),
		}
	}

	switch {

	case f.Type.Kind() == reflect.Ptr && !valField.IsNil():
		// Handle pointer fields
		if err := g.inspect(valField.Interface(), fieldName(f, tree), i, f, enumMap); err != nil {
			return err
		}
	case f.Type.Kind() == reflect.Struct:
		// Handle non-pointer struct fields
		if err := g.checkStruct(val, valField, fieldName(f, tree), enumMap); err != nil {
			return err
		}
	case f.Type.Kind() != reflect.Ptr:
		if !g.IgnoreRequired && isFieldRequired(f) && reflect.DeepEqual(valField.Interface(), reflect.Zero(f.Type).Interface()) {
			return RequiredFieldError(f, tree)
		}
	case !g.IgnoreRequired:
		if isFieldRequired(f) {
			if f.Type.Kind() == reflect.Ptr && valField.IsNil() {
				return RequiredFieldError(f, tree)
			}

		}

	case isPtr(v):
		if !v.IsValid() || v.IsNil() {
			return nil // nil pointer is valid
		}
		return g.inspect(v.Elem().Interface(), tree, i, v.Type().Field(i), enumMap)

	}
	if err := g.validateCondition(f, valField, tree, enumMap); err != nil {
		return err
	}
	if err := g.checkMinMax(f, valField, tree); err != nil {
		return err
	}
	if err := g.checkNumericConstraints(f, valField, tree); err != nil {
		return err
	}
	if err := g.checkDecimalConstraints(f, valField, tree); err != nil {
		return err
	}

	if err := g.regexPattern(f, valField, tree); err != nil {
		return err
	}
	if err := g.formatValidation(f, valField, tree); err != nil {
		return err
	}

	// Check for enum validation tags.
	if len(enums) > 0 {
		enumValues := strings.Split(strings.TrimSpace(enums), ",")
		err := g.strEnums(f, valField, tree, enumValues)

		if err != nil {
			return err
		}
	}
	if customValidator, ok := val.(ValidationPlugin); ok {
		if err := customValidator.Validate(); err != nil {
			return &Error{
				ErrType: err.ErrType,
				Message: err.Message,
				Path:    err.Path,
				err:     err,
			}
		}
	}
	if df, ok := val.(DynamicFieldsValidator); ok {
		if err := validateDynamicFields(df.GetValue(), df.GetAttribute(), df.GetValueType(), tree); err != nil {
			return err
		}
	}

	return nil
}

func (g *Validate) strEnums(f reflect.StructField, val reflect.Value, tree string, allowedValues []string) error {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	fieldValue := val.String()

	// Convert the allowedValues slice to a map for faster searching.
	allowedMap := make(map[string]bool, len(allowedValues))
	for _, allowedValue := range allowedValues {
		allowedMap[allowedValue] = true
	}

	if _, ok := allowedMap[fieldValue]; !ok {
		return &Error{
			ErrType: "INVALID_ENUM_ERR",
			Path:    fieldName(f, tree),
			Message: fmt.Sprintf("The field <%s> must have one of the following values: %s, '%s' was given",
				fieldName(f, tree), strings.Join(allowedValues, ", "), val.String()),
		}
	}

	return nil
}

func fieldName(f reflect.StructField, tree string) string {
	var name string
	if f.Tag != "" {
		jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
		name = jsonTag
		if len(tree) > 0 {
			name = tree + "." + jsonTag
		}
	}
	return name
}

func isFieldRequired(f reflect.StructField) bool {
	if f.Tag.Get("binding") == "required" {
		return true
	}
	return false
}

func RequiredFieldError(field reflect.StructField, tree string) error {
	return &Error{
		ErrType: "REQUIRED_FIELD_ERR",
		Path:    fieldName(field, tree),
		Message: fmt.Sprintf("The field <%s> is required", fieldName(field, tree)),
	}
}
