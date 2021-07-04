package gateway

import (
	"ad_service/dto"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func AddPostFromCampaign(ctx context.Context, createPost dto.CreatePostDTO) error {
	client := resty.New()
	userDomain := os.Getenv("POST_DOMAIN")
	if userDomain == "" {
		userDomain = "127.0.0.1"
	}

	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			SetBody(createPost).
			EnableTrace().
			Post("https://" + userDomain + ":8083/createPostFromCampaign")

		if resp.StatusCode() != 200 {
			return fmt.Errorf("error")
		}

		return nil
	} else {
		resp, _ := client.R().
			SetBody(createPost).
			EnableTrace().
			Post("https://" + userDomain + ":8083/createPostFromCampaign")

		if resp.StatusCode() != 200 {
			return fmt.Errorf("error")
		}

		return nil
	}
}
