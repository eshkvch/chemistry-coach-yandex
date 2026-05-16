package usecase

import "chemistry-coach/internal/catalog"

type CatalogUseCase struct{}

func NewCatalogUseCase() *CatalogUseCase { return &CatalogUseCase{} }

func (uc *CatalogUseCase) Goals(recommendedID string) []catalog.Goal {
	return catalog.GoalsWithRecommendation(recommendedID)
}

func (uc *CatalogUseCase) Personas() []catalog.Persona {
	return catalog.Personas
}
