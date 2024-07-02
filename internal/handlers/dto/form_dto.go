package dto

type LoginFormDto struct {
	Email     string `form:"email"`
	Passsword string `form:"password"`
}

type SettingsFormDto struct {
	Amount   uint `form:"amount"`
	SearchOn bool `form:"search-on"`
	AddNew   bool `form:"add-new"`
}
