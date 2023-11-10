package apis

import (
	"github.com/Kengathua/book-inventory-system/pkg/common/middlewares"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterInventoryRoutes(url fiber.Router, db *gorm.DB) {
	h := &Handler{
		DB: db,
	}

	userRoutes := url.Group("/users")
	userRoutes.Get("/", h.GetUsers)
	userRoutes.Get("/:id", h.GetUser)

	authorRoutes := url.Group("/authors")
	authorRoutes.Get("/", h.GetAuthors)
	authorRoutes.Post("/", h.AddAuthor)

	librarianRoutes := url.Group("/librarians")
	librarianRoutes.Get("/", h.GetLibrarians)
	librarianRoutes.Post("/", h.AddLibrarian)
	librarianRoutes.Post("/:id/review_book", h.LibrarianReviewBook)

	storeKeeperRoutes := url.Group("/store_keepers")
	storeKeeperRoutes.Get("/", h.GetStoreKeepers)
	storeKeeperRoutes.Post("/", h.AddStoreKeeper)
	storeKeeperRoutes.Post("/:id/review_book", h.StoreKeeperReviewBook)

	bookRoutes := url.Group("/books")
	bookRoutes.Get("/", h.GetBooks)
	bookRoutes.Get(":id", h.GetBook)
	bookRoutes.Post("/", h.AddBook)
	bookRoutes.Put(":id", h.UpdateBook)
	bookRoutes.Delete(":id", h.DeleteBook)
}

func RegisterV1Routes(url fiber.Router, db *gorm.DB) {
	inventoryURL := url.Group("/inventory", func(c *fiber.Ctx) error { // middleware for /api/v1/inventory
		c.Set("Version", "v1")
		return c.Next()
	})

	inventoryURL.Use(middlewares.AuthRequiredMiddleware(db))

	RegisterInventoryRoutes(inventoryURL, db)
}

func RegisterAPIsRoutes(url fiber.Router, db *gorm.DB) {
	v1 := url.Group("/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
		c.Set("Version", "v1")
		return c.Next()
	})

	RegisterV1Routes(v1, db)
}

func RegisterAuthRoutes(url fiber.Router, db *gorm.DB) {
	h := &Handler{
		DB: db,
	}

	auth := url.Group("/auth", func(c *fiber.Ctx) error { // middleware for /api/auth
		c.Set("Version", "auth")
		return c.Next()
	})

	auth.Post("/token", h.GetToken)

	userRoutes := auth.Group("/users")
	userRoutes.Post("/register", h.AddUser) // /apis/auth/users/register

}
