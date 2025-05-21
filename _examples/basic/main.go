package main

import (
	"context"
	"log"
	"os"
	"strings"

	kustomersdk "github.com/superpowerdotcom/kustomersdk"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("KUSTOMER_API_KEY")
	if apiKey == "" {
		log.Fatal("KUSTOMER_API_KEY environment variable is required")
	}

	// Initialize Kustomer client
	client, err := kustomersdk.New(&kustomersdk.Options{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Kustomer client: %v", err)
	}

	// Common parameters
	customerID := "68003f6360444743ff93a0b2" // Kingsley
	categoryTree := "5_system_generated.general1"
	
	// Test SendSMS
	testSendSMS(client, customerID, categoryTree)
	
	// Test SendEmail
	testSendEmail(client, customerID, categoryTree)
	
	// Test SendNote
	testSendNote(client, customerID, categoryTree)
}

func testSendSMS(client *kustomersdk.Kustomer, customerID, categoryTree string) {
	log.Println("Testing SendSMS...")
	
	// Example of a long message that would be chunked (more than 1600 characters)
	mesasge := strings.Repeat("This is a test message that will be chunked because it's very long. ", 40)
	// message := "Testing mesage"

	// Send a short test message
	err := client.SendSMS(context.Background(), customerID, &kustomersdk.SendSMSRequest{
		From:                 "14157422828",
		To:                   "+18253653818", // Kingsley
		Body:                 mesasge,
		AutoCloseConversation: true,
		CategoryTree:         categoryTree,
	})
	
	if err != nil {
		log.Fatalf("Failed to send SMS: %v", err)
	}
	
	log.Println("SMS sent successfully")
}

func testSendEmail(client *kustomersdk.Kustomer, customerID, categoryTree string) {
	log.Println("Testing SendEmail...")
	
	bodyText := `Hi Kingsley,

We are reaching out to confirm that your upcoming lab draw at 490 Post St Ste 419, San Francisco, CA 94102 at 10:45 AM - 11:00 AM PDT on Wednesday, May 28th has been cancelled. Life gets busy, and we're here to make getting back on track easy for you.

Your bloodwork is an essential step in unlocking a detailed picture of your health and helping us create a personalized plan tailored just for you.

Rebooking takes just a minute. Choose a time that works for you! You can rebook your appointment [here](https://app.superpower.com/services?tab=all).

If you need assistance or have questions, just reply to this email or send your concierge team an SMS.

Looking forward to helping you reach your health goals!`

	bodyHTML := `<p>Hi Kingsley,</p>

<p><br /></p>

<p>We are reaching out to confirm that your upcoming lab draw at 490 Post St Ste 419, San Francisco, CA 94102 at 10:45 AM - 11:00 AM PDT on Wednesday, May 28th has been cancelled. Life gets busy, and we're here to make getting back on track easy for you.</p>

<p><br /></p>

<p>Your bloodwork is an essential step in unlocking a detailed picture of your health and helping us create a personalized plan tailored just for you.</p>

<p><br /></p>

<p>Rebooking takes just a minute. Choose a time that works for you! You can rebook your appointment <a href="https://app.superpower.com/services?tab=all" target="_blank">here</a>.</p>

<p><br /></p>

<p>If you need assistance or have questions, just reply to this email or send your concierge team an SMS.</p>

<p><br /></p>

<p>Looking forward to helping you reach your health goals!</p>

<p><br /></p>

<p>Superpower Team</p>`

	err := client.SendEmail(context.Background(), customerID, &kustomersdk.SendEmailRequest{
		From: &kustomersdk.EmailContact{
			Name:  "Superpower Concierge",
			Email: "concierge@superpower.com",
		},
		To: []*kustomersdk.EmailContact{
			{
				Name:  "Kingsley Zhong",
				Email: "kingsley.zhong@gmail.com",
			},
		},
		Subject:               "Testing Subject",
		BodyText:              bodyText,
		BodyHTML:              bodyHTML,
		AutoCloseConversation: true,
		CategoryTree:          categoryTree,
	})
	
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	
	log.Println("Email sent successfully")
}

func testSendNote(client *kustomersdk.Kustomer, customerID, categoryTree string) {
	log.Println("Testing SendNote...")
	
	noteBody := `Internal note: Customer mentioned they're preparing for a marathon next month and might need some adjustments to their nutrition plan. 
Follow up about creating a custom meal plan for endurance training. They also expressed interest in recovery supplements.`

	err := client.SendNote(context.Background(), customerID, &kustomersdk.SendNoteRequest{
		Body:                 noteBody,
		Subject:              "Marathon Preparation",
		AutoCloseConversation: false,
		CategoryTree:         categoryTree,
	})
	
	if err != nil {
		log.Fatalf("Failed to send note: %v", err)
	}
	
	log.Println("Note sent successfully")
} 