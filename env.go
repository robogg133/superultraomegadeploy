package main

import (
	"os"
	"strconv"
)

var (
	DebugEnabled bool

	PostgresConnString string
)

func init() {
	PostgresConnString = os.Getenv("POSTGRES_CONN_STRING")
}

func init() {
	var err error
	DebugEnabled, err = strconv.ParseBool(os.Getenv("DEBUG_ENABLED"))
	if err != nil {
		panic(err)
	}
}
