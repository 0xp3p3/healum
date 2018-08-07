package mysql

var (
	mysqlSchema = `CREATE TABLE IF NOT EXISTS %[1]s (
id varchar(256) primary key,
created integer,
updated integer,
name varchar(255),
parameter1 varchar(255),
parameter2 varchar(10000),
parameter3 varchar(255),
metadata blob,
index(created),
index(updated),
index(name),
index(parameter1));
`

	migrations = []string{
		`ALTER TABLE %[1]s ADD lat real;`,
		`ALTER TABLE %[1]s ADD lng real;`,
		`ALTER TABLE %[1]s ADD SPATIAL INDEX lat, lng;`,
	}

)
