package contacts

import (
	entity "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

//UseCase структура
type UseCase struct {
	repo contactRepo
}

//contactRepo абстракция для определения методов репозитория
type contactRepo interface {
	AddContact(contact *entity.GlobalContact) error
	GetAllGlobalContacts() ([]entity.GlobalContact, error)
	UpdateGlobalContacts(contacts []entity.GlobalContact) error
	DeleteAccountContacts(accountID int) error
}

//New создает новый репозиторий
func New(r contactRepo) *UseCase {
	return &UseCase{repo: r}
}
