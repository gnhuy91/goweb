package main

const (
	DBDriver         = dbDriver
	MigrationsDir    = migrationsDir
	DBConnRetryCount = dbConnRetryCount
)

var (
	Port = configPort()
	DSN  = configDSN()
)
