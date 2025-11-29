#!/usr/bin/env python3
"""
Realistic German document corpus generator using Faker.
Generates invoices, tech logs, and emails for semantic search testing.
"""

import os
import random
from datetime import datetime, timedelta
from faker import Faker

# Initialize Faker with German locale
fake = Faker('de_DE')

# Output directory
OUTPUT_DIR = "/mnt/data/test_corpus"


def generate_invoice():
    """Generate a realistic German invoice."""
    invoice_number = f"RE-{fake.year()}-{random.randint(1000, 9999)}"
    company_name = fake.company()
    customer_name = fake.name()
    customer_address = f"{fake.street_address()}, {fake.postcode()} {fake.city()}"

    # Generate line items
    num_items = random.randint(1, 5)
    items = []
    total = 0

    for i in range(num_items):
        item_name = random.choice([
            "Server-Hosting Premium",
            "Cloud Storage 100GB",
            "SSL-Zertifikat",
            "Domain-Registrierung",
            "E-Mail-Postfach Business",
            "Backup-Service",
            "DDoS-Schutz",
            "Load Balancer",
            "Datenbank-Hosting",
            "CDN-Traffic"
        ])
        quantity = random.randint(1, 10)
        unit_price = round(random.uniform(5.99, 299.99), 2)
        item_total = round(quantity * unit_price, 2)
        total += item_total

        items.append(f"{i+1}. {item_name} - {quantity}x à {unit_price:.2f}€ = {item_total:.2f}€")

    vat = round(total * 0.19, 2)
    total_with_vat = round(total + vat, 2)

    invoice_date = fake.date_between(start_date='-1y', end_date='today')
    due_date = invoice_date + timedelta(days=14)

    iban = fake.iban()

    content = f"""RECHNUNG

Rechnungsnummer: {invoice_number}
Rechnungsdatum: {invoice_date.strftime('%d.%m.%Y')}
Fälligkeitsdatum: {due_date.strftime('%d.%m.%Y')}

Von:
{company_name}
{fake.street_address()}
{fake.postcode()} {fake.city()}
Steuernummer: {fake.vat_id().replace('DE', '')}

An:
{customer_name}
{customer_address}

Positionen:
{''.join(f'{item}' for item in items)}

Nettobetrag: {total:.2f}€
MwSt. (19%): {vat:.2f}€
━━━━━━━━━━━━━━━━━━━━━━━
Gesamtbetrag: {total_with_vat:.2f}€

Zahlungsinformationen:
IBAN: {iban}
BIC: {fake.swift()}
Verwendungszweck: {invoice_number}

Bitte überweisen Sie den Betrag bis zum Fälligkeitsdatum.

Mit freundlichen Grüßen
{company_name}
"""
    return content, f"rechnung_{invoice_number.replace('-', '_').lower()}.txt"


def generate_tech_log():
    """Generate realistic technical log entries."""
    server_name = f"{random.choice(['web', 'db', 'api', 'cache', 'worker'])}-{random.randint(1, 10)}"
    log_date = fake.date_time_between(start_date='-30d', end_date='now')

    log_entries = []
    num_entries = random.randint(10, 30)

    log_types = [
        ("INFO", "Service started successfully"),
        ("INFO", "Health check passed"),
        ("INFO", "Request processed in {}ms"),
        ("WARN", "High memory usage detected: {}%"),
        ("WARN", "Slow query detected: {} seconds"),
        ("ERROR", "Connection timeout to database"),
        ("ERROR", "Failed to write to disk: {}"),
        ("DEBUG", "Cache hit rate: {}%"),
        ("DEBUG", "Active connections: {}"),
    ]

    current_time = log_date
    for _ in range(num_entries):
        level, message = random.choice(log_types)

        # Fill in placeholders
        if '{}' in message:
            if 'ms' in message:
                message = message.format(random.randint(10, 5000))
            elif '%' in message:
                message = message.format(random.randint(50, 95))
            elif 'seconds' in message:
                message = message.format(round(random.uniform(1.0, 10.0), 2))
            elif 'connections' in message:
                message = message.format(random.randint(10, 500))
            else:
                message = message.format("unknown error")

        log_entries.append(
            f"[{current_time.strftime('%Y-%m-%d %H:%M:%S')}] [{level}] [{server_name}] {message}"
        )
        current_time += timedelta(seconds=random.randint(1, 300))

    content = f"""Server Log: {server_name}
Generated: {log_date.strftime('%Y-%m-%d')}

{''.join(f'{entry}' for entry in log_entries)}

Total Entries: {num_entries}
Server Uptime: {random.randint(1, 720)} hours
"""
    return content, f"log_{server_name}_{log_date.strftime('%Y%m%d')}.txt"


def generate_email():
    """Generate a realistic German business email."""
    sender = fake.name()
    sender_email = fake.email()
    recipient = fake.name()
    recipient_email = fake.email()

    subjects = [
        "Anfrage zu Serverkosten",
        "Rückfrage zur letzten Rechnung",
        "Terminvereinbarung",
        "Projektupdate",
        "Technisches Problem",
        "Angebot für Cloud-Infrastruktur",
        "Meeting-Zusammenfassung",
        "Frage zur Datenmigration",
        "Vertragsverlängerung",
        "Support-Ticket #{}".format(random.randint(10000, 99999))
    ]

    subject = random.choice(subjects)
    email_date = fake.date_time_between(start_date='-60d', end_date='now')

    bodies = [
        f"""Sehr geehrte/r {recipient.split()[0]},

ich habe eine Frage bezüglich der monatlichen Serverkosten. Könnten Sie mir eine Aufschlüsselung der Kosten für den letzten Monat zusenden?

Besonders interessieren mich:
- Hosting-Gebühren
- Traffic-Kosten
- Backup-Services

Vielen Dank im Voraus!

Mit freundlichen Grüßen
{sender}
{sender_email}""",

        f"""Hallo {recipient.split()[0]},

bezüglich unseres gestrigen Gesprächs möchte ich noch einmal zusammenfassen:

1. Migration der Datenbank auf PostgreSQL bis Ende des Monats
2. Setup der Backup-Strategie (täglich inkrementell)
3. Performance-Tests der API durchführen

Bitte bestätigen Sie, dass Sie mit diesen Punkten einverstanden sind.

Beste Grüße
{sender}""",

        f"""Guten Tag,

wir haben ein technisches Problem mit unserem Server festgestellt. Die Response-Zeiten haben sich in den letzten 24 Stunden deutlich verschlechtert.

Symptome:
- API antwortet mit >2 Sekunden Verzögerung
- Datenbank-Queries sind langsam
- CPU-Auslastung bei konstant 85%

Könnten Sie das bitte dringend prüfen?

Danke
{sender}"""
    ]

    body = random.choice(bodies)

    content = f"""Von: {sender} <{sender_email}>
An: {recipient} <{recipient_email}>
Betreff: {subject}
Datum: {email_date.strftime('%d.%m.%Y %H:%M')}

{body}
"""
    return content, f"email_{email_date.strftime('%Y%m%d_%H%M%S')}.txt"


def main():
    """Generate corpus of 50 realistic German documents."""
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    print(f"Generating realistic German document corpus in {OUTPUT_DIR}...")
    print("=" * 60)

    generated = 0
    target = 50

    # Distribution: 40% invoices, 30% logs, 30% emails
    num_invoices = int(target * 0.4)
    num_logs = int(target * 0.3)
    num_emails = target - num_invoices - num_logs

    # Generate invoices
    print(f"\n📄 Generating {num_invoices} invoices...")
    for i in range(num_invoices):
        content, filename = generate_invoice()
        filepath = os.path.join(OUTPUT_DIR, filename)
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        generated += 1
        if (i + 1) % 5 == 0:
            print(f"  ✓ {i + 1}/{num_invoices} invoices generated")

    # Generate tech logs
    print(f"\n📊 Generating {num_logs} tech logs...")
    for i in range(num_logs):
        content, filename = generate_tech_log()
        filepath = os.path.join(OUTPUT_DIR, filename)
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        generated += 1
        if (i + 1) % 5 == 0:
            print(f"  ✓ {i + 1}/{num_logs} logs generated")

    # Generate emails
    print(f"\n📧 Generating {num_emails} emails...")
    for i in range(num_emails):
        content, filename = generate_email()
        filepath = os.path.join(OUTPUT_DIR, filename)
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        generated += 1
        if (i + 1) % 5 == 0:
            print(f"  ✓ {i + 1}/{num_emails} emails generated")

    print("\n" + "=" * 60)
    print(f"✅ Successfully generated {generated} documents!")
    print(f"📁 Location: {OUTPUT_DIR}")
    print("\nBreakdown:")
    print(f"  • Invoices: {num_invoices}")
    print(f"  • Tech Logs: {num_logs}")
    print(f"  • Emails: {num_emails}")


if __name__ == "__main__":
    main()
