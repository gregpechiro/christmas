package main

import "github.com/cagnosolutions/adb"

type User struct {
	Id        string `json:"id"`
	Role      string `json:"role, omitempty"`
	Active    bool   `json:"active" auth:"active"`
	Email     string `json:"email, omitempty" auth:"username" required:"register, login"`
	Password  string `json:"pass, omitempty" auth:"password" required:"register, login"`
	FirstName string `json:"firstName, omitempty" required:"register"`
	LastName  string `json:"lastName, omitempty" required:"register"`
	Age       int    `json:"age, omitempty"`
}

type Item struct {
	Id        string `json:"id"`
	UserId    string `json:"userId, omitempty"`
	Name      string `json:"name, omitempty" required:"all"`
	Notes     string `json:"notes, omitempty"`
	Quantity  int    `json:"quantity, omitempty" required:"all"`
	Link      string `json:"link, omitempty"`
	Purchased int    `json:"purchased, omitempty"`
	Fulfilled bool   `json:"fulfilled"`
}

type Invite struct {
	Id     string `json:"id"`
	ToId   string `json:"toId"`
	FromId string `json:"fromId"`
	Secret bool   `json:"secret"`
}

type Exchange struct {
	Id      string `json:"id"`
	GiverId string `json:"giverId, omitempty"`
	RecId   string `json:"recId, omitempty"`
	Secret  bool   `json:"secret"`
}

type ExchangeView struct {
	Exchange
	User
}

func GetGivingTo(userId string) []ExchangeView {
	var exchanges []Exchange
	db.TestQuery("exchange", &exchanges, adb.Eq("giverId", `"`+userId+`"`))
	var exchangeViews []ExchangeView
	for _, exchange := range exchanges {
		var user User
		db.Get("user", exchange.RecId, &user)
		exchangeViews = append(exchangeViews, ExchangeView{
			Exchange: exchange,
			User:     user,
		})
	}
	return exchangeViews
}

func GetSharedWith(userId string) []ExchangeView {
	var exchanges []Exchange
	db.TestQuery("exchange", &exchanges, adb.Eq("recId", `"`+userId+`"`), adb.Eq("secret", "false"))
	var exchangeViews []ExchangeView
	for _, exchange := range exchanges {
		var user User
		db.Get("user", exchange.GiverId, &user)
		exchangeViews = append(exchangeViews, ExchangeView{
			Exchange: exchange,
			User:     user,
		})
	}
	return exchangeViews
}
