package helper

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    int = 1
		args []interface{}
	)

	for k, v := range params {
		if k != "" {
			namedQuery = strings.ReplaceAll(namedQuery, ":"+k, "$"+strconv.Itoa(i))

			args = append(args, v)
			i++
		}
	}

	return namedQuery, args
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

const otpChars = "1234567890"

func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

func Difference(a, b []int32) []int32 {
	mb := make(map[int32]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []int32
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func ValMultipleQuery(query string, vals []int32) (string, []interface{}) {
	params := []interface{}{}

	for i, id := range vals {
		query += fmt.Sprintf("$%d,", i+1)
		params = append(params, id)
	}

	query = query[:len(query)-1] // remove trailing ","
	query += ")"

	return query, params
}

func InsertMultiple(queryInsert string, id string, vals []string) (string, []interface{}) {
	insertparams := []interface{}{}

	for i, d := range vals {
		p1 := i * 2 // starting position for insert params
		queryInsert += fmt.Sprintf("($%d, $%d),", p1+1, p1+2)
		insertparams = append(insertparams, id, d)
	}

	queryInsert = queryInsert[:len(queryInsert)-1] // remove trailing ","

	return queryInsert, insertparams
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullBool(s bool) sql.NullBool {
	if !s {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  s,
		Valid: true,
	}
}

// ConvertMapToStruct converts a map[string]interface{} to the specified struct type.
func ConvertMapToStruct(data map[string]interface{}, output interface{}) error {
	if len(data) == 0 {
		return errors.New("empty map provided")
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Convert JSON to struct
	if err := json.Unmarshal(jsonData, output); err != nil {
		return err
	}

	return nil
}

// StructAsMap converts a struct into map[string]interface{}.
func StructAsMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	value := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typ.Field(i).Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = typ.Field(i).Name
		}

		if field.IsValid() {
			result[fieldName] = field.Interface()
		}
	}

	return result
}
func SortAndJoinTags(tags []string) string {
	// Sort tags alphabetically
	sort.Strings(tags)
	// Join with commas (e.g., "golang,backend,api")
	return strings.Join(tags, ",")
}
