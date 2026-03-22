package service

import (
	"github.com/dawnstack/shop-go/internal/cache"
	"github.com/dawnstack/shop-go/internal/config"
	"github.com/dawnstack/shop-go/internal/pkg/auth"
	"github.com/dawnstack/shop-go/internal/pkg/snowflake"
	"github.com/dawnstack/shop-go/internal/repository"
)

type Services struct {
	Auth    *AuthService
	User    *UserService
	Address *AddressService
	Product *ProductService
	Cart    *CartService
	Order   *OrderService
	Home    *HomeService
	Video   *VideoService
}

func NewServices(
	repos *repository.Repositories,
	txManager *repository.TxManager,
	productCache *cache.ProductCache,
	homeCache *cache.HomeCache,
	cfg config.Config,
) *Services {
	jwtManager := auth.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.AccessExpireMin,
		cfg.JWT.RefreshExpireHour,
	)
	idGen := snowflake.New(1)

	return &Services{
		Auth:    NewAuthService(repos.Users, idGen, jwtManager),
		User:    NewUserService(repos.Users),
		Address: NewAddressService(repos.Addresses, idGen),
		Product: NewProductService(repos.Products, productCache),
		Cart:    NewCartService(repos.Carts, repos.Products, idGen),
		Order:   NewOrderService(repos.Orders, repos.Carts, repos.Products, txManager, idGen),
		Home:    NewHomeService(repos.Homes, homeCache),
		Video:   NewVideoService(repos.Homes, homeCache),
	}
}
