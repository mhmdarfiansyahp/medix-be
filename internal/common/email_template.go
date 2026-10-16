package common

import (
	"fmt"
)

type EmailTemplate interface {
	Subject() string
	HtmlBody() string
	PlainTextBody() string
}

type ResetPasswordEmailTemplate struct {
	RecipientName string
	ResetLink     string
	ExpireMinutes int
}

func NewResetPasswordEmailTemplate(name, resetLink string) *ResetPasswordEmailTemplate {
	return &ResetPasswordEmailTemplate{
		RecipientName: name,
		ResetLink:     resetLink,
		ExpireMinutes: 30,
	}
}

func (t *ResetPasswordEmailTemplate) Subject() string {
	return "Reset Your Medix Account Password"
}

func (t *ResetPasswordEmailTemplate) PlainTextBody() string {
	return fmt.Sprintf(
		"Hello %s,\n\nClick the link below to reset your Medix account password:\n%s\n\nThis link is valid for %d minutes and can only be used once.",
		t.RecipientName,
		t.ResetLink,
		t.ExpireMinutes,
	)
}

func (t *ResetPasswordEmailTemplate) HtmlBody() string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Reset Password - Medix</title>
    <style type="text/css">
        body, table, td, a { -webkit-text-size-adjust: 100%%; -ms-text-size-adjust: 100%%; }
        table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
        img { -ms-interpolation-mode: bicubic; border: 0; outline: none; text-decoration: none; }
        body { margin: 0; padding: 0; width: 100%% !important; min-width: 100%%; height: 100%% !important; background-color: #f8fafc; }
        @media only screen and (max-width: 600px) {
            .email-container { width: 100%% !important; max-width: 100%% !important; border-radius: 0 !important; }
            .content-padding { padding: 24px 20px !important; }
            .header-padding { padding: 28px 20px !important; }
            .mobile-btn { width: 100%% !important; max-width: 100%% !important; display: block !important; box-sizing: border-box !important; }
            .btn-text { font-size: 16px !important; }
        }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; color: #1e293b;">
    <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f8fafc; padding: 20px 0;">
        <tr>
            <td align="center" style="padding: 10px;">
                <!-- Main Card Container -->
                <table class="email-container" border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 540px; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03); border: 1px solid #e2e8f0;">
                    
                    <!-- Header -->
                    <tr>
                        <td align="center" class="header-padding" style="background: linear-gradient(135deg, #2563eb, #1d4ed8); background-color: #2563eb; padding: 32px 24px;">
                            <div style="display: inline-block; background-color: rgba(255, 255, 255, 0.2); border-radius: 50%%; padding: 12px; margin-bottom: 12px;">
                                <span style="font-size: 28px; line-height: 1;">💊</span>
                            </div>
                            <div style="color: #ffffff; font-size: 20px; font-weight: 700; letter-spacing: -0.02em;">
                                Medix Pharmacy
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Content Body -->
                    <tr>
                        <td class="content-padding" style="padding: 36px 32px; text-align: left;">
                            <h1 style="margin: 0 0 16px 0; font-size: 20px; font-weight: 700; color: #0f172a; letter-spacing: -0.01em;">
                                Permintaan Reset Password
                            </h1>
                            <p style="margin: 0 0 14px 0; font-size: 15px; line-height: 1.6; color: #475569;">
                                Halo <strong style="color: #0f172a;">%s</strong>,
                            </p>
                            <p style="margin: 0 0 24px 0; font-size: 14px; line-height: 1.6; color: #475569;">
                                Kami menerima permintaan untuk mengatur ulang kata sandi akun Medix Anda. Klik tombol di bawah untuk melanjutkan:
                            </p>
                            
                            <!-- Action Button -->
                            <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="margin: 28px 0;">
                                <tr>
                                    <td align="center">
                                        <table class="mobile-btn" border="0" cellpadding="0" cellspacing="0" style="border-radius: 8px; background-color: #2563eb; width: auto;">
                                            <tr>
                                                <td align="center" style="border-radius: 8px;">
                                                    <a href="%s" target="_blank" class="mobile-btn btn-text" style="background-color: #2563eb; color: #ffffff; font-size: 15px; font-weight: 600; text-decoration: none; padding: 14px 28px; display: inline-block; border-radius: 8px; text-align: center; border: 1px solid #2563eb;">
                                                        Reset Password Sekarang
                                                    </a>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>

                            <!-- Alternative Link for mobile copy -->
                            <p style="margin: 0 0 24px 0; font-size: 12px; line-height: 1.5; color: #94a3b8; word-break: break-all;">
                                Atau salin tautan berikut ke browser Anda:<br>
                                <a href="%s" style="color: #2563eb; text-decoration: underline;">%s</a>
                            </p>

                            <!-- Security Note Box -->
                            <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #eff6ff; border-radius: 8px; border: 1px solid #dbeafe;">
                                <tr>
                                    <td style="padding: 14px 16px; font-size: 13px; color: #1e40af; line-height: 1.5;">
                                        ⏱️ Tautan ini hanya berlaku selama <strong>%d menit</strong>. Abaikan email ini bila Anda tidak pernah meminta reset kata sandi.
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td align="center" style="background-color: #f8fafc; padding: 20px 24px; border-top: 1px solid #f1f5f9;">
                            <p style="margin: 0 0 6px 0; font-size: 12px; color: #64748b; font-weight: 600;">
                                &copy; 2026 Medix Pharmacy System
                            </p>
                            <p style="margin: 0; font-size: 11px; color: #94a3b8; line-height: 1.4;">
                                Email dikirim secara otomatis. Mohon tidak membalas email ini.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`, t.RecipientName, t.ResetLink, t.ResetLink, t.ResetLink, t.ExpireMinutes)
}
