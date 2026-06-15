package gmailer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Gmailer struct {
	gmailService *gmail.Service
}

// Retrieve a token, saves the token, then returns the generated client.
func (g *Gmailer) getClient(config *oauth2.Config, tokenPath string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := g.tokenFromFile(tokenPath)
	if err != nil {
		tok = g.getTokenFromWeb(config)
		g.saveToken(tokenPath, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func (g *Gmailer) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func (g *Gmailer) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (g *Gmailer) saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (g *Gmailer) CreateGmailer(credentialPath string, tokenPath string) Gmailer {
	selectedCredentialPath := "credentials.json"
	if credentialPath == "" {
		selectedCredentialPath = credentialPath
	}

	selectedTokenPath := "token.json"
	if tokenPath == "" {
		selectedTokenPath = tokenPath
	}

	ctx := context.Background()
	b, err := os.ReadFile(selectedCredentialPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := g.getClient(config, selectedTokenPath)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return Gmailer{gmailService: srv}
}

func (g *Gmailer) sendEmail(gmailMessage GmailMessage) error {
	if g.gmailService == nil {
		log.Fatal("Gmail Service - 'Gmailer' is not initialized.")
		return errors.New("Error: Gmail Service - 'Gmailer' is not initialized.")
	}

	// Encode message to base64url format
	msg, err := gmailMessage.CreateGmailMessage()
	if err != nil {
		return err
	}
	_, err = g.gmailService.Users.Messages.Send("me", &msg).Do()
	if err != nil {
		log.Fatalf("Unable to send message: %v", err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}

type GmailMessage struct {
	to      string
	subject string
	body    string
	from    string
}

func (gm GmailMessage) CreateGmailMessage() (gmail.Message, error) {
	if gm.to == "" || gm.subject == "" || gm.body == "" || gm.from == "" {
		return gmail.Message{}, errors.New("GmailMessage missing elements.")
	}
	msgStr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", gm.from, gm.to, gm.subject, gm.body)
	return gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msgStr)),
	}, nil
}
