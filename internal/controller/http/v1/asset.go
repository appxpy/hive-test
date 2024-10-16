package v1

import (
	"net/http"
	"strconv"

	"github.com/appxpy/hive-test/internal/entity"
	"github.com/appxpy/hive-test/internal/middleware"
	"github.com/appxpy/hive-test/internal/usecase"
	"github.com/appxpy/hive-test/pkg/logger"
	"github.com/gin-gonic/gin"
)

type assetRoutes struct {
	a usecase.AssetUseCase
	l logger.Interface
}

func newAssetRoutes(handler *gin.RouterGroup, a usecase.AssetUseCase, l logger.Interface, jwtSecret string) {
	r := &assetRoutes{a, l}

	h := handler.Group("/assets")
	h.Use(middleware.JWTAuth(jwtSecret))
	{
		h.POST("/", r.addAsset)
		h.DELETE("/:id", r.removeAsset)
		h.POST("/purchase/:id", r.purchaseAsset)
		h.GET("/", r.getUserAssets)
	}
}

// @Security    BearerAuth
// @Summary     Add Asset
// @Description Adds a new asset for the user
// @Tags        assets
// @Accept      json
// @Produce     json
// @Param       asset body entity.Asset true "Asset Data"
// @Success     201
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /assets [post]
func (r *assetRoutes) addAsset(c *gin.Context) {
	var asset entity.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		r.l.Error(err, "http - v1 - addAsset")
		errorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := c.GetInt64("userID")
	asset.UserID = userID

	err := r.a.AddAsset(c.Request.Context(), &asset)
	if err != nil {
		r.l.Error(err, "http - v1 - addAsset")
		errorResponse(c, http.StatusInternalServerError, "Could not add asset")
		return
	}

	c.Status(http.StatusCreated)
}

// @Security    BearerAuth
// @Summary     Remove Asset
// @Description Removes an asset owned by the user
// @Tags        assets
// @Produce     json
// @Param       id   path     int true "Asset ID"
// @Success     200
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /assets/{id} [delete]
func (r *assetRoutes) removeAsset(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseInt(assetIDStr, 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - removeAsset")
		errorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	userID := c.GetInt64("userID")

	err = r.a.RemoveAsset(c.Request.Context(), assetID, userID)
	if err != nil {
		r.l.Error(err, "http - v1 - removeAsset")
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// @Security    BearerAuth
// @Summary     Purchase Asset
// @Description Allows a user to purchase an asset
// @Tags        assets
// @Produce     json
// @Param       id   path     int true "Asset ID"
// @Success     200
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /assets/purchase/{id} [post]
func (r *assetRoutes) purchaseAsset(c *gin.Context) {
	assetIDStr := c.Param("id")
	assetID, err := strconv.ParseInt(assetIDStr, 10, 64)
	if err != nil {
		r.l.Error(err, "http - v1 - purchaseAsset")
		errorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	userID := c.GetInt64("userID")

	err = r.a.PurchaseAsset(c.Request.Context(), assetID, userID)
	if err != nil {
		r.l.Error(err, "http - v1 - purchaseAsset")
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// @Security    BearerAuth
// @Summary     Get User Assets
// @Description Retrieves all assets owned by the user
// @Tags        assets
// @Produce     json
// @Success     200 {array} entity.Asset
// @Failure     500 {object} response
// @Router      /assets [get]
func (r *assetRoutes) getUserAssets(c *gin.Context) {
	userID := c.GetInt64("userID")

	assets, err := r.a.GetAssetsByUser(c.Request.Context(), userID)
	if err != nil {
		r.l.Error(err, "http - v1 - getUserAssets")
		errorResponse(c, http.StatusInternalServerError, "Could not retrieve assets")
		return
	}

	c.JSON(http.StatusOK, assets)
}
