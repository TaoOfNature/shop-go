package repository

import "database/sql"

type Repositories struct {
	Users     *UserRepository
	Addresses *AddressRepository
	Products  *ProductRepository
	Carts     *CartRepository
	Orders    *OrderRepository
	Homes     *HomeRepository
}

func NewRepositories(db *sql.DB, txManager *TxManager) *Repositories {
	return &Repositories{
		Users:     NewUserRepository(db),
		Addresses: NewAddressRepository(db),
		Products:  NewProductRepository(db),
		Carts:     NewCartRepository(db),
		Orders:    NewOrderRepository(db),
		Homes:     NewHomeRepository(db),
	}
}
