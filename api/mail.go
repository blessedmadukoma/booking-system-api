package api

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type mailStructure struct {
	employeeName   string
	employeeEmail  string
	visitorName    string
	visitorEmail   string
	visitorPicture string
}

func mailEmployee(mailInfo mailStructure) []byte {
	address := "helpdesk@visitorsapi.com"
	name := "visitorsapi Helpdesk"
	from := mail.NewEmail(name, address)
	subject := "Visit Request"
	address = mailInfo.employeeEmail
	name = mailInfo.employeeName

	to := mail.NewEmail(name, address)
	// content := mail.NewContent("text/plain", fmt.Sprintf("Hello %s, \n\n%s has requested to see you image<%s>. Kindly log in to your dashboard to approve, or deny this visit request. \n\nRegards, \nvisitorsapi Helpdesk", mailInfo.employeeName, mailInfo.visitorName, mailInfo.visitorPicture))
	content := mail.NewContent("text/html", fmt.Sprintf("Hello %s, <br/><br/>%s has requested to see you (<a href='%s' target='_blank'>Visitor Image</a>). <br/>Kindly <a href='https://approval.visitors.visitorsapi.com' target='_blank'>log in</a> to your dashboard to approve, or deny this visit request with a reason why. <br/><br/>Regards, <br/>visitorsapi Helpdesk", mailInfo.employeeName, mailInfo.visitorName, mailInfo.visitorPicture))
	m := mail.NewV3MailInit(from, subject, to, content)

	// add attachment
	// imgName := strings.SplitAfter(mailInfo.visitorPicture, "visitors/")[1]
	// imgName = strings.Split(imgName, "%20")[1]
	// format := strings.Split(imgName, ".")[1]

	// a1 := mail.NewAttachment()
	// a1.SetContent(mailInfo.visitorPicture)

	// encodedText := util.Base64Encode(mailInfo.visitorPicture)
	// decodedText, _ := util.Base64Decode(encodedText)

	// fmt.Println("Picture: => ", mailInfo.visitorPicture)
	// fmt.Println("Encoded: ", encodedText)
	// fmt.Println("Decoded: ", decodedText)

	// a1.SetFilename(fmt.Sprintf("%s", imgName))
	// a1.SetType(fmt.Sprintf("image/%s", format))
	// a1.SetDisposition("attachment")
	// a1.Content = mailInfo.visitorPicture
	// a1.SetContent(encodedText)
	// a1.SetDisposition("inline")
	// a1.SetContentID("Visitor Picture")

	// m.AddAttachment(a1)

	return mail.GetRequestBody(m)
}

func (srv *Server) sendMailEmployee(mailinfo mailStructure) error {
	request := sendgrid.GetRequest(srv.config.SendGridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mailEmployee(mailinfo)
	request.Body = Body
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println("error sending email:", err)
		return err
	}

	// log.Println("Response code:", response.StatusCode)
	log.Println("Mail sent to employee successfully")
	return nil
}

func mailVisitor(mailInfo mailStructure) []byte {
	address := "helpdesk@visitorsapi.com"
	name := "visitorsapi Helpdesk"
	fmt.Println(name, address)
	from := mail.NewEmail(name, address)
	subject := "Thank	you for visiting visitorsapi"
	address = mailInfo.visitorEmail
	name = mailInfo.visitorName
	to := mail.NewEmail(name, address)
	// fmt.Println(name, address)
	content := mail.NewContent("text/plain", fmt.Sprintf("Dear %s, \n\nThank you for taking the time to visit 21st Century Technologies. We appreciate your interest in our products and services and hope that you found the tour of our facilities informative and engaging. \nWe would love to hear about your experience at our company. Your feedback is valuable to us, and it helps us to improve our services and provide better experiences for our visitors, kindly check other services we offer on https://www.visitorsapi.com \nOnce again, thank you for visiting us, and we look forward to hearing from you soon. \n\nRegards, \nvisitorsapi Helpdesk", mailInfo.visitorName))
	m := mail.NewV3MailInit(from, subject, to, content)

	return mail.GetRequestBody(m)
}

func (srv *Server) sendMailVisitor(mailinfo mailStructure) (string, error) {
	request := sendgrid.GetRequest(srv.config.SendGridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mailVisitor(mailinfo)
	request.Body = Body
	_, err := sendgrid.API(request)
	if err != nil {
		log.Println("Error sending mail to visitor:", err)
		return "", err
	}
	// log.Println(response.StatusCode)
	// log.Println(response.Body)
	// log.Println(response.Headers)
	log.Println("Mail sent to visitor successfully")
	return "Mail sent successfully", nil
}
