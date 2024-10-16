package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/appxpy/hive-test/internal/entity"
	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/jmoiron/sqlx"
)

type AssetRepoImpl struct {
	db sqlx.ExtContext
}

// NewAssetRepo creates a new AssetRepo with a database connection.
func NewAssetRepo(db *sqlx.DB) usecase.AssetRepo {
	return &AssetRepoImpl{
		db: db,
	}
}

func (r *AssetRepoImpl) CreateAsset(ctx context.Context, asset *entity.Asset) error {
	query := `
        INSERT INTO assets (user_id, name, description, price)
        VALUES ($1, $2, $3, $4)
        RETURNING id`
	return sqlx.GetContext(ctx, r.db, &asset.ID, query, asset.UserID, asset.Name, asset.Description, asset.Price)
}

func (r *AssetRepoImpl) DeleteAsset(ctx context.Context, assetID, userID int64) error {
	query := `DELETE FROM assets WHERE id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, assetID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("asset not found or not owned by user")
	}

	return nil
}

func (r *AssetRepoImpl) GetAssetByID(ctx context.Context, assetID int64, forUpdate bool) (*entity.Asset, error) {
	asset := &entity.Asset{}
	query := `SELECT id, user_id, name, description, price FROM assets WHERE id = $1`
	if forUpdate {
		query += ` FOR UPDATE`
	}
	err := sqlx.GetContext(ctx, r.db, asset, query, assetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return asset, nil
}

func (r *AssetRepoImpl) UpdateAssetOwner(ctx context.Context, assetID, newOwnerID int64) error {
	query := `UPDATE assets SET user_id = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, newOwnerID, assetID)
	return err
}

func (r *AssetRepoImpl) GetAssetsByUserID(ctx context.Context, userID int64) ([]*entity.Asset, error) {
	var assets []*entity.Asset
	query := `SELECT id, user_id, name, description, price FROM assets WHERE user_id = $1`
	err := sqlx.SelectContext(ctx, r.db, &assets, query, userID)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *AssetRepoImpl) ExecuteTx(ctx context.Context, fn func(repo usecase.AssetRepo) error) error {
	db, ok := r.db.(*sqlx.DB)
	if !ok {
		return errors.New("ExecuteTx: cannot start a transaction within an existing transaction")
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	txRepo := &AssetRepoImpl{
		db: tx,
	}

	err = fn(txRepo)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
