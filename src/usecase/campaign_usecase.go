package usecase

import (
	"ad_service/domain"
	"ad_service/dto"
	"ad_service/gateway"
	"ad_service/repository"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
	"time"
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
	CreateApiToken(ctx context.Context, userId string) (string, error)
	SeeIfTokenExists(ctx context.Context, token string) (string, error)
	GenerateStatisticsReport(ctx context.Context, agentId string) ([]domain.StatisticsReport, error)

}

type campaignUseCase struct {
	campaignRepository repository.CampaignRepo
	adUseCase AdPostUseCase
	advertiseUseCase AdvertiseUseCase
}

func (c campaignUseCase) GenerateStatisticsReport(ctx context.Context, agentId string) ([]domain.StatisticsReport, error) {
	campaigns, err := c.GetAllMultipleCampaignsForAgent(context.Background(), agentId)
	if err != nil {
		fmt.Println(err)
	}
	var retVal []domain.StatisticsReport
	for _,campaign := range campaigns {
		report, err := c.advertiseUseCase.GenerateStatisticsReport(context.Background(), agentId, campaign.ID)
		if err != nil {
			continue
		}
		adFull, err := c.adUseCase.GetAdById(context.Background(), campaign.AgentId.ID, campaign.Post[0].ID)
		report.Description = adFull.Description
		var links []string
		for i, _ := range campaign.Post {
			adToGenerate, err := c.adUseCase.GetAdById(context.Background(), campaign.AgentId.ID, campaign.Post[i].ID)
			if err != nil {
				continue
			}
			links = append(links, adToGenerate.Link)
		}
		report.AdvertisedLinks = links
		retVal = append(retVal, report)
	}
	campaignsD, err := c.campaignRepository.GetAllDisposableCampaignsForAgent(context.Background(), agentId)
	if err != nil {
		fmt.Println(err)
	}
	for _,campaign := range campaignsD {
		report, err := c.advertiseUseCase.GenerateStatisticsReport(context.Background(), agentId, campaign.ID)
		if err != nil {
			continue
		}
		adFull, err := c.adUseCase.GetAdById(context.Background(), campaign.AgentId.ID, campaign.Post[0].ID)
		report.Description = adFull.Description
		var links []string
		for i, _ := range campaign.Post {
			adToGenerate, err := c.adUseCase.GetAdById(context.Background(), campaign.AgentId.ID, campaign.Post[i].ID)
			if err != nil {
				continue
			}
			links = append(links, adToGenerate.Link)
		}
		report.AdvertisedLinks = links
		retVal = append(retVal, report)
	}
	return retVal, nil
}

func (c campaignUseCase) CreateApiToken(ctx context.Context, userId string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = uuid.NewString()
	atClaims["refresh_uuid"] = uuid.NewString()
	atClaims["role"] = "AGENT"
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Second * 3600).Unix()


	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	err = c.campaignRepository.InsertIntoTokenTable(context.Background(), token, userId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c campaignUseCase) SeeIfTokenExists(ctx context.Context, token string) (string, error) {
	return c.campaignRepository.GetUserIdByToken(context.Background(), token)
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
	disposableCampaigns, err := c.campaignRepository.GetAllDisposableCampaignsForAgent(ctx, agentId)
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
	campaign.ID = uuid.NewString()
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

			createPostDto.Caption = ad.Description + " " + ad.Link
			createPostDto.Location = ad.Location
			createPostDto.Hashtags = ad.HashTags
			createPostDto.CampaignId = campaign.ID
			createPostDto.Link = ad.Link
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
			createStory.CampaignId = campaign.ID
			createStory.Link = ad.Link
			createStory.Location = domain.Location{Location: ad.Location}
			createStory.CloseFriends = false

			err := gateway.AddStoryFromCampaign(ctx, createStory)
			if err != nil {
				return err
			}
		}

	}
	err = c.advertiseUseCase.AddDisposableCampaignToAdvertisementTable(context.Background(), campaign)
	if err != nil {
		return err
	}
	return nil
}

func (c campaignUseCase) CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	campaign.ID = uuid.NewString()
	err := c.campaignRepository.CreateMultipleCampaign(ctx, campaign)
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
	err = c.advertiseUseCase.AddMultipleCampaignToAdvertisementTable(context.Background(), campaign)
	if err != nil {
		return err
	}
	return nil
}

func NewCampaignUseCase(campaignRepo repository.CampaignRepo, adUseCase AdPostUseCase, advertiseUseCase AdvertiseUseCase) CampaignUseCase {
	return &campaignUseCase{campaignRepository: campaignRepo, adUseCase: adUseCase, advertiseUseCase: advertiseUseCase}
}