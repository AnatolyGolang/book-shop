package handlers

import "book-shop/internal/app/services"

type HttpServer struct {
	bookService     services.BookService
	categoryService services.CategoryService
	userService     services.UserService
	cartService     services.CartService
	jwtService      services.JWTService
}

// NewHttpServer creates a new HTTP server for ports
func NewHttpServer(bs services.BookService,
	cs services.CategoryService,
	us services.UserService,
	carts services.CartService,
	jwts services.JWTService) HttpServer {
	return HttpServer{
		bookService:     bs,
		categoryService: cs,
		userService:     us,
		cartService:     carts,
		jwtService:      jwts,
	}
}
