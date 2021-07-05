package usecase

import (
	"ad_service/domain"
	"ad_service/dto"
	"ad_service/gateway"
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
	GetDisposableCampaign(ctx context.Context, campaignId string, agentId string) (domain.DisposableCampaign, error)
	GetMultipleCampaign(ctx context.Context, campaignId string, agentId string) (domain.MultipleCampaign, error)
}

type campaignUseCase struct {
	campaignRepository repository.CampaignRepo
	adUseCase AdPostUseCase
	advertiseUseCase AdvertiseUseCase
}

func (c campaignUseCase) GetDisposableCampaign(ctx context.Context, campaignId string, agentId string) (domain.DisposableCampaign, error) {
	campaign, err := c.campaignRepository.GetDisposableCampaign(ctx, campaignId, agentId)
	if err != nil {
		return domain.DisposableCampaign{}, err
	}
	var adPosts []domain.AdPost
	for _, ad := range campaign.Post {
		encodedAd, err := c.adUseCase.GetAdById(ctx, agentId, ad.ID)
		if err != nil {
			continue
		}
		adPosts = append(adPosts, encodedAd)
	}

		campaign.Post = adPosts


	return campaign, nil
}

func (c campaignUseCase) GetMultipleCampaign(ctx context.Context, campaignId string, agentId string) (domain.MultipleCampaign, error) {
	campaign, err := c.campaignRepository.GetMultipleCampaign(ctx, campaignId, agentId)
	if err != nil {
		return domain.MultipleCampaign{}, err
	}
	var adPosts []domain.AdPost
	for _, ad := range campaign.Post {
		encodedAd, err := c.adUseCase.GetAdById(ctx, agentId, ad.ID)
		if err != nil {
			continue
		}
		adPosts = append(adPosts, encodedAd)
	}

	campaign.Post = adPosts


	return campaign, nil
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
	err := c.campaignRepository.CreateDisposableCampaign(ctx, campaign)
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
			var media []string
			media = append(media, ad.Path)
			createPostDto.Media = media
			createPostDto.Caption = ad.Description + " " + ad.Link
			createPostDto.Location = ad.Location
			createPostDto.Hashtags = ad.HashTags
			createPostDto.UserId = dto.UserTag{UserId: campaign.AgentId.ID}
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
			createStory.UserId = campaign.AgentId.ID
			createStory.Story = ad.Path
			createStory.Location = domain.Location{Location: ad.Location}
			createStory.CloseFriends = false

			err := gateway.AddStoryFromCampaign(ctx, createStory)
			if err != nil {
				return err
			}
		}

	}
	err = c.advertiseUseCase.AddDisposableCampaignToAdvertisementTable(ctx, campaign)
	if err != nil {
		return err
	}
	return nil
}

func (c campaignUseCase) CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	err := c.campaignRepository.CreateMultipleCampaign(ctx, campaign)
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
			createPostDto.UserId = dto.UserTag{UserId: campaign.AgentId.ID}
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
			createStory.UserId = campaign.AgentId.ID
			createStory.Story = ad.Path
			createStory.Location = domain.Location{Location: ad.Location}
			createStory.CloseFriends = false

			err := gateway.AddStoryFromCampaign(ctx, createStory)
			if err != nil {
				return err
			}
		}

	}
	err = c.advertiseUseCase.AddMultipleCampaignToAdvertisementTable(ctx, campaign)
	if err != nil {
		return err
	}
	return nil
}

func NewCampaignUseCase(campaignRepo repository.CampaignRepo, adUseCase AdPostUseCase, advertiseUseCase AdvertiseUseCase) CampaignUseCase {
	return &campaignUseCase{campaignRepository: campaignRepo, adUseCase: adUseCase, advertiseUseCase: advertiseUseCase}
}