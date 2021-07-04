package repository

import (
	"ad_service/domain"
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)
const (
	CreateDisposableRequestTable = "CREATE TABLE IF NOT EXISTS adpost_keyspace.DisposableCampaignRequest (id text, influencer_id text, agent_id text, campaign_id text, status int, PRIMARY KEY(influencer_id, id));"
	CreateMultipleRequestTable = "CREATE TABLE IF NOT EXISTS adpost_keyspace.MultipleCampaignRequest (id text, influencer_id text, agent_id text, campaign_id text, status int, PRIMARY KEY(influencer_id, id));"
	InsertIntoDisposableRequestTable = "INSERT INTO adpost_keyspace.DisposableCampaignRequest (id, influencer_id, agent_id, campaign_id, status) VALUES (?, ?, ?, ?, ?) IF NOT EXISTS;"
	InsertIntoMultipleRequestTable = "INSERT INTO adpost_keyspace.MultipleCampaignRequest (id, influencer_id, agent_id, campaign_id, status) VALUES (?, ?, ?, ?, ?) IF NOT EXISTS;"
	UpdateDisposableRequestStatus = "UPDATE adpost_keyspace.DisposableCampaignRequest SET status = ? WHERE id = ? AND influencer_id = ?;"
	UpdateMultipleRequestStatus = "UPDATE adpost_keyspace.MultipleCampaignRequest SET status = ? WHERE id = ? AND influencer_id = ?;"
	GetAllDisposableCampaignRequestsByInfluencer = "SELECT id, campaign_id, status, agent_id FROM adpost_keyspace.DisposableCampaignRequest WHERE influencer_id = ?;"
	GetAllMultipleCampaignRequestsByInfluencer = "SELECT id, campaign_id, status, agent_id FROM adpost_keyspace.MultipleCampaignRequest WHERE influencer_id = ?;"


)
type CampaignRequestRepo interface {
	CreateDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error
	CreateMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error
	ApproveDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error
	ApproveMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error
	RejectDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error
	RejectMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error
	GetAllDisposableCampaignRequests(ctx context.Context, userId string) ([]domain.DisposableCampaignRequest, error)
	GetAllMultipleCampaignRequests(ctx context.Context, userId string) ([]domain.MultipleCampaignRequest, error)
}



type campaignRequestRepository struct {
	cassandraClient *gocql.Session
}

func (c campaignRequestRepository) CreateDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error {
	return c.cassandraClient.Query(InsertIntoDisposableRequestTable, uuid.NewString(), request.InfluencerId, request.AgentId, request.DisposableCampaign.ID, 0).Exec()
}

func (c campaignRequestRepository) CreateMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	return c.cassandraClient.Query(InsertIntoMultipleRequestTable, uuid.NewString(), request.InfluencerId, request.AgentId, request.MultipleCampaign.ID, 0).Exec()
}

func (c campaignRequestRepository) ApproveDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error {
	return c.cassandraClient.Query(UpdateDisposableRequestStatus, 1, request.Id, request.InfluencerId).Exec()
}

func (c campaignRequestRepository) ApproveMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	return c.cassandraClient.Query(UpdateMultipleRequestStatus, 1, request.Id, request.InfluencerId).Exec()
}

func (c campaignRequestRepository) RejectDisposableCampaignRequest(ctx context.Context, request domain.DisposableCampaignRequest) error {
	return c.cassandraClient.Query(UpdateDisposableRequestStatus, 2, request.Id, request.InfluencerId).Exec()
}

func (c campaignRequestRepository) RejectMultipleCampaignRequest(ctx context.Context, request domain.MultipleCampaignRequest) error {
	return c.cassandraClient.Query(UpdateMultipleRequestStatus, 2, request.Id, request.InfluencerId).Exec()
}

func (c campaignRequestRepository) GetAllDisposableCampaignRequests(ctx context.Context, userId string) ([]domain.DisposableCampaignRequest, error) {
	var id, campaignId, agentId string
	var status int
	iter := c.cassandraClient.Query(GetAllDisposableCampaignRequestsByInfluencer, userId).Iter().Scanner()
	var retVal []domain.DisposableCampaignRequest
	for iter.Next() {
		err := iter.Scan(&id, &campaignId, &status, &agentId)
		if err != nil {
			continue
		}
		retVal = append(retVal, domain.DisposableCampaignRequest{Id: id, InfluencerId: userId, AgentId: agentId, DisposableCampaign: domain.DisposableCampaign{ID: campaignId}, RequestStatus: domain.RequestStatus(status)})
	}

	return retVal, nil
}

func (c campaignRequestRepository) GetAllMultipleCampaignRequests(ctx context.Context, userId string) ([]domain.MultipleCampaignRequest, error) {
	var id, campaignId, agentId string
	var status int
	iter := c.cassandraClient.Query(GetAllMultipleCampaignRequestsByInfluencer, userId).Iter().Scanner()
	var retVal []domain.MultipleCampaignRequest
	for iter.Next() {
		err := iter.Scan(&id, &campaignId, &status, &agentId)
		if err != nil {
			continue
		}
		retVal = append(retVal, domain.MultipleCampaignRequest{Id: id, InfluencerId: userId, AgentId: agentId, MultipleCampaign: domain.MultipleCampaign{ID: campaignId}, RequestStatus: domain.RequestStatus(status)})
	}

	return retVal, nil
}

func NewCampaignRequestRepository(cassandraClient *gocql.Session) CampaignRequestRepo {
	err := cassandraClient.Query(CreateDisposableRequestTable).Exec()
	if err != nil {
		return nil
	}

	err = cassandraClient.Query(CreateMultipleRequestTable).Exec()
	if err != nil {
		return nil
	}

	return &campaignRequestRepository{cassandraClient: cassandraClient}
}