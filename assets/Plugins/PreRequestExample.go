// type HandlerConfig struct {
// 	Prefix   string
// 	PathDir  string
// 	Username string
// 	Password string
// 	ReadOnly bool
// }

// type PreRequestResult struct {
// 	ReadOnly      bool
// 	Authenticated bool
// }

import "gowebdav"
import "net/http"
import "fmt"

var count = 0

func PreRequest(cfg *gowebdav.HandlerConfig, r *http.Request) *gowebdav.PreRequestResult {
	count++
	fmt.Println("PreRequest", count) // This will be printed to the console when the hook is called
	fmt.Printf("cfg: %+v\n", cfg)

	if r.Method == "OPTIONS" {
		return &gowebdav.PreRequestResult{
			ReadOnly: false,
			Authed:   true,
		}
	}

	switch cfg.Prefix {
	case "/dav1":
		username, password, ok := r.BasicAuth()
		fmt.Println("username=", username, "password=", password, "ok=", ok)
		fmt.Println("Authorization=", r.Header.Get("Authorization"))
		if !ok {
			return &gowebdav.PreRequestResult{
				ReadOnly: true,
				Authed:   false,
			}
		}
		if _, ok := users[username]; !ok {
			return &gowebdav.PreRequestResult{
				ReadOnly: true,
				Authed:   false,
			}
		}
		if users[username] != password {
			return &gowebdav.PreRequestResult{
				ReadOnly: true,
				Authed:   false,
			}
		}

		if username == "user1" {
			return &gowebdav.PreRequestResult{
				ReadOnly: false,
				Authed:   true,
			}
		} else if username == "user2" {
			return &gowebdav.PreRequestResult{
				ReadOnly: true,
				Authed:   true,
			}
		} else {
			return &gowebdav.PreRequestResult{
				ReadOnly: true,
				Authed:   false,
			}
		}

	case "/dav2":
		return &gowebdav.PreRequestResult{
			ReadOnly: false,
			Authed:   true,
		}
	}

	return &gowebdav.PreRequestResult{
		ReadOnly: false,
		Authed:   true,
	}
}

// USERNAME: PASSWORD
// user1:pass1, user2:pass2, user3:pass3

var users = map[string]string{
	"user1": "pass1",
	"user2": "pass2",
	"user3": "pass3",
}

// DAVS
// dav1, dav2
var davs = []string{"dav1", "dav2"}

// PERMISSIONS
// dav1: user1 can read/write, user2 can read, user3 can't access, annonymous can't access
// dav2: all users and annonymous can read/write