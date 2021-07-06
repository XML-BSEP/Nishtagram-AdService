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
}
type Interactor interface {
	NewAdPostRepo() repository.AdPostRepo
	NewCampaignRepo() repository.CampaignRepo
	NewCampaignRequestRepo() repository.CampaignRequestRepo
	NewAdvertisementRepo() repository.AdvertisementRepo


	NewAdPostUseCase() usecase.AdPostUseCase
	NewCampaignUseCase() usecase.CampaignUseCase
	NewCampaignRequestUseCase() usecase.CampaignRequestUseCase
	NewAdvertisementUseCase() usecase.AdvertiseUseCase

	NewAdPostHandler() handler.AdHandler
	NewCampaignHandler() handler.CampaignHandler
	NewCampaignRequestHandler() handler.CampaignRequestHandler
	NewAdvertisementHandler() handler.AdvertiseHandler

	NewAppHandler() handler.AppHandler
}

type interctor struct {
	cassandraClient *gocql.Session
}

func (i interctor) NewAdvertisementRepo() repository.AdvertisementRepo {
	return repository.NewAdvertisementRepository(i.cassandraClient)
}

func (i interctor) NewAdvertisementUseCase() usecase.AdvertiseUseCase {
	return usecase.NewAdvertiseUseCase(i.NewAdPostUseCase(), i.NewAdvertisementRepo())
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

	return appHandler
}

func NewInteractor(cassandraClient *gocql.Session) Interactor {
	return &interctor{cassandraClient: cassandraClient}
}
