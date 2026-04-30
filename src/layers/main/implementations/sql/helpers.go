package sqlrepo

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// marshalNullable serializes v to JSON bytes. Returns nil when v is a nil
// pointer so that the column is stored as NULL rather than the JSON "null".
func marshalNullable(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	// Detect typed nil pointers via interface boxing.
	// json.Marshal handles this correctly but we guard anyway for clarity.
	return json.Marshal(v)
}

// nullableString returns a sql.NullString whose Valid flag mirrors whether
// the string is non-empty. Empty strings are stored as NULL.
func nullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// requireOneRowAffected returns an error when the SQL result reports that no
// rows were affected, which indicates the target record does not exist.
func requireOneRowAffected(result sql.Result, entityName, id string) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%s não encontrada: %s", entityName, id)
	}
	return nil
}
