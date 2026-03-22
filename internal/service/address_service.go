package service

import (
	"context"

	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/pkg/snowflake"
	"github.com/dawnstack/shop-go/internal/repository"
)

type AddressService struct {
	addresses *repository.AddressRepository
	idGen     *snowflake.Generator
}

type UpsertAddressRequest struct {
	ReceiverName string `json:"receiver_name" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	District     string `json:"district" binding:"required"`
	Detail       string `json:"detail" binding:"required"`
	IsDefault    bool   `json:"is_default"`
}

func NewAddressService(addresses *repository.AddressRepository, idGen *snowflake.Generator) *AddressService {
	return &AddressService{addresses: addresses, idGen: idGen}
}

func (s *AddressService) List(ctx context.Context, userID int64) ([]model.Address, error) {
	return s.addresses.ListByUserID(ctx, userID)
}

func (s *AddressService) Create(ctx context.Context, userID int64, req UpsertAddressRequest) error {
	return s.addresses.Create(ctx, &model.Address{
		ID:           s.idGen.NextID(),
		UserID:       userID,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Detail:       req.Detail,
		IsDefault:    req.IsDefault,
	})
}

func (s *AddressService) Update(ctx context.Context, userID int64, id int64, req UpsertAddressRequest) error {
	return s.addresses.Update(ctx, &model.Address{
		ID:           id,
		UserID:       userID,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Detail:       req.Detail,
		IsDefault:    req.IsDefault,
	})
}

func (s *AddressService) Delete(ctx context.Context, userID int64, id int64) error {
	return s.addresses.Delete(ctx, userID, id)
}
