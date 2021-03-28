package cassandra_table_managment

import (
	"strings"
)

var CassandraAllowedSimpleDataType = []string{
	"ascii",
	"bigint",
	"blob",
	"boolean",
	"counter",
	"date",
	"decimal",
	"double",
	"duration",
	"float",
	"inet",
	"int",
	"smallint",
	"text",
	"time",
	"timestamp",
	"timeuuid",
	"tinyint",
	"uuid",
	"varchar",
	"varint",
}

type TableMetaData struct {
	Name            string
	ColumnsAndTypes map[string]string
	PartKey         []string
	SortKey         []string
}

func (t *TableMetaData) BuildCreateTableQuery() string {

	var query strings.Builder

	query.WriteString("CREATE TABLE ")

	query.WriteString(t.Name)

	query.WriteString("\n")

	var b strings.Builder

	if len(t.PartKey) > 2 {
		b.WriteString("(")
		for index_partKey := range t.PartKey {
			b.WriteString(t.PartKey[index_partKey])
			if !(index_partKey == len(t.PartKey)-1) {
				b.WriteString(", ")
			}
		}
		b.WriteString(")")
	} else {
		b.WriteString(t.PartKey[0])
	}

	if len(t.SortKey) >= 1 {
		b.WriteString(", ")
		for index_sortKey := range t.SortKey {
			b.WriteString(t.SortKey[index_sortKey])
			if !(index_sortKey == len(t.PartKey)-1) {
				b.WriteString(", ")
			}
		}
	}

	for column, columnType := range t.ColumnsAndTypes {
		query.WriteString(column)
		query.WriteString(" ")
		query.WriteString(columnType)

		query.WriteString(",\n")
	}

	query.WriteString("PRIMARY KEY(")
	query.WriteString(b.String())
	query.WriteString(")")

	return query.String()
}

func (t *TableMetaData) BuildUpdateQueryFrom(oldColumnAndTypes map[string]string) (string, []string) {

	var addColumnQuery strings.Builder

	addColumnQuery.WriteString("ALTER TABLE ")

	addColumnQuery.WriteString(t.Name)

	addColumnQuery.WriteString(" ADD (")

	hasColumnUpdate := false

	updateColumnTypes := []string{}

	for columnName, desireType := range t.ColumnsAndTypes {

		extractedType, ok := oldColumnAndTypes[columnName]

		if !ok {
			hasColumnUpdate = true
			addColumnQuery.WriteString(columnName)
			addColumnQuery.WriteString(" ")
			addColumnQuery.WriteString(desireType)
			addColumnQuery.WriteString(", ")
			continue
		}

		if extractedType != desireType {
			var updateColumm strings.Builder
			updateColumm.WriteString("ALTER TABLE ")
			updateColumm.WriteString(t.Name)
			updateColumm.WriteString(" ALTER ")
			updateColumm.WriteString(columnName)
			updateColumm.WriteString(" TYPE ")
			updateColumm.WriteString(desireType)
			updateColumm.WriteString(";")

			updateColumnTypes = append(updateColumnTypes, updateColumm.String())
		}

	}

	addColumnQueryString := addColumnQuery.String()

	addColumnQueryString = addColumnQueryString[:len(addColumnQueryString)-2] //erase ', '

	addColumnQueryString += ")"

	if !hasColumnUpdate {
		return "", updateColumnTypes
	}

	return addColumnQueryString, updateColumnTypes
}

func (t *TableMetaData) BuildDeleteQuery(oldColumnAndTypes map[string]string) string {

	var deleteQuery strings.Builder

	deleteQuery.WriteString("ALTER TABLE ")
	deleteQuery.WriteString(t.Name)
	deleteQuery.WriteString(" DROP ")

	hasDropColumn := false

	for columnName := range oldColumnAndTypes {

		_, ok := t.ColumnsAndTypes[columnName]

		if !ok {
			hasDropColumn = true
			deleteQuery.WriteString(columnName)
			deleteQuery.WriteString(", ")
		}
	}

	deleteQuery.WriteString(";")

	if !hasDropColumn {
		return ""
	}

	deleteQueryString := deleteQuery.String()

	return deleteQueryString[:len(deleteQueryString)-2] //erase ', '
}
