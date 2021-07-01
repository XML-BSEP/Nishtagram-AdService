package repository

import (
	"ad_service/domain"
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"time"
)

const (
	CreateDisposableCampaign = "CREATE TABLE IF NOT EXISTS adpost_keyspace.DisposableCampaigns (id text, agent_id text, exposure_date time, status int, timestamp timestamp, ad_id list<text>, type int, PRIMARY KEY (agent_id, id));"
	CreateMultipleCampaign = "CREATE TABLE IF NOT EXISTS adpost_keyspace.MultipleCampaigns (id text, agent_id text, start_date date, end_date date, frequency int, status int, timestamp, timestamp, ad_id list<text>, type int, PRIMARY KEY (agent_id, id));"
	InsertIntoDisposableCampaign = "INSERT INTO adpost_keyspace.DisposableCampaigns (id, agent_id, exposure_date, status, timestamp, ad_id, type) VALUES (?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	InsertIntoMultipleCampaign = "INSERT INTO adpost_keyspace.MultipleCampaigns (id, agent_id, start_date, end_date, frequency, timestamp, ad_id, type) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"

)

type CampaignRepo interface {
	CreateDisposableCampaign(ctx context.Context, campaign domain.DisposableCampaign) error
	CreateMultipleCampaign(ctx context.Context, campaign domain.MultipleCampaign) error
}

type campaignRepository struct {
	cassandraClient *gocql.Session
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
