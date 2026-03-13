package controller

import (
	"bufio"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

func readInput(log zerolog.Logger) string {
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		log.Error().Err(err).Msg("err in reading stdin")
	}

	return strings.TrimSpace(text)
}
