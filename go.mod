module demo1

go 1.17

require (
	github.com/go-sql-driver/mysql v1.6.0
	orm v0.0.0
)

replace gee => ./gee

replace orm => ./orm
