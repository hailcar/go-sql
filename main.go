package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	input    string
	sql      bool
	template bool
)

func main() {
	flag.BoolVar(&sql, "s", true, "sql to php dao")
	flag.BoolVar(&template, "t", false, "generate template")
	flag.Parse()
	inputs := flag.Args()
	if len(inputs) <= 0 {
		return
	}
	input := inputs[0]
	if template {
		GenerateTemplate(input)
	} else if sql {
		fmt.Println(input)
		file, err := os.Open(input)
		if err != nil {
			log.Fatal(err)
			return
		}
		result := ParseSql(file)
		Sql2Dao(result)
		Sql2Validator(result)
		Sql2Editor(result)
	}
}
