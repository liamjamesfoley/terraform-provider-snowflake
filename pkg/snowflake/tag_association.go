package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/validation"
	"github.com/jmoiron/sqlx"
)

// TagAssociationBuilder abstracts the creation of SQL queries for a Snowflake tag.
type TagAssociationBuilder struct {
	databaseName     string
	objectIdentifier string
	objectType       string
	schemaName       string
	tagName          string
	tagValue         string
}

type tagAssociation struct {
	TagValue sql.NullString `db:"TAG_VALUE"`
}

// WithObjectIdentifier adds the name of the schema to the TagAssociationBuilder.
func (tb *TagAssociationBuilder) WithObjectIdentifier(objectIdentifier string) *TagAssociationBuilder {
	tb.objectIdentifier = objectIdentifier
	return tb
}

// WithObjectType adds the object type of the resource to add tag attachement to the TagAssociationBuilder.
func (tb *TagAssociationBuilder) WithObjectType(objectType string) *TagAssociationBuilder {
	tb.objectType = objectType
	return tb
}

// WithTagValue adds the name of the tag value to the TagAssociationBuilder.
func (tb *TagAssociationBuilder) WithTagValue(tagValue string) *TagAssociationBuilder {
	tb.tagValue = tagValue
	return tb
}

// GetTagDatabase returns the value of the tag database of TagAssociationBuilder.
func (tb *TagAssociationBuilder) GetTagDatabase() string {
	return tb.databaseName
}

// GetTagName returns the value of the tag name of TagAssociationBuilder.
func (tb *TagAssociationBuilder) GetTagName() string {
	return tb.schemaName
}

// GetTagSchema returns the value of the tag schema of TagAssociationBuilder.
func (tb *TagAssociationBuilder) GetTagSchema() string {
	return tb.schemaName
}

// TagAssociation returns a pointer to a Builder that abstracts the DDL operations for a tag sssociation.
//
// Supported DDL operations are:
//   - ALTER <object_type> SET TAG
//   - ALTER <object_type> UNSET TAG
//   - SYSTEM$GET_TAG (get current tag value)
//
// [Snowflake Reference](https://docs.snowflake.com/en/user-guide/object-tagging.html)
func TagAssociation(tagID string) *TagAssociationBuilder {
	databaseName, schemaName, tagName := validation.ParseFullyQualifiedObjectID(tagID)
	return &TagAssociationBuilder{
		databaseName: databaseName,
		schemaName:   schemaName,
		tagName:      tagName,
	}
}

// Create returns the SQL query that will set the tag on an object.
func (tb *TagAssociationBuilder) Create() string {
	return fmt.Sprintf(`ALTER %v %v SET TAG "%v"."%v"."%v" = '%v'`, tb.objectType, tb.objectIdentifier, tb.databaseName, tb.schemaName, tb.tagName, EscapeString(tb.tagValue))
}

// Drop returns the SQL query that will remove a tag from an object.
func (tb *TagAssociationBuilder) Drop() string {
	return fmt.Sprintf(`ALTER %v %v UNSET TAG "%v"."%v"."%v"`, tb.objectType, tb.objectIdentifier, tb.databaseName, tb.schemaName, tb.tagName)
}

// Show returns the SQL query that will show the current tag value on an object.
func (tb *TagAssociationBuilder) Show() string {
	return fmt.Sprintf(`SELECT SYSTEM$GET_TAG('"%v"."%v"."%v"', '%v', '%v') TAG_VALUE WHERE TAG_VALUE IS NOT NULL`, tb.databaseName, tb.schemaName, tb.tagName, tb.objectIdentifier, tb.objectType)
}

func ScanTagAssociation(row *sqlx.Row) (*tagAssociation, error) {
	r := &tagAssociation{}
	err := row.StructScan(r)
	return r, err
}

func ListTagAssociations(tb *TagAssociationBuilder, db *sql.DB) ([]tagAssociation, error) {
	stmt := `SELECT SYSTEM$GET_TAG('"?"."?"."?"', '?', '?') TAG_VALUE WHERE TAG_VALUE IS NOT NULL`
	rows, err := db.Query(stmt,
		tb.databaseName, tb.schemaName, tb.tagName, tb.objectIdentifier, tb.objectType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tagAssociations := []tagAssociation{}
	log.Printf("[DEBUG] tagAssociations is %v", tagAssociations)
	if err := sqlx.StructScan(rows, &tagAssociations); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DEBUG] no tag associations found for tag %s", tb.tagName)
			return nil, err
		}
		return nil, fmt.Errorf("unable to scan row for %s err = %w", stmt, err)
	}

	return tagAssociations, nil
}