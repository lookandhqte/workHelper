package inmemory

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
)

// GetAllGlobalContacts возвращает все контакты из таблицы global_contacts
func (m *MemoryStorage) GetAllGlobalContacts() ([]entity.GlobalContact, error) {
	return nil, nil
}

// UpdateGlobalContacts обновляет существующие и добавляет новые контакты
func (m *MemoryStorage) UpdateGlobalContacts(contacts []entity.GlobalContact) error {
	return nil
}

// DeleteAccountContacts удаляет все контакты
func (dm *MemoryStorage) DeleteAccountContacts(accountID int) error {
	return nil
}
