#!/usr/bin/env python3
"""
Extended Family Document Generator - Creates realistic German document corpus.

Categories:
- Finanzen: Steuern, Krypto, Rechnungen, Kontoauszüge
- Kinder: Zeugnisse, Arztberichte, Schulbriefe
- Familie: Versicherungen, Verträge, Reisen
- Arbeit: Gehaltsabrechnungen, Arbeitsverträge

Usage:
    python generate_corpus.py --all --output /mnt/data
    python generate_corpus.py --finance 20 --children 15 --output ./docs
"""

import argparse
import json
import os
import random
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Any

# German names and data
FIRST_NAMES_MALE = ["Max", "Leon", "Paul", "Lukas", "Felix", "Noah", "Elias", "Ben", "Jonas", "Finn"]
FIRST_NAMES_FEMALE = ["Emma", "Mia", "Hannah", "Sofia", "Lina", "Marie", "Lea", "Anna", "Laura", "Lisa"]
LAST_NAMES = ["Müller", "Schmidt", "Schneider", "Fischer", "Weber", "Meyer", "Wagner", "Becker", "Schulz", "Hoffmann"]
CITIES = ["Berlin", "Hamburg", "München", "Köln", "Frankfurt", "Stuttgart", "Düsseldorf", "Leipzig", "Dresden", "Hannover"]
STREETS = ["Hauptstraße", "Bahnhofstraße", "Schulstraße", "Gartenstraße", "Bergstraße", "Lindenstraße", "Kirchstraße"]

# Companies
COMPANIES = [
    "Deutsche Telekom AG", "SAP SE", "Siemens AG", "Allianz SE", "BMW AG",
    "Amazon Deutschland", "Google Germany GmbH", "Microsoft Deutschland",
    "Bosch GmbH", "Continental AG", "BASF SE", "Volkswagen AG"
]

# Crypto exchanges and coins
CRYPTO_EXCHANGES = ["Binance", "Coinbase", "Kraken", "Bitpanda", "Trade Republic", "Bitstamp"]
CRYPTO_COINS = ["Bitcoin (BTC)", "Ethereum (ETH)", "Solana (SOL)", "Cardano (ADA)", "Polkadot (DOT)", "Ripple (XRP)"]

# School subjects
SCHOOL_SUBJECTS = {
    "Deutsch": ["Sprachverständnis", "Aufsatz", "Grammatik", "Rechtschreibung"],
    "Mathematik": ["Grundrechenarten", "Geometrie", "Algebra", "Textaufgaben"],
    "Englisch": ["Vokabeln", "Grammatik", "Leseverständnis", "Sprechen"],
    "Sachunterricht": ["Natur", "Gesellschaft", "Technik"],
    "Sport": ["Ausdauer", "Ballsport", "Turnen"],
    "Kunst": ["Malen", "Basteln", "Kreativität"],
    "Musik": ["Singen", "Rhythmus", "Instrumentenkunde"],
}

GRADES = ["1", "2", "3", "4", "5", "6"]
GRADE_DESCRIPTIONS = {
    "1": "sehr gut",
    "2": "gut",
    "3": "befriedigend",
    "4": "ausreichend",
    "5": "mangelhaft",
    "6": "ungenügend"
}

# Medical terms
DIAGNOSES_CHILDREN = [
    "Akute Bronchitis", "Otitis media (Mittelohrentzündung)", "Gastroenteritis",
    "Tonsillitis (Mandelentzündung)", "Allergische Rhinitis", "Windpocken (Varizellen)",
    "Scharlach", "Hautkontusion (Prellung)", "Fieberkrampf", "ADHS-Verdacht"
]

MEDICATIONS = [
    "Paracetamol-Saft 200mg", "Ibuprofen Junior", "Amoxicillin 250mg",
    "Cetirizin-Tropfen", "NaCl Nasenspray", "ACC Kindersaft"
]


class FamilyDocumentGenerator:
    """Generates realistic German family documents organized in folders."""
    
    def __init__(self, output_dir: str):
        self.output_dir = Path(output_dir)
        self.family_name = random.choice(LAST_NAMES)
        self.family_members = self._generate_family()
        self.generated_files = []
        
    def _generate_family(self) -> Dict:
        """Generate a realistic German family."""
        father_name = random.choice(FIRST_NAMES_MALE)
        mother_name = random.choice(FIRST_NAMES_FEMALE)
        
        children = []
        num_children = random.randint(1, 3)
        for i in range(num_children):
            gender = random.choice(["male", "female"])
            name = random.choice(FIRST_NAMES_MALE if gender == "male" else FIRST_NAMES_FEMALE)
            age = random.randint(6, 16)
            children.append({
                "name": name,
                "gender": gender,
                "birth_year": datetime.now().year - age,
                "age": age,
                "school_class": min(age - 5, 10)
            })
        
        return {
            "last_name": self.family_name,
            "address": f"{random.choice(STREETS)} {random.randint(1, 150)}, {random.randint(10000, 99999)} {random.choice(CITIES)}",
            "father": {"name": father_name, "birth_year": random.randint(1975, 1990)},
            "mother": {"name": mother_name, "birth_year": random.randint(1977, 1992)},
            "children": children
        }

    def _ensure_dir(self, subdir: str) -> Path:
        """Ensure directory exists and return path."""
        path = self.output_dir / subdir
        path.mkdir(parents=True, exist_ok=True)
        return path

    def _save_document(self, subdir: str, filename: str, content: str) -> str:
        """Save document and track it."""
        dir_path = self._ensure_dir(subdir)
        filepath = dir_path / filename
        
        # Ensure unique filename
        counter = 1
        base, ext = os.path.splitext(filename)
        while filepath.exists():
            filename = f"{base}_{counter}{ext}"
            filepath = dir_path / filename
            counter += 1
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        
        self.generated_files.append({
            "path": str(filepath),
            "category": subdir,
            "filename": filename
        })
        
        return str(filepath)

    # ========== FINANCE DOCUMENTS ==========
    
    def generate_tax_document(self) -> str:
        """Generate a tax-related document."""
        year = random.randint(2021, 2024)
        doc_types = ["steuerbescheid", "vorauszahlung", "einkommensteuererklarung"]
        doc_type = random.choice(doc_types)
        
        gross_income = random.randint(45000, 120000)
        tax_amount = int(gross_income * random.uniform(0.25, 0.35))
        
        content = f"""FINANZAMT {random.choice(CITIES).upper()}
Steuernummer: {random.randint(100, 999)}/{random.randint(1000, 9999)}/{random.randint(10000, 99999)}

EINKOMMENSTEUERBESCHEID {year}

Steuerpflichtiger: {self.family_members['father']['name']} {self.family_name}
Anschrift: {self.family_members['address']}

Ermittlung des zu versteuernden Einkommens:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Einkünfte aus nichtselbständiger Arbeit:    {gross_income:,.2f} EUR
Werbungskosten (Pauschbetrag):              -1.200,00 EUR
Sonderausgaben:                             -{random.randint(1000, 5000):,.2f} EUR
Vorsorgeaufwendungen:                       -{random.randint(2000, 8000):,.2f} EUR
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Zu versteuerndes Einkommen:                 {gross_income - random.randint(10000, 20000):,.2f} EUR

Festgesetzte Einkommensteuer:               {tax_amount:,.2f} EUR
Solidaritätszuschlag:                       {int(tax_amount * 0.055):,.2f} EUR
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Gesamtbetrag:                               {int(tax_amount * 1.055):,.2f} EUR

Bereits gezahlte Lohnsteuer:                {int(tax_amount * random.uniform(0.9, 1.1)):,.2f} EUR

{'Nachzahlung' if random.random() > 0.5 else 'Erstattung'}: {abs(random.randint(500, 3000)):,.2f} EUR

Zahlungsfrist: {(datetime.now() + timedelta(days=30)).strftime('%d.%m.%Y')}

Rechtsbehelfsbelehrung:
Gegen diesen Bescheid kann innerhalb eines Monats nach Bekanntgabe 
Einspruch eingelegt werden.

gez. Sachbearbeiter/in
Finanzamt {random.choice(CITIES)}
"""
        
        filename = f"steuerbescheid_{year}_{self.family_name.lower()}.txt"
        return self._save_document("finanzen/steuern", filename, content)

    def generate_crypto_statement(self) -> str:
        """Generate crypto trading statement."""
        exchange = random.choice(CRYPTO_EXCHANGES)
        year = random.randint(2022, 2024)
        month = random.randint(1, 12)
        
        trades = []
        total_profit = 0
        for _ in range(random.randint(5, 15)):
            coin = random.choice(CRYPTO_COINS)
            buy_price = random.uniform(100, 50000)
            amount = random.uniform(0.01, 2.0)
            sell_price = buy_price * random.uniform(0.7, 1.5)
            profit = (sell_price - buy_price) * amount
            total_profit += profit
            trades.append({
                "coin": coin,
                "amount": amount,
                "buy": buy_price,
                "sell": sell_price,
                "profit": profit
            })
        
        trade_lines = ""
        for t in trades:
            trade_lines += f"{t['coin']:20} | {t['amount']:.4f} | Kauf: {t['buy']:,.2f}€ | Verkauf: {t['sell']:,.2f}€ | {'+' if t['profit'] > 0 else ''}{t['profit']:,.2f}€\n"
        
        content = f"""{exchange.upper()} - KONTOAUSZUG KRYPTOWÄHRUNGEN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Kontoinhaber: {self.family_members['father']['name']} {self.family_name}
Kunden-ID: {random.randint(100000, 999999)}
Zeitraum: {month:02d}/{year}

HANDELSÜBERSICHT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Kryptowährung        | Menge    | Kaufkurs        | Verkaufskurs    | Gewinn/Verlust
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{trade_lines}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
GESAMTERGEBNIS: {'+' if total_profit > 0 else ''}{total_profit:,.2f} EUR

HINWEIS ZUR STEUERPFLICHT:
Kryptowährungsgewinne sind nach § 23 EStG steuerpflichtig, wenn zwischen
Anschaffung und Veräußerung weniger als ein Jahr liegt. Bitte konsultieren
Sie Ihren Steuerberater für die korrekte Angabe in Ihrer Steuererklärung.

Spekulationsfrist beachten: Gewinne nach 1 Jahr Haltedauer sind steuerfrei.

Generiert am: {datetime.now().strftime('%d.%m.%Y %H:%M')}
{exchange} GmbH - Kryptobörse Deutschland
"""
        
        filename = f"krypto_auszug_{exchange.lower()}_{year}_{month:02d}.txt"
        return self._save_document("finanzen/krypto", filename, content)

    def generate_invoice(self) -> str:
        """Generate a business invoice."""
        company = random.choice(COMPANIES)
        invoice_num = f"RE-{random.randint(2024000, 2024999)}"
        amount = random.randint(50, 2000)
        vat = amount * 0.19
        
        items = [
            ("Software-Lizenz Premium Paket", random.randint(200, 800)),
            ("Cloud-Speicher 1TB monatlich", random.randint(10, 50)),
            ("IT-Support Servicepauschale", random.randint(100, 300)),
            ("Hardware Upgrade Kit", random.randint(150, 500)),
            ("Schulung Online-Kurs", random.randint(80, 250)),
        ]
        selected_items = random.sample(items, random.randint(1, 3))
        
        item_lines = ""
        subtotal = 0
        for item, price in selected_items:
            item_lines += f"  {item:40} {price:>10.2f} EUR\n"
            subtotal += price
        
        content = f"""{company}
Kundenservice | Rechnungsabteilung
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

RECHNUNG
Rechnungsnummer: {invoice_num}
Rechnungsdatum: {datetime.now().strftime('%d.%m.%Y')}
Fälligkeitsdatum: {(datetime.now() + timedelta(days=14)).strftime('%d.%m.%Y')}

Rechnungsempfänger:
{self.family_members['father']['name']} {self.family_name}
{self.family_members['address']}

LEISTUNGSÜBERSICHT:
───────────────────────────────────────────────────────────────
{item_lines}───────────────────────────────────────────────────────────────
  Zwischensumme (netto):                    {subtotal:>10.2f} EUR
  MwSt. 19%:                                {subtotal * 0.19:>10.2f} EUR
  ═══════════════════════════════════════════════════════════════
  GESAMTBETRAG:                             {subtotal * 1.19:>10.2f} EUR

Zahlungsart: Überweisung
IBAN: DE89 3704 0044 0532 0130 00
BIC: COBADEFFXXX
Verwendungszweck: {invoice_num}

Vielen Dank für Ihr Vertrauen!
{company}
"""
        
        filename = f"rechnung_{invoice_num.lower().replace('-', '_')}.txt"
        return self._save_document("finanzen/rechnungen", filename, content)

    def generate_bank_statement(self) -> str:
        """Generate a bank account statement."""
        balance_start = random.randint(2000, 15000)
        month = random.randint(1, 12)
        year = 2024
        
        transactions = []
        balance = balance_start
        for _ in range(random.randint(10, 20)):
            is_credit = random.random() > 0.4
            if is_credit:
                amount = random.choice([
                    (random.randint(2500, 5000), "Gehalt " + random.choice(COMPANIES)),
                    (random.randint(50, 200), "Überweisung von " + random.choice(FIRST_NAMES_MALE) + " " + random.choice(LAST_NAMES)),
                    (random.randint(100, 500), "Rückerstattung"),
                ])
            else:
                amount = (
                    -random.randint(10, 500),
                    random.choice([
                        "REWE Supermarkt", "Amazon EU", "Netflix", "Spotify",
                        "Deutsche Bahn", "Miete Wohnung", "Stadtwerke Strom",
                        "Vodafone Mobilfunk", "ADAC Mitgliedschaft"
                    ])
                )
            
            balance += amount[0] if isinstance(amount, tuple) else amount
            transactions.append({
                "date": f"{random.randint(1, 28):02d}.{month:02d}.{year}",
                "desc": amount[1] if isinstance(amount, tuple) else "Transaktion",
                "amount": amount[0] if isinstance(amount, tuple) else amount,
                "balance": balance
            })
        
        trans_lines = ""
        for t in transactions:
            sign = "+" if t['amount'] > 0 else ""
            trans_lines += f"{t['date']} | {t['desc']:35} | {sign}{t['amount']:>10.2f} EUR | {t['balance']:>12.2f} EUR\n"
        
        content = f"""SPARKASSE {random.choice(CITIES).upper()}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

KONTOAUSZUG
Kontonummer: DE{random.randint(10, 99)} {random.randint(1000, 9999)} {random.randint(1000, 9999)} {random.randint(1000, 9999)} {random.randint(1000, 9999)} {random.randint(10, 99)}
Kontoinhaber: {self.family_members['father']['name']} {self.family_name}
Auszugszeitraum: {month:02d}/{year}

Anfangssaldo: {balance_start:,.2f} EUR

UMSÄTZE:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Datum      | Buchungstext                         | Betrag         | Saldo
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{trans_lines}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ENDSALDO: {balance:,.2f} EUR

Sparkasse {random.choice(CITIES)} | IBAN: DE89 3704 0044 0532 0130 00
"""
        
        filename = f"kontoauszug_{year}_{month:02d}.txt"
        return self._save_document("finanzen/kontoauszuege", filename, content)

    # ========== CHILDREN DOCUMENTS ==========

    def generate_school_report(self, child: Dict) -> str:
        """Generate a school report card (Zeugnis)."""
        school_year = f"{random.randint(2022, 2024)}/{random.randint(2023, 2025)}"
        semester = random.choice(["1. Halbjahr", "2. Halbjahr (Jahreszeugnis)"])
        
        grades = {}
        for subject in SCHOOL_SUBJECTS.keys():
            grades[subject] = random.choice(["1", "2", "2", "3", "3", "3", "4"])
        
        grade_lines = ""
        for subject, grade in grades.items():
            grade_lines += f"  {subject:25} {grade} ({GRADE_DESCRIPTIONS[grade]})\n"
        
        behavior = random.choice([
            "zeigt vorbildliches Sozialverhalten",
            "arbeitet konzentriert und gewissenhaft",
            "beteiligt sich aktiv am Unterricht",
            "zeigt gute Teamfähigkeit",
            "ist hilfsbereit und freundlich"
        ])
        
        content = f"""GRUNDSCHULE AM STADTPARK
{random.choice(CITIES)}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

ZEUGNIS
Schuljahr {school_year} - {semester}

Schüler/in: {child['name']} {self.family_name}
Klasse: {child['school_class']}a
Klassenlehrer/in: Frau {random.choice(LAST_NAMES)}

BEURTEILUNG:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
{grade_lines}
ARBEITS- UND SOZIALVERHALTEN:
{child['name']} {behavior}. Die Hausaufgaben werden 
{random.choice(['regelmäßig', 'meist vollständig', 'zuverlässig'])} erledigt.
{child['name']} {random.choice(['sollte mehr Selbstvertrauen entwickeln', 'zeigt großes Interesse am Lernen', 'ist ein/e beliebte/r Mitschüler/in'])}.

BEMERKUNGEN:
Versäumte Tage: {random.randint(0, 10)}
davon unentschuldigt: 0

{random.choice(['Versetzung in die nächste Klasse', 'Erfolgreich in die nächste Klassenstufe versetzt'])}.

{random.choice(CITIES)}, den {datetime.now().strftime('%d.%m.%Y')}

_____________________          _____________________
Klassenlehrer/in               Schulleitung
"""
        
        filename = f"zeugnis_{child['name'].lower()}_{school_year.replace('/', '_')}.txt"
        return self._save_document(f"kinder/{child['name'].lower()}/schule", filename, content)

    def generate_medical_report_child(self, child: Dict) -> str:
        """Generate a pediatric medical report."""
        diagnosis = random.choice(DIAGNOSES_CHILDREN)
        medication = random.sample(MEDICATIONS, random.randint(1, 2))
        
        content = f"""KINDERARZTPRAXIS DR. MED. {random.choice(LAST_NAMES).upper()}
Facharzt für Kinder- und Jugendmedizin
{random.choice(STREETS)} {random.randint(1, 100)}, {random.choice(CITIES)}
Tel: 0{random.randint(30, 89)}-{random.randint(100000, 999999)}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

ÄRZTLICHER BERICHT

Patient: {child['name']} {self.family_name}
Geburtsdatum: {child['birth_year']}.{random.randint(1,12):02d}.{random.randint(1,28):02d}
Alter: {child['age']} Jahre
Vorstellungsdatum: {datetime.now().strftime('%d.%m.%Y')}

ANAMNESE:
Das Kind wurde von der Mutter vorgestellt mit folgenden Beschwerden:
- {random.choice(['Fieber seit 2 Tagen', 'Husten und Schnupfen', 'Ohrenschmerzen', 'Bauchschmerzen', 'Hautausschlag'])}
- {random.choice(['Appetitlosigkeit', 'Schlafstörungen', 'allgemeines Unwohlsein'])}

BEFUND:
- Allgemeinzustand: {random.choice(['leicht reduziert', 'altersentsprechend', 'gut'])}
- Temperatur: {random.uniform(37.5, 39.5):.1f}°C
- Gewicht: {random.randint(20, 50)} kg
- Größe: {random.randint(110, 160)} cm
- {random.choice(['Rachen gerötet', 'Trommelfell beidseits reizlos', 'Lunge frei', 'Abdomen weich'])}

DIAGNOSE:
{diagnosis}

THERAPIE:
{', '.join(medication)}
{random.choice(['Bettruhe empfohlen', 'Viel Flüssigkeit', 'Bei Verschlechterung Wiedervorstellung'])}

ARBEITSUNFÄHIGKEITSBESCHEINIGUNG:
Schulbefreiung für {random.randint(3, 7)} Tage bis {(datetime.now() + timedelta(days=random.randint(3,7))).strftime('%d.%m.%Y')}

Mit freundlichen Grüßen

Dr. med. {random.choice(FIRST_NAMES_MALE)} {random.choice(LAST_NAMES)}
Facharzt für Kinder- und Jugendmedizin
"""
        
        filename = f"arztbericht_{child['name'].lower()}_{datetime.now().strftime('%Y%m%d')}.txt"
        return self._save_document(f"kinder/{child['name'].lower()}/gesundheit", filename, content)

    def generate_school_letter(self, child: Dict) -> str:
        """Generate a letter from school."""
        topics = [
            ("Elternabend", f"Einladung zum Elternabend am {(datetime.now() + timedelta(days=random.randint(7, 30))).strftime('%d.%m.%Y')} um 19:00 Uhr"),
            ("Klassenfahrt", f"Information zur geplanten Klassenfahrt nach {random.choice(['Sylt', 'Bayern', 'Schwarzwald', 'Ostsee'])}"),
            ("Schulausflug", f"Tagesausflug zum {random.choice(['Zoo', 'Museum', 'Schwimmbad', 'Theater'])}"),
            ("Projektwoche", f"Projektwoche zum Thema '{random.choice(['Umwelt', 'Digitalisierung', 'Sport', 'Kunst'])}'"),
        ]
        
        topic, description = random.choice(topics)
        
        content = f"""GRUNDSCHULE AM STADTPARK
{random.choice(STREETS)} {random.randint(1, 50)}
{random.choice(CITIES)}

An die Eltern der Klasse {child['school_class']}a
{self.family_members['address']}

{random.choice(CITIES)}, den {datetime.now().strftime('%d.%m.%Y')}

Betreff: {topic}

Liebe Eltern,

{description}.

{random.choice([
    'Bitte geben Sie die Einverständniserklärung bis zum',
    'Wir bitten um Rückmeldung bis zum',
    'Anmeldeschluss ist der'
])} {(datetime.now() + timedelta(days=14)).strftime('%d.%m.%Y')}.

Kosten: {random.randint(5, 50)},00 EUR
Bitte überweisen Sie den Betrag auf das Schulkonto oder geben Sie 
das Geld in einem verschlossenen Umschlag ab.

Bei Fragen stehe ich Ihnen gerne zur Verfügung.

Mit freundlichen Grüßen

{random.choice(FIRST_NAMES_FEMALE)} {random.choice(LAST_NAMES)}
Klassenlehrerin {child['school_class']}a

───────────────────────────────────────────────────────────
✂ Bitte hier abtrennen und ausgefüllt zurückgeben

[ ] Mein Kind {child['name']} nimmt teil
[ ] Mein Kind {child['name']} nimmt NICHT teil

Unterschrift Erziehungsberechtigte: _______________________
"""
        
        filename = f"schulbrief_{topic.lower().replace(' ', '_')}_{child['name'].lower()}.txt"
        return self._save_document(f"kinder/{child['name'].lower()}/schule", filename, content)

    # ========== FAMILY DOCUMENTS ==========

    def generate_insurance_document(self) -> str:
        """Generate an insurance policy document."""
        insurance_types = [
            ("Hausratversicherung", random.randint(80, 200)),
            ("Haftpflichtversicherung", random.randint(50, 100)),
            ("Unfallversicherung", random.randint(100, 250)),
            ("Rechtsschutzversicherung", random.randint(150, 300)),
            ("Krankenversicherung (Zusatz)", random.randint(30, 80)),
        ]
        
        ins_type, yearly_cost = random.choice(insurance_types)
        
        content = f"""ALLIANZ VERSICHERUNG
Kundenservice
Postfach 10 01 64
{random.choice(CITIES)}
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

VERSICHERUNGSPOLICE

Versicherungsnehmer:
{self.family_members['father']['name']} {self.family_name}
{self.family_members['address']}

Versicherungsart: {ins_type}
Vertragsnummer: VN-{random.randint(10000000, 99999999)}
Versicherungsbeginn: 01.01.{random.randint(2020, 2023)}
Laufzeit: jährlich, automatische Verlängerung

VERSICHERUNGSUMFANG:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Versicherungssumme: {random.randint(50, 500)}.000,00 EUR
Selbstbeteiligung: {random.choice([0, 150, 250, 500])},00 EUR
Zusatzleistungen: {random.choice(['Elementarschäden', 'Fahrraddiebstahl', 'Glasbruch'])}

JAHRESBEITRAG: {yearly_cost},00 EUR
Zahlweise: {random.choice(['monatlich', 'vierteljährlich', 'jährlich'])}

HINWEIS:
Bei Schadensfällen kontaktieren Sie bitte unsere 24-Stunden-Hotline:
0800-123 456 78

Diese Police ersetzt alle vorherigen Versicherungsscheine.

Mit freundlichen Grüßen
Allianz Versicherung AG
"""
        
        filename = f"versicherung_{ins_type.lower().replace(' ', '_').replace('(', '').replace(')', '')}.txt"
        return self._save_document("familie/versicherungen", filename, content)

    def generate_travel_booking(self) -> str:
        """Generate a travel booking confirmation."""
        destinations = [
            ("Mallorca, Spanien", "Flug", random.randint(200, 500)),
            ("Paris, Frankreich", "Zug", random.randint(100, 300)),
            ("München", "Zug", random.randint(50, 150)),
            ("Dubai, VAE", "Flug", random.randint(400, 900)),
            ("Rom, Italien", "Flug", random.randint(150, 400)),
        ]
        
        dest, transport, cost_per_person = random.choice(destinations)
        num_persons = len(self.family_members['children']) + 2
        total = cost_per_person * num_persons
        
        travelers = f"- {self.family_members['father']['name']} {self.family_name}\n"
        travelers += f"- {self.family_members['mother']['name']} {self.family_name}\n"
        for child in self.family_members['children']:
            travelers += f"- {child['name']} {self.family_name} (Kind)\n"
        
        dep_date = datetime.now() + timedelta(days=random.randint(30, 180))
        ret_date = dep_date + timedelta(days=random.randint(7, 14))
        
        content = f"""{'LUFTHANSA' if transport == 'Flug' else 'DEUTSCHE BAHN'} - BUCHUNGSBESTÄTIGUNG
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Buchungsnummer: {random.choice(['LH', 'DB', 'EW'])}{random.randint(100000, 999999)}
Buchungsdatum: {datetime.now().strftime('%d.%m.%Y')}

REISENDE:
{travelers}

REISEDETAILS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ziel: {dest}
Transportmittel: {transport}
Hinreise: {dep_date.strftime('%d.%m.%Y')} um {random.randint(6, 20):02d}:{random.choice(['00', '15', '30', '45'])}
Rückreise: {ret_date.strftime('%d.%m.%Y')} um {random.randint(6, 20):02d}:{random.choice(['00', '15', '30', '45'])}

KOSTENÜBERSICHT:
{num_persons} Personen x {cost_per_person},00 EUR = {total},00 EUR
Steuern und Gebühren:                    {int(total * 0.1)},00 EUR
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
GESAMTPREIS: {int(total * 1.1)},00 EUR

Zahlungsstatus: ✓ Bezahlt (Kreditkarte ****{random.randint(1000, 9999)})

WICHTIGE HINWEISE:
- Check-in ab {(dep_date - timedelta(days=1)).strftime('%d.%m.%Y')} möglich
- Gültiger Personalausweis erforderlich
- Gepäck: {random.choice(['1x Handgepäck', '1x Koffer bis 23kg', '2x Koffer'])}

Gute Reise wünscht
Ihr {'Lufthansa' if transport == 'Flug' else 'Deutsche Bahn'} Team
"""
        
        filename = f"reisebuchung_{dest.split(',')[0].lower().replace(' ', '_')}_{dep_date.strftime('%Y%m')}.txt"
        return self._save_document("familie/reisen", filename, content)

    def generate_employment_contract(self) -> str:
        """Generate an employment-related document."""
        employer = random.choice(COMPANIES)
        salary = random.randint(3500, 8000)
        
        content = f"""{employer.upper()}
Personalabteilung
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GEHALTSABRECHNUNG
Monat: {datetime.now().strftime('%B %Y')}

Mitarbeiter: {self.family_members['father']['name']} {self.family_name}
Personalnummer: {random.randint(10000, 99999)}
Steuerklasse: {random.choice([1, 3, 4])}
Kinderfreibeträge: {len(self.family_members['children'])}.0

BEZÜGE:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Grundgehalt:                             {salary:,.2f} EUR
{random.choice(['Überstundenvergütung', 'Leistungszulage', 'Bonuszahlung'])}: {random.randint(100, 500):,.2f} EUR
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Bruttolohn:                              {salary + random.randint(100, 500):,.2f} EUR

ABZÜGE:
Lohnsteuer:                              -{salary * 0.2:,.2f} EUR
Solidaritätszuschlag:                    -{salary * 0.011:,.2f} EUR
Rentenversicherung:                      -{salary * 0.093:,.2f} EUR
Krankenversicherung:                     -{salary * 0.073:,.2f} EUR
Pflegeversicherung:                      -{salary * 0.015:,.2f} EUR
Arbeitslosenversicherung:                -{salary * 0.012:,.2f} EUR
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
NETTOLOHN:                               {salary * 0.6:,.2f} EUR

Auszahlung auf: DE89 **** **** **** 0130 00

{employer}
Personalabteilung
"""
        
        filename = f"gehaltsabrechnung_{datetime.now().strftime('%Y_%m')}.txt"
        return self._save_document("arbeit/gehaltsabrechnungen", filename, content)

    # ========== MAIN GENERATION ==========

    def generate_full_corpus(
        self,
        num_tax: int = 5,
        num_crypto: int = 5,
        num_invoices: int = 10,
        num_bank: int = 5,
        num_school: int = None,
        num_medical: int = None,
        num_school_letters: int = None,
        num_insurance: int = 5,
        num_travel: int = 3,
        num_salary: int = 6
    ) -> Dict:
        """Generate full document corpus."""
        
        print(f"\n🏠 Generating documents for the {self.family_name} family")
        print(f"   Address: {self.family_members['address']}")
        print(f"   Parents: {self.family_members['father']['name']} & {self.family_members['mother']['name']}")
        print(f"   Children: {', '.join([c['name'] for c in self.family_members['children']])}")
        print("=" * 70)
        
        # Finance documents
        print("\n💰 Generating FINANCE documents...")
        for _ in range(num_tax):
            self.generate_tax_document()
        for _ in range(num_crypto):
            self.generate_crypto_statement()
        for _ in range(num_invoices):
            self.generate_invoice()
        for _ in range(num_bank):
            self.generate_bank_statement()
        
        # Children documents (per child)
        print("\n👶 Generating CHILDREN documents...")
        for child in self.family_members['children']:
            reports = num_school if num_school else random.randint(2, 4)
            medical = num_medical if num_medical else random.randint(2, 4)
            letters = num_school_letters if num_school_letters else random.randint(3, 6)
            
            for _ in range(reports):
                self.generate_school_report(child)
            for _ in range(medical):
                self.generate_medical_report_child(child)
            for _ in range(letters):
                self.generate_school_letter(child)
        
        # Family documents
        print("\n👨‍👩‍👧‍👦 Generating FAMILY documents...")
        for _ in range(num_insurance):
            self.generate_insurance_document()
        for _ in range(num_travel):
            self.generate_travel_booking()
        
        # Work documents
        print("\n💼 Generating WORK documents...")
        for _ in range(num_salary):
            self.generate_employment_contract()
        
        print("\n" + "=" * 70)
        print(f"✅ Generated {len(self.generated_files)} documents!")
        
        # Summary by category
        categories = {}
        for f in self.generated_files:
            cat = f['category'].split('/')[0]
            categories[cat] = categories.get(cat, 0) + 1
        
        print("\nSummary:")
        for cat, count in sorted(categories.items()):
            print(f"  📁 {cat}: {count} files")
        
        return {
            "family": self.family_members,
            "files": self.generated_files,
            "total": len(self.generated_files)
        }


def main():
    parser = argparse.ArgumentParser(
        description="Generate realistic German family document corpus"
    )
    
    parser.add_argument(
        "--output", "-o",
        type=str,
        default="/mnt/data",
        help="Output directory (default: /mnt/data)"
    )
    
    parser.add_argument(
        "--all",
        action="store_true",
        help="Generate a complete corpus with default counts"
    )
    
    parser.add_argument("--tax", type=int, default=5, help="Number of tax documents")
    parser.add_argument("--crypto", type=int, default=5, help="Number of crypto statements")
    parser.add_argument("--invoices", type=int, default=10, help="Number of invoices")
    parser.add_argument("--bank", type=int, default=5, help="Number of bank statements")
    parser.add_argument("--insurance", type=int, default=5, help="Number of insurance docs")
    parser.add_argument("--travel", type=int, default=3, help="Number of travel bookings")
    parser.add_argument("--salary", type=int, default=6, help="Number of salary slips")
    
    args = parser.parse_args()
    
    generator = FamilyDocumentGenerator(args.output)
    
    result = generator.generate_full_corpus(
        num_tax=args.tax,
        num_crypto=args.crypto,
        num_invoices=args.invoices,
        num_bank=args.bank,
        num_insurance=args.insurance,
        num_travel=args.travel,
        num_salary=args.salary
    )
    
    # Save manifest
    manifest_path = Path(args.output) / "corpus_manifest.json"
    with open(manifest_path, 'w', encoding='utf-8') as f:
        json.dump(result, f, ensure_ascii=False, indent=2, default=str)
    
    print(f"\n📋 Manifest saved to: {manifest_path}")
    print(f"📁 Documents in: {args.output}")


if __name__ == "__main__":
    main()
