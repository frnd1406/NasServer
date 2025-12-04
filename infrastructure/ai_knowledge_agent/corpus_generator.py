"""
Corpus Generator - Central orchestration class.
Manages document generation, rendering, and ground truth registry.
"""

import json
import os
import random
from datetime import datetime
from pathlib import Path
from typing import List, Optional, Tuple, Any

from data_models import (
    InvoiceData, LogData, EmailData,
    GroundTruth, GroundTruthEntry
)
from generators.invoice import create_invoice_data
from generators.tech_log import create_log_data
from generators.email import create_email_data
from renderers.text_renderer import TextRenderer
from renderers.html_renderer import HtmlRenderer
from renderers.pdf_renderer import PdfRenderer
from noise import inject_noise


class CorpusGenerator:
    """
    Orchestrates corpus generation with metadata tracking.
    
    Manages the full pipeline:
    1. Generate structured data
    2. Render to multiple formats
    3. Optionally inject noise
    4. Record ground truth metadata
    """

    def __init__(self, output_dir: str, noise_probability: float = 0.3):
        """
        Initialize the corpus generator.
        
        Args:
            output_dir: Base directory for output
            noise_probability: Probability of applying noise to a document
        """
        self.output_dir = Path(output_dir)
        self.docs_dir = self.output_dir / "docs"
        self.noise_probability = noise_probability
        
        # Ground truth registry
        self.ground_truth = GroundTruth(
            generated_at=datetime.now(),
            total_documents=0,
            documents=[],
        )
        
        # Renderers
        self.text_renderer = TextRenderer()
        self.html_renderer = HtmlRenderer()
        self.pdf_renderer = PdfRenderer()

    def _ensure_dirs(self):
        """Create output directories."""
        self.docs_dir.mkdir(parents=True, exist_ok=True)

    def _select_renderer(self, doc_type: str) -> Tuple[Any, str]:
        """
        Randomly select a renderer based on document type.
        
        Returns:
            Tuple of (renderer, format_name)
        """
        if doc_type == "invoice":
            # Invoices: all three formats
            choice = random.choices(
                [self.text_renderer, self.html_renderer, self.pdf_renderer],
                weights=[0.4, 0.3, 0.3]
            )[0]
        elif doc_type == "log":
            # Logs: mainly text, some HTML
            choice = random.choices(
                [self.text_renderer, self.html_renderer],
                weights=[0.7, 0.3]
            )[0]
        else:  # email
            # Emails: text and HTML
            choice = random.choices(
                [self.text_renderer, self.html_renderer],
                weights=[0.5, 0.5]
            )[0]
        
        return choice, choice.format_name

    def _generate_filename(self, doc_type: str, data: Any, extension: str) -> str:
        """Generate unique filename for document."""
        if doc_type == "invoice":
            base = f"rechnung_{data.invoice_number.replace('-', '_').lower()}"
        elif doc_type == "log":
            base = f"log_{data.server_name}_{data.log_date.strftime('%Y%m%d')}"
        else:  # email
            base = f"email_{data.sent_date.strftime('%Y%m%d_%H%M%S')}"
        
        return f"{base}{extension}"

    def generate_document(self, doc_type: str) -> Tuple[str, bool]:
        """
        Generate a single document.
        
        Args:
            doc_type: Type of document ('invoice', 'log', 'email')
            
        Returns:
            Tuple of (filename, has_noise)
        """
        # Generate data
        if doc_type == "invoice":
            data = create_invoice_data()
        elif doc_type == "log":
            data = create_log_data()
        else:
            data = create_email_data()

        # Select renderer
        renderer, format_name = self._select_renderer(doc_type)
        
        # Render
        content_bytes, extension = renderer.render(data)
        
        # Apply noise (only to text content, not PDFs which are binary)
        has_noise = False
        if format_name != "pdf" and random.random() < self.noise_probability:
            content_str = content_bytes.decode('utf-8')
            content_str = inject_noise(content_str, intensity=0.1)
            content_bytes = content_str.encode('utf-8')
            has_noise = True
        
        # Generate filename and save
        filename = self._generate_filename(doc_type, data, extension)
        filepath = self.docs_dir / filename
        
        # Ensure unique filename
        counter = 1
        while filepath.exists():
            base, ext = os.path.splitext(filename)
            filename = f"{base}_{counter}{ext}"
            filepath = self.docs_dir / filename
            counter += 1

        with open(filepath, 'wb') as f:
            f.write(content_bytes)

        # Record to ground truth
        entry = GroundTruthEntry(
            filename=filename,
            doc_type=doc_type,
            format=format_name,
            has_noise=has_noise,
            metadata=data.to_dict(),
        )
        self.ground_truth.documents.append(entry)
        self.ground_truth.total_documents += 1

        return filename, has_noise

    def generate_corpus(
        self,
        num_invoices: int = 20,
        num_logs: int = 15,
        num_emails: int = 15,
    ) -> None:
        """
        Generate complete document corpus.
        
        Args:
            num_invoices: Number of invoices to generate
            num_logs: Number of tech logs to generate
            num_emails: Number of emails to generate
        """
        self._ensure_dirs()
        
        print(f"Generating corpus in {self.output_dir}...")
        print("=" * 60)

        # Generate invoices
        print(f"\n📄 Generating {num_invoices} invoices...")
        for i in range(num_invoices):
            filename, has_noise = self.generate_document("invoice")
            if (i + 1) % 5 == 0:
                print(f"  ✓ {i + 1}/{num_invoices} invoices")

        # Generate logs
        print(f"\n📊 Generating {num_logs} tech logs...")
        for i in range(num_logs):
            filename, has_noise = self.generate_document("log")
            if (i + 1) % 5 == 0:
                print(f"  ✓ {i + 1}/{num_logs} logs")

        # Generate emails
        print(f"\n📧 Generating {num_emails} emails...")
        for i in range(num_emails):
            filename, has_noise = self.generate_document("email")
            if (i + 1) % 5 == 0:
                print(f"  ✓ {i + 1}/{num_emails} emails")

        print("\n" + "=" * 60)
        print(f"✅ Generated {self.ground_truth.total_documents} documents!")
        
        # Count formats
        formats = {}
        noisy_count = 0
        for doc in self.ground_truth.documents:
            formats[doc.format] = formats.get(doc.format, 0) + 1
            if doc.has_noise:
                noisy_count += 1
        
        print(f"\nFormats: {formats}")
        print(f"Documents with noise: {noisy_count}")

    def dump_ground_truth(self) -> str:
        """
        Export ground truth registry to JSON.
        
        Returns:
            Path to the ground truth file
        """
        filepath = self.output_dir / "ground_truth.json"
        
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(self.ground_truth.to_dict(), f, ensure_ascii=False, indent=2)
        
        print(f"\n📋 Ground truth saved to: {filepath}")
        return str(filepath)
