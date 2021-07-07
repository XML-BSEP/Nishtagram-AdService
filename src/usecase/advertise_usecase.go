package usecase

import (
	"ad_service/domain"
	"ad_service/domain/events"
	"ad_service/dto"
	"ad_service/gateway"
	"ad_service/repository"
	"context"
)

type AdvertiseUseCase interface {
	AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign) error
	AddMultipleCampaignToAdvertisementTable(ctx context.Context, multipleCampaign domain.MultipleCampaign) error
	GetAllPostAdsForUser(ctx context.Context, profileId string) ([]dto.ShowAdPost, error)
	GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]dto.StoryDTO, error)
	GenerateStatisticsReport(ctx context.Context, agentId string, campaignId string) (domain.StatisticsReport, error)
	AddClickEvent(ctx context.Context, event events.ClickEvent) error
}

type advertiseUseCase struct {
	adPostUseCase AdPostUseCase
	advertiseRepository repository.AdvertisementRepo
	likeRepository repository.LikeRepo
}

func (a advertiseUseCase) AddClickEvent(ctx context.Context, event events.ClickEvent) error {
	return a.advertiseRepository.InsertClickEvent(context.Background(), event.InfluencerId, event.CampaignId)
}

func (a advertiseUseCase) GenerateStatisticsReport(ctx context.Context, agentId string, campaignId string) (domain.StatisticsReport, error) {
	timesAdvertised, err := a.advertiseRepository.GetTimesAdvertised(context.Background(), campaignId, agentId)
	if err != nil {
		return domain.StatisticsReport{}, err
	}
	clicks, err := a.advertiseRepository.GetNumberOfClicks(context.Background(), campaignId)
	return domain.StatisticsReport{CampaignId: campaignId, Clicks: clicks, AdvertisingCount: timesAdvertised}, nil

}

func checkIfElementExists(list []string, element string) bool {
	for _, elem := range list {
		if elem == element {
			return true
		}
	}
	return false
}

func (a advertiseUseCase) AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign) error {
	var profileIds []string
	for _, ad := range disposableCampaign.Post {
		profileIdsLocation, err := gateway.GetProfilesByLocation(ctx, ad.Location)
		if err == nil {
			for _, id := range profileIdsLocation {
				if !checkIfElementExists(profileIds, id) {
					profileIds = append(profileIds, id)
				}
			}
		}


		for _, hashtag := range ad.HashTags {
			profileIdsHash, err := gateway.GetProfilesByHashtag(ctx, hashtag)
			if err != nil {
				continue
			}
			for _, id := range profileIdsHash {
				if !checkIfElementExists(profileIds, id) {
					profileIds = append(profileIds, id)
				}
			}

		}


	}
	return a.advertiseRepository.AddDisposableCampaignToAdvertisementTable(context.Background(), disposableCampaign, profileIds)
}

func (a advertiseUseCase) AddMultipleCampaignToAdvertisementTable(ctx context.Context, multipleCampaign domain.MultipleCampaign) error {
	var profileIds []string
	for _, ad := range multipleCampaign.Post {
		profileIdsLocation, err := gateway.GetProfilesByLocation(ctx, ad.Location)
		if err != nil {
			continue
		}
		for _, id := range profileIdsLocation {
			if !checkIfElementExists(profileIds, id) {
				profileIds = append(profileIds, id)
			}
		}

		for _, hashtag := range ad.HashTags {
			profileIdsHash, err := gateway.GetProfilesByHashtag(ctx, hashtag)
			if err != nil {
				continue
			}
			for _, id := range profileIdsHash {
				if !checkIfElementExists(profileIds, id) {
					profileIds = append(profileIds, id)
				}
			}

		}


	}
	return a.advertiseRepository.AddMultipleCampaignToAdvertisementTable(multipleCampaign, profileIds)
}

func (a advertiseUseCase) GetAllPostAdsForUser(ctx context.Context, profileId string) ([]dto.ShowAdPost, error) {
	var retVal []dto.ShowAdPost
	adsToShow, err := a.advertiseRepository.GetAllPostAdsForUser(context.Background(), profileId)
	if err != nil {
		return nil, err
	}

	for _, adToShow := range adsToShow {
		oneAd, err := a.adPostUseCase.GetAdById(ctx, adToShow.AgentId.ID, adToShow.ID)
		if err != nil {
			continue
		}
		adToAdd := dto.ShowAdPost{}
		if oneAd.Type == 0 {
			adToAdd.IsVideo = false
			adToAdd.IsAlbum = false
		} else {

			adToAdd.IsVideo = true
			adToAdd.IsAlbum = false
		}
		adToAdd.Location = oneAd.Location
		user, _ := gateway.GetUser(context.Background(), oneAd.AgentId.ID)
		adToAdd.User = domain.Profile{ProfileId : oneAd.AgentId.ID, ProfilePhoto: user.ProfilePhoto, Username: user.Username}
		media := make([]string, 1)
		media[0] = oneAd.Path
		adToAdd.Media = media
		adToAdd.NumOfComments = oneAd.NumOfComments
		adToAdd.NumOfLikes = oneAd.NumOfLikes
		adToAdd.NumOfDislikes = oneAd.NumOfDislikes
		adToAdd.Link = oneAd.Link
		adToAdd.Description = oneAd.Description
		adToAdd.IsCampaign = true
		adToAdd.Timestamp = oneAd.Timestamp
		adToAdd.PostBy = oneAd.AgentId.ID
		adToAdd.Id = oneAd.ID
		if a.likeRepository.SeeIfLikeExists(oneAd.ID, profileId, context.Background()) {
			adToAdd.IsLiked = true
		}

		if  a.likeRepository.SeeIfDislikeExists(oneAd.ID, profileId, context.Background()) {
			adToAdd.IsDisliked = true
		}

		retVal = append(retVal, adToAdd)
	}

	return retVal, nil
}

func (a advertiseUseCase) GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]dto.StoryDTO, error) {
	var retVal []dto.StoryDTO
	adsToShow, err := a.advertiseRepository.GetAllStoryAdsForUser(ctx, profileId)
	if err != nil {
		return nil, err
	}

	for _, adToShow := range adsToShow {
		oneAd, err := a.adPostUseCase.GetAdById(ctx, adToShow.AgentId.ID, adToShow.ID)
		if err != nil {
			continue
		}
		var story dto.StoryDTO
		if oneAd.Type == 0 {
			story.IsVideo = false
		} else {
			story.IsVideo = true
		}
		story.Story = oneAd.Path
		story.Link = oneAd.Link
		story.MediaPath.Path = oneAd.Path
		story.CampaignId = adToShow.CampaignId
		if oneAd.Type == 1 {
			story.StoryContent = dto.StoryContent{IsVideo: true, Content: story.MediaPath.Path}
		} else {
			story.StoryContent = dto.StoryContent{IsVideo: false, Content: story.MediaPath.Path}
		}
		user, _ := gateway.GetUser(context.Background(), oneAd.AgentId.ID)
		story.User = domain.Profile{ProfileId : oneAd.AgentId.ID, ProfilePhoto: user.ProfilePhoto, Username: user.Username}
		story.Timestamp = oneAd.Timestamp

		retVal = append(retVal, story)
	}

	return retVal, nil
}

func NewAdvertiseUseCase(adPostUseCase AdPostUseCase, advertiseRepository repository.AdvertisementRepo, likeRepo repository.LikeRepo) AdvertiseUseCase {
	return &advertiseUseCase{adPostUseCase: adPostUseCase, advertiseRepository: advertiseRepository, likeRepository: likeRepo}
}
