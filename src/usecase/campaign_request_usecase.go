package usecase

import (
	"ad_service/domain"
	"ad_service/dto"
	"ad_service/gateway"
	"ad_service/repository"
	"context"
	"fmt"
)

type CampaignRequestUseCase interface {
	CreateDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error
	CreateMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error
	ApproveDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error
	ApproveMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error
	RejectDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error
	RejectMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error
	GetAllDisposableCampaignRequests(ctx context.Context, userId string) ([]domain.DisposableCampaignRequest, error)
	GetAllMultipleCampaignRequests(ctx context.Context, userId string) ([]domain.MultipleCampaignRequest, error)
}

type campaignRequestUseCase struct {
	campaignRequestRepo repository.CampaignRequestRepo
	campaignUseCase CampaignUseCase
}

func (c campaignRequestUseCase) CreateDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error {
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(ctx, request.InfluencerId)
	if !isInfluencer.IsInfluencer {
		return fmt.Errorf("not an influencer")
	}
	if isInfluencer.IsPrivate {
		followDto := dto.FollowDTO{Follower: dto.ProfileDTO{ID: request.InfluencerId}, User: dto.ProfileDTO{ID: request.AgentId}}
		followResponse, _ := gateway.SeeIfAgentFollows(ctx, followDto)
		if followResponse.Message == "Allowed to follow" || followResponse.Message == "Request already sent" {
			return fmt.Errorf("user private and not following")
		}
	}

	return c.campaignRequestRepo.CreateDisposableCampaignRequest(ctx, request)
}

func (c campaignRequestUseCase) CreateMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(ctx, request.InfluencerId)
	if !isInfluencer.IsInfluencer {
		return fmt.Errorf("not an influencer")
	}
	if isInfluencer.IsPrivate {
		followDto := dto.FollowDTO{Follower: dto.ProfileDTO{ID: request.InfluencerId}, User: dto.ProfileDTO{ID: request.AgentId}}
		followResponse, _ := gateway.SeeIfAgentFollows(ctx, followDto)
		if followResponse.Message == "Allowed to follow" || followResponse.Message == "Request already sent" {
			return fmt.Errorf("user private and not following")
		}
	}

	return c.campaignRequestRepo.CreateMultipleCampaignRequest(ctx, request)
}

func (c campaignRequestUseCase) ApproveDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error {
	campaign, err := c.campaignUseCase.GetDisposableCampaign(ctx, request.DisposableCampaign.ID, request.AgentId)
	if err != nil {
		return err
	}
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(ctx, request.InfluencerId)
	if !isInfluencer.IsInfluencer {
		return fmt.Errorf("not an influencer")
	}
	if isInfluencer.IsPrivate {
		followDto := dto.FollowDTO{Follower: dto.ProfileDTO{ID: request.InfluencerId}, User: dto.ProfileDTO{ID: request.AgentId}}
		followResponse, _ := gateway.SeeIfAgentFollows(ctx, followDto)
		if followResponse.Message == "Allowed to follow" || followResponse.Message == "Request already sent" {
			return fmt.Errorf("user private and not following")
		}
	}

	err = c.campaignRequestRepo.ApproveDisposableCampaignRequest(ctx, request)
	if err != nil {
		return err
	}
	if campaign.Type == 1 {

		for _, ad := range campaign.Post {
			createPostDto := dto.CreatePostDTO{}
			if ad.Type == 0 {
				createPostDto.IsImage = true
				createPostDto.IsVideo = false
				createPostDto.IsAlbum = false
			} else {
				createPostDto.IsImage = false
				createPostDto.IsVideo = true
				createPostDto.IsAlbum = false
			}
			var media []string
			media = append(media, ad.Path)
			createPostDto.Media = media
			createPostDto.Caption = ad.Description + " " + ad.Link
			createPostDto.Location = ad.Location
			createPostDto.Hashtags = ad.HashTags
			createPostDto.UserId = dto.UserTag{UserId: request.InfluencerId}
			err := gateway.AddPostFromCampaign(ctx, createPostDto)
			if err != nil {
				return err
			}

		}
	} else {
		for _, ad := range campaign.Post {
			createStory := dto.StoryDTO{}
			if ad.Type == 0 {
				createStory.IsVideo = false
				createStory.Type = "PHOTO"
			} else {
				createStory.IsVideo = true
				createStory.Type = "VIDEO"
			}
			createStory.UserId = request.InfluencerId
			createStory.Story = ad.Path
			createStory.Location = domain.Location{Location: ad.Location}
			createStory.CloseFriends = false

			err := gateway.AddStoryFromCampaign(ctx, createStory)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (c campaignRequestUseCase) ApproveMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	campaign, err := c.campaignUseCase.GetMultipleCampaign(ctx, request.MultipleCampaign.ID, request.AgentId)
	if err != nil {
		return err
	}
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(ctx, request.InfluencerId)
	if !isInfluencer.IsInfluencer {
		return fmt.Errorf("not an influencer")
	}
	if isInfluencer.IsPrivate {
		followDto := dto.FollowDTO{Follower: dto.ProfileDTO{ID: request.InfluencerId}, User: dto.ProfileDTO{ID: request.AgentId}}
		followResponse, _ := gateway.SeeIfAgentFollows(ctx, followDto)
		if followResponse.Message == "Allowed to follow" || followResponse.Message == "Request already sent" {
			return fmt.Errorf("user private and not following")
		}
	}

	err = c.campaignRequestRepo.ApproveMultipleCampaignRequest(ctx, request)
	if err != nil {
		return err
	}
	if campaign.Type == 1 {
		for _, ad := range campaign.Post {
			createPostDto := dto.CreatePostDTO{}
			if ad.Type == 0 {
				createPostDto.IsImage = true
				createPostDto.IsVideo = false
				createPostDto.IsAlbum = false
			} else {
				createPostDto.IsImage = false
				createPostDto.IsVideo = true
				createPostDto.IsAlbum = false
			}
			var media []string
			media = append(media, ad.Path)
			createPostDto.Media = media
			createPostDto.Caption = ad.Description + " " + ad.Link
			createPostDto.Location = ad.Location
			createPostDto.Hashtags = ad.HashTags
			createPostDto.UserId = dto.UserTag{UserId: request.InfluencerId}
			err := gateway.AddPostFromCampaign(ctx, createPostDto)
			if err != nil {
				return err
			}

		}
	} else {
		for _, ad := range campaign.Post {
			createStory := dto.StoryDTO{}
			if ad.Type == 0 {
				createStory.IsVideo = false
				createStory.Type = "PHOTO"
			} else {
				createStory.IsVideo = true
				createStory.Type = "VIDEO"
			}
			createStory.UserId = request.InfluencerId
			createStory.Story = ad.Path
			createStory.Location = domain.Location{Location: ad.Location}
			createStory.CloseFriends = false

			err := gateway.AddStoryFromCampaign(ctx, createStory)
			if err != nil {
				return err
			}

		}
	}
	return nil

}

func (c campaignRequestUseCase) RejectDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error {
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(ctx, request.InfluencerId)
	if !isInfluencer.IsInfluencer {
		return fmt.Errorf("not an influencer")
	}
	if isInfluencer.IsPrivate {
		followDto := dto.FollowDTO{Follower: dto.ProfileDTO{ID: request.InfluencerId}, User: dto.ProfileDTO{ID: request.AgentId}}
		followResponse, _ := gateway.SeeIfAgentFollows(ctx, followDto)
		if followResponse.Message == "Allowed to follow" || followResponse.Message == "Request already sent" {
			return fmt.Errorf("user private and not following")
		}
	}

	return c.campaignRequestRepo.RejectDisposableCampaignRequest(ctx, request)
}

func (c campaignRequestUseCase) RejectMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(ctx, request.InfluencerId)
	if !isInfluencer.IsInfluencer {
		return fmt.Errorf("not an influencer")
	}
	if isInfluencer.IsPrivate {
		followDto := dto.FollowDTO{Follower: dto.ProfileDTO{ID: request.InfluencerId}, User: dto.ProfileDTO{ID: request.AgentId}}
		followResponse, _ := gateway.SeeIfAgentFollows(ctx, followDto)
		if followResponse.Message == "Allowed to follow" || followResponse.Message == "Request already sent" {
			return fmt.Errorf("user private and not following")
		}
	}

	return c.campaignRequestRepo.RejectMultipleCampaignRequest(ctx, request)
}

func (c campaignRequestUseCase) GetAllDisposableCampaignRequests(ctx context.Context, userId string) ([]domain.DisposableCampaignRequest, error) {
	requests, err := c.campaignRequestRepo.GetAllDisposableCampaignRequests(ctx, userId)
	if err != nil {
		return nil, err
	}
	var retVal []domain.DisposableCampaignRequest
	for _, r := range requests {
		campaign, err := c.campaignUseCase.GetDisposableCampaign(ctx, r.DisposableCampaign.ID, r.AgentId)
		if err != nil {
			continue
		}
		r.DisposableCampaign = campaign
		retVal = append(retVal, r)
	}
	return retVal, nil
}

func (c campaignRequestUseCase) GetAllMultipleCampaignRequests(ctx context.Context, userId string) ([]domain.MultipleCampaignRequest, error) {
	requests, err := c.campaignRequestRepo.GetAllMultipleCampaignRequests(ctx, userId)
	if err != nil {
		return nil, err
	}
	var retVal []domain.MultipleCampaignRequest
	for _, r := range requests {
		campaign, err := c.campaignUseCase.GetMultipleCampaign(ctx, r.MultipleCampaign.ID, r.AgentId)
		if err != nil {
			continue
		}
		r.MultipleCampaign = campaign
		retVal = append(retVal, r)
	}
	return retVal, nil
}

func NewCampaignRequestUseCase(campaignRequestRepo repository.CampaignRequestRepo, campaignUseCase CampaignUseCase) CampaignRequestUseCase {
	return &campaignRequestUseCase{campaignRequestRepo: campaignRequestRepo, campaignUseCase: campaignUseCase}
}

