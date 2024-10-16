package usecase_test

import (
	"context"
	"testing"

	"github.com/appxpy/hive-test/internal/entity"
	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AssetUseCaseSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	ctx  context.Context

	// Intermidiate variables
	someAsset *entity.Asset

	// Mocked units
	mockAssetRepo *MockAssetRepo

	// Tested usecase
	assetUseCase usecase.AssetUseCase
}

func (t *AssetUseCaseSuite) SetupSuite() {
	t.someAsset = &entity.Asset{
		UserID:      1,
		Name:        "Test Asset",
		Description: "A test asset",
		Price:       100.0,
	}
}

func (t *AssetUseCaseSuite) SetupTest() {
	t.ctx = context.Background()
	t.ctrl = gomock.NewController(t.T())
	t.mockAssetRepo = NewMockAssetRepo(t.ctrl)
	t.assetUseCase = usecase.NewAssetUseCase(t.mockAssetRepo)
}

func TestAssetUseCaseSuite(t *testing.T) {
	suite.Run(t, new(AssetUseCaseSuite))
}

func (t *AssetUseCaseSuite) TestAddAsset_GreenPath() {
	t.mockAssetRepo.EXPECT().CreateAsset(t.ctx, t.someAsset).Return(nil)

	err := t.assetUseCase.AddAsset(t.ctx, t.someAsset)

	t.NoError(err)
}

func (t *AssetUseCaseSuite) TestAddAsset_ReturnsError_WhenRepoReturnsError() {
	t.mockAssetRepo.EXPECT().CreateAsset(t.ctx, t.someAsset).Return(assert.AnError)

	err := t.assetUseCase.AddAsset(t.ctx, t.someAsset)

	t.ErrorIs(err, assert.AnError)
}

func (t *AssetUseCaseSuite) TestRemoveAsset_GreenPath() {
	t.mockAssetRepo.EXPECT().DeleteAsset(t.ctx, t.someAsset.ID, t.someAsset.UserID).Return(nil)

	err := t.assetUseCase.RemoveAsset(t.ctx, t.someAsset.ID, t.someAsset.UserID)

	t.NoError(err)
}

func (t *AssetUseCaseSuite) TestRemoveAsset_ReturnsError_WhenRepoReturnsError() {
	t.mockAssetRepo.EXPECT().DeleteAsset(t.ctx, t.someAsset.ID, t.someAsset.UserID).Return(assert.AnError)

	err := t.assetUseCase.RemoveAsset(t.ctx, t.someAsset.ID, t.someAsset.UserID)

	t.ErrorIs(err, assert.AnError)
}

func (t *AssetUseCaseSuite) TestPurchaseAsset_GreenPath() {
	assetID := int64(1)
	buyerID := int64(2)
	originalAsset := &entity.Asset{ID: assetID, UserID: 3} // Asset owned by userID 3

	t.mockAssetRepo.EXPECT().ExecuteTx(t.ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(usecase.AssetRepo) error) error {

			// Expected calls within the transaction
			t.mockAssetRepo.EXPECT().GetAssetByID(ctx, assetID, true).Return(originalAsset, nil)
			t.mockAssetRepo.EXPECT().UpdateAssetOwner(ctx, assetID, buyerID).Return(nil)

			return fn(t.mockAssetRepo)
		},
	)

	err := t.assetUseCase.PurchaseAsset(t.ctx, assetID, buyerID)

	t.NoError(err)
}

func (t *AssetUseCaseSuite) TestPurchaseAsset_ReturnsError_WhenAssetNotFound() {
	assetID := int64(1)
	buyerID := int64(2)

	t.mockAssetRepo.EXPECT().ExecuteTx(t.ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(usecase.AssetRepo) error) error {
			// Asset not found
			t.mockAssetRepo.EXPECT().GetAssetByID(ctx, assetID, true).Return(nil, nil)

			return fn(t.mockAssetRepo)
		},
	)

	err := t.assetUseCase.PurchaseAsset(t.ctx, assetID, buyerID)

	t.Error(err)
	t.Contains(err.Error(), "asset not found")
}

func (t *AssetUseCaseSuite) TestPurchaseAsset_ReturnsError_WhenBuyerIsOwner() {
	assetID := int64(1)
	buyerID := int64(2)
	originalAsset := &entity.Asset{ID: assetID, UserID: buyerID}

	t.mockAssetRepo.EXPECT().ExecuteTx(t.ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(usecase.AssetRepo) error) error {
			t.mockAssetRepo.EXPECT().GetAssetByID(ctx, assetID, true).Return(originalAsset, nil)

			return fn(t.mockAssetRepo)
		},
	)

	err := t.assetUseCase.PurchaseAsset(t.ctx, assetID, buyerID)

	t.Error(err)
	t.Contains(err.Error(), "cannot purchase your own asset")
}

func (t *AssetUseCaseSuite) TestGetAssetsByUser_GreenPath() {
	t.mockAssetRepo.EXPECT().GetAssetsByUserID(t.ctx, t.someAsset.UserID).Return([]*entity.Asset{t.someAsset, t.someAsset}, nil)

	res, err := t.assetUseCase.GetAssetsByUser(t.ctx, t.someAsset.UserID)

	t.NoError(err)
	t.ElementsMatch(res, []*entity.Asset{t.someAsset, t.someAsset})
}

func (t *AssetUseCaseSuite) TestGetAssetsByUser_ReturnsError_WhenRepoReturnsError() {
	t.mockAssetRepo.EXPECT().GetAssetsByUserID(t.ctx, t.someAsset.UserID).Return(nil, assert.AnError)

	res, err := t.assetUseCase.GetAssetsByUser(t.ctx, t.someAsset.UserID)

	t.ErrorIs(err, assert.AnError)
	t.Nil(res)
}
