#!/usr/bin/env python3
"""
Generate realistic German document corpus for semantic search testing.

Usage:
    python generate_corpus.py --count 50 --output ./output
    python generate_corpus.py --invoices 20 --logs 15 --emails 15 --output ./output
"""

import argparse
import sys
from pathlib import Path

# Add parent to path for imports
sys.path.insert(0, str(Path(__file__).parent))

from corpus_generator import CorpusGenerator


def main():
    parser = argparse.ArgumentParser(
        description="Generate realistic German document corpus"
    )
    
    parser.add_argument(
        "--output", "-o",
        type=str,
        default="./output",
        help="Output directory (default: ./output)"
    )
    
    parser.add_argument(
        "--count", "-c",
        type=int,
        default=None,
        help="Total documents (distributed 40%% invoices, 30%% logs, 30%% emails)"
    )
    
    parser.add_argument(
        "--invoices",
        type=int,
        default=None,
        help="Number of invoices to generate"
    )
    
    parser.add_argument(
        "--logs",
        type=int,
        default=None,
        help="Number of tech logs to generate"
    )
    
    parser.add_argument(
        "--emails",
        type=int,
        default=None,
        help="Number of emails to generate"
    )
    
    parser.add_argument(
        "--noise",
        type=float,
        default=0.3,
        help="Probability of noise injection (0.0-1.0, default: 0.3)"
    )
    
    args = parser.parse_args()
    
    # Determine document counts
    if args.count is not None:
        num_invoices = int(args.count * 0.4)
        num_logs = int(args.count * 0.3)
        num_emails = args.count - num_invoices - num_logs
    elif args.invoices is not None or args.logs is not None or args.emails is not None:
        num_invoices = args.invoices or 0
        num_logs = args.logs or 0
        num_emails = args.emails or 0
    else:
        # Default: 50 documents
        num_invoices = 20
        num_logs = 15
        num_emails = 15

    # Create generator and run
    generator = CorpusGenerator(
        output_dir=args.output,
        noise_probability=args.noise,
    )
    
    generator.generate_corpus(
        num_invoices=num_invoices,
        num_logs=num_logs,
        num_emails=num_emails,
    )
    
    generator.dump_ground_truth()
    
    print(f"\n📁 Output location: {args.output}")


if __name__ == "__main__":
    main()
