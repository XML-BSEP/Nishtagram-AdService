package usecase

import (
	"ad_service/domain"
	"ad_service/gateway"
	"ad_service/repository"
	"context"
)

type AdvertiseUseCase interface {
	AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign) error
	AddMultipleCampaignToAdvertisementTable(ctx context.Context, multipleCampaign domain.MultipleCampaign) error
	GetAllPostAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error)
	GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error)
}

type advertiseUseCase struct {
	adPostUseCase AdPostUseCase
	advertiseRepository repository.AdvertisementRepo
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
	return a.advertiseRepository.AddDisposableCampaignToAdvertisementTable(ctx, disposableCampaign, profileIds)
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
	return a.advertiseRepository.AddMultipleCampaignToAdvertisementTable(ctx, multipleCampaign, profileIds)
}

func (a advertiseUseCase) GetAllPostAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error) {
	var retVal []domain.AdPost
	adsToShow, err := a.GetAllPostAdsForUser(ctx, profileId)
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

func (a advertiseUseCase) GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error) {
	var retVal []domain.AdPost
	adsToShow, err := a.GetAllStoryAdsForUser(ctx, profileId)
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

func NewAdvertiseUseCase(adPostUseCase AdPostUseCase, advertiseRepository repository.AdvertisementRepo) AdvertiseUseCase {
	return &advertiseUseCase{adPostUseCase: adPostUseCase, advertiseRepository: advertiseRepository}
}
