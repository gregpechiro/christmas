package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var db *adb.DB
var tmpl *web.TmplCache

func init() {
	db = adb.NewDB()
	web.Funcs["sub"] = sub
	tmpl = web.NewTmplCache()
	db.AddStore("user")
	db.AddStore("item")
	db.AddStore("exchange")

}

func main() {
	mux := web.NewMux()

	mux.AddRoutes(login, register, loginPost, logout)

	mux.AddSecureRoutes(USER, myList, saveItem, giving, sharing, invite, list, purchase)

	http.ListenAndServe(":8080", mux)
}

var myList = web.Route{"GET", "/my-list", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	var user User
	if !db.Get("user", userId, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	var items []Item
	db.TestQuery("item", &items, adb.Eq("userId", `"`+user.Id+`"`))
	tmpl.Render(w, r, "my-list.tmpl", web.Model{
		"user":  user,
		"items": items,
	})
}}

var saveItem = web.Route{"POST", "/item", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	r.ParseForm()
	var item Item
	db.Get("item", r.FormValue("id"), &item)
	if errs, ok := FormToStruct(&item, r.Form, ""); !ok {
		e := ""
		for _, msg := range errs {
			e += msg + " "
		}
		fmt.Println()
		web.SetErrorRedirect(w, r, "/my-list", "Error saving Item<br>"+e)
		return
	}

	if item.Id == "" {
		item.Id = genId()
		item.UserId = userId
	}
	fmt.Println(item)
	db.Set("item", item.Id, item)
	web.SetSuccessRedirect(w, r, "/my-list", "Successfully saved item")

	return
}}

var giving = web.Route{"GET", "/giving", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")

		return
	}
	var user User
	if !db.Get("user", userId, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	givingTo := GetGivingTo(user.Id)
	tmpl.Render(w, r, "giving.tmpl", web.Model{

		"user":     user,
		"givingTo": givingTo,
	})
}}

var sharing = web.Route{"GET", "/sharing", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	var user User
	if !db.Get("user", userId, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	sharingWith := GetSharedWith(user.Id)
	tmpl.Render(w, r, "sharing.tmpl", web.Model{
		"user":        user,
		"sharingWith": sharingWith,
	})
}}

var intives = web.Route{"GET", "/invite", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	var user User
	if !db.Get("user", userId, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	exInvites := GetInvites(user.Id)
	tmpl.Render(w, r, "giving.tmpl", web.Model{

		"user":    user,
		"invites": exInvites,
	})
}}

var invite = web.Route{"POST", "/invite", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	email := r.FormValue("email")
	var user User
	if !db.TestQueryOne("user", &user, adb.Eq("email", email)) {
		web.SetErrorRedirect(w, r, "/sharing", "Sorry we could not found the user you were looking for")
		return
	}

	exchange := Exchange{
		Id:       genId(),
		RecId:    userId,
		GiverId:  user.Id,
		Secret:   false,
		Accepted: false,
	}

	db.Set("exchange", exchange.Id, exchange)
	web.SetSuccessRedirect(w, r, "/sharing", "Successfully shared Christmas list")
	return

}}

var list = web.Route{"GET", "/list/:id", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	var user User
	if !db.Get("user", userId, &user) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}
	recId := r.FormValue(":id")
	var recUser User
	if !db.Get("user", recId, &recUser) {
		web.SetErrorRedirect(w, r, "/giving", "Sorry we could not find that user")
		return
	}
	var exchange Exchange
	if !db.TestQueryOne("exchange", &exchange, adb.Eq("recId", `"`+recId+`"`), adb.Eq("giverId", `"`+userId+`"`)) {
		web.SetErrorRedirect(w, r, "/giving", "Sorry that user hasn't shared their list with you")
		return
	}
	var items []Item
	db.TestQuery("item", &items, adb.Eq("userId", `"`+recId+`"`), adb.Eq("fulfilled", "false"))
	tmpl.Render(w, r, "list.tmpl", web.Model{
		"user":  user,
		"items": items,
		"exchangeView": ExchangeView{
			Exchange: exchange,
			User:     recUser,
		},
	})
}}

var purchase = web.Route{"POST", "/purchase/:recId/:itemId", func(w http.ResponseWriter, r *http.Request) {
	userId, ok := web.GetSess(r, "id").(string)
	if !ok {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/", "Please login")
		return
	}

	recId := r.FormValue(":recId")

	var exchange Exchange
	if !db.TestQueryOne("exchange", &exchange, adb.Eq("recId", `"`+recId+`"`), adb.Eq("giverId", `"`+userId+`"`)) {
		web.SetErrorRedirect(w, r, "/giving", "Sorry that user hasn't shared their list with you")
		return
	}

	var item Item
	if !db.Get("item", r.FormValue(":itemId"), &item) {
		web.SetErrorRedirect(w, r, "/list/"+recId, "Sorry we could not fond that item")
		return
	}

	if item.UserId != recId {
		web.SetErrorRedirect(w, r, "/list/"+recId, "Sorry we could not fond that item")
		return
	}

	purQty, _ := strconv.Atoi(r.FormValue("purchased"))

	item.Purchased += purQty
	item.Fulfilled = item.Purchased >= item.Quantity
	db.Set("item", item.Id, item)
	web.SetSuccessRedirect(w, r, "/list/"+recId, "Successfully marked item purchased")
	return
}}

/*var invite = web.Route{"GET", "/invite", func(w http.ResponseWriter, r *http.Request) {

}}*/
