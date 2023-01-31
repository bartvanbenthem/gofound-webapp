package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bartvanbenthem/gofound-webapp/internal/models"
	"github.com/bartvanbenthem/gofound-webapp/internal/validator"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	blogposts, err := app.blogposts.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.BlogPosts = blogposts

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) blogpostView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	blogpost, err := app.blogposts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.BlogPost = blogpost

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) blogpostCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = blogpostCreateForm{}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

type blogpostCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Author              string `form:"author"`
	ImgURL              string `form:"img_url"`
	validator.Validator `form:"-"`
}

func (app *application) blogpostCreatePost(w http.ResponseWriter, r *http.Request) {
	var form blogpostCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Author), "author", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.blogposts.Insert(form.Title, form.Content, form.Author, form.ImgURL)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "BlogPost successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/blogpost/view/%d", id), http.StatusSeeOther)
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/blogpost/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	blogposts, err := app.blogposts.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.BlogPosts = blogposts

	app.render(w, http.StatusOK, "blog.tmpl", data)
}

// MailData holds email
type mailDataCreateForm struct {
	Name                string `form:"name"`
	From                string `form:"from"`
	Phone               string `form:"phone"`
	Subject             string `form:"subject"`
	Content             string `form:"content"`
	validator.Validator `form:"-"`
}

func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = mailDataCreateForm{}
	app.render(w, http.StatusOK, "contact.tmpl", data)
}

// PostContact is the handler for the send mail on contact page
func (app *application) contactPost(w http.ResponseWriter, r *http.Request) {

	var form mailDataCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 50), "name", "This field cannot be more than 50 characters long")
	form.CheckField(validator.Matches(form.From, validator.EmailRX), "from", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Subject), "subject", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Phone, 12), "subject", "This field cannot be more than 12 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "contact.tmpl", data)
		return
	}

	// add provided name and phone number to the content body
	body := fmt.Sprintf("%s \nNAME: %s \nPHONE: %s",
		form.Content, form.Name, form.Phone)

	msg := models.MailData{
		To:      app.mailAddress,
		From:    form.From,
		Phone:   form.Phone,
		Subject: form.Subject,
		Content: body,
	}

	app.mailChan <- msg

	// send confirmation mail back to sender
	msg = models.MailData{
		To:      form.From,
		From:    "noreply@gofound.nl",
		Subject: "contact confirmation",
		Content: `thanks for contacting us, we will reply as soon as possible`,
	}

	app.mailChan <- msg

	app.sessionManager.Put(r.Context(), "senderName", form.Name)
	app.sessionManager.Put(r.Context(), "senderMail", form.From)

	http.Redirect(w, r, "/contactresponse", http.StatusSeeOther)

}

func (app *application) contactResponse(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = mailDataCreateForm{
		Name: app.sessionManager.GetString(r.Context(), "senderName"),
		From: app.sessionManager.GetString(r.Context(), "senderMail"),
	}

	app.render(w, http.StatusOK, "contact.response.tmpl", data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
