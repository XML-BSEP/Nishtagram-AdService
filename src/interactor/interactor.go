package interactor

import (
	"ad_service/http/handler"
	"ad_service/repository"
	"ad_service/usecase"
	"github.com/gocql/gocql"
)
type appHandler struct {
	handler.AdHandler
	handler.CampaignHandler
	handler.CampaignRequestHandler
	handler.AdvertiseHandler
	handler.LikeHandler
	handler.CommentHandler
}
type Interactor interface {
	NewAdPostRepo() repository.AdPostRepo
	NewCampaignRepo() repository.CampaignRepo
	NewCampaignRequestRepo() repository.CampaignRequestRepo
	NewAdvertisementRepo() repository.AdvertisementRepo
	NewCommentRepo() repository.CommentRepo
	NewLikeRepo() repository.LikeRepo


	NewAdPostUseCase() usecase.AdPostUseCase
	NewCampaignUseCase() usecase.CampaignUseCase
	NewCampaignRequestUseCase() usecase.CampaignRequestUseCase
	NewAdvertisementUseCase() usecase.AdvertiseUseCase
	NewLikeUseCase() usecase.LikeUseCase
	NewCommentUseCase() usecase.CommentUseCase

	NewAdPostHandler() handler.AdHandler
	NewCampaignHandler() handler.CampaignHandler
	NewCampaignRequestHandler() handler.CampaignRequestHandler
	NewAdvertisementHandler() handler.AdvertiseHandler
	NewLikeHandler() handler.LikeHandler
	NewCommentHandler() handler.CommentHandler

	NewAppHandler() handler.AppHandler
}

type interctor struct {
	cassandraClient *gocql.Session
}

func (i interctor) NewCommentRepo() repository.CommentRepo {
	return repository.NewCommentRepository(i.cassandraClient)
}

func (i interctor) NewLikeRepo() repository.LikeRepo {
	return repository.NewLikeRepository(i.cassandraClient)
}

func (i interctor) NewLikeUseCase() usecase.LikeUseCase {
	return usecase.NewLikeUseCase(i.NewLikeRepo())
}

func (i interctor) NewCommentUseCase() usecase.CommentUseCase {
	return usecase.NewCommentUseCase(i.NewCommentRepo())
}

func (i interctor) NewLikeHandler() handler.LikeHandler {
	return handler.NewLikeHandler(i.NewLikeUseCase())
}

func (i interctor) NewCommentHandler() handler.CommentHandler {
	return handler.NewCommentHandler(i.NewCommentUseCase())
}

func (i interctor) NewAdvertisementRepo() repository.AdvertisementRepo {
	return repository.NewAdvertisementRepository(i.cassandraClient)
}

func (i interctor) NewAdvertisementUseCase() usecase.AdvertiseUseCase {
	return usecase.NewAdvertiseUseCase(i.NewAdPostUseCase(), i.NewAdvertisementRepo(), i.NewLikeRepo())
}

func (i interctor) NewAdvertisementHandler() handler.AdvertiseHandler {
	return handler.NewAdvertiseHandler(i.NewAdvertisementUseCase())
}

func (i interctor) NewCampaignRequestRepo() repository.CampaignRequestRepo {
	return repository.NewCampaignRequestRepository(i.cassandraClient)
}

func (i interctor) NewCampaignRequestUseCase() usecase.CampaignRequestUseCase {
	return usecase.NewCampaignRequestUseCase(i.NewCampaignRequestRepo(), i.NewCampaignUseCase())
}

func (i interctor) NewCampaignRequestHandler() handler.CampaignRequestHandler {
	return handler.NewCampaignRequestHandler(i.NewCampaignRequestUseCase())
}

func (i interctor) NewAdPostRepo() repository.AdPostRepo {
	return repository.NewAdPostRepo(i.cassandraClient)
}

func (i interctor) NewCampaignRepo() repository.CampaignRepo {
	return repository.NewCampaignRepo(i.cassandraClient)
}

func (i interctor) NewAdPostUseCase() usecase.AdPostUseCase {
	return usecase.NewAdPostUseCase(i.NewAdPostRepo())
}

func (i interctor) NewCampaignUseCase() usecase.CampaignUseCase {
	return usecase.NewCampaignUseCase(i.NewCampaignRepo(), i.NewAdPostUseCase(), i.NewAdvertisementUseCase())
}

func (i interctor) NewAdPostHandler() handler.AdHandler {
	return handler.NewAdHandler(i.NewAdPostUseCase())
}

func (i interctor) NewCampaignHandler() handler.CampaignHandler {
	return handler.NewCampaignHandler(i.NewCampaignUseCase())
}

func (i interctor) NewAppHandler() handler.AppHandler {
	appHandler := &appHandler{}
	appHandler.AdHandler = i.NewAdPostHandler()
	appHandler.CampaignHandler = i.NewCampaignHandler()
	appHandler.CampaignRequestHandler = i.NewCampaignRequestHandler()
	appHandler.AdvertiseHandler = i.NewAdvertisementHandler()
	appHandler.LikeHandler = i.NewLikeHandler()
	appHandler.CommentHandler = i.NewCommentHandler()

	return appHandler
}

func NewInteractor(cassandraClient *gocql.Session) Interactor {
	return &interctor{cassandraClient: cassandraClient}
}
