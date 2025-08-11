package provider

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider/amocrm"
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider/unisender"
)

type Provider struct {
	Amo amocrm.Provider
	Uni unisender.Provider
}

// New создает нового провайдера
func New() *Provider {
	return &Provider{
		Amo: *amocrm.New(),
		Uni: *unisender.New(),
	}
}
