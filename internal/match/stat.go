package match

import (
	// "encoding/json"
	// "fmt"
	"math"
	// "net/http"
	// "net/url"
	// "time"
)

type userListResponse struct {
	Results []struct {
		ID     string `json:"id"`
		Rating int    `json:"rating"`
	} `json:"results"`
}

// poka serv u nix ne rabotayet, tak chto otdayem customniye reytingi
func RequestRating(player *PlayerConn) (int, error) {
	return 1500, nil

	// endpoint := "http://37.9.53.248:8000/users/"
	// client := &http.Client{Timeout: 5 * time.Second}

	// params := url.Values{}
	// params.Add("ids", player.ID)

	// req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	// if err != nil {
	// 	return 0, fmt.Errorf("creating request: %w", err)
	// }

	// req.Header.Set("Authorization", player.AccessToken)

	// resp, err := client.Do(req)
	// if err != nil {
	// 	return 0, fmt.Errorf("sending request: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != 200 {
	// 	return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	// var parsed userListResponse
	// if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
	// 	return 0, fmt.Errorf("decoding response: %w", err)
	// }

	// if len(parsed.Results) == 0 {
	// 	return 0, fmt.Errorf("user not found")
	// }

	// return parsed.Results[0].Rating, nil
}

func GetRatingGain(ratingWinner, ratingLoser int) (int, int) {
	expWinner := 1.0 / (1.0 + math.Pow(10, float64(ratingLoser-ratingWinner)/400.0))
	expLoser := 1.0 - expWinner

	kWinner := getKFactor(ratingWinner)
	kLoser := getKFactor(ratingLoser)

	gainWinner := int(math.Round(kWinner * (1.0 - expWinner))) // actual=1 for winner
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
