// Code generated by entc, DO NOT EDIT.

package user

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldChatID holds the string denoting the chat_id field in the database.
	FieldChatID = "chat_id"
	// EdgeSet holds the string denoting the set edge name in mutations.
	EdgeSet = "set"
	// Table holds the table name of the user in the database.
	Table = "users"
	// SetTable is the table that holds the set relation/edge.
	SetTable = "sets"
	// SetInverseTable is the table name for the Set entity.
	// It exists in this package in order to avoid circular dependency with the "set" package.
	SetInverseTable = "sets"
	// SetColumn is the table column denoting the set relation/edge.
	SetColumn = "user_set"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldChatID,
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