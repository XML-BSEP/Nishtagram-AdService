package usecase

import (
	"ad_service/domain"
	"ad_service/repository"
	"context"
)

type CampaignUseCase interface {
	CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
	CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
}

type campaignUseCase struct {
	campaignRepository repository.CampaignRepo
}

func (c campaignUseCase) CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error {
	return c.campaignRepository.CreateDisposableCampaign(ctx, campaign)
}

func (c campaignUseCase) CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	return c.campaignRepository.CreateMultipleCampaign(ctx, campaign)
}

func NewCampaignUseCase(campaignRepo repository.CampaignRepo) CampaignUseCase {
	return &campaignUseCase{campaignRepository: campaignRepo}
}