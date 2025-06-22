package controller

import (
	"html/template"
	"net/http"
	"strings"
	"unicode/utf8"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func UserLoginPage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/user_login
	ctx := r.Context()
	session := getSession(r)
	queryParams := r.URL.Query()

	type FormData struct {
		Username   string
		Password   string
		Errors     []string
		RedirectTo string
	}
	formData := FormData{}

	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			logger.Errorf(ctx, "Error while parsing form upon user login: %s", err)
			return
		}
		formData.Username = strings.TrimSpace(r.Form.Get("username"))
		if formData.Username == "" {
			formData.Errors = append(formData.Errors, "You have specified an incorrect username. Please check your username and try again.")
		} else if utf8.RuneCountInString(formData.Username) < 4 {
			formData.Errors = append(formData.Errors, "The username you entered is too short.")
		} else if utf8.RuneCountInString(formData.Username) > 20 {
			formData.Errors = append(formData.Errors, "The username you entered is too long.")
		}
		formData.Password = strings.TrimSpace(r.Form.Get("password"))
		if formData.Password == "" {
			formData.Errors = append(formData.Errors, "You cannot login without a password.")
		} else if utf8.RuneCountInString(formData.Password) < 8 {
			formData.Errors = append(formData.Errors, "The password you entered is too short.")
		}
		formData.RedirectTo = strings.TrimSpace(r.Form.Get("redirect"))
		if formData.RedirectTo == "" {
			formData.RedirectTo = "./"
		}
		user := model.User{}
		if len(formData.Errors) == 0 {
			user, err = model.GetUserForLogin(ctx, formData.Username)
			if err != nil {
				logger.Errorf(ctx, "Error while retrieving user name '%s' for login: %s", formData.Username, err)
				formData.Errors = append(formData.Errors, "Invalid username or password")
			}
		}
		if len(formData.Errors) == 0 {
			if !helper.IsPasswordCorrect(formData.Password, user.UserPasswordHashed) {
				logger.Errorf(ctx, "Error while login for user name '%s': Password is incorrect", formData.Username)
				formData.Errors = append(formData.Errors, "Invalid username or password")
			}
		}
		if len(formData.Errors) == 0 {
			// Validation successful

			// fmt.Fprintf(w, "Form submitted successfully!\n")
			// fmt.Fprintf(w, "Username: %s\n", formData.Username)
			// fmt.Fprintf(w, "Password: %s\n", formData.Password)
			// fmt.Fprintf(w, "RedirectTo: %s\n", formData.RedirectTo)
			// TODO: Handle CSRF token validation

			// Create user session (for user registration and user login)
			ip, browser, forwardedFor := helper.ExtractUserFingerprint(r)
			session, err = createSession(ctx, user.UserId, user.UserName, ip, browser, forwardedFor)
			if err != nil {
				logger.Errorf(ctx, "Error while creating user session: %s", err)
				return
			}
			err = model.UpdateLastVisitTimeForUser(ctx, user.UserId)
			if err != nil {
				logger.Errorf(ctx, "Error while updating last visit time for user id %d: %s", user.UserId, err)
				return
			}

			// Redirect user to their last visited page
			http.Redirect(w, r, helper.UrlWithSID(formData.RedirectTo, session.SessionId), http.StatusFound)
			return
		}

		fallthrough

	case "GET":
		// Prepare template files
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/user_login.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing user login template files: %s", err)
			return
		}

		// Prepare data
		if formData.RedirectTo == "" {
			// If say, user entered a wrong password, then we shall keep the hidden input "redirect" intact.
			// Otherwise we record query parameter "redirect" as hidden input.
			if queryParams.Get("redirect") == "" {
				formData.RedirectTo = "./"
			} else {
				formData.RedirectTo = queryParams.Get("redirect")
			}
		}
		type UserLoginPageData struct {
			FormData                FormData
			Session                 model.Session
			RedirectURIForLoginPage string
			ForumNavTrails          []forumhelper.ForumNavTrail
		}
		userLoginPageData := UserLoginPageData{
			FormData:                formData,
			Session:                 session,
			RedirectURIForLoginPage: "./",
			ForumNavTrails:          []forumhelper.ForumNavTrail{},
		}

		// Execute template
		err = templateOutput.ExecuteTemplate(w, "overall", userLoginPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing user login template: %s", err)
			return
		}
	}
}
