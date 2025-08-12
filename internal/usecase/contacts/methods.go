package contacts

import "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"

//ReturnAll возвращает все контакты из хранилища
func (uc *UseCase) ReturnAll() ([]entity.GlobalContact, error) {
	return uc.repo.GetAllGlobalContacts()
}

//Update обновляет контакт в хранилище
func (uc *UseCase) Update(contacts []entity.GlobalContact) error {
	return uc.repo.UpdateGlobalContacts(contacts)
}

//Delete удаляет контакт в хранилище
func (uc *UseCase) Delete(id int) error {
	return uc.repo.DeleteAccountContacts(id)
}

func (uc *UseCase) Create(contact *entity.GlobalContact) error {
	return uc.repo.AddContact(contact)
}
