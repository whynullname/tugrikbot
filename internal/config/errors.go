package config

import "errors"

var ErrorEmptyBotToken error = errors.New("error empty bot token in env")
var ErrorEmptyDBDSN error = errors.New("error empty db dsn in env")
