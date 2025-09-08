package provider

import (
	"github.com/lookandhqte/workHelper/internal/provider/deepseek"
	"github.com/lookandhqte/workHelper/internal/provider/hh"
)

type Provider struct {
	HH       hh.Provider
	DeepSeek deepseek.Provider
}

// New создает нового провайдера
func New() *Provider {
	return &Provider{
		HH:       *hh.New(),
		DeepSeek: *deepseek.New(),
	}
}
