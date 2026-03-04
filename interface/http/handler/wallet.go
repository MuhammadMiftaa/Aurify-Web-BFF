package handler

import (
	"context"
	"fmt"

	logger "refina-web-bff/config/log"
	grpcClient "refina-web-bff/interface/grpc/client"
	"refina-web-bff/internal/types/dto"
	"refina-web-bff/internal/utils/data"

	wpb "github.com/MuhammadMiftaa/Refina-Protobuf/wallet"
	"github.com/gofiber/fiber/v2"
)

type walletHandler struct {
	wallet grpcClient.WalletClient
}

func NewWalletHandler(wc grpcClient.WalletClient) *walletHandler {
	return &walletHandler{
		wallet: wc,
	}
}

// GetUserWallets — GET /wallets
func (h *walletHandler) GetUserWallets(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(dto.UserData)
	requestID, _ := c.Locals(data.REQUEST_ID_LOCAL_KEY).(string)

	result, err := h.wallet.GetUserWallets(context.Background(), userData.ID)
	if err != nil {
		logger.Error(data.LogGetWalletsFailed, map[string]any{
			"service":    data.WalletService,
			"request_id": requestID,
			"user_id":    userData.ID,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 500,
			Message:    "Failed to get user wallets",
		})
	}

	return c.JSON(dto.APIResponse{
		Status:     true,
		StatusCode: 200,
		Message:    "User wallets retrieved successfully",
		Data:       result.GetWallets(),
	})
}

// GetWalletSummary — GET /wallets/summary
func (h *walletHandler) GetWalletSummary(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(dto.UserData)
	requestID, _ := c.Locals(data.REQUEST_ID_LOCAL_KEY).(string)

	result, err := h.wallet.GetWalletSummary(context.Background(), userData.ID)
	if err != nil {
		logger.Error(data.LogGetWalletSummaryFailed, map[string]any{
			"service":    data.WalletService,
			"request_id": requestID,
			"user_id":    userData.ID,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 500,
			Message:    "Failed to get wallet summary",
		})
	}

	return c.JSON(dto.APIResponse{
		Status:     true,
		StatusCode: 200,
		Message:    "Wallet summary retrieved successfully",
		Data:       result,
	})
}

// GetWalletByID — GET /wallets/:id
func (h *walletHandler) GetWalletByID(c *fiber.Ctx) error {
	walletID := c.Params("id")
	requestID, _ := c.Locals(data.REQUEST_ID_LOCAL_KEY).(string)

	if walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 400,
			Message:    "Wallet ID is required",
		})
	}

	result, err := h.wallet.GetWalletByID(context.Background(), walletID)
	if err != nil {
		logger.Error(data.LogGetWalletByIDFailed, map[string]any{
			"service":    data.WalletService,
			"request_id": requestID,
			"wallet_id":  walletID,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 500,
			Message:    "Failed to get wallet",
		})
	}

	return c.JSON(dto.APIResponse{
		Status:     true,
		StatusCode: 200,
		Message:    "Wallet retrieved successfully",
		Data:       result,
	})
}

// CreateWallet — POST /wallets
func (h *walletHandler) CreateWallet(c *fiber.Ctx) error {
	userData := c.Locals("user_data").(dto.UserData)
	requestID, _ := c.Locals(data.REQUEST_ID_LOCAL_KEY).(string)

	var req dto.CreateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 400,
			Message:    "Invalid request body",
		})
	}

	if req.Name == "" || req.WalletTypeID == "" || req.Number == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 400,
			Message:    "Name, wallet_type_id, and number are required",
		})
	}

	grpcReq := &wpb.CreateWalletRequest{
		UserId:       userData.ID,
		WalletTypeId: req.WalletTypeID,
		Name:         req.Name,
		Balance:      req.Balance,
		Number:       req.Number,
	}

	result, err := h.wallet.CreateWallet(context.Background(), grpcReq)
	if err != nil {
		logger.Error(data.LogCreateWalletFailed, map[string]any{
			"service":    data.WalletService,
			"request_id": requestID,
			"user_id":    userData.ID,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 500,
			Message:    fmt.Sprintf("Failed to create wallet: %s", err.Error()),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status:     true,
		StatusCode: 201,
		Message:    "Wallet created successfully",
		Data:       result,
	})
}

// UpdateWallet — PUT /wallets/:id
func (h *walletHandler) UpdateWallet(c *fiber.Ctx) error {
	walletID := c.Params("id")
	requestID, _ := c.Locals(data.REQUEST_ID_LOCAL_KEY).(string)

	if walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 400,
			Message:    "Wallet ID is required",
		})
	}

	var req dto.UpdateWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 400,
			Message:    "Invalid request body",
		})
	}

	grpcReq := &wpb.UpdateWalletRequest{
		Id:           walletID,
		Name:         req.Name,
		Number:       req.Number,
		WalletTypeId: req.WalletTypeID,
	}

	result, err := h.wallet.UpdateWallet(context.Background(), grpcReq)
	if err != nil {
		logger.Error(data.LogUpdateWalletFailed, map[string]any{
			"service":    data.WalletService,
			"request_id": requestID,
			"wallet_id":  walletID,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 500,
			Message:    fmt.Sprintf("Failed to update wallet: %s", err.Error()),
		})
	}

	return c.JSON(dto.APIResponse{
		Status:     true,
		StatusCode: 200,
		Message:    "Wallet updated successfully",
		Data:       result,
	})
}

// DeleteWallet — DELETE /wallets/:id
func (h *walletHandler) DeleteWallet(c *fiber.Ctx) error {
	walletID := c.Params("id")
	requestID, _ := c.Locals(data.REQUEST_ID_LOCAL_KEY).(string)

	if walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 400,
			Message:    "Wallet ID is required",
		})
	}

	_, err := h.wallet.DeleteWallet(context.Background(), walletID)
	if err != nil {
		logger.Error(data.LogDeleteWalletFailed, map[string]any{
			"service":    data.WalletService,
			"request_id": requestID,
			"wallet_id":  walletID,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 500,
			Message:    fmt.Sprintf("Failed to delete wallet: %s", err.Error()),
		})
	}

	return c.JSON(dto.APIResponse{
		Status:     true,
		StatusCode: 200,
		Message:    "Wallet deleted successfully",
	})
}

// GetWalletTypes — GET /wallet-types
func (h *walletHandler) GetWalletTypes(c *fiber.Ctx) error {
	requestID, _ := c.Locals(data.REQUEST_ID_LOCAL_KEY).(string)

	result, err := h.wallet.GetWalletTypes(context.Background())
	if err != nil {
		logger.Error(data.LogGetWalletTypesFailed, map[string]any{
			"service":    data.WalletService,
			"request_id": requestID,
			"error":      err.Error(),
		})
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:     false,
			StatusCode: 500,
			Message:    "Failed to get wallet types",
		})
	}

	return c.JSON(dto.APIResponse{
		Status:     true,
		StatusCode: 200,
		Message:    "Wallet types retrieved successfully",
		Data:       result.GetWalletTypes(),
	})
}
