package cassandra_table_managment

import (
	"errors"
	"strings"

	"github.com/gocql/gocql"
)

var (
	ErrTableDoesNotExist = errors.New("table does not exist")
)

type TableManagment struct {
	keyspace string
	session  *gocql.Session
	tables   []TableMetaData
}

func (t *TableManagment) SetSession(session *gocql.Session) {
	t.session = session
}

func (t *TableManagment) SetKeySpace(keySpace string) {
	t.keyspace = keySpace
}

func (t *TableManagment) AddTableMetaData(tableMeta TableMetaData) {
	t.tables = append(t.tables, tableMeta)
}

func (t *TableManagment) CheckExistanceOfTables() error {
	for tableMetaData_index := range t.tables {
		if !t.isTableMetaDataExist(&t.tables[tableMetaData_index]) {
			if err := t.createTableMetaData(&t.tables[tableMetaData_index]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TableManagment) DropOldColumnsOfTables() error {
	return t.dropOldColumnsOfTables(false)
}

func (t *TableManagment) DropOldColumnsOfTablesReturnOnFailure() error {
	return t.dropOldColumnsOfTables(true)
}

func (t *TableManagment) dropOldColumnsOfTables(returnOnFailure bool) error {
	for tableMetaData_index := range t.tables {
		if !t.isTableMetaDataExist(&t.tables[tableMetaData_index]) {
			return ErrTableDoesNotExist
		}

		currentColumns, err := t.collectTableMetaDataCurrentColumns(&t.tables[tableMetaData_index])

		if err != nil {
			return err
		}

		deleteQuery := t.tables[tableMetaData_index].BuildDeleteQuery(currentColumns)

		result := t.session.Query(deleteQuery).Exec()

		if result != nil && returnOnFailure {
			return result
		}

	}

	return nil
}

func (t *TableManagment) UpdateTables() error {
	return t.updateTable(false)
}

func (t *TableManagment) UpdateTablesReturnOnFailureOfDataType() error {
	return t.updateTable(true)
}

func (t *TableManagment) updateTable(returnOnFailure bool) error {
	for tableMetaData_index := range t.tables {
		if !t.isTableMetaDataExist(&t.tables[tableMetaData_index]) {
			return ErrTableDoesNotExist
		}

		currentColumns, err := t.collectTableMetaDataCurrentColumns(&t.tables[tableMetaData_index])

		if err != nil {
			return err
		}

		addColumnQuery, updateColumnsQuery := t.tables[tableMetaData_index].BuildUpdateQueryFrom(currentColumns)

		if addColumnQuery != "" {

			result := t.addColumnToTableMetaData(addColumnQuery)

			if result != nil {
				return result
			}
		}

		if len(updateColumnsQuery) != 0 {
			err = t.updateColumnsDataType(updateColumnsQuery, returnOnFailure)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TableManagment) isTableMetaDataExist(tableMeta *TableMetaData) bool {
	existanceQuery := tableMeta.BuildExistanceQuery(t.keyspace)

	scanner := t.session.Query(existanceQuery).Iter().Scanner()

	return scanner.Next()
}

func (t *TableManagment) createTableMetaData(tableMeta *TableMetaData) error {
	createQuery := tableMeta.BuildCreateTableQuery()

	return t.session.Query(createQuery).Exec()
}

func (t *TableManagment) addColumnToTableMetaData(addColumnQuery string) error {
	return t.session.Query(addColumnQuery).Exec()
}

func (t *TableManagment) updateColumnsDataType(updateColumnsQuery []string, returnOnFailure bool) error {

	for updateQuery_index := range updateColumnsQuery {
		result := t.session.Query(updateColumnsQuery[updateQuery_index]).Exec()

		if result != nil && returnOnFailure {
			return result
		}
	}

	return nil
}

func (t *TableManagment) collectTableMetaDataCurrentColumns(tableMeta *TableMetaData) (map[string]string, error) {

	sliceMap, err := t.session.Query(tableMeta.BuildCollectColumnsQuery(t.keyspace)).Iter().SliceMap()

	if err != nil {
		return nil, err
	}

	filedsTypes := make(map[string]string)
	for row := range sliceMap {
		retrRow := sliceMap[row]
		rowName := strings.ToLower(retrRow["column_name"].(string))
		filedsTypes[rowName] = retrRow["type"].(string)
	}

	return filedsTypes, nil
}
