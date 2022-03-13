package common

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	//DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"` // gorm soft-delete feature
}

func (m *BaseModel) SetInit() {
	m.ID = uuid.New()
	m.CreatedAt = time.Now()
}

func (m *BaseModel) SetUpdate() {
	m.UpdatedAt = time.Now()
}

type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
