package entities

import "kraken/utils"

// Certificate - Entity for Certificate
type Certificate struct {

	//Id Identifier attribute for internal control
	Id int64 `json:"id"`

	userID int64 `json:"userID"`

	//Code - Code of certificate, like EAN13
	ImageBase64 string `json:"imageBase64,omitempty"`

	ImageURL string `json:"imageURL"`

	//Name certificate name
	Name string `json:"name"`

	//IsActive indicates if the certificate is active
	IsActive bool `json:"isActive"`

	Address *Address `json:"address,omitempty"`

	CPF string `json:"cpf,omitempty"`

	CNPJ string `json:"cnpj,omitempty"`

	PHONE string `json:"phone,omitempty"`

	Email string `json:"email,omitempty"`

	LastVisitDate *utils.DateTime `json:"lastVisitDate,omitempty"`

	// StatusCode indicates if the certificate is active
	StatusCode StatusCode `json:"status_code,omitempty"`

	//ModifiedAt is the date when the certificate was modified for the last time.
	ModifiedAt *utils.DateTime `json:"lastChange"`

	//Shops used to list which shops are related with certificate
	CreatedAt *utils.DateTime `json:"created_at"`
}
