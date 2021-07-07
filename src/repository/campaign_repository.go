package repository

import (
	"ad_service/domain"
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"time"
)

const (
	CreateDisposableCampaign = "CREATE TABLE IF NOT EXISTS adpost_keyspace.DisposableCampaigns (id text, agent_id text, exposure_date timestamp, status int, timestamp timestamp, ad_id list<text>, type int, PRIMARY KEY (agent_id, id));"
	CreateMultipleCampaign = "CREATE TABLE IF NOT EXISTS adpost_keyspace.MultipleCampaigns (id text, agent_id text, start_date timestamp, end_date timestamp, frequency int, status int, timestamp timestamp, ad_id list<text>, type int, PRIMARY KEY (agent_id, id));"
	CreateTokenTable = "CREATE TABLE IF NOT EXISTS adpost_keyspace.NishtagramApiTokens (user_token text, user_id text, PRIMARY KEY (user_token));"
	InsertIntoTokenTable = "INSERT INTO adpost_keyspace.NishtagramApiTokens (user_token, user_id) VALUES (?, ?) IF NOT EXISTS;"
	SelectUserId = "SELECT FROM adpost_keyspace.NishtagramApiTokens user_id WHERE user_token = ?;"
	InsertIntoDisposableCampaign = "INSERT INTO adpost_keyspace.DisposableCampaigns (id, agent_id, exposure_date, status, timestamp, ad_id, type) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	InsertIntoMultipleCampaign = "INSERT INTO adpost_keyspace.MultipleCampaigns (id, agent_id, start_date, end_date, frequency, timestamp, ad_id, type) VALUES (?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	GetAllDisposableCampaigns = "SELECT id, exposure_date, status, timestamp, ad_id, type FROM adpost_keyspace.DisposableCampaigns WHERE agent_id = ?;"
	GetAllMultipleCampaigns = "SELECT id, start_date, end_date, frequency, status, timestamp, ad_id, type FROM adpost_keyspace.MultipleCampaigns WHERE agent_id = ?;"
	UpdateMultipleCampaign = "UPDATE adpost_keyspace.MultipleCampaigns SET start_date = ?, end_date = ?, frequency = ? WHERE id = ? AND agent_id = ?;"
	DeleteMultipleCampaign = "DELETE FROM adpost_keyspace.MultipleCampaigns WHERE id = ? AND agent_id = ?;"
	DeleteDisposableCampaign = "DELETE FROM adpost_keyspace.DisposableCampaigns WHERE id = ? AND agent_id = ?;"
	GetDisposableCampaign = "SELECT id, exposure_date, status, timestamp, ad_id, type FROM adpost_keyspace.DisposableCampaigns WHERE agent_id = ? AND id = ?;"
	GetMultipleCampaign = "SELECT id, start_date, end_date, frequency, status, timestamp, ad_id, type FROM adpost_keyspace.MultipleCampaigns WHERE agent_id = ? AND id = ?;"
	CreatePendingEditCampaign = "CREATE TABLE IF NOT EXISTS adpost_keyspace.PendingEditCampaign (campaign_id text, start_date timestamp, end_date timestamp, frequency int, timestamp timestamp, PRIMARY KEY(campaign_id));"
	InsertIntoPendingEditCampaign = "INSERT INTO adpost_keyspace.PendingEditCampaign (campaign_id, start_date, end_date, frequency, timestamp) VALUES (?, ?, ?, ?, ?);"
	SeeIfPendingEditExists = "SELECT count(*) FROM adpost_keyspace.PendingEditCampaign WHERE campaign_id = ?;"
	DeletePendingEdit = "DELETE FROM adpost_keyspace.PendingEditCampaign WHERE campaign_id = ?;"
	GetNewChanges = "SELECT start_date, end_date, frequency, timestamp FROM adpost_keyspace.PendingEditCampaign WHERE campaign_id = ?;"
	)

type CampaignRepo interface {
	CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
	CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	GetAllMultipleCampaignsForAgent(ctx context.Context, agentId string) ([]domain.MultipleCampaign, error)
	GetAllDisposableCampaignsForAgent(ctx context.Context, agentId string) ([]domain.DisposableCampaign, error)
	UpdateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	DeleteMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	DeleteDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
	GetDisposableCampaign(ctx context.Context, campaignId string, userId string) (domain.DisposableCampaign, error)
	GetMultipleCampaign(ctx context.Context, campaignId string, userId string) (domain.MultipleCampaign, error)
	InsertIntoTokenTable(ctx context.Context, token string, userId string) error
	GetUserIdByToken(ctx context.Context, token string) (string, error)
}

type campaignRepository struct {
	cassandraClient *gocql.Session
}

func (c campaignRepository) InsertIntoTokenTable(ctx context.Context, token string, userId string) error {
	return c.cassandraClient.Query(InsertIntoTokenTable, token, userId).Exec()
}

func (c campaignRepository) GetUserIdByToken(ctx context.Context, token string) (string, error) {
	var userId string
	c.cassandraClient.Query(SelectUserId, token).Iter().Scan(&userId)
	if userId == "" {
		return userId, fmt.Errorf("no such token")
	}
	return userId, nil
}

func (c campaignRepository) GetDisposableCampaign(ctx context.Context, campaignId string, userId string) (domain.DisposableCampaign, error) {
	var id string
	var exposureDate, timestamp time.Time
	var status, campaignType int
	var adIds []string
	var ads []domain.AdPost

	err := c.cassandraClient.Query(GetDisposableCampaign, userId, campaignId).Iter().Scan(&id, &exposureDate, &status, &timestamp, &adIds, &campaignType)
	if !err {
		return domain.DisposableCampaign{}, fmt.Errorf("no such campaign")
	}
	for _, a := range adIds {
		ads = append(ads, domain.AdPost{ID: a})
	}
	return domain.DisposableCampaign{ID: id, Post: ads, AgentId: domain.Profile{ID: userId}, ExposureDate: exposureDate, Type: domain.Type(campaignType)}, nil

}

func (c campaignRepository) GetMultipleCampaign(ctx context.Context, campaignId string, userId string) (domain.MultipleCampaign, error) {
	var id string
	var startDate, endDate, timestamp time.Time
	var status, campaignType, frequency int
	var adIds []string

	err := c.cassandraClient.Query(GetMultipleCampaign, userId, campaignId).Iter().Scan(&id, &startDate, &endDate, &frequency, &status, &timestamp, &adIds, &campaignType)
	if !err {
		return domain.MultipleCampaign{}, fmt.Errorf("no such campaign")
	}
	var ads []domain.AdPost
	for _, a := range adIds {
		ads = append(ads, domain.AdPost{ID: a})
	}
	return domain.MultipleCampaign{ID: id, Post: ads, AgentId: domain.Profile{ID: userId}, StartDate: startDate, EndDate: endDate, AdvertisementFrequency: frequency, Type: domain.Type(campaignType)}, nil
}

func (c campaignRepository) GetAllMultipleCampaignsForAgent(ctx context.Context, agentId string) ([]domain.MultipleCampaign, error) {
	iter := c.cassandraClient.Query(GetAllMultipleCampaigns, agentId).Iter().Scanner()
	var id string
	var startDate, endDate, timestamp time.Time
	var status, campaignType, frequency int
	var adIds []string
	var retVal []domain.MultipleCampaign
	timeString := time.Now().Format("02-01-2006 15:04:05")
	timeToCompare, _ := time.Parse("02-01-2006 15:04:05", timeString)
	for iter.Next() {
		err := iter.Scan(&id, &startDate, &endDate, &frequency, &status, &timestamp, &adIds, &campaignType)
		if err != nil {
			continue
		}
		var count int
		_ = c.cassandraClient.Query(SeeIfPendingEditExists, id).Iter().Scan(&count)
		if count > 0 {
			var newStartDate, newEndDate, timeOfChange time.Time
			var freq int
			c.cassandraClient.Query(GetNewChanges, id).Iter().Scan(&newStartDate, &newEndDate, &freq, &timeOfChange)
			isDayAfter := timeOfChange.Add(time.Hour*24)
			if isDayAfter.Before(timeToCompare) {
				c.cassandraClient.Query(UpdateMultipleCampaign, newStartDate, newEndDate, freq, id, agentId)
				startDate = newStartDate
				endDate = newEndDate
				frequency = freq
				err := c.cassandraClient.Query(DeletePendingEdit, id).Exec()
				fmt.Println(err)
			}
		}
		var ads []domain.AdPost
		for _, a := range adIds {
			ads = append(ads, domain.AdPost{ID: a})
		}
		retVal = append(retVal, domain.MultipleCampaign{ID: id, Post: ads, AgentId: domain.Profile{ID: agentId}, StartDate: startDate, EndDate: endDate, AdvertisementFrequency: frequency, Type: domain.Type(campaignType)})

	}
	return retVal, nil
}

func (c campaignRepository) GetAllDisposableCampaignsForAgent(ctx context.Context, agentId string) ([]domain.DisposableCampaign, error) {
	iter := c.cassandraClient.Query(GetAllDisposableCampaigns, agentId).Iter().Scanner()
	var id string
	var exposureDate, timestamp time.Time
	var status, campaignType int
	var adIds []string
	var retVal []domain.DisposableCampaign

	for iter.Next() {
		var ads []domain.AdPost
		err := iter.Scan(&id, &exposureDate, &status, &timestamp, &adIds, &campaignType)
		if err != nil {
			continue
		}

		for _, a := range adIds {
			ads = append(ads, domain.AdPost{ID: a})
		}
		retVal = append(retVal, domain.DisposableCampaign{ID: id, Post: ads, AgentId: domain.Profile{ID: agentId}, ExposureDate: exposureDate, Type: domain.Type(campaignType)})

	}
	return retVal, nil
}

func (c campaignRepository) UpdateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	var count int
	_ = c.cassandraClient.Query(SeeIfPendingEditExists, campaign.ID).Iter().Scan(&count)
	if count > 0 {
		err := c.cassandraClient.Query(DeletePendingEdit, campaign.ID).Exec()
		fmt.Println(err)
	}

	return c.cassandraClient.Query(InsertIntoPendingEditCampaign, campaign.ID, campaign.StartDate, campaign.EndDate, campaign.AdvertisementFrequency, time.Now()).Exec()
}

func (c campaignRepository) DeleteMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {

	return c.cassandraClient.Query(DeleteMultipleCampaign, campaign.ID, campaign.AgentId.ID).Exec()
}

func (c campaignRepository) DeleteDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error {
	return c.cassandraClient.Query(DeleteDisposableCampaign, campaign.ID, campaign.AgentId.ID).Exec()
}

func (c campaignRepository) CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error {
	id := uuid.NewString()
	timestamp := time.Now()
	var adIds []string
	for _, ad := range campaign.Post {
		adIds = append(adIds, ad.ID)
	}
	err := c.cassandraClient.Query(InsertIntoDisposableCampaign, id, campaign.AgentId.ID, campaign.ExposureDate, campaign.Status, timestamp, adIds, campaign.Type).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (c campaignRepository) CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	id := uuid.NewString()
	timestamp := time.Now()
	var adIds []string
	for _, ad := range campaign.Post {
		adIds = append(adIds, ad.ID)
	}
	err := c.cassandraClient.Query(InsertIntoMultipleCampaign, id, campaign.AgentId.ID, campaign.StartDate, campaign.EndDate, campaign.AdvertisementFrequency, timestamp, adIds, campaign.Type).Exec()
	if err != nil {
		return err
	}
	return nil
}

func NewCampaignRepo(cassandraClient *gocql.Session) CampaignRepo {
	err := cassandraClient.Query(CreateMultipleCampaign).Exec()
	if err != nil {
		return nil
	}
	err = cassandraClient.Query(CreatePendingEditCampaign).Exec()
	if err != nil {
		return nil
	}
	err = cassandraClient.Query(CreateDisposableCampaign).Exec()

	if err != nil {
		return nil
	}

	err = cassandraClient.Query(CreateTokenTable).Exec()
	if err != nil {
		return nil
	}

	return &campaignRepository{cassandraClient: cassandraClient}
}
