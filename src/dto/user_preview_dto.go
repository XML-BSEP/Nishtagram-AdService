package dto

type UserPreviewDTO struct {
	Username string
	Name string
	Surname string
	ProfilePicture string
	Id string
}

func NewUserPreviewDTO(id string) UserPreviewDTO {
	return UserPreviewDTO{Id: id}
}
