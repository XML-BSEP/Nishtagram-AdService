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

	return c.campaignRequestRepo.CreateDisposableCampaignRequest(context.Background(), request)
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

	return c.campaignRequestRepo.CreateMultipleCampaignRequest(context.Background(), request)
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

	c.campaignRequestRepo.ApproveDisposableCampaignRequest(context.Background(), request)

	if campaign.Type == 1 {

		for _, ad := range campaign.Post {
			createPostDto := dto.CreatePostDTO{}
			if ad.Type == 0 {
				createPostDto.Image = ad.Path
				createPostDto.IsImage = true
				createPostDto.IsVideo = false
				createPostDto.IsAlbum = false
			} else {
				createPostDto.Video = ad.Path
				createPostDto.IsImage = false
				createPostDto.IsVideo = true
				createPostDto.IsAlbum = false
			}

			createPostDto.Caption = ad.Description + " " + ad.Link
			createPostDto.Location = ad.Location
			createPostDto.Hashtags = ad.HashTags
			createPostDto.CampaignId = campaign.ID
			createPostDto.Link = ad.Link
			createPostDto.UserId = dto.UserTag{UserId: request.InfluencerId}
			createPostDto.CampaignId = request.DisposableCampaign.ID
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
			createStory.CampaignId = campaign.ID
			createStory.Link = ad.Link
			err := gateway.AddStoryFromCampaign(context.Background(), createStory)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (c campaignRequestUseCase) ApproveMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	campaign, err := c.campaignUseCase.GetMultipleCampaign(context.Background(), request.MultipleCampaign.ID, request.AgentId)
	if err != nil {
		return err
	}
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(context.Background(), request.InfluencerId)
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

	err = c.campaignRequestRepo.ApproveMultipleCampaignRequest(context.Background(), request)
	if err != nil {
		return err
	}
	if campaign.Type == 1 {
		for _, ad := range campaign.Post {
			createPostDto := dto.CreatePostDTO{}
			if ad.Type == 0 {
				createPostDto.Image = ad.Path
				createPostDto.IsImage = true
				createPostDto.IsVideo = false
				createPostDto.IsAlbum = false
			} else {
				createPostDto.Video = ad.Path
				createPostDto.IsImage = false
				createPostDto.IsVideo = true
				createPostDto.IsAlbum = false
			}

			createPostDto.Caption = ad.Description + " " + ad.Link
			createPostDto.Location = ad.Location
			createPostDto.Hashtags = ad.HashTags
			createPostDto.CampaignId = campaign.ID
			createPostDto.Link = ad.Link
			createPostDto.UserId = dto.UserTag{UserId: request.InfluencerId}
			createPostDto.CampaignId = request.MultipleCampaign.ID
			err := gateway.AddPostFromCampaign(context.Background(), createPostDto)
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
			createStory.CampaignId = campaign.ID
			createStory.Link = ad.Link
			createStory.Location = domain.Location{Location: ad.Location}
			createStory.CloseFriends = false

			err := gateway.AddStoryFromCampaign(context.Background(), createStory)
			if err != nil {
				return err
			}

		}
	}
	return nil

}

func (c campaignRequestUseCase) RejectDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error {
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(context.Background(), request.InfluencerId)
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

	return c.campaignRequestRepo.RejectDisposableCampaignRequest(context.Background(), request)
}

func (c campaignRequestUseCase) RejectMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	isInfluencer, _ := gateway.CheckIsUserIsInfluencer(context.Background(), request.InfluencerId)
	if !isInfluencer.IsInfluencer {
		return fmt.Errorf("not an influencer")
	}
	if isInfluencer.IsPrivate {
		followDto := dto.FollowDTO{Follower: dto.ProfileDTO{ID: request.InfluencerId}, User: dto.ProfileDTO{ID: request.AgentId}}
		followResponse, _ := gateway.SeeIfAgentFollows(context.Background(), followDto)
		if followResponse.Message == "Allowed to follow" || followResponse.Message == "Request already sent" {
			return fmt.Errorf("user private and not following")
		}
	}

	return c.campaignRequestRepo.RejectMultipleCampaignRequest(context.Background(), request)
}

func (c campaignRequestUseCase) GetAllDisposableCampaignRequests(ctx context.Context, userId string) ([]domain.DisposableCampaignRequest, error) {
	requests, err := c.campaignRequestRepo.GetAllDisposableCampaignRequests(ctx, userId)
	if err != nil {
		return nil, err
	}
	var retVal []domain.DisposableCampaignRequest
	for _, r := range requests {
		campaign, err := c.campaignUseCase.GetDisposableCampaign(context.Background(), r.DisposableCampaign.ID, r.AgentId)
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
		campaign, err := c.campaignUseCase.GetMultipleCampaign(context.Background(), r.MultipleCampaign.ID, r.AgentId)
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

