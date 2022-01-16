package orm

import "database/sql"

type Orm struct {
	db      *sql.DB
	dialect Dialect
}

func NewOrm(driver, source string) (orm *Orm, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		Error(err)
		return
	}
	d, ok := GetDialect(driver)
	if !ok {
		Error("diver did not exist in dialectmap")
		return
	}
	orm = &Orm{db: db, dialect: d}
	Info("database connect success")
	return
}

func (o *Orm) Close() {
	if err := o.db.Close(); err != nil {
		Error("database close error")
		return
	}
	Info("database close success")
}

func (o *Orm) NewSession() *Session {
	return NewSession(o.db, o.dialect)
}
