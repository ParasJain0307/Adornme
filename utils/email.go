// // // package utils

// // // import (
// // // 	"net/smtp"
// // // )

// // // // func SendResetEmail(to, link string) error {

// // // // 	from := "rj0648665@gmail.com"
// // // // 	password := "bpkobbdxjclbhwep" // use environment variables or secrets manager to store sensitive info

// // // // 	smtpHost := "smtp.gmail.com"
// // // // 	smtpPort := "587"

// // // // 	auth := smtp.PlainAuth("", from, password, smtpHost)

// // // // 	subject := "Reset Password"

// // // // 	body := fmt.Sprintf(`
// // // // Click here to reset your password:

// // // // %s

// // // // This link expires in 15 minutes.
// // // // `, link)

// // // // 	msg := []byte("Subject: " + subject + "\r\n" +
// // // // 		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
// // // // 		body)

// // // // 	err := smtp.SendMail(
// // // // 		smtpHost+":"+smtpPort,
// // // // 		auth,
// // // // 		from,
// // // // 		[]string{to},
// // // // 		msg,
// // // // 	)

// // // // 	return err

// // // func SendResetEmail(to, link string) error {

// // // 	from := "noreply.adornme@gmail.com"
// // // 	password := "rcnsshwdrxmasata"    // app password of this Gmail

// // // 	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

// // // 	msg := []byte(
// // // 		"From: Adornme Support <" + from + ">\r\n" +
// // // 			"To: " + to + "\r\n" +
// // // 			"Subject: Reset Password\r\n" +
// // // 			"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
// // // 			"Click here to reset your password:\n" + link,
// // // 	)

// // // 	err := smtp.SendMail(
// // // 		"smtp.gmail.com:587",
// // // 		auth,
// // // 		from,
// // // 		[]string{to},
// // // 		msg,
// // // 	)

// // // 	return err
// // // }

// // package utils

// // import (
// // 	"bytes"
// // 	"fmt"
// // 	"html/template"
// // 	"net/smtp"
// // )

// // // 🔹 HTML Template
// // const resetTemplate = `
// // <!DOCTYPE html>
// // <html>
// // <body style="font-family: Arial, sans-serif; background:#f4f4f4; padding:20px;">
// //   <div style="max-width:600px;margin:auto;background:#fff;padding:20px;border-radius:8px;">

// //     <h2 style="color:#333;">Reset your password</h2>

// //     <p>Hi {{.Name}},</p>

// //     <p>We received a request to reset your password for your <b>Adornme</b> account.</p>

// //     <p style="text-align:center;">
// //       <a href="{{.Link}}"
// //          style="background:#000;color:#fff;padding:12px 20px;text-decoration:none;border-radius:5px;">
// //          Reset Password
// //       </a>
// //     </p>

// //     <p>If the button doesn’t work, use this link:</p>
// //     <p>{{.Link}}</p>

// //     <p>This link will expire in <b>15 minutes</b>.</p>

// //     <hr/>

// //     <p style="font-size:12px;color:#777;">
// //       If you didn’t request this, you can safely ignore this email.
// //     </p>

// //     <p style="font-size:12px;">
// //       Need help? support@adornme.com
// //     </p>

// //     <p>— Adornme Team</p>
// //   </div>
// // </body>
// // </html>
// // `

// // // 🔹 Template Data
// // type emailData struct {
// // 	Name string
// // 	Link string
// // }

// // // 🔹 Config
// // type EmailConfig struct {
// // 	FromEmail string
// // 	Password  string
// // 	SMTPHost  string
// // 	SMTPPort  string
// // }

// // // 🔹 Load Config (for now hardcoded — move to env later)
// // func loadEmailConfig() EmailConfig {
// // 	return EmailConfig{
// // 		FromEmail: "rj0648665@gmail.com",
// // 		Password:  "bpkobbdxjclbhwep", // ⚠️ replace with env variable later
// // 		SMTPHost:  "smtp.gmail.com",
// // 		SMTPPort:  "587",
// // 	}
// // }

// // // 🔹 Send Reset Email
// // func SendResetEmail(to, name, link string) error {

// // 	config := loadEmailConfig()

// // 	// 🔹 Parse template
// // 	tmpl, err := template.New("reset").Parse(resetTemplate)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	var body bytes.Buffer
// // 	err = tmpl.Execute(&body, emailData{
// // 		Name: name,
// // 		Link: link,
// // 	})
// // 	if err != nil {
// // 		return err
// // 	}

// // 	subject := "Reset your Adornme password"

// // 	// 🔥 Correct SMTP message format (IMPORTANT FIX)
// // 	msg := fmt.Sprintf(
// // 		"From: Adornme Support <%s>\r\n"+
// // 			"To: %s\r\n"+
// // 			"Subject: %s\r\n"+
// // 			"MIME-Version: 1.0\r\n"+
// // 			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
// // 			"\r\n"+
// // 			"%s",
// // 		config.FromEmail,
// // 		to,
// // 		subject,
// // 		body.String(),
// // 	)

// // 	// 🔐 Auth
// // 	auth := smtp.PlainAuth("", config.FromEmail, config.Password, config.SMTPHost)

// // 	// 🚀 Send email
// // 	err = smtp.SendMail(
// // 		config.SMTPHost+":"+config.SMTPPort,
// // 		auth,
// // 		config.FromEmail,
// // 		[]string{to},
// // 		[]byte(msg),
// // 	)

// // 	return err
// // }

// package utils

// import (
// 	"bytes"
// 	"encoding/base64"
// 	"fmt"
// 	"html/template"
// 	"os"

// 	"net/smtp"
// )

// // 🔹 HTML TEMPLATE (Premium UI)
// const resetTemplate = `
// <!DOCTYPE html>
// <html>
// <body style="margin:0; padding:0; background:#f5f3ef; font-family: Arial, sans-serif;">

//   <table width="100%" style="padding:30px 0;">
//     <tr>
//       <td align="center">

//         <table width="600" style="background:#ffffff; border-radius:12px; padding:30px;">

//           <!-- 🔥 LOGO -->
//           <tr>
//             <td align="center" style="padding-bottom:20px;">
//               <img src="cid:logo" width="140" />
//             </td>
//           </tr>

//           <!-- Title -->
//           <tr>
//             <td align="center">
//               <h2 style="margin:0; color:#111;">Reset your password</h2>
//             </td>
//           </tr>

//           <!-- Content -->
//           <tr>
//             <td style="padding:20px; color:#555; text-align:center;">
//               Hi {{.Name}},<br><br>
//               We received a request to reset your password for your <b>Adornme</b> account.
//             </td>
//           </tr>

//           <!-- Button -->
//           <tr>
//             <td align="center" style="padding:20px;">
//               <a href="{{.Link}}"
//                  style="background:#000; color:#fff; padding:12px 25px; text-decoration:none; border-radius:6px; display:inline-block;">
//                  Reset Password
//               </a>
//             </td>
//           </tr>

//           <!-- Footer -->
//           <tr>
//             <td style="padding-top:20px; font-size:12px; color:#888; text-align:center;">
//               This link expires in 15 minutes.<br><br>
//               If you didn’t request this, you can ignore this email.
//             </td>
//           </tr>

//         </table>

//       </td>
//     </tr>
//   </table>

// </body>
// </html>
// `

// // 🔹 Template Data
// type emailData struct {
// 	Name string
// 	Link string
// }

// // 🔹 Config (Use ENV in production)
// type EmailConfig struct {
// 	FromEmail string
// 	Password  string
// 	SMTPHost  string
// 	SMTPPort  string
// }

// func loadEmailConfig() EmailConfig {
// 	return EmailConfig{
// 		FromEmail: "rj0648665@gmail.com",
// 		Password:  "bpkobbdxjclbhwep", // ⚠️ replace with env variable later
// 		SMTPHost:  "smtp.gmail.com",
// 		SMTPPort:  "587",
// 	}
// }

// // 🚀 MAIN FUNCTION
// func SendResetEmail(to, name, link string) error {

// 	config := loadEmailConfig()

// 	// 🔹 Parse template
// 	tmpl, err := template.New("reset").Parse(resetTemplate)
// 	if err != nil {
// 		return err
// 	}

// 	var body bytes.Buffer
// 	err = tmpl.Execute(&body, emailData{
// 		Name: name,
// 		Link: link,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// 🔥 Load logo
// 	logoBytes, err := os.ReadFile("../assets/logo.png")
// 	if err != nil {
// 		return err
// 	}

// 	logoBase64 := base64.StdEncoding.EncodeToString(logoBytes)

// 	boundary := "BOUNDARY-ADORNME-123"

// 	// 🔥 MIME EMAIL (HTML + IMAGE)
// 	msg := fmt.Sprintf(
// 		"From: Adornme Support <%s>\r\n"+
// 			"To: %s\r\n"+
// 			"Subject: Reset your password\r\n"+
// 			"MIME-Version: 1.0\r\n"+
// 			"Content-Type: multipart/related; boundary=%s\r\n\r\n"+

// 			"--%s\r\n"+
// 			"Content-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n\r\n"+

// 			"--%s\r\n"+
// 			"Content-Type: image/png\r\n"+
// 			"Content-Transfer-Encoding: base64\r\n"+
// 			"Content-ID: <logo>\r\n\r\n%s\r\n\r\n"+

// 			"--%s--",
// 		config.FromEmail,
// 		to,
// 		boundary,
// 		boundary,
// 		body.String(),
// 		boundary,
// 		logoBase64,
// 		boundary,
// 	)

// 	// 🔐 Auth
// 	auth := smtp.PlainAuth("", config.FromEmail, config.Password, config.SMTPHost)

// 	// 📤 Send
// 	return smtp.SendMail(
// 		config.SMTPHost+":"+config.SMTPPort,
// 		auth,
// 		config.FromEmail,
// 		[]string{to},
// 		[]byte(msg),
// 	)
// }

// package utils

// import (
// 	"bytes"
// 	"encoding/base64"
// 	"fmt"
// 	"html/template"
// 	"net/smtp"
// 	"os"
// )

// // 💎 PREMIUM HTML TEMPLATE
// const resetTemplate = `
// <!DOCTYPE html>
// <html>
// <head>
// <meta charset="UTF-8">
// <meta name="viewport" content="width=device-width, initial-scale=1.0">
// </head>

// <body style="margin:0; padding:0; background:#f5f3ef; font-family: 'Segoe UI', Arial, sans-serif;">

//   <table width="100%" cellpadding="0" cellspacing="0" style="padding:40px 0;">
//     <tr>
//       <td align="center">

//         <!-- Main Card -->
//         <table width="600" cellpadding="0" cellspacing="0"
//           style="background:#ffffff; border-radius:16px; overflow:hidden; box-shadow:0 10px 40px rgba(0,0,0,0.08);">

//           <!-- 🔥 Header -->
//           <tr>
//             <td align="center" style="padding:30px 20px; background:#f8f5f1;">
//               <img src="cid:logo" width="130" style="display:block; margin-bottom:10px;" />
//               <div style="font-size:12px; letter-spacing:3px; color:#999;">ADORNME</div>
//             </td>
//           </tr>

//           <!-- Content -->
//           <tr>
//             <td style="padding:40px 30px; text-align:center;">

//               <h2 style="margin:0; font-size:24px; color:#111;">
//                 Reset your password
//               </h2>

//               <p style="margin:20px 0; font-size:15px; color:#555; line-height:1.6;">
//                 Hi <b>{{.Name}}</b>,<br><br>
//                 We received a request to reset your password for your <b>Adornme</b> account.
//               </p>

//               <!-- 🔘 PREMIUM BUTTON -->
//               <a href="{{.Link}}"
//                  style="
//                    display:inline-block;
//                    margin-top:15px;
//                    background: linear-gradient(135deg, #000000, #333333);
//                    color:#ffffff;
//                    padding:14px 32px;
//                    border-radius:999px;
//                    text-decoration:none;
//                    font-size:14px;
//                    font-weight:600;
//                    letter-spacing:0.5px;
//                    box-shadow: 0 8px 20px rgba(0,0,0,0.25);
//                    border:1px solid rgba(0,0,0,0.2);
//                  ">
//                  Reset Password →
//               </a>

//               <!-- Divider -->
//               <div style="margin:30px 0; height:1px; background:#eee;"></div>

//               <p style="font-size:13px; color:#777;">
//                 Or copy and paste this link into your browser:
//               </p>

//               <p style="word-break:break-all; font-size:12px; color:#999;">
//                 {{.Link}}
//               </p>

//               <p style="margin-top:20px; font-size:13px; color:#777;">
//                 ⏳ This link expires in <b>15 minutes</b>.
//               </p>

//             </td>
//           </tr>

//           <!-- Footer -->
//           <tr>
//             <td style="background:#fafafa; padding:25px; text-align:center;">

//               <p style="margin:0; font-size:12px; color:#888;">
//                 If you didn’t request this, you can safely ignore this email.
//               </p>

//               <p style="margin:10px 0 0; font-size:12px; color:#aaa;">
//                 Need help? support@adornme.com
//               </p>

//               <p style="margin-top:15px; font-size:11px; color:#bbb;">
//                 © 2026 Adornme. All rights reserved.
//               </p>

//             </td>
//           </tr>

//         </table>

//       </td>
//     </tr>
//   </table>

// </body>
// </html>
// `

// // 🔹 Template Data
// type emailData struct {
// 	Name string
// 	Link string
// }

// // 🔹 Config (use ENV in real apps)
// type EmailConfig struct {
// 	FromEmail string
// 	Password  string
// 	SMTPHost  string
// 	SMTPPort  string
// }

// func loadEmailConfig() EmailConfig {
// 	return EmailConfig{
// 		FromEmail: "rj0648665@gmail.com",
// 		Password:  "bpkobbdxjclbhwep", // ⚠️ replace with env variable later
// 		SMTPHost:  "smtp.gmail.com",
// 		SMTPPort:  "587",
// 	}
// }

// // 🚀 MAIN FUNCTION
// func SendResetEmail(to, name, link string) error {

// 	config := loadEmailConfig()

// 	// Parse template
// 	tmpl, err := template.New("reset").Parse(resetTemplate)
// 	if err != nil {
// 		return err
// 	}

// 	var body bytes.Buffer
// 	err = tmpl.Execute(&body, emailData{
// 		Name: name,
// 		Link: link,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// 🔥 Load logo (IMPORTANT PATH)
// 	logoBytes, err := os.ReadFile("assets/logo.png")
// 	if err != nil {
// 		return err
// 	}

// 	logoBase64 := base64.StdEncoding.EncodeToString(logoBytes)

// 	boundary := "BOUNDARY-ADORNME-123"

// 	// 🔥 MIME EMAIL
// 	msg := fmt.Sprintf(
// 		"From: Adornme Support <%s>\r\n"+
// 			"To: %s\r\n"+
// 			"Subject: Reset your password\r\n"+
// 			"MIME-Version: 1.0\r\n"+
// 			"Content-Type: multipart/related; boundary=%s\r\n\r\n"+

// 			"--%s\r\n"+
// 			"Content-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n\r\n"+

// 			"--%s\r\n"+
// 			"Content-Type: image/png\r\n"+
// 			"Content-Transfer-Encoding: base64\r\n"+
// 			"Content-ID: <logo>\r\n\r\n%s\r\n\r\n"+

// 			"--%s--",
// 		config.FromEmail,
// 		to,
// 		boundary,
// 		boundary,
// 		body.String(),
// 		boundary,
// 		logoBase64,
// 		boundary,
// 	)

// 	auth := smtp.PlainAuth("", config.FromEmail, config.Password, config.SMTPHost)

// 	return smtp.SendMail(
// 		config.SMTPHost+":"+config.SMTPPort,
// 		auth,
// 		config.FromEmail,
// 		[]string{to},
// 		[]byte(msg),
// 	)
// }

package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

// 💎 FINAL HTML TEMPLATE (EMAIL SAFE)
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
            <td align="center" style="padding-bottom:20px;">
              <img src="cid:logo" width="130" />
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

          <!-- 🔘 BULLETPROOF BUTTON -->
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
                         letter-spacing:0.3px;
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

          <!-- Footer -->
          <tr>
            <td style="padding-top:25px; text-align:center; font-size:12px; color:#777;">
              ⏳ Link expires in <b>15 minutes</b><br><br>
              If you didn’t request this, ignore this email.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`

// 🔹 Template Data
type emailData struct {
	Name string
	Link string
}

// 🔹 Config
type EmailConfig struct {
	FromEmail string
	Password  string
	SMTPHost  string
	SMTPPort  string
}

func loadEmailConfig() EmailConfig {
	return EmailConfig{
		FromEmail: "rj0648665@gmail.com",
		Password:  "bpkobbdxjclbhwep", // ⚠️ replace with env variable later
		SMTPHost:  "smtp.gmail.com",
		SMTPPort:  "587",
	}
}

// 🚀 SEND EMAIL FUNCTION
func SendResetEmail(to, name, link string) error {

	config := loadEmailConfig()

	// Parse template
	tmpl, err := template.New("reset").Parse(resetTemplate)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, emailData{
		Name: name,
		Link: link,
	})
	if err != nil {
		return err
	}

	// 🔥 Load logo
	logoBytes, err := os.ReadFile("assets/logo.png")
	if err != nil {
		return err
	}

	logoBase64 := base64.StdEncoding.EncodeToString(logoBytes)

	boundary := "BOUNDARY-ADORNME-123"

	// 🔥 MIME MESSAGE
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
