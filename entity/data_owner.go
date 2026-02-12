package entity

import (
	"encoding/json"

	"github.com/google/uuid"
)

type DataOwner struct {
	CreatedBy   *uuid.UUID `json:"created_by,omitempty" gorm:"type:varchar(36)" swaggerignore:"true"` // created by
	UpdatedBy   *uuid.UUID `json:"updated_by,omitempty" gorm:"type:varchar(36)" swaggerignore:"true"` // updated by
	CreatorName *string    `json:"creator_name,omitempty" gorm:"-"`
	UpdaterName *string    `json:"updater_name,omitempty" gorm:"-"`
}

type DataUser struct {
	UserID *uuid.UUID `json:"user_id,omitempty"`
}

func (d *DataOwner) AssignCreator(v any) error {
	by, err := json.Marshal(v)
	if err != nil {
		return err
	}

	dataUser := DataUser{}
	if err := json.Unmarshal(by, &dataUser); err != nil {
		return err
	}

	d.CreatedBy = dataUser.UserID
	d.UpdatedBy = dataUser.UserID

	return nil
}

func (d *DataOwner) AssignUpdater(v any) error {
	by, err := json.Marshal(v)
	if err != nil {
		return err
	}

	dataUser := DataUser{}
	if err := json.Unmarshal(by, &dataUser); err != nil {
		return err
	}
	d.UpdatedBy = dataUser.UserID

	return nil
}
