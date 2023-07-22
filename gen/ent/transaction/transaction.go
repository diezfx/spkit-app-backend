// Code generated by ent, DO NOT EDIT.

package transaction

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the transaction type in the database.
	Label = "transaction"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldAmount holds the string denoting the amount field in the database.
	FieldAmount = "amount"
	// FieldSourceID holds the string denoting the source_id field in the database.
	FieldSourceID = "source_id"
	// FieldTransactionType holds the string denoting the transaction_type field in the database.
	FieldTransactionType = "transaction_type"
	// FieldTargetIds holds the string denoting the target_ids field in the database.
	FieldTargetIds = "target_ids"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// Table holds the table name of the transaction in the database.
	Table = "transactions"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "projects"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "transaction_project"
)

// Columns holds all SQL columns for transaction fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldAmount,
	FieldSourceID,
	FieldTransactionType,
	FieldTargetIds,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "transactions"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"project_transactions",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// TransactionType defines the type for the "transaction_type" enum field.
type TransactionType string

// TransactionType values.
const (
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

func (tt TransactionType) String() string {
	return string(tt)
}

// TransactionTypeValidator is a validator for the "transaction_type" field enum values. It is called by the builders before save.
func TransactionTypeValidator(tt TransactionType) error {
	switch tt {
	case TransactionTypeExpense, TransactionTypeTransfer:
		return nil
	default:
		return fmt.Errorf("transaction: invalid enum value for transaction_type field: %q", tt)
	}
}

// OrderOption defines the ordering options for the Transaction queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByAmount orders the results by the amount field.
func ByAmount(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAmount, opts...).ToFunc()
}

// BySourceID orders the results by the source_id field.
func BySourceID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSourceID, opts...).ToFunc()
}

// ByTransactionType orders the results by the transaction_type field.
func ByTransactionType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTransactionType, opts...).ToFunc()
}

// ByProjectCount orders the results by project count.
func ByProjectCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProjectStep(), opts...)
	}
}

// ByProject orders the results by project terms.
func ByProject(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProjectStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProjectStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProjectInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ProjectTable, ProjectColumn),
	)
}
