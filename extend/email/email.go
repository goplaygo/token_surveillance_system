/*
@Time : 2020/7/6 6:47 下午
@Author : L
@File : email.go
@Software: GoLand
*/
package email

import (
	"net/smtp"
	"strconv"
	"token_surveillance_system/extend/conf"
)

func SendEmail(subject string,recvEmail string,emailContent string) error  {
	auth := smtp.PlainAuth("",conf.EmailConf.UserName,conf.EmailConf.Password,conf.EmailConf.Host,)

	msg := []byte("To: " + recvEmail + "\r\n" +
	"From: " + conf.EmailConf.ServName + "<" + conf.EmailConf.UserName + ">\r\n" +
	"Subject " + subject + "\r\n" + "MIME-version: 1.0;\nContent-Type: " +
	conf.EmailConf.ContentTypeHTML + ";charset=\"UTF-8\";\t\n\r\n" + emailContent,
	)
	err := smtp.SendMail(
	conf.EmailConf.Host + ":" +strconv.Itoa(conf.EmailConf.Port),
	auth,
	conf.EmailConf.UserName,
	[]string{recvEmail},
	msg,
	)
	if err != nil {
		return err
	}
	return nil
}
