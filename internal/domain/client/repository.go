package client

import (
	"github.com/chack93/scrumpoker_api/internal/service/database"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func CreateClient(model *Client) (err error) {
	model.SetInit()
	if err = database.Get().Create(model).Error; err != nil {
		logrus.Errorf("failed, err: %v", err)
	}
	return
}

func ListClient(modelList *[]Client) (err error) {
	if err = database.Get().Find(modelList).Error; err != nil {
		logrus.Errorf("failed, err: %v", err)
	}
	return
}

func ListClientOfSession(sessionID uuid.UUID, modelList *[]Client) (err error) {
	if err = database.Get().Where("session_id = ?", sessionID).Find(modelList).Error; err != nil {
		logrus.Errorf("failed, err: %v", err)
	}
	return
}

func ReadClient(id uuid.UUID, model *Client) (err error) {
	if err = database.Get().First(&model, id).Error; err != nil {
		logrus.Errorf("failed, id: %s, err: %v", id.String(), err)
	}
	return
}

func UpdateClient(id uuid.UUID, model *Client) (err error) {
	model.SetUpdate()
	if err = database.Get().Save(model).Error; err != nil {
		logrus.Errorf("failed, id: %s, err: %v", model.ID.String(), err)
	}
	return
}

func DeleteClient(id uuid.UUID, model *Client) (err error) {
	if err = ReadClient(id, model); err != nil {
		logrus.Errorf("read failed, id: %s, err: %v", model.ID.String(), err)
		return
	}
	if err = database.Get().Delete(model).Error; err != nil {
		logrus.Errorf("failed, id: %s, err: %v", model.ID.String(), err)
	}
	return
}
