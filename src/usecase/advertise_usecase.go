package usecase

import (
	"ad_service/domain"
	"ad_service/dto"
	"ad_service/gateway"
	"ad_service/repository"
	"context"
)

type AdvertiseUseCase interface {
	AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign) error
	AddMultipleCampaignToAdvertisementTable(ctx context.Context, multipleCampaign domain.MultipleCampaign) error
	GetAllPostAdsForUser(ctx context.Context, profileId string) ([]dto.ShowAdPost, error)
	GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error)
}

type advertiseUseCase struct {
	adPostUseCase AdPostUseCase
	advertiseRepository repository.AdvertisementRepo
	likeRepository repository.LikeRepo
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

func (a advertiseUseCase) GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error) {
	var retVal []domain.AdPost
	adsToShow, err := a.advertiseRepository.GetAllStoryAdsForUser(ctx, profileId)
	if err != nil {
		return nil, err
	}

	for _, adToShow := range adsToShow {
		oneAd, err := a.adPostUseCase.GetAdById(ctx, adToShow.AgentId.ID, adToShow.ID)
		if err != nil {
			continue
		}
		retVal = append(retVal, oneAd)
	}

	return retVal, nil
}

func NewAdvertiseUseCase(adPostUseCase AdPostUseCase, advertiseRepository repository.AdvertisementRepo, likeRepo repository.LikeRepo) AdvertiseUseCase {
	return &advertiseUseCase{adPostUseCase: adPostUseCase, advertiseRepository: advertiseRepository, likeRepository: likeRepo}
}
