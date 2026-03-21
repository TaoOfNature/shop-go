package model

import "time"

type User struct {
	ID           int64      `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Nickname     string     `json:"nickname" db:"nickname"`
	AvatarURL    string     `json:"avatar_url" db:"avatar_url"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Address struct {
	ID           int64      `json:"id" db:"id"`
	UserID       int64      `json:"user_id" db:"user_id"`
	ReceiverName string     `json:"receiver_name" db:"receiver_name"`
	Phone        string     `json:"phone" db:"phone"`
	Province     string     `json:"province" db:"province"`
	City         string     `json:"city" db:"city"`
	District     string     `json:"district" db:"district"`
	Detail       string     `json:"detail" db:"detail"`
	IsDefault    bool       `json:"is_default" db:"is_default"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Category struct {
	ID        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	IconURL   string     `json:"icon_url" db:"icon_url"`
	SortOrder int        `json:"sort_order" db:"sort_order"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Product struct {
	ID          int64      `json:"id" db:"id"`
	CategoryID  int64      `json:"category_id" db:"category_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Price       int64      `json:"price" db:"price"`
	Stock       int        `json:"stock" db:"stock"`
	CoverURL    string     `json:"cover_url" db:"cover_url"`
	Sales       int        `json:"sales" db:"sales"`
	Tags        []string   `json:"tags" db:"-"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CartItem struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Checked   bool      `json:"checked" db:"checked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	ID             int64       `json:"id" db:"id"`
	UserID         int64       `json:"user_id" db:"user_id"`
	AddressID      int64       `json:"address_id" db:"address_id"`
	OrderNo        string      `json:"order_no" db:"order_no"`
	Status         string      `json:"status" db:"status"`
	TotalAmount    int64       `json:"total_amount" db:"total_amount"`
	IdempotencyKey string      `json:"idempotency_key" db:"idempotency_key"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
	Items          []OrderItem `json:"items"`
}

type OrderItem struct {
	ID           int64  `json:"id" db:"id"`
	OrderID      int64  `json:"order_id" db:"order_id"`
	ProductID    int64  `json:"product_id" db:"product_id"`
	ProductName  string `json:"product_name" db:"product_name"`
	ProductPrice int64  `json:"product_price" db:"product_price"`
	Quantity     int    `json:"quantity" db:"quantity"`
}

type Banner struct {
	ID       int64  `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	ImageURL string `json:"image_url" db:"image_url"`
	LinkURL  string `json:"link_url" db:"link_url"`
}

type Video struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	CoverURL    string `json:"cover_url" db:"cover_url"`
	PlaybackURL string `json:"playback_url" db:"playback_url"`
	ProductID   int64  `json:"product_id" db:"product_id"`
}
