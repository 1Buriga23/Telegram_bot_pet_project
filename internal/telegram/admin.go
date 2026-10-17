package telegram

import "os"

func GetAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
