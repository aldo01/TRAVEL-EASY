package notifier

// bookingEmailTemplate is the HTML body for a booking confirmation email.
const bookingEmailTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; color: #333;">
  <div style="text-align: center; padding: 20px 0; border-bottom: 2px solid #2563eb;">
    <h1 style="color: #2563eb; margin: 0;">Travel Easy</h1>
    <p style="color: #666; margin: 4px 0 0;">Luggage storage & locker booking</p>
  </div>

  <div style="padding: 24px 0;">
    <h2 style="color: #111; margin: 0 0 4px;">Booking Confirmed!</h2>
    <p style="color: #666; margin: 0;">Hi {{.UserName}}, your luggage storage booking is confirmed.</p>
  </div>

  <div style="background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 8px; padding: 20px; margin-bottom: 20px;">
    <table style="width: 100%; border-collapse: collapse;">
      <tr>
        <td style="padding: 8px 0; color: #666; width: 40%;">Booking ID</td>
        <td style="padding: 8px 0; font-weight: bold; color: #111;">{{.BookingNumber}}</td>
      </tr>
      <tr>
        <td style="padding: 8px 0; color: #666;">Location</td>
        <td style="padding: 8px 0; color: #111;">{{.LocationName}}</td>
      </tr>
      {{if .LocationAddr}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Address</td>
        <td style="padding: 8px 0; color: #111;">{{.LocationAddr}}</td>
      </tr>
      {{end}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Drop-off</td>
        <td style="padding: 8px 0; color: #111;">{{.StartTime}}</td>
      </tr>
      {{if .EndTime}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Pick-up</td>
        <td style="padding: 8px 0; color: #111;">{{.EndTime}}</td>
      </tr>
      {{end}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Bags</td>
        <td style="padding: 8px 0; color: #111;">{{.Bags}}</td>
      </tr>
      <tr>
        <td style="padding: 8px 0; color: #666;">Rate type</td>
        <td style="padding: 8px 0; color: #111;">{{.RateType}}</td>
      </tr>
      <tr style="border-top: 1px solid #e5e7eb;">
        <td style="padding: 12px 0 8px; font-weight: bold; color: #111;">Total</td>
        <td style="padding: 12px 0 8px; font-weight: bold; color: #111; font-size: 18px;">₹{{printf "%.2f" .TotalPrice}}</td>
      </tr>
    </table>
  </div>

  {{if .QRCode}}
  <div style="background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 8px; padding: 16px; margin-bottom: 20px; text-align: center;">
    <p style="margin: 0 0 4px; font-weight: bold; color: #1e40af;">Your QR Code</p>
    <p style="margin: 0; color: #3b82f6; font-family: monospace; font-size: 14px;">{{.QRCode}}</p>
    <p style="margin: 8px 0 0; color: #666; font-size: 12px;">Show this code at the location to unlock your locker.</p>
  </div>
  {{end}}

  <div style="padding: 16px 0; border-top: 1px solid #e5e7eb; color: #999; font-size: 12px; text-align: center;">
    <p>Thank you for using Travel Easy.</p>
    <p>If you need to cancel, please do so before midnight the day before your drop-off.</p>
  </div>
</body>
</html>`
