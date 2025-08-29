package entity

type App struct { // СТРУКТУРА ПРИЛОЖЕНИЯ СОЗДАЕТСЯ ОДНА И МОДИФИЦИРУЕТСЧ ОДНА!! ни одной более
	ClientID     string `json:"client_id"`     //обязательное поле
	ClientSecret string `json:"client_secret"` //обязательное поле
	RedirectURI  string `json:"redirect_uri"`  //обязательное поле приложения
}
