package database

import (
	"errors"
	"log"
	"migrator/config"
	"sort"
	"strings"
)

type Table struct {
	Schema   string
	Name     string
	Priority int
}

type Tables []*Table

func (t Tables) Len() int           { return len(t) }
func (t Tables) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tables) Less(i, j int) bool { return t[i].Priority > t[j].Priority }

func Analize() (xs *Tables, err error) {
	if dstDB == nil {
		err = errors.New("Database not connect yet")
		return
	}

	cfg := config.GetConfig()

	xs = new(Tables)
	*xs = make(Tables, 0)

	for _, v := range cfg.Tables {
		if len(strings.Split(v, ".")) == 1 {
			v = "public." + v
		}
		splt := strings.Split(v, ".")
		if err = analizeTableDepencies(splt[1], splt[0], xs, 1); err != nil {
			return
		}
	}

	sort.Sort(xs)
	for _, t := range *xs {
		log.Println(t)
	}

	return
}

func analizeTableDepencies(t, s string, m *Tables, n int) (err error) {
	log.Println("Analizing table", t, "from schema", s, "with priority", n)

	table := &Table{
		Name:     t,
		Schema:   s,
		Priority: n,
	}

	if i := checkIn(table, m); i > 0 {
		if (*m)[i].Priority < n {
			(*m)[i].Priority = n
		}
	}

	(*m) = append(*m, table)

	rows, err := srcDB.Query(`
		SELECT
			CCU.table_schema AS foreign_table_schema,
			CCU.table_name AS foreign_table_name
		FROM information_schema.table_constraints AS TC
		INNER JOIN information_schema.key_column_usage AS KCU ON TC.constraint_name = KCU.constraint_name AND TC.table_schema = KCU.table_schema
		INNER JOIN information_schema.constraint_column_usage AS CCU ON CCU.constraint_name = TC.constraint_name AND CCU.table_schema = TC.table_schema
		WHERE TC.constraint_type = 'FOREIGN KEY' AND TC.table_name=$1 AND TC.table_schema = $2;
	`, t, s)
	if err != nil {
		return
	}

	for rows.Next() {
		if err = rows.Scan(&s, &t); err != nil {
			return
		}

		analizeTableDepencies(t, s, m, n+1)
	}

	return
}

func checkIn(t *Table, m *Tables) int {
	for i, v := range *m {
		if v.Name == t.Name && v.Schema == t.Schema {
			return i
		}
	}

	return -1
}
