package database

type FakeRecord interface {
	Fakeable
	Insertable
}

type Insertable interface {
	GenerateInsertSql() string
}

type Fakeable interface {
	NewFakeRecord() (Insertable, error)
}
