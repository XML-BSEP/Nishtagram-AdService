package gateway

import (
	"ad_service/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func GetProfilesByLocation(ctx context.Context, location string) ([]string, error) {
	client := resty.New()
	userDomain := os.Getenv("SEARCH_DOMAIN")
	if userDomain == "" {
		userDomain = "127.0.0.1"
	}

	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			EnableTrace().
			Get("http://" + userDomain + ":8087/getPostLocationsByLocationContaining?searchedLocation=" + location)

		if resp.StatusCode() != 200 {
			return nil, fmt.Errorf("error")
		}
		var postLocation []dto.PostLocationsDTO
		err := json.Unmarshal(resp.Body(), &postLocation)
		if err != nil {
			fmt.Println(err)
		}
		var retVal []string
		for _, s := range postLocation {
			for _, id := range s.PostProfileId {
				retVal = append(retVal, id.ProfileId)
			}
		}
		return retVal, nil
	} else {
		resp, _ := client.R().
			EnableTrace().
			Get("http://" + userDomain + ":8087/getPostLocationsByLocationContaining?searchedLocation=" + location)

		if resp.StatusCode() != 200 {
			return nil, fmt.Errorf("error")
		}

		var postLocation []dto.PostLocationsDTO
		err := json.Unmarshal(resp.Body(), &postLocation)
		if err != nil {
			fmt.Println(err)
		}
		var retVal []string
		for _, s := range postLocation {
			for _, id := range s.PostProfileId {
				retVal = append(retVal, id.ProfileId)
			}
		}
		return retVal, nil
	}
}


func GetProfilesByHashtag(ctx context.Context, hashtag string) ([]string, error) {
	client := resty.New()
	userDomain := os.Getenv("SEARCH_DOMAIN")
	if userDomain == "" {
		userDomain = "127.0.0.1"
	}

	if os.Getenv("DOCKER_ENV") == "" {
		resp, _ := client.R().
			EnableTrace().
			Get("http://" + userDomain + ":8087/getPostsByTag?searchedTag=" + hashtag)

		if resp.StatusCode() != 200 {
			return nil, fmt.Errorf("error")
		}
		var postLocation []dto.PostTagsDTO
		err := json.Unmarshal(resp.Body(), &postLocation)
		if err != nil {
			fmt.Println(err)
		}
		var retVal []string
		for _, s := range postLocation {
			for _, id := range s.PostProfileId {
				retVal = append(retVal, id.ProfileId)
			}
		}
		return retVal, nil
	} else {
		resp, _ := client.R().
			EnableTrace().
			Get("http://" + userDomain + ":8087/getPostsByTag?searchedTag=" + hashtag)

		if resp.StatusCode() != 200 {
			return nil, fmt.Errorf("error")
		}

		var postLocation []dto.PostTagsDTO
		err := json.Unmarshal(resp.Body(), &postLocation)
		if err != nil {
			fmt.Println(err)
		}
		var retVal []string
		for _, s := range postLocation {
			for _, id := range s.PostProfileId {
				retVal = append(retVal, id.ProfileId)
			}
		}
		return retVal, nil
	}
}

