package model

import (
	"database/sql"
	"fmt"
	"log"
)

type Table struct {
	TableName  string `json:"table_name"`
	LocalKey   string `json:"local_key"`
	ForeignKey string `json:"foreign_key"`
}
type TableList struct {
	RootTable  string  `json:"root_table"`
	TableArray []Table `json:"table_array"`
}

type Entry struct {
	Date string `json:"date"`
	Day  string `json:"day"`
	Task string `json:"task"`
}
type Pair struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func GetQuery(v TableList) string {
	query := "SELECT * FROM " + v.RootTable
	leftTable := v.RootTable
	foreignKey := leftTable + "." + v.TableArray[0].ForeignKey

	for i := 0; i < len(v.TableArray); i++ {
		if i > 0 {
			leftTable = v.TableArray[i-1].TableName
			foreignKey = leftTable + "." + v.TableArray[i-1].ForeignKey
		}
		rightTable := v.TableArray[i].TableName
		localKey := rightTable + "." + v.TableArray[i].LocalKey

		query = query + " INNER JOIN " + rightTable + " ON " + localKey + " = " + foreignKey
		fmt.Println(query)
	}
	return query
}

func JoinTable(query string, db *sql.DB) [][]Pair {

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	var results [][]Pair
	for rows.Next() {
		rowValues := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range rowValues {
			valuePtrs[i] = &rowValues[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}
		var rowData []Pair
		for i, col := range columns {
			ele := string(rowValues[i].([]byte))
			rowData = append(rowData, Pair{Field: col, Value: ele})
		}

		results = append(results, rowData)
	}
	return results
}
