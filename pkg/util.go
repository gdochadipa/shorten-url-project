package pkg

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate *validator.Validate
var randomizer *rand.Rand
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var lengthLetters = len(letters)
func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	randomizer = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b{
		b[i] = letters[randomizer.Intn(lengthLetters)]
	}
	return string(b)
}

func ValidateOne(data string, rules string) error {
	err := validate.Var(data, rules)
	if err != nil {
		return fmt.Errorf("validation failed: %s", rules)
	}
	return nil
}

func IsAlphaNumeric(input string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(input)
}

func ValidateURL(uri string) error {
	u, err := url.Parse(uri)

	if err != nil {
		zap.L().Error("url.malformed", zap.Error(err))
		return fmt.Errorf("malformed URL")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must start with http or https")
	}

	if u.Host == "" {
        return fmt.Errorf("URL is missing a host")
    }

    return nil
}
