package controllerUtils

import (
	"sort"
	"time"

	"backend/models"
)

func FilterByFutureScheduleTimeOfMessages(messages []models.Message) []models.Message {
	result := make([]models.Message, 0)
	nowDate := time.Now()
	for _, m := range messages {
		if m.ScheduleTime.After(nowDate) {
			continue
		}
		result = append(result, m)
	}
	return result
}

func UpdateCreatedAt(messages []models.Message) []models.Message {
	for i, m := range messages {
		if m.ScheduleTime.After(m.CreatedAt) {
			messages[i].CreatedAt = m.ScheduleTime
		}
	}
	return messages
}

func SortMessageByCreatedAt(messages []models.Message) []models.Message {
	sort.Slice(messages, func(i, j int) bool { return messages[i].CreatedAt.After(messages[j].CreatedAt) })
	return messages
}
