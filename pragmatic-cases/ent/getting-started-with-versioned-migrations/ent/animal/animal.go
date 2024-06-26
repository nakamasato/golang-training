// Code generated by ent, DO NOT EDIT.

package animal

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the animal type in the database.
	Label = "animal"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldColor holds the string denoting the color field in the database.
	FieldColor = "color"
	// FieldGender holds the string denoting the gender field in the database.
	FieldGender = "gender"
	// Table holds the table name of the animal in the database.
	Table = "animals"
)

// Columns holds all SQL columns for animal fields.
var Columns = []string{
	FieldID,
	FieldColor,
	FieldGender,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Animal queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByColor orders the results by the color field.
func ByColor(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldColor, opts...).ToFunc()
}

// ByGender orders the results by the gender field.
func ByGender(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGender, opts...).ToFunc()
}
