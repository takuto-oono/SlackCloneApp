package models

import "fmt"

func Create(x int) error {
	cmd := fmt.Sprintf("INSERT INTO %s (num, word) VALUES (?, ?)", "test_db1")
	_, err := DbConnection.Exec(cmd, x, "Hello Golang sqlite3")
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func CreateTestData() {
	for i := 0; i < 10; i ++ {
		Create(i)
	}
}

func GetTestData(x int) string {
	tableName := "test_db1"
	cmd := fmt.Sprintf("SELECT num, word FROM %s WHERE num = ?", tableName)
	row := DbConnection.QueryRow(cmd, x)
	fmt.Println(row)
	var num int
	var word string
	row.Scan(&num, &word)
	fmt.Println(num, word)
	return word
}