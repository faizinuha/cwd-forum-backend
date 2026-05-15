package mailtraps

import (
	"fmt"
	"gin-quickstart/pkg/email"
)

// This is a test/example file to demonstrate email functionality
// To run: go run main.go mcp.go

func TestSendWelcomeEmail() {
	client := email.NewEmailClient()

	err := client.SendWelcomeEmail("test@example.com", "JohnDoe")
	if err != nil {
		fmt.Printf("Error sending welcome email: %v\n", err)
		return
	}

	fmt.Println("Welcome email sent successfully!")
}

func TestSendForgotPasswordEmail() {
	client := email.NewEmailClient()

	resetLink := "https://cwdforum.com/reset-password?token=abc123"
	err := client.SendForgotPasswordEmail("test@example.com", resetLink)
	if err != nil {
		fmt.Printf("Error sending forgot password email: %v\n", err)
		return
	}

	fmt.Println("Forgot password email sent successfully!")
}

func TestSendNotificationEmail() {
	client := email.NewEmailClient()

	message := "Ada balasan baru pada thread Anda"
	err := client.SendNotificationEmail("test@example.com", "Thread Reply", message)
	if err != nil {
		fmt.Printf("Error sending notification email: %v\n", err)
		return
	}

	fmt.Println("Notification email sent successfully!")
}
