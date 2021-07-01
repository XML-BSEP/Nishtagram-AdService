package repository

import (
	"ad_service/domain"
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	time2 "time"
)

const (
	CreateAdPostRepo = "CREATE TABLE IF NOT EXISTS adpost_keyspace.AdPosts (id text, agent_id text, media text, description text, " +
		"timestamp timestamp, link string, hashtags list<text>, location string, type int, num_of_likes int, num_of_dislikes int, PRIMARY KEY(agent_id, id));"
	InsertIntoAdPostRepo = "INSERT INTO adpost_keyspace.AdPosts (id, agent_id, media, description, timestamp, link, hashtags, location, type, num_of_likes, num_of_dislikes) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS;"
	DeleteFromAdPostRepo = "DELETE FROM adpost_keyspace.AdPosts WHERE id = ? AND agent_id = ?;"
	GetAdsByAgent = "SELECT id, agent_id, media, description, timestamp, link, hashtags, location, type, num_of_likes, num_of_dislikes FROM adpost_keyspace.AdPosts WHERE agent_id = ?;"
	GetAdsByAgentAndId = "SELECT id, agent_id, media, description, timestamp, link, hashtags, location, type, num_of_likes, num_of_dislikes FROM adpost_keyspace.AdPosts WHERE agent_id = ? AND id = ?;"

	)
type AdPostRepo interface {
	CreateAd(ctx context.Context, ad domain.AdPost) error
	EditAd(ctx context.Context, ad domain.AdPost)
	DeleteAd(ctx context.Context, ad domain.AdPost)
	GetAdsByAgent(ctx context.Context, agentId string) ([]domain.AdPost, error)
	GetAdByAgentIdAndId(ctx context.Context, agentId string, id string) (domain.AdPost, error)
}

type adPostRepository struct {
	cassandraClient *gocql.Session
}

func (a adPostRepository) GetAdByAgentIdAndId(ctx context.Context, agent string, adId string) (domain.AdPost, error) {
	var id, agentId, media, description, link, location string
	var adType, numOfLikes, numOfDislikes int
	var hashtags []string
	var timestamp time2.Time

	iter := a.cassandraClient.Query(GetAdsByAgentAndId, agent, adId).Iter().Scan(&id, &agentId, &media, &description, &timestamp, &link, &hashtags, &location, &adType, &numOfLikes, &numOfDislikes)
	if !iter {
		return domain.AdPost{}, fmt.Errorf("error while unmarshaling ad")
	}

	return domain.AdPost{ID: id, AgentId: domain.Profile{ID: agentId}, Description: description, Timestamp: timestamp, HashTags: hashtags, Location: location,
		NumOfLikes: numOfLikes, NumOfDislikes: numOfDislikes, Link: link, Path: media}, nil

}

func (a adPostRepository) GetAdsByAgent(ctx context.Context, agent string) ([]domain.AdPost, error) {
	var id, agentId, media, description, link, location string
	var adType, numOfLikes, numOfDislikes int
	var hashtags []string
	var timestamp time2.Time

	var retVal []domain.AdPost
	iter := a.cassandraClient.Query(GetAdsByAgent, agent).Iter().Scanner()

	for iter.Next() {
		err := iter.Scan(&id, &agentId, &media, &description, &timestamp, &link, &hashtags, &location, &adType, &numOfLikes, &numOfDislikes)
		if err != nil {
			continue
		}

		retVal = append(retVal, domain.AdPost{ID: id, AgentId: domain.Profile{ID: agentId}, Description: description, Timestamp: timestamp, HashTags: hashtags, Location: location,
			NumOfLikes: numOfLikes, NumOfDislikes: numOfDislikes, Link: link, Path: media})

	}

	return retVal, nil
}

func (a adPostRepository) CreateAd(ctx context.Context, ad domain.AdPost) error {
	id := uuid.NewString()
	time := time2.Now()
	err := a.cassandraClient.Query(InsertIntoAdPostRepo, id, ad.AgentId.ID, ad.Path, ad.Description, time, ad.Link, ad.HashTags, ad.Location, ad.Type, 0, 0).Exec()

	if err != nil {
		return err
	}

	return nil
}

func (a adPostRepository) EditAd(ctx context.Context, ad domain.AdPost) {
	panic("implement me")
}

func (a adPostRepository) DeleteAd(ctx context.Context, ad domain.AdPost) {
	panic("implement me")
}

func NewAdPostRepo(cassandraClient *gocql.Session) AdPostRepo {
	err := cassandraClient.Query(CreateAdPostRepo).Exec()
	if err != nil {
		return nil
	}
	return &adPostRepository{cassandraClient: cassandraClient}
}