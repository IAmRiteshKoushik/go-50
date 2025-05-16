package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendGmailGoMail(to []string, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	gmailUsername := os.Getenv("GMAIL_USERNAME")
	gmailAppPassword := os.Getenv("GMAIL_APP_PASSWORD")

	if smtpHost == "" || smtpPortStr == "" || gmailUsername == "" || gmailAppPassword == "" {
		return fmt.Errorf("Gmail SMTP configuration environment variables not set")
	}

	smtpPort, err := fmt.Sscanf(smtpPortStr, "%d")
	if err != nil {
		log.Printf("Warning: Failed to parse SMTP_PORT '%s' as int using Sscanf, attempting with Atoi...", smtpPortStr)
		port, parseErr := strconv.Atoi(smtpPortStr)
		if parseErr != nil {
			return fmt.Errorf("invalid SMTP_PORT '%s'", parseErr)
		}
		smtpPort = port
	}

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", gmailUsername)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	// Set the body - supports text/plain or text/html
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(smtpHost, smtpPort, gmailUsername, gmailAppPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email via gomail: %w", err)
	}

	log.Println("Email sent successfully via gomail!")
	return nil
}

func main() {
	recipients := []string{"riteshkoushik39@gmail.com"}
	emailSubject := "Test Email from Go via Gmail (gomail)"
	emailBody := "Hello from Go!\n\nThis is a test email sent via Gmail using the go-gomail package."

	err := SendGmailGoMail(recipients, emailSubject, emailBody)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}
}
