package usecase

import (
	"ad_service/domain"
	"ad_service/repository"
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"strings"
)

type AdPostUseCase interface {
	CreateAdPost(ctx context.Context, post domain.AdPost) error
	GetAdsByAgent(ctx context.Context, agentId string) ([]domain.AdPost, error)
	GetAdById(ctx context.Context, agentId string, id string) (domain.AdPost, error)
	EncodeBase64(media string, agentId string, ctx context.Context) (string, error)
	DecodeBase64(media string, agentId string, ctx context.Context) (string, error)

}

type adPostUseCase struct {
	adPostRepository repository.AdPostRepo
}

func (a adPostUseCase) EncodeBase64(media string, agentId string, ctx context.Context) (string, error) {
	workingDirectory, _ := os.Getwd()
	if !strings.HasSuffix(workingDirectory, "src") {
		firstPart := strings.Split(workingDirectory, "src")
		value := firstPart[0] + "/src"
		workingDirectory = value
		os.Chdir(workingDirectory)

	}
	path1 := "./assets/images/"
	err := os.Chdir(path1)

	err = os.Mkdir(agentId, 0755)
	fmt.Println(err)

	err = os.Chdir(agentId)
	fmt.Println(err)


	s := strings.Split(media, ",")
	part := strings.Split(s[0], "/")
	format := strings.Split(part[1], ";")

	dec, err := base64.StdEncoding.DecodeString(s[1])

	uuid := uuid.NewString()
	f, err := os.Create(uuid + "." + format[0])


	defer f.Close()

	if _, err := f.Write(dec); err != nil {

	}
	if err := f.Sync(); err != nil {

	}

	os.Chdir(workingDirectory)
	return agentId + "/" + uuid + "." + format[0], nil
}

func (a adPostUseCase) DecodeBase64(media string, agentId string, ctx context.Context) (string, error) {
	workingDirectory, _ := os.Getwd()
	if !strings.HasSuffix(workingDirectory, "src") {
		firstPart := strings.Split(workingDirectory, "src")
		value := firstPart[0] + "/src"
		workingDirectory = value
		os.Chdir(workingDirectory)
	}

	path1 := "./assets/images/"
	err := os.Chdir(path1)
	fmt.Println(err)
	spliced := strings.Split(media, "/")
	var f *os.File
	if len(spliced) > 1 {
		err = os.Chdir(agentId)
		f, _ = os.Open(spliced[1])
	} else {
		f, _ = os.Open(spliced[0])
	}

	defer f.Close()
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	encoded := base64.StdEncoding.EncodeToString(content)

	fmt.Println("ENCODED: " + encoded)
	os.Chdir(workingDirectory)

	return "data:image/jpg;base64," + encoded, nil
}

func (a adPostUseCase) CreateAdPost(ctx context.Context, post domain.AdPost) error {
	encoded, err := a.EncodeBase64(post.Path, post.AgentId.ID, ctx)
	if err != nil {
		return err
	}

	post.Path = encoded
	return a.adPostRepository.CreateAd(context.Background(), post)
}

func (a adPostUseCase) GetAdsByAgent(ctx context.Context, agentId string) ([]domain.AdPost, error) {

	ads, err := a.adPostRepository.GetAdsByAgent(context.Background(), agentId)

	if err != nil {
		return nil, err
	}
	var retVal []domain.AdPost

	for _, ad := range ads {
		ad.Path, _ = a.DecodeBase64(ad.Path, agentId, ctx)
		retVal = append(retVal, ad)
	}
	return retVal, nil
}

func (a adPostUseCase) GetAdById(ctx context.Context, agentId string, id string) (domain.AdPost, error) {
	ad, err := a.adPostRepository.GetAdByAgentIdAndId(ctx, agentId, id)
	if err != nil {
		return domain.AdPost{}, err
	}
	ad.Path, err = a.DecodeBase64(ad.Path, agentId, ctx)
	return ad, err
}

func NewAdPostUseCase(adPostRepo repository.AdPostRepo) AdPostUseCase {
	return &adPostUseCase{adPostRepository: adPostRepo}
}