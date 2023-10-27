package mail_consumer

import (
	"encoding/json"
	"fmt"
	base_consumer "main/consumer"
	"main/models"
	mail_repository "main/repositories/email"
)

func New() {
	event := base_consumer.New("user_registered")
	event.SubcribeEvent(sendMailListener)
}

func sendMailListener(recieved []byte) {
	var data models.User
	err := json.Unmarshal(recieved, &data)
	if err != nil {
		fmt.Printf("error when decode json data is: %v\n", err)
		return
	}
	if err := mail_repository.SendMailToUser(&models.Email{
		SendTo:  data.Email,
		Content: "Crm verify user code: " + data.OtpCode,
		Subject: "",
	}); err != nil {
		fmt.Printf("Error when send mail is: %v %v\n", err, data)
	}
}
