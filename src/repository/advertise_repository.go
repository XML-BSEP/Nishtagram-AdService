package repository

import (
	"ad_service/domain"
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

const (
	CreateAdvertisementTable = "CREATE TABLE IF NOT EXIST apdpost_keyspace.AdvertisementTable (profile_id text, ad_id text, agent_id text, advertisement_time timestamp, type text, seen bool, PRIMARY KEY (profile_id, advertisement_time, type));"
	InsertIntoAdvertisementTable = "INSERT INTO adpost_keyspace.AdvertisementTable (profile_id, ad_id, agent_id, advertisement_time, type, seen) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	SelectAllAdsForAdvertisement = "SELECT ad_id, agent_id, advertisement_time FROM adpost_keyspace.AdvertisementTable WHERE profile_id = ? AND advertisement_time >= ? AND advertisement_time <= ? AND type = ?;"
)

type AdvertisementRepo interface {
	AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign, userIds []string) error
	AddMultipleCampaignToAdvertisementTable(ctx context.Context, multipleCampaign domain.MultipleCampaign, userIds []string) error
	GetAllPostAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error)
	GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error)
}

type advertisementRepository struct {
	cassandraClient *gocql.Session
}

func (a advertisementRepository) AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign, userIds []string) error {
	for _, userId := range userIds {
		for _, ad := range disposableCampaign.Post {
			if disposableCampaign.Type == 0 {
				a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, disposableCampaign.AgentId, disposableCampaign.ExposureDate, "STORY", false)
			} else {
				a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, disposableCampaign.AgentId, disposableCampaign.ExposureDate, "POST", false)
			}
		}
	}
	return nil
}


func (a advertisementRepository) AddMultipleCampaignToAdvertisementTable(ctx context.Context, multipleCampaign domain.MultipleCampaign, userIds []string) error {
	var exposureDates []time.Time
	startTime :=  time.Date(multipleCampaign.StartDate.Year(), multipleCampaign.StartDate.Month(), multipleCampaign.StartDate.Day(), 8, 0, 0, multipleCampaign.StartDate.Nanosecond(), multipleCampaign.StartDate.Location())
	endTime :=  time.Date(multipleCampaign.EndDate.Year(), multipleCampaign.EndDate.Month(), multipleCampaign.EndDate.Day(), 0, 0, 0, multipleCampaign.StartDate.Nanosecond(), multipleCampaign.StartDate.Location())

	if multipleCampaign.AdvertisementFrequency <= 0 {
		return fmt.Errorf("error while adding adv")
	}

	hoursToAdd := 16 / multipleCampaign.AdvertisementFrequency

	for i := startTime; i.After(endTime) == false; i.AddDate(0, 0, 1) {
		for hour := 0; hour <= hoursToAdd; hour += 1 {
			exposureDates = append(exposureDates, i.Add(time.Hour * time.Duration(hour)))
		}
	}

	for _, userId := range userIds {
		for _, ad := range multipleCampaign.Post {
			for _, date := range exposureDates {
				if multipleCampaign.Type == 0 {
					a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, multipleCampaign.AgentId, date, "STORY", false)
				} else {
					a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, multipleCampaign.AgentId, date, "POST", false)
				}
			}

		}
	}
	return nil
}

func (a advertisementRepository) GetAllPostAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error) {
	var adId, agentId string
	var advertisementTime time.Time
	timeToCompare := time.Now()
	iter := a.cassandraClient.Query(SelectAllAdsForAdvertisement, profileId, timeToCompare.Add(-time.Minute * 15), timeToCompare.Add(time.Minute * 15), "POST").Iter().Scanner()
	var retVal []domain.AdPost
	for iter.Next() {
		err := iter.Scan(&adId, &agentId, &advertisementTime)
		if err != nil {
			continue
		}
		retVal = append(retVal, domain.AdPost{ID: adId, AgentId: domain.Profile{ID: agentId}})

	}
	return retVal, nil
}

func (a advertisementRepository) GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error) {
	var adId, agentId string
	var advertisementTime time.Time
	timeToCompare := time.Now()
	iter := a.cassandraClient.Query(SelectAllAdsForAdvertisement, profileId, timeToCompare.Add(-time.Minute * 15), timeToCompare.Add(time.Minute * 15), "STORY").Iter().Scanner()
	var retVal []domain.AdPost
	for iter.Next() {
		err := iter.Scan(&adId, &agentId, &advertisementTime)
		if err != nil {
			continue
		}
		retVal = append(retVal, domain.AdPost{ID: adId, AgentId: domain.Profile{ID: agentId}})

	}
	return retVal, nil
}

func NewAdvertisementRepository(cassandraClient *gocql.Session) AdvertisementRepo {
	err := cassandraClient.Query(CreateAdvertisementTable).Exec()
	if err != nil {
		return nil
	}
	return &advertisementRepository{cassandraClient: cassandraClient}
}


