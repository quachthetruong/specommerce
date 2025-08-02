package customer

type Customer struct {
	Id       string `json:"id" bun:"id,pk,skipupdate"`
	FullName string `json:"full_name" validate:"required"`
}
