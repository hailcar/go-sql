package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type Field struct {
	Name    string
	Type    string
	Comment string
}
type Table struct {
	Name   string
	Fields []Field
}

func toClassName(tableName string) string {
	names := strings.Split(tableName, "_")
	output := ""
	for _, v := range names {
		output += strings.Title(v)
	}
	return output
}
func ParseSql(r io.Reader) []Table {
	scanner := bufio.NewScanner(r)
	startString := regexp.MustCompile("CREATE.*`(.*)`")
	endString := regexp.MustCompile("^\\).*")
	fieldString := regexp.MustCompile("(?U)^\\s*`(.*)` (.*) .*(COMMENT\\s+'(.*)')?,$")
	tableFlag := 1
	tables := make([]Table, 1)
	mapper := map[string]string{
		"varchar":  "string",
		"int":      "int",
		"tinyint":  "int",
		"smallint": "int",
		"float":    "float",
		"double":   "double",
		"decimal":  "float",
		"char":     "string",
		"text":     "string",
		"enum":     "string",
	}
	tableName := ""
	var currentTable Table
	for scanner.Scan() {
		str := scanner.Text()
		switch tableFlag {
		case 1:
			ggroup := startString.FindAllStringSubmatch(str, -1)
			if len(ggroup) > 0 {
				sub := ggroup[0]
				if len(sub) > 0 {
					tableName = sub[1]
					currentTable = Table{
						Name:   tableName,
						Fields: make([]Field, 1),
					}
					fmt.Printf(`
					class %s {
					`, toClassName(tableName))
					tableFlag = 2
				}
			}
			break
		case 2:
			gp := endString.FindAllString(str, -1)
			if len(gp) > 0 {
				tableFlag = 1
				fmt.Printf(`
			}
				`)
				tables = append(tables, currentTable)
				break
			}
			ggroup := fieldString.FindAllStringSubmatch(str, -1)
			if len(ggroup) > 0 {
				sub := ggroup[0]
				fieldName := sub[1]
				typeArray := strings.Split(sub[2], "(")
				fieldType := mapper[typeArray[0]]
				comment := sub[4]

				if len(sub) > 3 {
					fmt.Printf(`
					/**
					* @var %s %s
					*/
					`, fieldType, comment)
				}
				fmt.Printf(`public %s $%s;`, fieldType, fieldName)
				field := Field{
					Name:    fieldName,
					Type:    fieldType,
					Comment: comment,
				}
				currentTable.Fields = append(currentTable.Fields, field)
			}
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return tables
}

func Sql2Dao(tables []Table) {
	for _, v := range tables {
		if len(v.Fields) == 0 {
			continue
		}
		fmt.Printf(`
		class %s {
		`, toClassName(v.Name))
		for _, field := range v.Fields {
			if field.Name == "" {
				continue
			}
			fmt.Printf(`
					/**
					* @var %s %s
					*/
					`, field.Type, field.Comment)
			fmt.Printf(`public %s $%s;`, field.Type, field.Name)
		}
		fmt.Printf(`
			}
				`)
	}
}

func Sql2Validator(tables []Table) {

	mapper := map[string]string{
		"varchar":  "string",
		"int":      "number",
		"tinyint":  "number",
		"smallint": "number",
		"float":    "float",
		"double":   "float",
		"decimal":  "float",
		"char":     "string",
		"text":     "string",
	}

	for _, table := range tables {
		if len(table.Fields) == 0 {
			continue
		}
		fmt.Printf(`
	use think\Validate;

	class %s extends Validate
	{
	protected $rule=[
	`, toClassName(table.Name))

		for _, field := range table.Fields {
			if field.Name == "" {
				continue
			}
			fmt.Printf(`
			"%s"=>[ "require", "%s"],`, field.Name, mapper[field.Type])
		}

		fmt.Printf(`
	];
	`)
		fmt.Printf(`
	protected $message=[
	`)

		for _, field := range table.Fields {
			if field.Name == "" {
				continue
			}
			fmt.Printf(`
			"%s.require"=>"%s必传",`, field.Name, field.Comment)
		}
		fmt.Println("\n")
		for _, field := range table.Fields {
			if field.Name == "" {
				continue
			}
			fmt.Printf(`
			"%s.%s"=>"%s格式错误",`, field.Name, mapper[field.Type], field.Comment)
		}

		fmt.Printf(`
	];
	`)
	}

}

func Sql2Editor(tables []Table) {

	for _, table := range tables {
		if len(table.Fields) == 0 {
			continue
		}
		objectName := strcase.ToLowerCamel(table.Name)
		objectNameUpperCase := strcase.ToCamel(table.Name)
		fmt.Printf("$%s=new %s()\n", objectName, objectNameUpperCase)
		for _, field := range table.Fields {
			if field.Name == "" {
				continue
			}
			fmt.Printf(`
				$%[1]s->%[2]s=$request->%[2]s;
			`, objectName, field.Name)
		}
		fmt.Print("---------------\n\n")
		fmt.Printf(`
		%s 
		`, objectName)
		fmt.Printf(`
		$%[1]s= $%[1]sManager->findById($request->id);
		if($%[1]s==null){
			throw new NotFoundException("指定的商品未找到");
		}
		`, objectName)
		for _, field := range table.Fields {
			if field.Name == "" {
				continue
			}
			fmt.Printf(`
			if(isset($%[1]s->%[2]s) and $%[1]s->%[2]s!=$request->%[2]s){
				$%[1]s->%[2]s=$request->%[2]s;
			}
			`, objectName, field.Name)
		}
	}
}
