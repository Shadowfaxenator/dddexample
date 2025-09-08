package ddd

import (
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
	_ "github.com/tursodatabase/go-libsql"
)

func ParseToSql(t reflect.Type) string {
	sts := &StructToSql{
		Tables: make(map[string][]struct {
			Name string
			Type string
		})}

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct type, got %s", t.Kind()))
	}
	sts.parse(t, t.Name(), "")
	return fmt.Sprintf("%+v", sts.Tables)

}

type StructToSql struct {
	Tables map[string][]struct {
		Name string
		Type string
	}
}

func (s *StructToSql) parse(t reflect.Type, tname, fname string) {
	tname = strcase.ToSnake(tname)
	fname = strcase.ToSnake(fname)
	switch t.Kind() {
	case reflect.Slice:
		s.parse(t.Elem(), tname, t.Elem().Name())
	case reflect.String:
		s.Tables[tname] = append(s.Tables[tname], struct {
			Name string
			Type string
		}{
			Name: fname,
			Type: "TEXT",
		})
	case reflect.Array:
		s.Tables[tname] = append(s.Tables[tname], struct {
			Name string
			Type string
		}{
			Name: fname,
			Type: "BLOB",
		})
	case reflect.Struct:
		str_name := tname
		for i := range t.NumField() {
			if t.Field(i).Type.Kind() == reflect.Slice {
				tname = fmt.Sprintf("%s_%s", str_name, t.Field(i).Name)
			}
			s.parse(t.Field(i).Type, tname, t.Field(i).Name)
		}
	default:
		panic(fmt.Sprintf("unsupported field type: %s", t.Kind()))
	}

}

func db() {
	// path := filepath.Join(".", "projections", a.boundedCtx)
	// if err := os.MkdirAll(path, os.ModePerm); err != nil {
	// 	panic(err)
	// }
	// dbpath := fmt.Sprintf("file:%s/%s.db", path, strings.ToLower(a.atype))
	// db, err := sql.Open("libsql", dbpath)
	// if err != nil {
	// 	panic(err)
	// }
	// a.proj = db
	// db.ExecContext(ctx, "PRAGMA journal_mode=WAL")
}
