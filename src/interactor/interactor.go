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
}
type Interactor interface {
	NewAdPostRepo() repository.AdPostRepo
	NewCampaignRepo() repository.CampaignRepo


	NewAdPostUseCase() usecase.AdPostUseCase
	NewCampaignUseCase() usecase.CampaignUseCase

	NewAdPostHandler() handler.AdHandler
	NewCampaignHandler() handler.CampaignHandler

	NewAppHandler() handler.AppHandler


}

type interctor struct {
	cassandraClient *gocql.Session
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
	return usecase.NewCampaignUseCase(i.NewCampaignRepo(), i.NewAdPostUseCase())
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
	return appHandler
}

func NewInteractor(cassandraClient *gocql.Session) Interactor {
	return &interctor{cassandraClient: cassandraClient}
}
