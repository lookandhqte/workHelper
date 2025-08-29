package provider

import (
	"git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/provider/hh"
)

type Provider struct {
	HH hh.Provider
}

// New создает нового провайдера
func New() *Provider {
	return &Provider{
		HH: *hh.New(),
	}
}
