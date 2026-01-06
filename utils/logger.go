package utils

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

var output = zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: "15:04:05",
	PartsOrder: []string{"time", "level", "caller", "message"},

	FormatLevel: func(i interface{}) string {
		var levelColor string
		switch i.(string) {
		case "debug":
			levelColor = "\x1b[36m" // Циан
		case "info":
			levelColor = "\x1b[32m" // Зелёный
		case "warn":
			levelColor = "\x1b[33m" // Жёлтый
		case "error":
			levelColor = "\x1b[31m" // Красный
		case "fatal":
			levelColor = "\x1b[35m" // Маджента
		case "panic":
			levelColor = "\x1b[35m" // Маджента
		default:
			levelColor = "\x1b[0m" // Сброс
		}
		return "| " + levelColor + "\x1b[1m" + fmt.Sprintf("%-8s", strings.ToUpper(i.(string))) + "\x1b[0m" + " |"
	},
	FormatCaller: func(i interface{}) string {
		caller := i.(string)
		parts := strings.Split(caller, "/")
		file := parts[len(parts)-1]
		return file + " > "
	},
}

func GetLogger() {
	log.Logger = log.Output(output).With().Caller().Logger()
}
