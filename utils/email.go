package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

// 💎 FINAL PREMIUM TEMPLATE
const resetTemplate = `
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>

<body style="margin:0; padding:0; background:#f5f3ef; font-family: Arial, sans-serif;">

  <table width="100%" cellpadding="0" cellspacing="0" style="padding:40px 0;">
    <tr>
      <td align="center">

        <!-- Card -->
        <table width="600" cellpadding="0" cellspacing="0"
          style="background:#ffffff; border-radius:12px; padding:30px;">

          <!-- Logo -->
          <tr>
            <td align="center" style="padding-bottom:15px;">
              <img src="cid:logo" width="150" style="display:block;" />
              <p style="font-size:11px; letter-spacing:3px; color:#999; margin-top:5px;">
                ADORNME
              </p>
            </td>
          </tr>

          <!-- Title -->
          <tr>
            <td align="center">
              <h2 style="margin:0; color:#111;">Reset your password</h2>
            </td>
          </tr>

          <!-- Text -->
          <tr>
            <td style="padding:20px; text-align:center; color:#555;">
              Hi <b>{{.Name}}</b>,<br><br>
              We received a request to reset your password for your <b>Adornme</b> account.
            </td>
          </tr>

          <!-- Button -->
          <tr>
            <td align="center">
              <table cellpadding="0" cellspacing="0" border="0" style="margin-top:10px;">
                <tr>
                  <td align="center" bgcolor="#000000" 
                      style="border-radius:10px; border:1px solid #222;">
                    
                    <a href="{{.Link}}"
                       style="
                         display:inline-block;
                         padding:14px 32px;
                         font-size:14px;
                         font-weight:600;
                         color:#ffffff;
                         text-decoration:none;
                         border-radius:10px;
                       ">
                       Reset Password →
                    </a>

                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Link fallback -->
          <tr>
            <td style="padding:25px 20px 10px; text-align:center; font-size:12px; color:#888;">
              Or copy this link:
            </td>
          </tr>

          <tr>
            <td style="word-break:break-all; text-align:center; font-size:11px; color:#999;">
              {{.Link}}
            </td>
          </tr>

          <!-- 🔥 INTERACTIVE SOCIAL SECTION -->
          <tr>
            <td align="center" style="padding-top:30px;">

              <div style="width:40px;height:2px;background:#ddd;margin-bottom:20px;"></div>

              <p style="font-size:13px; color:#555; margin-bottom:15px;">
                Connect with us
              </p>

              <table cellpadding="0" cellspacing="0" border="0">
                <tr>

                  <!-- Instagram -->
                  <td style="padding:0 6px;">
                    <a href="https://instagram.com/adornme">
                      <div style="background:#000;border-radius:50%;width:36px;height:36px;line-height:36px;text-align:center;">
                        <img src="https://cdn-icons-png.flaticon.com/512/2111/2111463.png" width="16" />
                      </div>
                    </a>
                  </td>

                  <!-- X -->
                  <td style="padding:0 6px;">
                    <a href="https://twitter.com/adornme">
                      <div style="background:#000;border-radius:50%;width:36px;height:36px;line-height:36px;text-align:center;">
                        <img src="https://cdn-icons-png.flaticon.com/512/5968/5968830.png" width="16" />
                      </div>
                    </a>
                  </td>

                  <!-- Facebook -->
                  <td style="padding:0 6px;">
                    <a href="https://facebook.com/adornme">
                      <div style="background:#000;border-radius:50%;width:36px;height:36px;line-height:36px;text-align:center;">
                        <img src="https://cdn-icons-png.flaticon.com/512/733/733547.png" width="16" />
                      </div>
                    </a>
                  </td>

                  <!-- LinkedIn -->
                  <td style="padding:0 6px;">
                    <a href="https://linkedin.com/company/adornme">
                      <div style="background:#000;border-radius:50%;width:36px;height:36px;line-height:36px;text-align:center;">
                        <img src="https://cdn-icons-png.flaticon.com/512/145/145807.png" width="16" />
                      </div>
                    </a>
                  </td>

                  <!-- YouTube -->
                  <td style="padding:0 6px;">
                    <a href="https://youtube.com/@adornme">
                      <div style="background:#000;border-radius:50%;width:36px;height:36px;line-height:36px;text-align:center;">
                        <img src="https://cdn-icons-png.flaticon.com/512/1384/1384060.png" width="16" />
                      </div>
                    </a>
                  </td>

                  <!-- Website -->
                  <td style="padding:0 6px;">
                    <a href="https://adornme.com">
                      <div style="background:#000;border-radius:50%;width:36px;height:36px;line-height:36px;text-align:center;">
                        <img src="https://cdn-icons-png.flaticon.com/512/841/841364.png" width="16" />
                      </div>
                    </a>
                  </td>

                </tr>
              </table>

              <p style="font-size:12px; color:#888; margin-top:15px;">
                Follow us for updates, offers & new arrivals ✨
              </p>

            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="padding-top:15px; text-align:center; font-size:12px; color:#777;">
              ⏳ Link expires in <b>15 minutes</b><br><br>
              If you didn’t request this, ignore this email.
              <br><br>
              <span style="font-size:11px; color:#aaa;">
                © 2026 Adornme. All rights reserved.
              </span>
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`

type emailData struct {
	Name string
	Link string
}

type EmailConfig struct {
	FromEmail string
	Password  string
	SMTPHost  string
	SMTPPort  string
}

func loadEmailConfig() EmailConfig {
	return EmailConfig{
		FromEmail: "rj0648665@gmail.com",
		Password:  "bpkobbdxjclbhwep", // ⚠️ move to ENV
		SMTPHost:  "smtp.gmail.com",
		SMTPPort:  "587",
	}
}

func SendResetEmail(to, name, link string) error {

	config := loadEmailConfig()

	tmpl, err := template.New("reset").Parse(resetTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, emailData{Name: name, Link: link})
	if err != nil {
		return err
	}

	logoBytes, err := os.ReadFile("assets/logo.png")
	if err != nil {
		return err
	}

	logoBase64 := base64.StdEncoding.EncodeToString(logoBytes)
	boundary := "BOUNDARY-ADORNME-123"

	msg := fmt.Sprintf(
		"From: Adornme Support <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: Reset your password\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/related; boundary=%s\r\n\r\n"+

			"--%s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n\r\n"+

			"--%s\r\n"+
			"Content-Type: image/png\r\n"+
			"Content-Transfer-Encoding: base64\r\n"+
			"Content-ID: <logo>\r\n\r\n%s\r\n\r\n"+

			"--%s--",
		config.FromEmail,
		to,
		boundary,
		boundary,
		body.String(),
		boundary,
		logoBase64,
		boundary,
	)

	auth := smtp.PlainAuth("", config.FromEmail, config.Password, config.SMTPHost)

	return smtp.SendMail(
		config.SMTPHost+":"+config.SMTPPort,
		auth,
		config.FromEmail,
		[]string{to},
		[]byte(msg),
	)
}
