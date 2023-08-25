package user

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	usersPasswords map[string][]byte
}

func Create() *User {
	users := &User{}

	users.usersPasswords = map[string][]byte{
		"asd": []byte(`$2a$10$o5uBq878SF50DzFmcVyB.elVkH461ObMOnC9pu2pxMKDyuVPKXW8C`), //pwd: asd                                                         //pwd: 123
		"qwe": []byte(`$2a$10$hBRu.UVzmzJXoGvJxWSxAOvRk/dAmmrMey6JrYwXgqKrG6IS9UD6O`), //pwd: qwe
	}

	return users
}

func (u *User) verifyUserPass(username, password string) bool {
	userPass, ok := u.usersPasswords[username]
	if !ok {
		return false
	}

	if err := bcrypt.CompareHashAndPassword(userPass, []byte(password)); err == nil {
		return true
	}

	return false
}

func (u *User) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok && u.verifyUserPass(user, pass) {
			h.ServeHTTP(w, r)
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
