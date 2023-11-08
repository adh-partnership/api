/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package certifications

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
)

type CertificationDTO struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Order       uint   `json:"order"`
	Hidden      bool   `json:"hidden"`
}

// Get certification types
// @Summary Get certification types
// @Description Get certification types
// @Tags certifications
// @Success 200 {object} []string
// @Failure 500 {object} response.R
// @Router /v1/certifications [get]
func getCertifications(c *gin.Context) {
	response.Respond(c, http.StatusOK, database.GetCertifications())
}

// Create a new certification type
// @Summary Create a new certification type
// @Description Create a new certification type
// @Tags certifications
// @Param body body CertificationDTO true "Certification type data"
// @Success 204
// @Failure 401 {object} response.R
// @Failure 409 {object} response.R "Conflict - certification name already in use"
// @Failure 500 {object} response.R
// @Router /v1/certifications [post]
func postCertifications(c *gin.Context) {
	var certificationDTO CertificationDTO

	if err := c.ShouldBind(&certificationDTO); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	if database.DB.Where(models.Certification{Name: certificationDTO.Name}).First(&models.Certification{}).Error == nil {
		response.RespondError(c, http.StatusConflict, "Certification name already in use")
		return
	}

	if err := database.DB.Create(&models.Certification{
		DisplayName: certificationDTO.DisplayName,
		Name:        certificationDTO.Name,
		Order:       certificationDTO.Order,
		Hidden:      certificationDTO.Hidden,
	}).Error; err != nil {
		log.Errorf("Failed to create certification: %+v", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	database.InvalidateCertCache()

	response.Respond(c, http.StatusNoContent, nil)
}

// Bulk Reorder Certifications
// @Summary Bulk Reorder Certifications
// @Description Bulk Reorder Certifications
// @Tags certifications
// @Param body body []int true "Array of certification IDs in order"
// @Success 204
// @Failure 401 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/certifications/bulk-order [patch]
func patchBulkOrder(c *gin.Context) {
	var certifications []int

	if err := c.ShouldBind(&certifications); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for i, certID := range certifications {
			if err := tx.Model(&models.Certification{}).Where(&models.Certification{
				ID: uint(certID),
			}).Update("order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to reorder certifications: %+v", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	database.InvalidateCertCache()

	response.Respond(c, http.StatusNoContent, nil)
}

// Update a certification type
// @Summary Update a certification type
// @Description Update a certification type
// @Tags certifications
// @Param name path string true "Current certification name"
// @Param body body CertificationDTO true "Certification type data"
// @Success 204
// @Failure 401 {object} response.R
// @Failure 404 {object} response.R
// @Failure 409 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/certifications/{name} [put]
func putCertifications(c *gin.Context) {
	var certificationDTO CertificationDTO
	if err := c.ShouldBind(&certificationDTO); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	// If changing names and the new name is already in use, return a conflict
	if c.Param("name") != certificationDTO.Name &&
		database.DB.Where(models.Certification{Name: certificationDTO.Name}).First(&models.Certification{}).Error == nil {
		response.RespondError(c, http.StatusConflict, "Certification name already in use")
		return
	}

	cert := &models.Certification{}
	if err := database.DB.Where(models.Certification{Name: c.Param("name")}).First(cert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RespondError(c, http.StatusNotFound, "Certification not found")
			return
		}
		log.Errorf("Failed to find certification: %+v", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	cert.Name = certificationDTO.Name
	cert.DisplayName = certificationDTO.DisplayName
	cert.Order = certificationDTO.Order
	cert.Hidden = certificationDTO.Hidden

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(cert).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.UserCertification{}).Where(&models.UserCertification{
			Name: c.Param("name"),
		}).Update("name", certificationDTO.Name).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("Failed to update certification: %+v", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	database.InvalidateCertCache()

	response.Respond(c, http.StatusNoContent, nil)
}

// Delete a certification type
// @Summary Delete a certification type
// @Description Delete a certification type
// @Tags certifications
// @Param name path string true "Certification name"
// @Success 204
// @Failure 401 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/certifications/{name} [delete]
func deleteCertifications(c *gin.Context) {
	cert := &models.Certification{}
	if err := database.DB.Where(models.Certification{Name: c.Param("name")}).First(cert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.RespondError(c, http.StatusNotFound, "Certification not found")
			return
		}
		log.Errorf("Failed to find certification: %+v", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(cert).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.UserCertification{}).Where(&models.Certification{
			Name: c.Param("name"),
		}).Delete(&models.UserCertification{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("Failed to delete certification: %+v", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	database.InvalidateCertCache()

	response.Respond(c, http.StatusNoContent, nil)
}
