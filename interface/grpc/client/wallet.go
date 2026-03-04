package client

import (
	"context"
	"time"

	wpb "github.com/MuhammadMiftaa/Refina-Protobuf/wallet"
)

type WalletClient interface {
	GetUserWallets(ctx context.Context, userID string) (*wpb.GetUserWalletsResponse, error)
	GetWalletByID(ctx context.Context, walletID string) (*wpb.Wallet, error)
	CreateWallet(ctx context.Context, req *wpb.CreateWalletRequest) (*wpb.Wallet, error)
	UpdateWallet(ctx context.Context, req *wpb.UpdateWalletRequest) (*wpb.Wallet, error)
	DeleteWallet(ctx context.Context, walletID string) (*wpb.Wallet, error)
	GetWalletTypes(ctx context.Context) (*wpb.GetWalletTypesResponse, error)
	GetWalletSummary(ctx context.Context, userID string) (*wpb.WalletSummary, error)
}

type walletClientImpl struct {
	client wpb.WalletServiceClient
}

func NewWalletClient(grpcClient wpb.WalletServiceClient) WalletClient {
	return &walletClientImpl{
		client: grpcClient,
	}
}

func (w *walletClientImpl) GetUserWallets(ctx context.Context, userID string) (*wpb.GetUserWalletsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	return w.client.GetUserWallets(ctx, &wpb.UserID{Id: userID})
}

func (w *walletClientImpl) GetWalletByID(ctx context.Context, walletID string) (*wpb.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return w.client.GetWalletByID(ctx, &wpb.WalletID{Id: walletID})
}

func (w *walletClientImpl) CreateWallet(ctx context.Context, req *wpb.CreateWalletRequest) (*wpb.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	return w.client.CreateWallet(ctx, req)
}

func (w *walletClientImpl) UpdateWallet(ctx context.Context, req *wpb.UpdateWalletRequest) (*wpb.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return w.client.UpdateWallet(ctx, req)
}

func (w *walletClientImpl) DeleteWallet(ctx context.Context, walletID string) (*wpb.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return w.client.DeleteWallet(ctx, &wpb.WalletID{Id: walletID})
}

func (w *walletClientImpl) GetWalletTypes(ctx context.Context) (*wpb.GetWalletTypesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return w.client.GetWalletTypes(ctx, &wpb.Empty{})
}

func (w *walletClientImpl) GetWalletSummary(ctx context.Context, userID string) (*wpb.WalletSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return w.client.GetWalletSummary(ctx, &wpb.UserID{Id: userID})
}
