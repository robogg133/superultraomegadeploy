package main

import (
	"os"
	"strconv"
)

var (
	DebugEnabled bool

	PostgresConnString string
	JWTSecret          string
)

func init() {
	PostgresConnString = os.Getenv("POSTGRES_CONN_STRING")
	JWTSecret = os.Getenv("JWT_SECRET")
}

func init() {
	var err error
	DebugEnabled, err = strconv.ParseBool(os.Getenv("DEBUG_ENABLED"))
	if err != nil && os.Getenv("DEBUG_ENABLED") != "" {
		panic(err)
	}
}
