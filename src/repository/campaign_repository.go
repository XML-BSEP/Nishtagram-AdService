package repository

import (
	"ad_service/domain"
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"time"
)

const (
	CreateDisposableCampaign = "CREATE TABLE IF NOT EXISTS adpost_keyspace.DisposableCampaigns (id text, agent_id text, exposure_date timestamp, status int, timestamp timestamp, ad_id list<text>, type int, PRIMARY KEY (agent_id, id));"
	CreateMultipleCampaign = "CREATE TABLE IF NOT EXISTS adpost_keyspace.MultipleCampaigns (id text, agent_id text, start_date timestamp, end_date timestamp, frequency int, status int, timestamp timestamp, ad_id list<text>, type int, PRIMARY KEY (agent_id, id));"
	InsertIntoDisposableCampaign = "INSERT INTO adpost_keyspace.DisposableCampaigns (id, agent_id, exposure_date, status, timestamp, ad_id, type) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	InsertIntoMultipleCampaign = "INSERT INTO adpost_keyspace.MultipleCampaigns (id, agent_id, start_date, end_date, frequency, timestamp, ad_id, type) VALUES (?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	GetAllDisposableCampaigns = "SELECT id, exposure_date, status, timestamp, ad_id, type FROM adpost_keyspace.DisposableCampaigns WHERE agent_id = ?;"
	GetAllMultipleCampaigns = "SELECT id, start_date, end_date, frequency, status, timestamp, ad_id, type FROM adpost_keyspace.MultipleCampaigns WHERE agent_id;"
	UpdateMultipleCampaign = "UPDATE adpost_keyspace.MultipleCampaigns SET start_date = ?, end_date = ?, frequency = ? WHERE id = ? AND agent_id = ?;"
	DeleteMultipleCampaign = "DELETE FROM adpost_keyspace.MultipleCampaigns WHERE id = ? AND agent_id = ?;"
	DeleteDisposableCampaign = "DELETE FROM adpost_keyspace.DisposableCampaigns WHERE id = ? AND agent_id = ?;"
)

type CampaignRepo interface {
	CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
	CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	GetAllMultipleCampaignsForAgent(ctx context.Context, agentId string) ([]domain.MultipleCampaign, error)
	GetAllDisposableCampaignsForAgenyt(ctx context.Context, agentId string) ([]domain.DisposableCampaign, error)
	UpdateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	DeleteMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
	DeleteDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
}

type campaignRepository struct {
	cassandraClient *gocql.Session
}

func (c campaignRepository) GetAllMultipleCampaignsForAgent(ctx context.Context, agentId string) ([]domain.MultipleCampaign, error) {
	iter := c.cassandraClient.Query(GetAllMultipleCampaigns, agentId).Iter().Scanner()
	var id string
	var startDate, endDate, timestamp time.Time
	var status, campaignType, frequency int
	var adIds []string
	var retVal []domain.MultipleCampaign
	for iter.Next() {
		err := iter.Scan(&id, &startDate, &endDate, &frequency, &status, &timestamp, &adIds, &campaignType)
		if err != nil {
			continue
		}
		var ads []domain.AdPost
		for _, a := range adIds {
			ads = append(ads, domain.AdPost{ID: a})
		}
		retVal = append(retVal, domain.MultipleCampaign{ID: id, Post: ads, AgentId: domain.Profile{ID: agentId}, StartDate: startDate, EndDate: endDate, AdvertisementFrequency: frequency, Type: campaignType})

	}
	return retVal, nil
}

func (c campaignRepository) GetAllDisposableCampaignsForAgenyt(ctx context.Context, agentId string) ([]domain.DisposableCampaign, error) {
	iter := c.cassandraClient.Query(GetAllDisposableCampaigns, agentId).Iter().Scanner()
	var id string
	var exposureDate, timestamp time.Time
	var status, campaignType int
	var adIds []string
	var retVal []domain.DisposableCampaign
	var ads []domain.AdPost
	for iter.Next() {
		err := iter.Scan(&id, &exposureDate, &status, &timestamp, &adIds, &campaignType)
		if err != nil {
			continue
		}
		for _, a := range adIds {
			ads = append(ads, domain.AdPost{ID: a})
		}
		retVal = append(retVal, domain.DisposableCampaign{ID: id, Post: ads, AgentId: domain.Profile{ID: agentId}, ExposureDate: exposureDate, Type: campaignType})

	}
	return retVal, nil
}

func (c campaignRepository) UpdateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error {
	return c.cassandraClient.Query(UpdateMultipleCampaign, campaign.StartDate, campaign.EndDate, campaign.AdvertisementFrequency, campaign.ID, campaign.AgentId.ID ).Exec()
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

	err = cassandraClient.Query(CreateDisposableCampaign).Exec()

	if err != nil {
		return nil
	}

	return &campaignRepository{cassandraClient: cassandraClient}
}
