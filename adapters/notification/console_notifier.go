package notification

import "fmt"

type ConsoleNotifier struct{}

func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

func (c *ConsoleNotifier) Notify(userID, message string) error {
	fmt.Printf("Notification to user %s: %s\n", userID, message)
	return nil
}