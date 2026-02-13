package model

import (
	"api/entity"
	"strings"
)

type SensorCategory struct {
	entity.Base
	SensorCategoryAPI
}

type SensorCategoryAPI struct {
	Code        string  `json:"code" validate:"required" gorm:"type:varchar(255);uniqueIndex;not null"`
	Name        string  `json:"name" validate:"required" gorm:"type:varchar(255);not null"`
	Description *string `json:"description,omitempty" gorm:"type:text"`
	Units       string  `json:"units" gorm:"type:varchar(255)"`
	DefaultUnit string  `json:"default_unit" validate:"required" gorm:"type:varchar(255);not null"`
}

type SensorCategoryRequest struct {
	Code        *string  `json:"code,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Units       []string `json:"units,omitempty"`
	DefaultUnit *string  `json:"default_unit,omitempty"`
}

type SensorCategoryResponse struct {
	SensorCategory
	Units []Unit `json:"units"`
}

type Unit struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	IsDefault bool   `json:"is_default"`
}

func (s *SensorCategory) GetUnitNames() []string {
	if s == nil || s.Units == "" {
		return nil
	}
	units := strings.Split(s.Units, "|")
	var unitNames []string
	for _, unit := range units {
		name := unit
		if idx := strings.Index(unit, "("); idx != -1 {
			name = unit[:idx]
		}
		unitNames = append(unitNames, name)
	}
	return unitNames
}

func (s *SensorCategory) ToSensorCategoryResponse() *SensorCategoryResponse {
	if s == nil {
		return nil
	}
	units := strings.Split(s.Units, "|")
	var unitsResponse []Unit
	for _, unit := range units {
		name := unit
		symbol := unit

		// Parse format: "Name(Symbol)" or just "Name" if no parentheses
		if idx := strings.Index(unit, "("); idx != -1 {
			if idx2 := strings.Index(unit, ")"); idx2 != -1 && idx2 > idx {
				name = unit[:idx]
				symbol = unit[idx+1 : idx2]
			}
		}

		unitsResponse = append(unitsResponse, Unit{
			Name:      name,
			Symbol:    symbol,
			IsDefault: name == s.DefaultUnit,
		})
	}
	return &SensorCategoryResponse{
		SensorCategory: *s,
		Units:          unitsResponse,
	}
}

func (s *SensorCategory) Seed() *[]SensorCategory {
	return &[]SensorCategory{
		{
			SensorCategoryAPI: SensorCategoryAPI{
				Code:        "temperature",
				Name:        "Temperature",
				Description: nil,
				Units:       "Celcius(°C)|Fahrenheit(°F)",
				DefaultUnit: "Celcius",
			},
		},
		{
			SensorCategoryAPI: SensorCategoryAPI{
				Code:        "humidity",
				Name:        "Humidity",
				Description: nil,
				Units:       "Percentage(%)",
				DefaultUnit: "Percentage",
			},
		},
		{
			SensorCategoryAPI: SensorCategoryAPI{
				Code:        "pressure",
				Name:        "Pressure",
				Description: nil,
				Units:       "Pa(Pa)|Bar(Bar)|Psi(Psi)",
				DefaultUnit: "Pa",
			},
		},
		{
			SensorCategoryAPI: SensorCategoryAPI{
				Code:        "voltage",
				Name:        "Voltage",
				Description: nil,
				Units:       "V(V)",
				DefaultUnit: "V",
			},
		},
	}
}
