package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	usecase "github.com/open-wm/blockehr/pkg/domain/usecases"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type AuthHandler interface {
	WithUser(auth AuthenticatedHandler, roles ...entities.Role) Handler
	LoginJSON(w http.ResponseWriter, r *http.Request)
	LoginForm(w http.ResponseWriter, r *http.Request)
	LoginFormData(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}
type authHandler struct {
	uc usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) AuthHandler {
	return &authHandler{uc}
}

// // mapa//claveMapa//valores de mapa tipo session

type AuthenticatedHandler func(user *entities.User, rw http.ResponseWriter, r *http.Request)

func unauth(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		JSON(rw, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	http.Redirect(rw, r, "/login?next="+r.URL.Path, http.StatusFound)
}
func (h *authHandler) WithUser(next AuthenticatedHandler, roles ...entities.Role) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user from session
		sessCookie, err := r.Cookie(usecase.SESSION_ID) //obtiene el cookie "session_token"
		if err != nil {
			if err == http.ErrNoCookie { //si hay problemas
				log.Println("No cookie found")
				// add a query parameter to the url with the "next" page
				unauth(w, r)
				return
			}
			return
		}

		sessionToken := sessCookie.Value //almacena el valor de la cookie

		userSession := h.uc.GetSession(sessionToken)
		if userSession == nil { //si no existe
			log.Println("No session found", sessionToken)
			unauth(w, r)
			return
		}
		if userSession.IsExpired() { //si esta experiado, se borra del mapa sessions
			log.Println("expired session found")
			h.uc.DeleteSession(sessionToken)
			unauth(w, r)
			return
		}
		userSession.User.Session = userSession

		if len(roles) > 0 {
			allowed := false
			for _, role := range roles {
				if userSession.User.Profile.Role == role {
					allowed = true
					break
				}
			}
			if !allowed {
				unauth(w, r)
				return
			}
		}

		next(userSession.User, w, r)
	}
}

func (h *authHandler) LoginForm(w http.ResponseWriter, r *http.Request) {
	// if it is running from tests then change the baseDir
	tmpl := template.New("login.html.tmpl")
	tmpl, err := tmpl.ParseFiles(basePath + "login.html.tmpl")
	if err != nil {
		panic(err)
	}

	// get next page from query parameter
	next := r.URL.Query().Get("next")

	err = tmpl.Execute(w, map[string]interface{}{
		"next": next,
	})
	if err != nil {
		log.Println(err)
	}
}
func (h *authHandler) LoginFormData(w http.ResponseWriter, r *http.Request) {
	var creds usecase.LoginStruct
	// read form data
	err := r.ParseForm()
	if err != nil {
		log.Println("Bad request", err)
		// w.WriteHeader(http.StatusBadRequest)
		// http.Redirect(w, r, "/login", http.StatusBadRequest)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	creds.Email = r.FormValue("email")
	creds.Password = r.FormValue("password")

	sess, err := h.uc.Login(creds)
	if err != nil {
		slog.Error("Error logging in", "error", err.Error())
		// w.WriteHeader(http.StatusUnauthorized) //401 no autorizado
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{ //configuracion de cookie
		Name:    usecase.SESSION_ID,
		Value:   sess.Token,
		Expires: sess.CreatedAt.Add(24 * time.Hour),
	})
	next := r.FormValue("next")
	if next == "" {
		if sess.User.Profile.Role == entities.DOCTOR {
			next = "/appointments"
		} else if sess.User.Profile.Role == entities.STAFF {
			log.Println(sess.User.Profile.Role)
			next = "/doctors"
		}
	}
	log.Println("redirecting to to" + next)
	http.Redirect(w, r, next, http.StatusFound)
}

func (h *authHandler) LoginJSON(w http.ResponseWriter, r *http.Request) {
	var creds usecase.LoginStruct
	err := json.NewDecoder(r.Body).Decode(&creds) // JSON a Credentials, error se alm en err
	if err != nil {                               //si no funciona, manda el status
		w.WriteHeader(http.StatusBadRequest) //bad request
		return
	}

	sess, err := h.uc.Login(creds)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized) //401 no autorizado
		return
	}

	http.SetCookie(w, &http.Cookie{ //configuracion de cookie
		Name:    usecase.SESSION_ID,
		Value:   sess.Token,
		Expires: sess.CreatedAt.Add(24 * time.Hour),
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//obtener el cookie de la sesion
	cookie, err := r.Cookie(usecase.SESSION_ID)

	if err != nil { //si no encuentra el cookie de la sesion, ya esta borrado
		w.WriteHeader(http.StatusOK)
		return
	}
	err = h.uc.DeleteSession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    usecase.SESSION_ID,
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})

	// if its a json request, return a json response
	if r.Header.Get("Content-Type") == "application/json" ||
		r.Header.Get("Accept") == "application/json" {
		JSON(w, http.StatusOK, "Logged out", nil)
	} else {
		// if its a form request, redirect to login
		http.Redirect(w, r, "/login", http.StatusFound)
	}

}
