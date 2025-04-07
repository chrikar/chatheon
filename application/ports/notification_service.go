package ports

type NotificationService interface {
	Notify(userID, message string) error
}