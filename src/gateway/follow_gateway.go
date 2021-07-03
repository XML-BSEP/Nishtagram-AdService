package gateway

import (
	"ad_service/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func SeeIfAgentFollows(ctx context.Context, followDto dto.FollowDTO) (dto.FollowResponseDto, error) {
	client := resty.New()

	domain := os.Getenv("FOLLOW_DOMAIN")
	if domain == "" {
		domain = "127.0.0.1"
	}
	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			SetBody(followDto).
			EnableTrace().
			Post("https://" + domain + ":8089/isAllowedToFollow")

		var responseDTO dto.FollowResponseDto
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			fmt.Println(err)
		}

		return responseDTO, nil
	} else {
		resp, _ := client.R().
			SetBody(followDto).
			EnableTrace().
			Post("http://" + domain + ":8089/isAllowedToFollow")

		var responseDTO dto.FollowResponseDto
		err := json.Unmarshal(resp.Body(), &responseDTO)
		if err != nil {
			fmt.Println(err)
		}

		return responseDTO, nil
	}
}
