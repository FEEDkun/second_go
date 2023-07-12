package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Go MySQL Tutorial")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/sql5")
	fmt.Println("连接成功")
	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("异常")
	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// perform a db.Query insert
	insert, err := db.Query("INSERT INTO excel01 VALUES ( '050', '第二列','第三列','第四列' )")
	fmt.Println("插入成功")

	// 查询数据
	rows, err := db.Query("SELECT * FROM excel01 where item1='050'")
	//checkErr(err)
	for rows.Next() {
		var item1 string
		var item2 string
		var item3 string
		var item4 string
		err = rows.Scan(&item1, &item2, &item3, &item4)
		//checkErr(err)
		fmt.Println(item1, item2, item3, item4)
		fmt.Println("查询成功")
	}

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	//defer Select.Close()
}
