package usecase

import (
	"ad_service/domain"
	"ad_service/repository"
	"context"
)

type CampaignUseCase interface {
	CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
	CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	GetAllDisposableCampaignsForAgent(ctx context.Context, agentId string) ([]domain.DisposableCampaign, error)
	GetAllMultipleCampaignsForAgent(ctx context.Context, agentId string) ([]domain.MultipleCampaign, error)
	UpdateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	DeleteMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	DeleteDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
}

type campaignUseCase struct {
	campaignRepository repository.CampaignRepo
	adUseCase AdPostUseCase
}

func (c campaignUseCase) GetAllDisposableCampaignsForAgent(ctx context.Context, agentId string) ([]domain.DisposableCampaign, error) {
	disposableCampaigns, err := c.campaignRepository.GetAllDisposableCampaignsForAgenyt(ctx, agentId)
	if err != nil {
		return nil, err
	}
	var retVal []domain.DisposableCampaign
	for _, campaign := range disposableCampaigns {
		var adPosts []domain.AdPost
		for _, ad := range campaign.Post {
			encodedAd, err := c.adUseCase.GetAdById(ctx, agentId, ad.ID)
			if err != nil {
				continue
			}
			adPosts = append(adPosts, encodedAd)
		}

		campaign.Post = adPosts
		retVal = append(retVal, campaign)
	}

	return retVal, nil
}

func (c campaignUseCase) GetAllMultipleCampaignsForAgent(ctx context.Context, agentId string) ([]domain.MultipleCampaign, error) {
	disposableCampaigns, err := c.campaignRepository.GetAllMultipleCampaignsForAgent(ctx, agentId)
	if err != nil {
		return nil, err
	}
	var retVal []domain.MultipleCampaign
	for _, campaign := range disposableCampaigns {
		var adPosts []domain.AdPost
		for _, ad := range campaign.Post {
			encodedAd, err := c.adUseCase.GetAdById(ctx, agentId, ad.ID)
			if err != nil {
				continue
			}
			adPosts = append(adPosts, encodedAd)
		}

		campaign.Post = adPosts
		retVal = append(retVal, campaign)
	}

	return retVal, nil
}

func (c campaignUseCase) UpdateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	return c.campaignRepository.UpdateMultipleCampaign(ctx, campaign)
}

func (c campaignUseCase) DeleteMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	return c.campaignRepository.DeleteMultipleCampaign(ctx, campaign)
}

func (c campaignUseCase) DeleteDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error {
	return c.campaignRepository.DeleteDisposableCampaign(ctx, campaign)
}

func (c campaignUseCase) CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error {
	return c.campaignRepository.CreateDisposableCampaign(ctx, campaign)
}

func (c campaignUseCase) CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	return c.campaignRepository.CreateMultipleCampaign(ctx, campaign)
}

func NewCampaignUseCase(campaignRepo repository.CampaignRepo, adUseCase AdPostUseCase) CampaignUseCase {
	return &campaignUseCase{campaignRepository: campaignRepo, adUseCase: adUseCase}
}