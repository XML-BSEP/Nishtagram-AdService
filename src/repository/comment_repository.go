package repository

import (
	"ad_service/domain"
	"ad_service/dto"
	"ad_service/gateway"
	"context"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"time"
)
const (
	CreateCommentTable = "CREATE TABLE if not exists adpost_keyspace.Comments (id text, comment text, ad_id text, comment_by text, timestamp timestamp, mentions list<text>, " +
		"PRIMARY KEY (ad_id, id));"
	InsertComment = "INSERT INTO adpost_keyspace.Comments (id, comment, ad_id, comment_by, timestamp, mentions) VALUES (?, ? , ?, ?, ?, ?);"
	DeleteComment = "DELETE FROM adpost_keyspace.Comments WHERE ad_id = ? AND comment_by = ?;"
	GetComments = "SELECT id, comment, ad_id, comment_by, timestamp, mentions FROM adpost_keyspace.Comments WHERE ad_id = ?;"
	AddCommentToPost = "UPDATE adpost_keyspace.AdPosts SET num_of_comments = ?  WHERE id = ? and agent_id = ?;"

)
type CommentRepo interface {
	CommentPost(comment dto.CommentDTO, ctx context.Context) error
	DeleteComment(comment dto.CommentDTO, ctx context.Context) error
	GetComments(postId string, ctx context.Context) ([]dto.CommentDTO, error)
}

type commentRepository struct {
	cassandraSession *gocql.Session
}

func (c commentRepository) GetComments(postId string, ctx context.Context) ([]dto.CommentDTO, error) {
	//id, comment, post_id, comment_by, timestamp, mentions
	var id, comment, post_id, comment_by string
	var timestamp time.Time
	var mentions []string
	var retVal []dto.CommentDTO
	iter := c.cassandraSession.Query(GetComments, postId).Iter().Scanner()

	for iter.Next() {
		iter.Scan(&id, &comment, &post_id, &comment_by, &timestamp, &mentions)
		dto := dto.CommentDTO{}
		profile, err := gateway.GetUser(context.Background(), comment_by)
		if err != nil {
		}
		dto.CommentBy = domain.Profile{ID: comment_by, ProfilePhoto: profile.ProfilePhoto, Username: profile.Username}
		dto.Comment = comment
		dto.PostId = postId
		retVal = append(retVal, dto)
	}

	return retVal, nil
}

func (c commentRepository) CommentPost(comment dto.CommentDTO, ctx context.Context) error {
	mentions := make([]string, 1)
	id := uuid.NewString()



	err := c.cassandraSession.Query(InsertComment, id, comment.Comment, comment.PostId, comment.CommentBy.ID, time.Now(), mentions).Exec()

	if err != nil {
		return err
	}

	var numOfCom int

	c.cassandraSession.Query(GetNumOfCommentsForPost, comment.PostId, comment.PostBy).Iter().Scan(&numOfCom)
	numOfCom = numOfCom + 1
	c.cassandraSession.Query(AddCommentToPost, numOfCom, comment.PostId, comment.PostBy).Exec()
	return nil
}

func (c commentRepository) DeleteComment(comment dto.CommentDTO, ctx context.Context) error {
	err := c.cassandraSession.Query(DeleteComment, comment.PostId, comment.CommentBy).Exec()
	if err != nil {
		return err
	}

	return nil
}

func NewCommentRepository(cassandraSession *gocql.Session) CommentRepo {
	var c = &commentRepository{
		cassandraSession: cassandraSession,
	}
	err := c.cassandraSession.Query(CreateCommentTable).Exec()
	if err != nil {
		return nil
	}
	return c
}
