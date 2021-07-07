package repository

import (
	"ad_service/domain"
	"ad_service/domain/events"
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

const (
	CreateAdvertisementTable = "CREATE TABLE IF NOT EXISTS adpost_keyspace.AdvertisementTable (profile_id text, ad_id text, agent_id text, advertisement_time timestamp, type text, seen boolean, campaign_id text, PRIMARY KEY ((profile_id), advertisement_time))  WITH CLUSTERING ORDER BY (advertisement_time DESC);;"
	InsertIntoAdvertisementTable = "INSERT INTO adpost_keyspace.AdvertisementTable (profile_id, ad_id, agent_id, advertisement_time, type, seen, campaign_id) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	SelectAllAdsForAdvertisement = "SELECT ad_id, agent_id, advertisement_time, campaign_id, type FROM adpost_keyspace.AdvertisementTable WHERE profile_id = ? AND advertisement_time >= ? AND advertisement_time <= ?;"
	CreateStatisticsTable = "CREATE TABLE IF NOT EXISTS adpost_keyspace.StatisticsTable (campaign_id text, influencer_id  text, clicks int, PRIMARY KEY (campaign_id, influencer_id)); "
	InsertIntoStatisticsTable = "INSERT INTO adpost_keyspace.StatisticsTable(campaign_id, influencer_id, clicks) VALUES (?, ?, ?) IF NOT EXISTS;"
	UpdateClick = "UPDATE adpost_keyspace.StatisticsTable SET clicks = ? WHERE campaign_id = ? AND influencer_id = ?;"
	GetNumberOfClicksForAll = "SELECT influencer_id, clicks FROM adpost_keyspace.StatisticsTable WHERE campaign_id = ?;"
	GetNumberOfClick = "SELECT clicks FROM adpost_keyspace.StatisticsTable WHERE campaign_id = ? AND influencer_id = ?;"
	SeeIfExists = "SELECT count(*) FROM adpost_keyspace.StatisticsTable WHERE campaign_id = ? AND influencer_id = ?;"
	CreateAdvertisingCountTable = "CREATE TABLE IF NOT EXISTS adpost_keyspace.AdvertisingCount (campaign_id text, agent_id text, times_advertised int, PRIMARY KEY (agent_id, campaign_id));"
	InsertIntoAdvertisingCountTable = "INSERT INTO adpost_keyspace.AdvertisingCount (campaign_id, agent_id, times_advertised) VALUES (?, ?, ?) IF NOT EXISTS;"
	UpdateAdvertisingTime = "UPDATE adpost_keyspace.AdvertisingCount SET times_advertised = ? WHERE campaign_id = ? AND agent_id = ?;"
	GetAllAdvertisingTimesForAgent = "SELECT times_advertised, campaign_id FROM adpost_keyspace.AdvertisingCount WHERE agent_id = ?;"
	GetTimesAdvertised = "SELECT times_advertised FROM adpost_keyspace.AdvertisingCount WHERE campaign_id = ? AND agent_id = ?;"
	)

type AdvertisementRepo interface {
	AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign, userIds []string) error
	AddMultipleCampaignToAdvertisementTable(multipleCampaign domain.MultipleCampaign, userIds []string) error
	GetAllPostAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error)
	GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error)
	InsertClickEvent(ctx context.Context, influencerId string, campaignId string) error
	GetNumberOfClicks(ctx context.Context, campaignId string) ([]events.ClickEvent, error)
	GetTimesAdvertised(ctx context.Context, campaignId string, agentId string) (domain.AdvertisingCount, error)
}

type advertisementRepository struct {
	cassandraClient *gocql.Session
}

func (a advertisementRepository) GetTimesAdvertised(ctx context.Context, campaignId string, agentId string) (domain.AdvertisingCount, error) {
	var timesAdvertised int
	a.cassandraClient.Query(GetTimesAdvertised, campaignId, agentId).Iter().Scan(&timesAdvertised)
	return domain.AdvertisingCount{AdvertisedCount: timesAdvertised, CampaignId: campaignId, AgentId: agentId}, nil
}

func (a advertisementRepository) InsertClickEvent(ctx context.Context, influencerId string, campaignId string) error {
	var count int
	a.cassandraClient.Query(SeeIfExists, campaignId, influencerId).Iter().Scan(&count)
	if count > 0 {
		var numOfClicks int
		a.cassandraClient.Query(GetNumberOfClick, campaignId, influencerId).Iter().Scan(&numOfClicks)
		numOfClicks += 1
		err := a.cassandraClient.Query(UpdateClick, numOfClicks, campaignId, influencerId).Exec()
		if err != nil {
			return err
		}
	}
	return a.cassandraClient.Query(InsertIntoStatisticsTable, campaignId, influencerId, 1).Exec()
}

func (a advertisementRepository) GetNumberOfClicks(ctx context.Context, campaignId string) ([]events.ClickEvent, error) {
	var retVal []events.ClickEvent
	var influencerId string
	var clicks int
	iter := a.cassandraClient.Query(GetNumberOfClicksForAll, campaignId).Iter().Scanner()
	for iter.Next() {
		err := iter.Scan(&influencerId, &clicks)
		if err != nil {
			continue
		}
		retVal = append(retVal, events.ClickEvent{Clicks: clicks, InfluencerId: influencerId, CampaignId: campaignId})

	}
	return retVal, nil
}

func (a advertisementRepository) AddDisposableCampaignToAdvertisementTable(ctx context.Context, disposableCampaign domain.DisposableCampaign, userIds []string) error {
	for _, userId := range userIds {
		for _, ad := range disposableCampaign.Post {
			if disposableCampaign.Type == 0 {
				err := a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, disposableCampaign.AgentId.ID, disposableCampaign.ExposureDate, "STORY", false, disposableCampaign.ID).Exec()
				if err != nil {
					continue
				}
			} else {
				err := a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, disposableCampaign.AgentId.ID, disposableCampaign.ExposureDate, "POST", false, disposableCampaign.ID).Exec()
				if err != nil {
					continue
				}
			}
		}
	}
	return a.cassandraClient.Query(InsertIntoAdvertisingCountTable, disposableCampaign.ID, disposableCampaign.AgentId.ID, 0).Exec()

}


func (a advertisementRepository) AddMultipleCampaignToAdvertisementTable(multipleCampaign domain.MultipleCampaign, userIds []string) error {
	var exposureDates []time.Time
	startTime :=  time.Date(multipleCampaign.StartDate.Year(), multipleCampaign.StartDate.Month(), multipleCampaign.StartDate.Day(), 8, 0, 0, multipleCampaign.StartDate.Nanosecond(), multipleCampaign.StartDate.Location())
	endTime :=  time.Date(multipleCampaign.EndDate.Year(), multipleCampaign.EndDate.Month(), multipleCampaign.EndDate.Day(), 0, 0, 0, multipleCampaign.StartDate.Nanosecond(), multipleCampaign.StartDate.Location())

	if multipleCampaign.AdvertisementFrequency <= 0 {
		return fmt.Errorf("error while adding adv")
	}

	hoursToAdd := 16 / multipleCampaign.AdvertisementFrequency

	for i := startTime; i.After(endTime) == false; i = i.Add(time.Hour*24) {
		for hour := 0; hour <= 16; hour += hoursToAdd {
			date := i
			exposureDates = append(exposureDates, date.Add(time.Hour * time.Duration(hour)))
		}
	}

	for _, userId := range userIds {
		for _, ad := range multipleCampaign.Post {
			for _, date := range exposureDates {
				if multipleCampaign.Type == 0 {
					err :=a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, multipleCampaign.AgentId.ID, date, "STORY", false, multipleCampaign.ID).Exec()
					if err != nil {
						continue
					}
				} else {
					err := a.cassandraClient.Query(InsertIntoAdvertisementTable, userId, ad.ID, multipleCampaign.AgentId.ID, date, "POST", false, multipleCampaign.ID).Exec()
					if err != nil {
						continue
					}
				}
			}

		}
	}
	return a.cassandraClient.Query(InsertIntoAdvertisingCountTable, multipleCampaign.ID, multipleCampaign.AgentId.ID, 0).Exec()

	return nil
}

func (a advertisementRepository) GetAllPostAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error) {
	var adId, agentId, campaignId, typeAd string
	var advertisementTime time.Time
	timeString := time.Now().Format("02-01-2006 15:04:05")
	timeToCompare, _ := time.Parse("02-01-2006 15:04:05", timeString)
	fmt.Println(timeToCompare.Add(-time.Minute * 15).String())
	fmt.Println(timeToCompare.Add(time.Minute * 15).String())
	minTime := timeToCompare.Add(-time.Minute * 15)
	maxTime := timeToCompare.Add(time.Minute * 15)
	iter := a.cassandraClient.Query(SelectAllAdsForAdvertisement, profileId, minTime, maxTime).Iter().Scanner()
	var retVal []domain.AdPost
	for iter.Next() {
		err := iter.Scan(&adId, &agentId, &advertisementTime, &campaignId, &typeAd)
	/*	if advertisementTime.Before(minTime) {
		continue
		}
		if advertisementTime.After(maxTime) {
			continue
		}*/
		if err != nil {
			continue
		}
		if typeAd != "POST" {
			continue
		}
		retVal = append(retVal, domain.AdPost{ID: adId, AgentId: domain.Profile{ID: agentId}})

		var timesAdvertised int
		a.cassandraClient.Query(GetTimesAdvertised, campaignId, agentId).Iter().Scan(&timesAdvertised)
		timesAdvertised += 1
		a.cassandraClient.Query(UpdateAdvertisingTime, timesAdvertised, campaignId, agentId).Exec()

	}
	return retVal, nil
}

func (a advertisementRepository) GetAllStoryAdsForUser(ctx context.Context, profileId string) ([]domain.AdPost, error) {
	var adId, agentId, campaignId, typeAd string
	var advertisementTime time.Time
	timeString := time.Now().Format("02-01-2006 15:04:05")
	timeToCompare, _ := time.Parse("02-01-2006 15:04:05", timeString)
	fmt.Println(timeToCompare.Add(-time.Minute * 15).String())
	fmt.Println(timeToCompare.Add(time.Minute * 15).String())
	minTime := timeToCompare.Add(-time.Minute * 15)
	maxTime := timeToCompare.Add(time.Minute * 15)
	iter := a.cassandraClient.Query(SelectAllAdsForAdvertisement, profileId, minTime, maxTime).Iter().Scanner()
	var retVal []domain.AdPost
	for iter.Next() {
		err := iter.Scan(&adId, &agentId, &advertisementTime, &campaignId, &typeAd)
		if err != nil {
			continue
		}

		if typeAd != "STORY" {
			continue
		}
		retVal = append(retVal, domain.AdPost{ID: adId, AgentId: domain.Profile{ID: agentId}})
		var timesAdvertised int
		a.cassandraClient.Query(GetTimesAdvertised, campaignId, agentId).Iter().Scan(&timesAdvertised)
		timesAdvertised += 1
		a.cassandraClient.Query(UpdateAdvertisingTime, timesAdvertised, campaignId, agentId).Exec()
	}
	return retVal, nil
}

func NewAdvertisementRepository(cassandraClient *gocql.Session) AdvertisementRepo {
	err := cassandraClient.Query(CreateAdvertisementTable).Exec()
	if err != nil {
		return nil
	}
	err = cassandraClient.Query(CreateStatisticsTable).Exec()
	if err != nil {
		return nil
	}
	err = cassandraClient.Query(CreateAdvertisingCountTable).Exec()
	if err != nil {
		return nil
	}
	return &advertisementRepository{cassandraClient: cassandraClient}
}


