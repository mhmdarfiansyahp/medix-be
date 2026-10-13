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
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Reset Password - Medix</title>
</head>
<body bgcolor="#f1f5f9" style="margin: 0; padding: 0; background-color: #f1f5f9; font-family: Arial, sans-serif;">
    <table border="0" cellpadding="0" cellspacing="0" width="100%%" bgcolor="#f1f5f9" style="background-color: #f1f5f9; padding: 40px 0;">
        <tr>
            <td align="center" style="padding: 0;">
                <table border="0" cellpadding="0" cellspacing="0" width="600" bgcolor="#ffffff" style="background-color: #ffffff;">
                    <!-- Header -->
                    <tr>
                        <td align="center" bgcolor="#2563eb" style="background-color: #2563eb; padding: 36px 24px;">
                            <table border="0" cellpadding="0" cellspacing="0">
                                <tr>
                                    <td align="center" bgcolor="#1d4ed8" width="48" height="48" style="background-color: #1d4ed8; font-size: 24px; color: #ffffff; text-align: center;">
                                        💊
                                    </td>
                                </tr>
                            </table>
                            <div style="color: #ffffff; font-size: 22px; font-weight: bold; margin-top: 12px; text-align: center;">
                                Medix Pharmacy System
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Content Body -->
                    <tr>
                        <td style="padding: 40px 32px 32px 32px; text-align: left;">
                            <div style="margin: 0 0 24px 0; font-size: 20px; font-weight: bold; color: #0f172a; text-align: center;">
                                Password Reset Request
                            </div>
                            <p style="margin: 0 0 16px 0; font-size: 15px; line-height: 20px; color: #475569;">
                                Hello <strong style="color: #1e293b;">%s</strong>,
                            </p>
                            <p style="margin: 0 0 24px 0; font-size: 15px; line-height: 20px; color: #475569;">
                                We received a request to reset your Medix account password. If you did not make this request, you can safely ignore this email.
                            </p>
                            
                            <!-- Action Button -->
                            <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="margin: 32px 0;">
                                <tr>
                                    <td align="center">
                                        <table border="0" cellpadding="0" cellspacing="0">
                                            <tr>
                                                <td align="center" bgcolor="#2563eb" style="background-color: #2563eb; padding: 14px 28px;">
                                                    <a href="%s" style="color: #ffffff; font-size: 15px; font-weight: bold; text-decoration: none; text-align: center; display: block;">
                                                        Reset Password
                                                    </a>
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                            </table>

                            <!-- Security Note Box -->
                            <table border="0" cellpadding="0" cellspacing="0" width="100%%" bgcolor="#eff6ff" style="background-color: #eff6ff;">
                                <tr>
                                    <td style="padding: 16px; font-size: 13px; color: #1e40af; line-height: 18px;">
                                        <strong>Security Note:</strong> This link is valid for <strong>%d minutes</strong> and can only be used once. Never share this link with anyone.
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td align="center" bgcolor="#f8fafc" style="background-color: #f8fafc; padding: 24px 32px;">
                            <p style="margin: 0 0 8px 0; font-size: 12px; color: #64748b; font-weight: bold;">
                                &copy; 2026 Medix Pharmacy System. All rights reserved.
                            </p>
                            <p style="margin: 0; font-size: 11px; color: #94a3b8;">
                                This email was sent automatically by the system. Please do not reply to this email.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`, t.RecipientName, t.ResetLink, t.ExpireMinutes)
}
