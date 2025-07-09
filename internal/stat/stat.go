package stat

import (
	"math"

	"github.com/lesta-battleship/server-core/internal/match"
)


func RequestRating(player *match.PlayerConn) (int, error) {
	// TODO, написать обращение к сервису статистики, чтобы получить рейтинг
	// сейчас он фулово лежит на любые команды поиска, так что не понимаю как тестить

	// В player лежит id, login и токен пользователя -- добавил все, так как хз что конкретно понадобится
	return 1500, nil
}

func GetRatingGain(ratingWinner, ratingLoser int) (int, int) {
	expWinner := 1.0 / (1.0 + math.Pow(10, float64(ratingLoser-ratingWinner)/400.0))
	expLoser := 1.0 - expWinner

	kWinner := getKFactor(ratingWinner)
	kLoser := getKFactor(ratingLoser)

	
	gainWinner := int(math.Round(kWinner * (1.0 - expWinner)))  // actual=1 for winner
	gainLoser := int(math.Round(kLoser * (0.0 - expLoser)))    // actual=0 for loser

	return gainWinner - gainLoser, gainLoser - gainWinner
}

// getKFactor determines the K-factor based on player rating
func getKFactor(rating int) float64 {
	switch {
	case rating < 2100:
		return 32.0
	case rating >= 2100 && rating < 2400:
		return 24.0
	default: // 2400 and above
		return 16.0
	}
}
