package handlers

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/qingstor/openvpn-warder/config"
	mail "github.com/xhit/go-simple-mail/v2"
)

func connectMailServer(m *config.Email) (smtpClient *mail.SMTPClient, err error) {
	server := mail.NewSMTPClient()
	server.Host = *m.Host
	server.Port = *m.Port
	server.Username = *m.User
	server.Password = *m.Password
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	switch *m.Encryption {
	case "tls":
		server.Encryption = mail.EncryptionTLS
	case "ssl":
		server.Encryption = mail.EncryptionSSL
	case "none":
		server.Encryption = mail.EncryptionNone
	default:
		return nil, errors.New("the email encryption type not support ")
	}
	switch *m.Auth {
	case "plain":
		server.Authentication = mail.AuthPlain
	case "login":
		server.Authentication = mail.AuthLogin
	case "CRAMMD5":
		server.Authentication = mail.AuthCRAMMD5
	default:
		return nil, errors.New("the email auth type not support ")
	}

	smtpClient, err = server.Connect()
	if err != nil {
		return
	}

	return
}

func sendEmail(m *config.Email, sendTo string, content string) (err error) {
	smtpClient, err := connectMailServer(m)
	if err != nil {
		return
	}

	subject := "Auto Change Openvpn Password Notification"
	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("openvpn_warder <%s>", *m.User)).AddTo(sendTo).SetSubject(subject)
	email.SetBody(mail.TextPlain, content)

	err = email.Send(smtpClient)
	if err != nil {
		return
	}

	return
}

func makeMailContent(user, newPassword, updateAt string) string {
	return fmt.Sprintf("Your password of openvpn account %s has been changed to %s\n\nUpdate at: %s", user, newPassword, updateAt)
}

func checkMail(m *config.Email) (err error) {
	_, err = connectMailServer(m)
	return
}

func checkUserCycle(user *config.User, cycle int) (needChange bool, err error) {
	timeUpdatedAt, err := time.Parse("2006-01-02", *user.UpdatedAt)
	if err != nil {
		return
	}

	timeNow := time.Now()
	if int((timeNow.Sub(timeUpdatedAt)).Hours()/24) >= cycle {
		return true, nil
	}

	return false, nil
}

func generatePassword() (password string) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	var length int = 32
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func changeUserPassword(m *config.Email, cycle int) (err error) {
	db, err := readEntireDB()
	if err != nil {
		return err
	}

	for index, user := range db.Users {
		if *user.Ignore {
			continue
		}
		needChange, err := checkUserCycle(user, cycle)
		if err != nil {
			return err
		}
		if !needChange {
			continue
		}

		newPass := generatePassword()
		updatedAt := time.Now().Format("2006-01-02")
		user.UpdatedAt = &updatedAt
		user.Password = &newPass
		err = checkMail(m)
		if err != nil {
			return err
		}
		db.Users[index] = user
		content := makeMailContent(*user.Name, *user.Password, updatedAt)
		err = sendEmail(m, "chartoldong@yunify.com", content)
		if err != nil {
			return err
		}

		err = writeEntireDB(db)
		if err != nil {
			return err
		}
	}
	return
}

// CycleChangePassword will change user password cyclly.
func CycleChangePassword(m *config.Email, cycle int) (err error) {
	err = checkMail(m)
	if err != nil {
		return
	}
	go func() {
		for true {
			changeUserPassword(m, cycle)
			time.Sleep(24 * time.Hour)
		}
	}()
	return
}
