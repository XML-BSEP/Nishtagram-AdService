package gateway

import (
	"ad_service/domain"
	"ad_service/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func CheckIsUserIsInfluencer(ctx context.Context, userId string) (dto.InfluencerPrivateDTO, error) {
	client := resty.New()
	userDomain := os.Getenv("USER_DOMAIN")
	if userDomain == "" {
		userDomain = "127.0.0.1"
	}
	user := domain.Profile{ID: userId}

	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			SetBody(user).
			EnableTrace().
			Post("https://" + userDomain + ":8082/IsInfluencerAndPrivate")

		var responseDTO dto.InfluencerPrivateDTO
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			fmt.Println(err)
		}

		return responseDTO, nil
	} else {
		resp, _ := client.R().
			SetBody(user).
			EnableTrace().
			Post("https://" + userDomain + ":8082/IsInfluencerAndPrivate")

		var responseDTO dto.InfluencerPrivateDTO
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			fmt.Println(err)
		}

		return responseDTO, nil
	}
}


func GetUser(ctx context.Context, userId string) (dto.ProfileUsernameImageDTO, error) {
	client := resty.New()
	domain := os.Getenv("USER_DOMAIN")
	if domain == "" {
		domain = "127.0.0.1"
	}


	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			EnableTrace().
			Get("https://" + domain + ":8082/getProfileUsernameImageById?userId=" + userId)

		var responseDTO dto.ProfileUsernameImageDTO
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			fmt.Println(err)
		}

		return responseDTO, nil
	} else {
		resp, _ := client.R().
			EnableTrace().
			Get("http://" + domain + ":8082/getProfileUsernameImageById?userId=" + userId)

		var responseDTO dto.ProfileUsernameImageDTO
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			fmt.Println(err)
		}

		return responseDTO, nil
	}

}

