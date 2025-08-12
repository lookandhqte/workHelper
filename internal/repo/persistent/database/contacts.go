package database

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/entity"
	"gorm.io/gorm"
)

func (d *Storage) AddContact(contact *entity.GlobalContact) error {
	return d.DB.Create(contact).Error
}

// GetAllGlobalContacts возвращает все контакты из таблицы global_contacts
func (d *Storage) GetAllGlobalContacts() ([]entity.GlobalContact, error) {
	var contacts []entity.GlobalContact
	result := d.DB.Find(&contacts)
	return contacts, result.Error
}

// UpdateGlobalContacts обновляет существующие и добавляет новые контакты
func (d *Storage) UpdateGlobalContacts(contacts []entity.GlobalContact) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		for _, contact := range contacts {
			if contact.ID != 0 {
				if err := tx.Model(&entity.GlobalContact{}).
					Where("account_id = ?", contact.AccountID).
					Updates(map[string]interface{}{
						"email":  contact.Email,
						"status": contact.Status,
					}).Error; err != nil {
					return err
				}
			}
		}
		var newContacts []entity.GlobalContact
		for _, contact := range contacts {
			if contact.ID == 0 {
				newContacts = append(newContacts, contact)
			}
		}

		if len(newContacts) > 0 {
			if err := tx.Create(&newContacts).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteAccountContacts удаляет все контакты аккаунта из таблицы global_contacts
func (d *Storage) DeleteAccountContacts(accountID int) error {
	return d.DB.Where("account_id = ?", accountID).Delete(&entity.GlobalContact{}).Error
}
