#!/usr/bin/env python3
"""
Batch indexing script for test corpus documents.
Sends all files in test_corpus directory to the AI agent for processing.
"""

import os
import json
import requests
import time
from pathlib import Path

CORPUS_DIR = "/mnt/data/test_corpus"
AI_AGENT_URL = "http://localhost:5000/process"


def index_file(filepath):
    """Index a single file by sending it to the AI agent."""
    filename = os.path.basename(filepath)

    payload = {
        "file_path": filepath,
        "file_id": filename,
        "mime_type": "text/plain"
    }

    try:
        response = requests.post(
            AI_AGENT_URL,
            json=payload,
            timeout=30
        )

        if response.status_code == 200:
            data = response.json()
            return True, data
        else:
            return False, f"HTTP {response.status_code}: {response.text}"

    except Exception as e:
        return False, str(e)


def main():
    """Batch index all files in the corpus directory."""
    if not os.path.exists(CORPUS_DIR):
        print(f"❌ Corpus directory not found: {CORPUS_DIR}")
        return

    files = sorted([f for f in os.listdir(CORPUS_DIR) if f.endswith('.txt')])

    if not files:
        print(f"❌ No .txt files found in {CORPUS_DIR}")
        return

    print(f"📚 Batch Indexing {len(files)} documents...")
    print("=" * 60)

    success_count = 0
    fail_count = 0
    total_chars = 0

    start_time = time.time()

    for i, filename in enumerate(files, 1):
        filepath = os.path.join(CORPUS_DIR, filename)

        # Show progress
        print(f"\n[{i}/{len(files)}] Processing: {filename}")

        success, result = index_file(filepath)

        if success:
            content_len = result.get('content_length', 0)
            embedding_dim = result.get('embedding_dim', 0)
            print(f"  ✓ Success - {content_len} chars, {embedding_dim}D embedding")
            success_count += 1
            total_chars += content_len
        else:
            print(f"  ✗ Failed - {result}")
            fail_count += 1

        # Small delay to avoid overwhelming the AI agent
        if i < len(files):
            time.sleep(0.2)

    elapsed = time.time() - start_time

    print("\n" + "=" * 60)
    print(f"✅ Batch indexing complete!")
    print(f"\nResults:")
    print(f"  • Total files: {len(files)}")
    print(f"  • Successful: {success_count}")
    print(f"  • Failed: {fail_count}")
    print(f"  • Total content: {total_chars:,} characters")
    print(f"  • Time elapsed: {elapsed:.1f}s")
    print(f"  • Average: {elapsed/len(files):.2f}s per file")


if __name__ == "__main__":
    main()
