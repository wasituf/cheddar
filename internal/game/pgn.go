package game

import (
	"fmt"
	"log"
	"strings"

	chess "github.com/corentings/chess/v2"
)

type gameUpdater func(*chess.Game)

func GetPGNFromString(pgnStr string) gameUpdater {
	r := strings.NewReader(pgnStr)

	pgn, err := chess.PGN(r)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error importing from PGN")
	}

	return pgn
}
