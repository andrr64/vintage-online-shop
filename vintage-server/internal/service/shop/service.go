package shop

import (
	"context"
	"vintage-server/internal/domain/shop"
	"vintage-server/internal/model"
	"vintage-server/internal/repository"
	db_error "vintage-server/pkg"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/controller"
	"vintage-server/pkg/uploader"

	"github.com/google/uuid"
)

type shopService struct {
	store    repository.ShopStore
	jwt      auth.JWTService
	uploader uploader.Uploader
}

func NewShopService(
	store repository.ShopStore,
	jwt auth.JWTService,
	uploader uploader.Uploader,
) shop.ShopService {
	return &shopService{
		store:    store,
		jwt:      jwt,
		uploader: uploader,
	}
}

// CreateShop implements shop.ShopService.
func (s *shopService) CreateShop(ctx context.Context, accountID uuid.UUID, req shop.ShopRequest) (shop.ShopDetail, error) {
	newShop := model.Shop{
		AccountID:   accountID,
		Name:        req.Name,
		Description: &req.Description,
		Summary:     &req.Summary,
		Active:      true,
	}
	var savedShop model.Shop

	ctx, cancel := controller.WithTxTimeout(ctx)
	defer cancel()

	err := s.store.ExecTx(ctx, func(sr shop.ShopRepository) error {
		qctx, qcancel := controller.WithQueryTimeout(ctx)
		defer qcancel()
		roleID, err := sr.GetRoleIDByName(qctx, "seller")
		if err != nil {
			return err
		}

		qctx, qcancel = controller.WithQueryTimeout(ctx)
		defer qcancel()
		savedShop, err = sr.CreateShop(qctx, newShop)
		if err != nil {
			return err
		}

		return sr.InsertAccountRole(qctx, accountID, roleID)
	})
	if err != nil {
		return shop.ShopDetail{}, db_error.HandlePgError(err)
	}

	return shop.ConvertShopModelToShopDetail(savedShop), nil
}

func (s *shopService) UpdateShop(ctx context.Context, accountID uuid.UUID, shopID uuid.UUID, req shop.ShopRequest) (shop.ShopDetail, error) {

	return shop.ShopDetail{}, nil
}
