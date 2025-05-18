//go:generate go run github.com/yeisme/notevault/sql/gen
package main

// 如果你的 SQL 文件包含特定数据库语法，可能还需要导入对应的数据库驱动
import (
	"fmt"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

// SQL 文件路径
var (
	schema_path          = `../schema/`
	schema_file_name_set = []string{
		"001-file.sql",
	}
)

func main() {

	schema_file_set := make([]string, len(schema_file_name_set))
	for i, schema_file_name := range schema_file_name_set {
		schema_file_set[i] = schema_path + schema_file_name
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "../../pkg/storage/repository/dao", // data access object 目录
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// 使用 rawsql 驱动初始化一个 *gorm.DB 实例
	// rawsql.New 的 Config 中可以指定 SQL 字符串或 SQL 文件路径
	gormdb, err := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: schema_file_set,
	}), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("failed to open rawsql db: %w", err))
	}

	g.UseDB(gormdb)

	g.ApplyBasic(
		g.GenerateAllTable()...,
	)

	g.Execute()
}
