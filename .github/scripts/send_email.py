import smtplib
import os
import sys
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.base import MIMEBase
from email import encoders

def send_email(server, port, username, password, to, subject, body_file, attachment):
    msg = MIMEMultipart()
    msg["Subject"] = subject
    msg["From"] = username
    msg["To"] = to

    with open(body_file, "r") as f:
        msg.attach(MIMEText(f.read(), "plain", "utf-8"))

    if os.path.isfile(attachment):
        with open(attachment, "rb") as f:
            part = MIMEBase("application", "octet-stream")
            part.set_payload(f.read())
            encoders.encode_base64(part)
            part.add_header("Content-Disposition", "attachment", filename=os.path.basename(attachment))
            msg.attach(part)
    else:
        print(f"Warning: attachment not found: {attachment}", file=sys.stderr)

    server_conn = smtplib.SMTP(server, port)
    server_conn.ehlo()
    server_conn.starttls()
    server_conn.login(username, password)
    server_conn.send_message(msg)
    server_conn.quit()
    print("Email sent successfully.")

if __name__ == "__main__":
    send_email(
        server=os.environ["MAIL_SERVER_ADDRESS"],
        port=int(os.environ["MAIL_SERVER_PORT"]),
        username=os.environ["MAIL_USERNAME"],
        password=os.environ["MAIL_PASSWORD"],
        to=os.environ["MAIL_TO"],
        subject=os.environ["MAIL_SUBJECT"],
        body_file=os.environ["MAIL_BODY_FILE"],
        attachment=os.environ.get("MAIL_ATTACHMENT", ""),
    )
