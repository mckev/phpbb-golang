package controller

import (
	"html/template"
	"net/http"
	"strings"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func UserRegisterPage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/user_register
	ctx := r.Context()
	session := GetSession(r)

	type FormData struct {
		Username        string
		NewPassword     string
		PasswordConfirm string
		Email           string
		Errors          []string
	}
	formData := FormData{}

	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			logger.Errorf(ctx, "Error while parsing form upon user registration: %s", err)
			return
		}
		formData.Username = strings.TrimSpace(r.Form.Get("username"))
		if len(formData.Username) < 4 {
			formData.Errors = append(formData.Errors, "The username you entered is too short.")
		} else if len(formData.Username) > 20 {
			formData.Errors = append(formData.Errors, "The username you entered is too long.")
		}
		// TODO: Check that username does not have invalid characters, such as space
		formData.NewPassword = strings.TrimSpace(r.Form.Get("new_password"))
		if len(formData.NewPassword) < 8 {
			formData.Errors = append(formData.Errors, "The password you entered is too short.")
		}
		if !helper.IsPasswordValid(formData.NewPassword) {
			formData.Errors = append(formData.Errors, "Password must be at least 8 characters long, must contain letters in mixed case and must contain numbers.")
		}
		formData.PasswordConfirm = strings.TrimSpace(r.Form.Get("password_confirm"))
		if formData.PasswordConfirm != formData.NewPassword {
			formData.Errors = append(formData.Errors, "Password and confirmation do not match.")
		}
		formData.Email = strings.TrimSpace(r.Form.Get("email"))
		if formData.Email != "" && !helper.IsEmailValid(formData.Email) {
			formData.Errors = append(formData.Errors, "The email address format is invalid.")
		}
		user := model.User{}
		if len(formData.Errors) == 0 {
			// Insert user into database
			user.UserId, err = model.InsertUser(ctx, formData.Username, formData.NewPassword, formData.Email, "")
			if err != nil {
				if strings.Contains(err.Error(), model.DB_ERROR_UNIQUE_CONSTRAINT) {
					formData.Errors = append(formData.Errors, "This username is already taken. Please choose a different one.")
				} else {
					logger.Errorf(ctx, "Error while inserting user: %s", err)
					formData.Errors = append(formData.Errors, "The system is currently experiencing an issue. Please try again later.")
				}
			}
		}
		if len(formData.Errors) == 0 {
			// Validation successful
			user.UserName = formData.Username

			// fmt.Fprintf(w, "Form submitted successfully!\n")
			// fmt.Fprintf(w, "Username: %s\n", formData.Username)
			// fmt.Fprintf(w, "Password: %s\n", formData.NewPassword)
			// fmt.Fprintf(w, "Confirm password: %s\n", formData.PasswordConfirm)
			// fmt.Fprintf(w, "Email address: %s\n", formData.Email)
			// TODO: Handle CSRF token validation

			// Create user session (for user registration and user login)
			ip, browser, forwardedFor := helper.ExtractUserFingerprint(r)
			session, err = model.CreateSession(ctx, user.UserId, user.UserName, ip, browser, forwardedFor)
			if err != nil {
				logger.Errorf(ctx, "Error while creating user session: %s", err)
				return
			}
			err = model.UpdateLastVisitTimeForUser(ctx, user.UserId)
			if err != nil {
				logger.Errorf(ctx, "Error while updating last visit time for user id %d: %s", user.UserId, err)
				return
			}

			templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/user_register_created.html")
			if err != nil {
				logger.Errorf(ctx, "Error while parsing user register created template files: %s", err)
				return
			}
			type UserRegisterPageData struct {
				Session                 model.Session
				RedirectURIForLoginPage string
				ForumNavTrails          []forumhelper.ForumNavTrail
			}
			userRegisterPageData := UserRegisterPageData{
				Session:                 session,
				RedirectURIForLoginPage: "./",
				ForumNavTrails:          []forumhelper.ForumNavTrail{},
			}
			err = templateOutput.ExecuteTemplate(w, "overall", userRegisterPageData)
			if err != nil {
				logger.Errorf(ctx, "Error while executing user register created template: %s", err)
				return
			}
			return
		}

		fallthrough

	case "GET":
		// Prepare template files
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/user_register.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing user register template files: %s", err)
			return
		}
		// Prepare data
		type UserRegisterPageData struct {
			FormData                FormData
			Session                 model.Session
			RedirectURIForLoginPage string
			ForumNavTrails          []forumhelper.ForumNavTrail
		}
		userRegisterPageData := UserRegisterPageData{
			FormData:                formData,
			Session:                 session,
			RedirectURIForLoginPage: "./",
			ForumNavTrails:          []forumhelper.ForumNavTrail{},
		}
		// Execute template
		err = templateOutput.ExecuteTemplate(w, "overall", userRegisterPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing user register template: %s", err)
			return
		}
	}
}
