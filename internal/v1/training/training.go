package training

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/network/vatusa"
)

// Get Training Records for cid
// @Summary Get Training Records for cid
// @Tags training
// @Param cid path string true "CID"
// @Success 200 {object} []models.TrainingNote
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/:cid [GET]
func getTraining(c *gin.Context) {
	var notes []models.TrainingNote

	user := c.MustGet("x-user").(*models.User)

	if !auth.InGroup(user, "training") && fmt.Sprint(user.CID) != c.Param("cid") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if err := database.DB.Where(models.TrainingNote{ControllerID: database.Atou(c.Param("cid"))}).Find(&notes).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, notes)
}

// Create Training Record
// @Summary Create Training Record
// @Tags training
// @Param cid path string true "CID"
// @Param training body dto.TrainingNoteRequest true "Training Note"
// @Success 201 {object} models.TrainingNote
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/:cid [POST]
func postTraining(c *gin.Context) {
	var trainingRequest dto.TrainingNoteRequest
	if err := c.ShouldBind(&trainingRequest); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	user := c.MustGet("x-user").(*models.User)

	student, err := database.FindUserByCID(c.Param("cid"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if student == nil {
		response.RespondError(c, http.StatusNotFound, "Student Not Found")
		return
	}

	training := models.TrainingNote{
		Controller:  student,
		Instructor:  user,
		Position:    trainingRequest.Position,
		Type:        trainingRequest.Type,
		Duration:    trainingRequest.Duration,
		Comments:    trainingRequest.Comments,
		SessionDate: &trainingRequest.SessionDate,
	}

	status, id, err := vatusa.SubmitTrainingNote(
		fmt.Sprint(student.CID),
		fmt.Sprint(user.CID),
		training.Position,
		*training.SessionDate,
		training.Duration,
		training.Comments,
		training.Type,
	)
	if err != nil {
		log.Errorf("VATUSA returned status code %d: %+v (%+v)", status, training, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if status > 299 {
		log.Errorf("VATUSA returned status code %d: %+v", status, training)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	training.VATUSAID = uint(id)

	if err := database.DB.Create(&training).Error; err != nil {
		log.Errorf("Failed to create training note: %+v (%+v)", training, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusCreated, training)
}

// Update Training Record
// @Summary Update Training Record
// @Tags training
// @Param cid path string true "CID"
// @Param id path string true "ID"
// @Param training body dto.TrainingNoteRequest true "Training Note"
// @Success 200 {object} models.TrainingNote
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/:cid/:id [PUT]
func putTraining(c *gin.Context) {
	var trainingRequest dto.TrainingNoteRequest
	if err := c.ShouldBind(&trainingRequest); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	training := &models.TrainingNote{}
	if err := database.DB.Find(training, c.Param("id")).Error; err != nil {
		response.RespondError(c, http.StatusNotFound, "Training Note Not Found")
		return
	}

	training.Position = trainingRequest.Position
	training.Type = trainingRequest.Type
	training.Duration = trainingRequest.Duration
	training.Comments = trainingRequest.Comments
	training.SessionDate = &trainingRequest.SessionDate

	status, err := vatusa.EditTrainingNote(
		fmt.Sprint(training.VATUSAID),
		fmt.Sprint(training.Controller.CID),
		fmt.Sprint(training.Instructor.CID),
		training.Position,
		*training.SessionDate,
		training.Duration,
		training.Comments,
		training.Type,
	)
	if err != nil || status > 299 {
		log.Errorf("VATUSA returned status code %d: %+v (%+v)", status, training, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := database.DB.Save(training).Error; err != nil {
		log.Errorf("Failed to update training note: %+v (%+v)", training, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, training)
}

// Delete Training Record
// @Summary Delete Training Record
// @Tags training
// @Param cid path string true "CID"
// @Param id path string true "ID"
// @Success 204
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/:cid/:id [DELETE]
func deleteTraining(c *gin.Context) {
	user := c.MustGet("x-user").(*models.User)

	if !auth.InGroup(user, "training") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	training := &models.TrainingNote{}
	if err := database.DB.Find(training, c.Param("id")).Error; err != nil {
		response.RespondError(c, http.StatusNotFound, "Training Note Not Found")
		return
	}

	status, err := vatusa.DeleteTrainingNote(fmt.Sprint(training.VATUSAID))
	if err != nil || status > 299 {
		log.Errorf("VATUSA returned status code %d: %+v (%+v)", status, training, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := database.DB.Delete(training).Error; err != nil {
		log.Errorf("Failed to delete training note: %+v (%+v)", training, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondBlank(c, http.StatusNoContent)
}
