package services

import (
	"backend/internal/datasources/records"
	"backend/internal/datasources/repositories"
	"backend/internal/http/data_transfers"
	"backend/pkg/convert"
	"errors"
	"github.com/jinzhu/copier"
	"net/http"
)

type NutritionsRepository interface {
	FindAllByOwnerID(ownerID int) ([]records.Nutritions, error)
	FindByID(id int) (records.Nutritions, error)
	Save(nutrition records.Nutritions) (int, error)
	Update(id int, nutrition map[string]interface{}) error
	Delete(id int) error
}

type NutritionsService struct {
	repository NutritionsRepository
}

func NewNutritionsService(repository NutritionsRepository) *NutritionsService {
	return &NutritionsService{repository}
}

func (n *NutritionsService) FindAllByOwnerID(ownerID int) ([]data_transfers.NutritionsResponse, int, error) {
	var nutritionsResponse []data_transfers.NutritionsResponse

	nutritions, err := n.repository.FindAllByOwnerID(ownerID)
	if err != nil {
		return nutritionsResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&nutritionsResponse, &nutritions)
	if err != nil {
		return nutritionsResponse, http.StatusInternalServerError, err
	}

	return nutritionsResponse, http.StatusOK, nil
}

func (n *NutritionsService) FindByID(id int) (data_transfers.NutritionsResponse, int, error) {
	var nutritionResponse data_transfers.NutritionsResponse

	nutrition, err := n.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrorRowNotFound) {
			return nutritionResponse, http.StatusNotFound, err
		}

		return nutritionResponse, http.StatusInternalServerError, err
	}

	err = copier.Copy(&nutritionResponse, &nutrition)
	if err != nil {
		return nutritionResponse, http.StatusInternalServerError, err
	}

	return nutritionResponse, http.StatusOK, nil
}

func (n *NutritionsService) Save(nutritionRequest data_transfers.CreateNutritionsRequest) (int, int, error) {
	var nutrition records.Nutritions

	err := copier.Copy(&nutrition, &nutritionRequest)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	id, err := n.repository.Save(nutrition)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return id, http.StatusCreated, nil
}

func (n *NutritionsService) Update(id int, updateNutritionRequest data_transfers.UpdateNutritionsRequest) (int, error) {
	nutritionMap, err := convert.StructToMap(updateNutritionRequest)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = n.repository.Update(id, nutritionMap)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (n *NutritionsService) Delete(id int) error {
	err := n.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
