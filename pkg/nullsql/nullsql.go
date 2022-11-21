package nullsql

import "gopkg.in/guregu/null.v4"

// ValidateIntSQLValid returns a null.Int type, but distinguish a 0 as null value
func ValidateIntSQLValid(id int64) null.Int {
	if id == 0 {
		return null.NewInt(0, false)
	}
	return null.IntFrom(id)
}

func ValidateStringSQLValid(value string) null.String {
	if value == "" {
		return null.NewString(value, false)
	}
	return null.NewString(value, true)
}
