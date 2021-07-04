package gateway

import (
	"ad_service/dto"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func AddStoryFromCampaign(ctx context.Context, createPost dto.StoryDTO) error {
	client := resty.New()
	userDomain := os.Getenv("STORY_DOMAIN")
	if userDomain == "" {
		userDomain = "127.0.0.1"
	}

	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			SetBody(createPost).
			EnableTrace().
			Post("https://" + userDomain + ":8084/story/createStoryFromCampaign")


		if resp.StatusCode() != 200 {
			return fmt.Errorf("error")
		}

		return nil
	} else {
		resp, _ := client.R().
			SetBody(createPost).
			EnableTrace().
			Post("https://" + userDomain + ":8084/story/createStoryFromCampaign")


		if resp.StatusCode() != 200 {
			return fmt.Errorf("error")
		}

		return nil
	}
}

