package database

import (
	"database/sql"
	"errors"
	"log"
	"migrator/config"
	"reflect"
	"strconv"
	"strings"
)

func Migrate(tables *Tables) (err error) {
	cfg := config.GetConfig()
	limit := 15

	if srcDB == nil || dstDB == nil {
		err = errors.New("Database is not connected")
	}

	if cfg.RecordLimit > 0 {
		limit = cfg.RecordLimit
	}

	for i := 0; i < len(*tables); i++ {
		log.Println((*tables)[i])

		rows, err := srcDB.Query(`
			SELECT *
			FROM `+(*tables)[i].Schema+"."+(*tables)[i].Name+`
			LIMIT $1
		`, limit)
		if err != nil {
			return err
		}

		types, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		cols, err := rows.Columns()
		if err != nil {
			return err
		}

		for rows.Next() {
			vals := getVals(types)
			err := rows.Scan(vals...)
			if err != nil {
				return err
			}

			var id int64
			if err = dstDB.QueryRow(`
				INSERT INTO `+(*tables)[i].Schema+"."+(*tables)[i].Name+`
				(`+strings.Join(cols, ",")+`)
				VALUES (`+generatePlaceholder(len(cols))+`)
				RETURNING id
			`, vals...).Scan(&id); err != nil {
				return err
			}
		}
	}

	return
}

func getVals(t []*sql.ColumnType) []interface{} {
	r := make([]interface{}, len(t))

	for i := range t {
		tmp := reflect.New(reflect.PtrTo(t[i].ScanType()))
		tmp.Elem().Set(reflect.New(t[i].ScanType()))
		r[i] = tmp.Elem().Interface()
	}

	return r
}

func generatePlaceholder(i int) (p string) {
	xs := make([]string, i)

	for j := 0; j < i; j++ {
		xs[j] = "$" + strconv.Itoa(j+1)
	}

	return strings.Join(xs, ",")
}
