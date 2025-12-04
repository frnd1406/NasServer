"""
Invoice data generator.
Creates structured invoice data using Faker.
"""

import random
from datetime import timedelta
from faker import Faker

from data_models import (
    InvoiceData, CompanyData, CustomerData, LineItem
)

fake = Faker('de_DE')

# Available product/service items
INVOICE_ITEMS = [
    "Server-Hosting Premium",
    "Cloud Storage 100GB",
    "SSL-Zertifikat",
    "Domain-Registrierung",
    "E-Mail-Postfach Business",
    "Backup-Service",
    "DDoS-Schutz",
    "Load Balancer",
    "Datenbank-Hosting",
    "CDN-Traffic",
]


def create_invoice_data() -> InvoiceData:
    """
    Generate structured invoice data.
    
    Returns:
        InvoiceData: Complete invoice with all fields populated.
    """
    # Generate line items
    num_items = random.randint(1, 5)
    items = []
    net_total = 0.0

    for i in range(num_items):
        item_name = random.choice(INVOICE_ITEMS)
        quantity = random.randint(1, 10)
        unit_price = round(random.uniform(5.99, 299.99), 2)
        item_total = round(quantity * unit_price, 2)
        net_total += item_total

        items.append(LineItem(
            position=i + 1,
            name=item_name,
            quantity=quantity,
            unit_price=unit_price,
            total=item_total,
        ))

    net_total = round(net_total, 2)
    vat = round(net_total * 0.19, 2)
    gross_total = round(net_total + vat, 2)

    invoice_date = fake.date_between(start_date='-1y', end_date='today')
    due_date = invoice_date + timedelta(days=14)

    return InvoiceData(
        invoice_number=f"RE-{fake.year()}-{random.randint(1000, 9999)}",
        invoice_date=invoice_date,
        due_date=due_date,
        company=CompanyData(
            name=fake.company(),
            street=fake.street_address(),
            postcode=fake.postcode(),
            city=fake.city(),
            tax_id=fake.vat_id().replace('DE', ''),
        ),
        customer=CustomerData(
            name=fake.name(),
            street=fake.street_address(),
            postcode=fake.postcode(),
            city=fake.city(),
        ),
        items=items,
        net_total=net_total,
        vat=vat,
        gross_total=gross_total,
        iban=fake.iban(),
        bic=fake.swift(),
    )
