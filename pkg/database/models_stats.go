package database

import (
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
)

func GetStatsForUserAndMonth(user *models.User, month int, year int) (float32, float32, float32, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := startDate.AddDate(0, 1, -1)
	endDate := time.Date(year, time.Month(month), lastDayOfMonth.Day(), 23, 59, 59, 999, time.UTC)

	cab, err := GetCabStatsForUser(user, startDate, endDate)
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	terminal, err := GetTerminalStatsForUser(user, startDate, endDate)
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	enroute, err := GetEnrouteStatsForUser(user, startDate, endDate)
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}

	return cab, terminal, enroute, nil
}

func GetCabStatsForUser(user *models.User, startDate, endDate time.Time) (float32, error) {
	type result struct {
		Total float32
	}
	res := &result{}
	if err := DB.Model(&models.ControllerStat{}).Where(
		"user_id = ? AND (position LIKE ? OR position LIKE ? OR position LIKE ?) AND logon_time BETWEEN ? AND ?",
		user.CID,
		"%_TWR",
		"%_GND",
		"%_DEL",
		startDate,
		endDate,
	).Select("SUM(duration) AS total").First(&res).Error; err != nil {
		return 0.0, err
	}
	return res.Total, nil
}

func GetTerminalStatsForUser(user *models.User, startDate, endDate time.Time) (float32, error) {
	type result struct {
		Total float32
	}
	res := &result{}
	if err := DB.Model(&models.ControllerStat{}).Where(
		"user_id = ? AND (position LIKE ? OR position LIKE ?) AND logon_time BETWEEN ? AND ?",
		user.CID,
		"%_APP",
		"%_DEP",
		startDate,
		endDate,
	).Select("SUM(duration) AS total").First(&res).Error; err != nil {
		return 0.0, err
	}
	return res.Total, nil
}

func GetEnrouteStatsForUser(user *models.User, startDate, endDate time.Time) (float32, error) {
	type result struct {
		Total float32
	}
	res := &result{}
	if err := DB.Model(&models.ControllerStat{}).Where(
		"user_id = ? AND (position LIKE ? OR position LIKE ?) AND logon_time BETWEEN ? AND ?",
		user.CID,
		"%_CTR",
		"%_FSS",
		startDate,
		endDate,
	).Select("SUM(duration) AS total").First(&res).Error; err != nil {
		return 0.0, err
	}
	return res.Total, nil
}
