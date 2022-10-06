package storage

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/gin/response"
	storagePackage "github.com/adh-partnership/api/pkg/storage"
	"github.com/adh-partnership/api/pkg/utils"
)

// Get Storage Listing
// @Summary Get Storage Listing
// @Tags storage
// @Param category path string false "Category, if applicable"
// @Success 200 {object} []models.Document
// @Failure 500 {object} response.R
// @Router /v1/storage/:category [GET]
func getStorage(c *gin.Context) {
	storage := []models.Document{}

	if c.Param("category") != "" {
		if err := database.DB.Where(models.Document{Category: c.Param("Category")}).Find(&storage).Error; err != nil {
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	} else {
		if err := database.DB.Find(&storage).Error; err != nil {
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	response.Respond(c, http.StatusOK, storage)
}

// Create storage object
// @Summary Create storage object
// @Tags storage
// @Param storage body dto.StorageRequest true "Storage Object"
// @Success 201 {object} models.Document
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/storage [POST]
func postStorage(c *gin.Context) {
	var storageRequest dto.StorageRequest
	if err := c.ShouldBind(&storageRequest); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	s := &models.Document{
		Category:    storageRequest.Category,
		Name:        storageRequest.Name,
		Description: storageRequest.Description,
		CreatedBy:   *c.MustGet("x-user").(*models.User),
		UpdatedBy:   *c.MustGet("x-user").(*models.User),
	}

	if err := database.DB.Create(&s).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response.Respond(c, http.StatusCreated, s)
}

// Update storage object
// @Summary Update storage object
// @Tags storage
// @Param storage body dto.StorageRequest true "Storage Object"
// @Param id path int true "Storage ID"
// @Success 200 {object} models.Document
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/storage/:id [PUT]
func putStorage(c *gin.Context) {
	var storageRequest dto.StorageRequest
	if err := c.ShouldBind(&storageRequest); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	s := &models.Document{}
	if err := database.DB.First(&s, c.Param("id")).Error; err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	if err := database.DB.Model(&s).Updates(models.Document{
		Category:    storageRequest.Category,
		Name:        storageRequest.Name,
		Description: storageRequest.Description,
		UpdatedBy:   *c.MustGet("x-user").(*models.User),
	}).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response.Respond(c, http.StatusOK, s)
}

// Delete storage object
// @Summary Delete storage object
// @Tags storage
// @Param id path int true "Storage ID"
// @Success 204 {object} response.R
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/storage/:id [DELETE]
func deleteStorage(c *gin.Context) {
	s := &models.Document{}
	if err := database.DB.First(&s, c.Param("id")).Error; err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	if err := database.DB.Delete(&s).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	go func(url string) {
		if url != "" {
			slug := GetSlugFromURL(url)
			err := storagePackage.Storage("uploads").DeleteObject(slug)
			if err != nil {
				log.Errorf("Error deleting object from storage: %s", err.Error())
				_ = discord.SendWebhookMessage("uploads", "API", fmt.Sprintf("Error deleting object %s from uploads storage: %v", slug, err))
			}
		}
	}(s.URL)

	response.Respond(c, http.StatusNoContent, nil)
}

// Upload file data to storage object
// @Summary Upload file data to storage object
// @Tags storage
// @Param id path int true "Storage ID"
// @Param file formData file true "File"
// @Success 204
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/storage/:id/file [PUT]
func putStorageFile(c *gin.Context) {
	s := &models.Document{}
	if err := database.DB.First(&s, c.Param("id")).Error; err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}
	// Capped at 100MB
	if file.Size > 100*1024*1024 {
		response.RespondError(c, http.StatusBadRequest, "File too large")
		return
	}
	if file.Size == 0 {
		response.RespondError(c, http.StatusBadRequest, "File is empty")
		return
	}

	if s.URL != "" {
		slug := GetSlugFromURL(s.URL)
		err := storagePackage.Storage("uploads").DeleteObject(slug)
		if err != nil {
			log.Errorf("Error deleting object from storage: %s", err.Error())
			_ = discord.SendWebhookMessage("uploads", "API", fmt.Sprintf("Error deleting object %s from uploads storage: %v", slug, err))
		}
	}
	fileSlug := fmt.Sprintf("%s.%s", utils.StringToSlug(s.Name), filepath.Ext(file.Filename))

	mtype, err := mimetype.DetectFile(file.Filename)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	err = storagePackage.Storage("uploads").PutObject(fileSlug, file.Filename, false, file.Size, mtype.String())
	if err != nil {
		log.Errorf("Error uploading file to storage: %s", err.Error())
		_ = discord.SendWebhookMessage("uploads", "API", fmt.Sprintf("Error uploading file %s to uploads storage: %v", fileSlug, err))
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	s.UpdatedBy = *c.MustGet("x-user").(*models.User)
	s.URL = GenerateURL(fileSlug)
	if err := database.DB.Save(&s).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	_ = discord.SendWebhookMessage("uploads", "API", fmt.Sprintf("Uploaded file %s to uploads storage", fileSlug))

	response.RespondBlank(c, http.StatusNoContent)
}
