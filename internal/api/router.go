package api

import (
	"net/http"

	"github.com/dawnstack/shop-go/internal/api/handler"
	"github.com/dawnstack/shop-go/internal/api/middleware"
	"github.com/dawnstack/shop-go/internal/config"
	"github.com/dawnstack/shop-go/internal/pkg/response"
	"github.com/dawnstack/shop-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(services *service.Services, cfg config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		middleware.Metrics(),
		middleware.Tracing(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.RateLimit(cfg.App.RateLimitRPS, cfg.App.RateLimitBurst),
	)

	authHandler := handler.NewAuthHandler(services.Auth)
	userHandler := handler.NewUserHandler(services.User)
	addressHandler := handler.NewAddressHandler(services.Address)
	productHandler := handler.NewProductHandler(services.Product)
	cartHandler := handler.NewCartHandler(services.Cart)
	orderHandler := handler.NewOrderHandler(services.Order)
	homeHandler := handler.NewHomeHandler(services.Home)
	videoHandler := handler.NewVideoHandler(services.Video)
	recommendHandler := handler.NewRecommendHandler(services.Home)
	seckillHandler := handler.NewSeckillHandler(services.Seckill)

	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{"service": cfg.App.Name, "env": cfg.App.Env})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	apiGroup := r.Group("/api")
	{
		authGroup := apiGroup.Group("/auth")
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)

		apiGroup.GET("/products/search", productHandler.Search)
		apiGroup.GET("/products/:id", productHandler.Detail)
		apiGroup.GET("/home/banners", homeHandler.Banners)
		apiGroup.GET("/home/categories", homeHandler.Categories)
		apiGroup.GET("/home/recommend", homeHandler.Recommend)
		apiGroup.GET("/videos/recommend", videoHandler.Recommend)
	}

	secured := apiGroup.Group("")
	secured.Use(middleware.AuthRequired(services.Auth.TokenManager()))
	{
		secured.GET("/user/profile", userHandler.Profile)
		secured.PUT("/user/profile", userHandler.UpdateProfile)

		secured.GET("/user/addresses", addressHandler.List)
		secured.POST("/user/address", addressHandler.Create)
		secured.PUT("/user/address/:id", addressHandler.Update)
		secured.DELETE("/user/address", addressHandler.Delete)

		secured.GET("/cart", cartHandler.List)
		secured.POST("/cart", cartHandler.Add)
		secured.PUT("/cart/:id", cartHandler.Update)
		secured.DELETE("/cart", cartHandler.Delete)

		secured.GET("/orders", orderHandler.List)
		secured.GET("/orders/:id", orderHandler.Detail)
		secured.POST("/orders", orderHandler.Create)

		secured.GET("/home/recommend/personalized", recommendHandler.Personalized)
		secured.POST("/seckill/orders", seckillHandler.CreateOrder)
	}

	r.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, 4040, "route not found")
	})

	return r
}
