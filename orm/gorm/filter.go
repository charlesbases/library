package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type columns map[string]interface{}

type filters struct {
	eq      columns // a = b
	ne      columns // a != b
	gt      columns // a > b
	lt      columns // a < b
	bt      columns // a BETWEEN b AND c
	like    columns // a LIKE b
	in      columns // a IN (b, c)
	notin   columns // a NOT IN (b, c)
	null    columns // a IS NULL
	notnull columns // a IS NOT NULL
}

type Filter func(f *filters)

func (c *columns) bind(col string, arg interface{}) {
	if c.len() == 0 {
		*c = make(map[string]interface{}, 0)
	}
	(*c)[col] = arg
}

func (c *columns) len() int {
	return len(*c)
}

// build .
func (fs *filters) build() *[]clause.Expression {
	var expres = make([]clause.Expression, 0)

	if fs.eq.len() != 0 {
		for col, arg := range fs.eq {
			expres = append(
				expres,
				clause.Eq{
					Column: col,
					Value:  arg,
				},
			)
		}
	}

	if fs.ne.len() != 0 {
		for col, arg := range fs.ne {
			expres = append(
				expres,
				clause.Neq{
					Column: col,
					Value:  arg,
				})
		}
	}

	if fs.gt.len() != 0 {

	}

	if fs.lt.len() != 0 {

	}

	if fs.bt.len() != 0 {
		for col, arg := range fs.bt {
			switch arg.(type) {
			case []interface{}:
				if args := arg.([]interface{}); len(args) == 2 {
					expres = append(expres,
						clause.Gte{
							Column: col,
							Value:  args[0],
						},
						clause.Lte{
							Column: col,
							Value:  args[1],
						})
				}
			}
		}
	}

	return &expres
}

// Eq col = arg
func Eq(col string, arg interface{}) Filter {
	return func(f *filters) {
		f.eq.bind(col, arg)
	}
}

// Ne col != arg
func Ne(col string, arg interface{}) Filter {
	return func(f *filters) {
		f.ne.bind(col, arg)
	}
}

// Bt col BETWEEN a AND b
func Bt(col string, a, b interface{}) Filter {
	return func(f *filters) {
		f.bt.bind(col, []interface{}{a, b})
	}
}

// Bind .
func Bind(tableName string, args ...Filter) *gorm.DB {
	var filters = new(filters)
	for _, f := range args {
		f(filters)
	}

	return db.Table(tableName).Where(clause.Where{*filters.build()})
}
