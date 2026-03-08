package routes

import (
	grpcClient "refina-web-bff/interface/grpc/client"
	"refina-web-bff/interface/http/handler"
	"refina-web-bff/interface/http/middleware"
	"refina-web-bff/internal/cache"

	"github.com/gofiber/fiber/v2"
)

func InvestmentRoutes(app *fiber.App, ic grpcClient.InvestmentClient, c cache.Cache) {
	h := handler.NewInvestmentHandler(ic, c)

	investments := app.Group("/investments")
	investments.Use(middleware.AuthMiddleware())

	investments.Get("/", h.GetUserInvestments)
	investments.Get("/summary", h.GetInvestmentSummary)
	investments.Get("/asset-codes", h.GetAssetCodes)
	investments.Post("/", h.CreateInvestment)
	investments.Post("/sell", h.SellInvestment)
	investments.Get("/:id", h.GetInvestmentDetail)
}
