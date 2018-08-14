package main

import (

    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
    "github.com/GitbookIO/go-github-webhook"
    "strconv"
)


func main() {

    // Your GitHub secret should be added to repo->webhook and in export in bash_profile
    
    secret := os.Getenv("GITHUB_SECRET_KEY")
    
    if err := http.ListenAndServe(":8000", WebhookLog(secret)); err != nil {
        fmt.Errorf("Error: %s", err)
    }
}

func sendMail(subscriber_mailid string, pusher_email *string,  push_message []github.GitHubCommit, pusher string, repo string) {

    from := mail.NewEmail(pusher, *pusher_email)
    
    subject := pusher+" pushed into repo "+repo

    to := mail.NewEmail("subscriber", subscriber_mailid)

    //read from multiple commit messages

    var htmlContent string
    var plainTextContent string =" "

    for i, commit_message := range push_message {
       htmlContent = htmlContent+"<strong>"+"commit : "+strconv.Itoa(i+1)+" Message: "+commit_message.Message+" by "+*commit_message.Author.Name+"</strong><br/>"
       htmlContent = htmlContent+"<a href="+commit_message.URL+">link text</a><br/>"
    }
    
    //create the email
    message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

    //get an API key from sendgrid account and export it in bash_profile
    client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))


    //send mail
    response, err := client.Send(message)
    if err != nil {
        fmt.Println("email not sent")
        log.Println(err)
    } else {
        fmt.Println("email sent")
        fmt.Println(response.StatusCode)
        fmt.Println(response.Body)
        fmt.Println(response.Headers)
    }
}

func WebhookLog(secret string) http.Handler {
    return github.Handler(secret, func(event string, payload *github.GitHubPayload, req *http.Request) error {

        // Log webhook
        fmt.Println("Received", event, "for ", payload.Repository.Name)

        //search for repo name
        //get each subscriber and send messages
        // send message to all subscribers
        sendMail("devkhv129@gmail.com",payload.Pusher.Email,payload.Commits,*payload.Pusher.Name,payload.Repository.Name)

        // All is good (return an error to fail)
        return nil
    })
}