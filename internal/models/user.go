package models

import "time"

type User struct {
    ID             int       `json:"id"`
    PassportNumber string    `json:"passportNumber"`
    Surname        string    `json:"surname,omitempty"`
    Name           string    `json:"name,omitempty"`
    Patronymic     string    `json:"patronymic,omitempty"`
    Address        string    `json:"address,omitempty"`
    CreatedAt      time.Time `json:"createdAt"`
    UpdatedAt      time.Time `json:"updatedAt"`
}
