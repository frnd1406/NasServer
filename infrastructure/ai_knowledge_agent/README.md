# AI Knowledge Agent - Corpus Generator

Generate realistic German document corpora for semantic search testing and OCR validation.

## Features

- **Multi-format output**: PDF, HTML, TXT
- **Document types**: Invoices, Tech Logs, Business Emails
- **Ground Truth Registry**: JSON metadata for search validation
- **Noise Injection**: OCR-style corruption for robustness testing

## Quick Start

```bash
# Setup
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

# Generate 50 documents
python generate_corpus.py --count 50 --output ./output
```

## Usage

```bash
# Default: 50 docs (40% invoices, 30% logs, 30% emails)
python generate_corpus.py --count 50 --output ./output

# Custom distribution
python generate_corpus.py --invoices 30 --logs 10 --emails 10 --output ./output

# Adjust noise probability (0.0-1.0)
python generate_corpus.py --count 50 --noise 0.5 --output ./output
```

## Output Structure

```
output/
├── docs/
│   ├── rechnung_re_2023_1234.pdf
│   ├── rechnung_re_2023_5678.html
│   ├── log_web-3_20231015.txt
│   └── email_20231018_143021.html
└── ground_truth.json
```

### Ground Truth Schema

```json
{
  "generated_at": "2025-12-04T15:42:07Z",
  "total_documents": 50,
  "documents": [
    {
      "filename": "rechnung_re_2023_1234.pdf",
      "type": "invoice",
      "format": "pdf",
      "has_noise": false,
      "metadata": {
        "invoice_number": "RE-2023-1234",
        "gross_total": 142.80,
        "customer": "Max Mustermann",
        "company": "TechCorp GmbH"
      }
    }
  ]
}
```

## Architecture

```
├── generate_corpus.py      # CLI entry point
├── corpus_generator.py     # Orchestrator class
├── data_models.py          # Typed dataclasses
├── noise.py                # OCR-style noise injection
├── generators/
│   ├── invoice.py          # → InvoiceData
│   ├── tech_log.py         # → LogData
│   └── email.py            # → EmailData
├── renderers/
│   ├── text_renderer.py    # Jinja2 → .txt
│   ├── html_renderer.py    # Jinja2 → .html
│   └── pdf_renderer.py     # ReportLab → .pdf
└── templates/
    ├── invoice_clean.j2
    ├── invoice_table.html.j2
    ├── log.j2 / log.html.j2
    └── email.j2 / email.html.j2
```

## Noise Injection

OCR-style corruption for testing search robustness:

| Type | Example |
|------|---------|
| Character swap | `i` → `1`, `l` → `I`, `O` → `0` |
| Keyword typos | `Rechnung` → `Rechugn` |
| Letter swaps | `Server` → `Serevr` |

Noise is applied **only to rendered content**, never to ground truth metadata.

## Dependencies

- `faker` - German locale data generation
- `jinja2` - Template rendering
- `reportlab` - PDF generation
