package item

import (
	"github.com/Yuriekokubu/workflow/internal/constant"
	"github.com/Yuriekokubu/workflow/internal/model"

	"gorm.io/gorm"
)

type Service struct {
	Repository Repository
}

func NewService(db *gorm.DB) Service {
	return Service{
		Repository: NewRepository(db),
	}
}

func (service Service) Create(req model.RequestItem) (model.Item, error) {
	item := model.Item{
		Title:    req.Title,
		Amount:   req.Amount,
		Quantity: req.Quantity,
		Status:   constant.ItemPendingStatus,
		OwnerID:  req.OwnerID,
	}

	if err := service.Repository.Create(&item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (service Service) Find(query model.RequestFindItem) ([]model.Item, error) {
	return service.Repository.Find(query)
}

func (service Service) UpdateStatus(id uint, status constant.ItemStatus) (model.Item, error) {
	item, err := service.Repository.FindByID(id)
	if err != nil {
		return model.Item{}, err
	}

	item.Status = status

	if err := service.Repository.Replace(item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (service Service) FindByID(id uint) (model.Item, error) {
	item, err := service.Repository.FindByID(id)
	if err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (service Service) UpdateItemByID(id uint, req model.RequestItem) (model.Item, error) {
	item, err := service.Repository.FindByID(id)
	if err != nil {
		return model.Item{}, err
	}

	item.Title = req.Title
	item.Amount = req.Amount
	item.Quantity = req.Quantity

	item.Status = constant.ItemPendingStatus

	if err := service.Repository.Update(&item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (service Service) UpdateItemStatusByID(id uint, status constant.ItemStatus) (model.Item, error) {
	item, err := service.Repository.FindByID(id)
	if err != nil {
		return model.Item{}, err
	}

	item.Status = status

	if err := service.Repository.Update(&item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (service Service) DeleteItemByID(id uint) error {
	item, err := service.Repository.FindByID(id)
	if err != nil {
		return err
	}

	if err := service.Repository.Delete(&item); err != nil {
		return err
	}

	return nil
}