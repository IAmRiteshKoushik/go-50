package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendGmailSMTP(to []string, subject, body string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	smtpHost := os.Getenv("SMTP_HOST")                  // Should be "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT")                  // Should be "587" or "465"
	gmailUsername := os.Getenv("GMAIL_USERNAME")        // Your Gmail address
	gmailAppPassword := os.Getenv("GMAIL_APP_PASSWORD") // Your App Password

	if smtpHost == "" || smtpPort == "" || gmailUsername == "" || gmailAppPassword == "" {
		return fmt.Errorf("Gmail SMTP configuration environment variables not set")
	}

	// Authentication: Use the Gmail address and the App Password
	// For port 587 (STARTTLS), PlainAuth is standard.
	// The 'identity' param is usually "" or the username.
	auth := smtp.PlainAuth("", gmailUsername, gmailAppPassword, smtpHost)

	// Build the email message content (including headers)
	message := fmt.Sprintf("From: %s\r\n", gmailUsername)
	message += fmt.Sprintf("To: %s\r\n", joinStrings(to, ","))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	message += "\r\n" // Essential blank line between headers and body
	message += body

	smtpAddress := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Send the email
	// For port 587, SendMail automatically handles STARTTLS.
	// For port 465 (implicit TLS), SendMail *might* work depending on Go version/server config,
	// but sometimes requires manually establishing a TLS connection first. Port 587 is recommended.
	err = smtp.SendMail(smtpAddress, auth, gmailUsername, to, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email via Gmail SMTP: %w", err)
	}

	log.Println("Email sent successfully via Gmail SMTP!")
	return nil
}

// Helper function for joining recipients (same as before)
func joinStrings(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0]
	}
	result := s[0]
	for _, str := range s[1:] {
		result += sep + str
	}
	return result
}

func main() {
	// Ensure your environment variables GMAIL_USERNAME and GMAIL_APP_PASSWORD are set
	recipients := []string{"riteshkoushik39@gmail.com"} // Replace
	emailSubject := "Test Email from Go via Gmail (net/smtp)"
	emailBody := "Hello from Go!\n\nThis is a test email sent via Gmail using the net/smtp package."

	err := SendGmailSMTP(recipients, emailSubject, emailBody)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}
}
