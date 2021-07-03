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

	return c.campaignRequestRepo.ApproveDisposableCampaignRequest(ctx, request)
}

func (c campaignRequestUseCase) ApproveMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
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

	return c.campaignRequestRepo.ApproveMultipleCampaignRequest(ctx, request)
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
		campaign, err := c.campaignUseCase.GetDisposableCampaign(ctx, r.DisposableCampaign.ID, userId)
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
		campaign, err := c.campaignUseCase.GetMultipleCampaign(ctx, r.MultipleCampaign.ID, userId)
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

