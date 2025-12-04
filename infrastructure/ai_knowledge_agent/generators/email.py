"""
Email data generator.
Creates structured email data using Faker.
"""

import random
from faker import Faker

from data_models import EmailData

fake = Faker('de_DE')

# Email subject templates
SUBJECTS = [
    "Anfrage zu Serverkosten",
    "Rückfrage zur letzten Rechnung",
    "Terminvereinbarung",
    "Projektupdate",
    "Technisches Problem",
    "Angebot für Cloud-Infrastruktur",
    "Meeting-Zusammenfassung",
    "Frage zur Datenmigration",
    "Vertragsverlängerung",
]

# Email body templates (will be formatted with names)
BODY_TEMPLATES = [
    """Sehr geehrte/r {recipient_first},

ich habe eine Frage bezüglich der monatlichen Serverkosten. Könnten Sie mir eine Aufschlüsselung der Kosten für den letzten Monat zusenden?

Besonders interessieren mich:
- Hosting-Gebühren
- Traffic-Kosten
- Backup-Services

Vielen Dank im Voraus!

Mit freundlichen Grüßen
{sender}
{sender_email}""",

    """Hallo {recipient_first},

bezüglich unseres gestrigen Gesprächs möchte ich noch einmal zusammenfassen:

1. Migration der Datenbank auf PostgreSQL bis Ende des Monats
2. Setup der Backup-Strategie (täglich inkrementell)
3. Performance-Tests der API durchführen

Bitte bestätigen Sie, dass Sie mit diesen Punkten einverstanden sind.

Beste Grüße
{sender}""",

    """Guten Tag,

wir haben ein technisches Problem mit unserem Server festgestellt. Die Response-Zeiten haben sich in den letzten 24 Stunden deutlich verschlechtert.

Symptome:
- API antwortet mit >2 Sekunden Verzögerung
- Datenbank-Queries sind langsam
- CPU-Auslastung bei konstant 85%

Könnten Sie das bitte dringend prüfen?

Danke
{sender}""",
]


def create_email_data() -> EmailData:
    """
    Generate structured email data.
    
    Returns:
        EmailData: Complete email with all fields populated.
    """
    sender_name = fake.name()
    sender_email = fake.email()
    recipient_name = fake.name()
    recipient_email = fake.email()

    # Get first name for greeting
    recipient_first = recipient_name.split()[0]

    # Select and format body
    body_template = random.choice(BODY_TEMPLATES)
    body = body_template.format(
        recipient_first=recipient_first,
        sender=sender_name,
        sender_email=sender_email,
    )

    # Subject with optional ticket number
    subject = random.choice(SUBJECTS)
    if random.random() < 0.2:  # 20% chance of ticket reference
        subject = f"Support-Ticket #{random.randint(10000, 99999)}"

    return EmailData(
        sender_name=sender_name,
        sender_email=sender_email,
        recipient_name=recipient_name,
        recipient_email=recipient_email,
        subject=subject,
        sent_date=fake.date_time_between(start_date='-60d', end_date='now'),
        body=body,
    )
