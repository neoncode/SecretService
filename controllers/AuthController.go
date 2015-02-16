package Controllers
import(
	"github.com/martini-contrib/auth"
)
func IsAuthorized(username, password string) bool {
    	return auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "admin")
}