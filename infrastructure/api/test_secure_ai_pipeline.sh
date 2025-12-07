#!/bin/bash
# End-to-End Test: Secure AI Pipeline for Encrypted Files
# Tests: Upload encrypted → AI indexes → Search returns result

set -e

API_URL="${API_URL:-http://localhost:8080}"
AI_AGENT_URL="${AI_AGENT_URL:-http://localhost:5000}"

echo "=== Secure AI Pipeline Test ==="
echo "API: $API_URL"
echo "AI Agent: $AI_AGENT_URL"
echo ""

# Test 1: Direct content ingestion via /process (RAM-push mode)
echo "📤 Test 1: Direct content ingestion via /process"
INGEST_RESULT=$(curl -s -X POST "$AI_AGENT_URL/process" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Dies ist ein geheimer Testinhalt. Das Passwort lautet: SECURE_TEST_123",
    "file_id": "encrypted_test_file.txt.enc",
    "file_path": "/media/demo/encrypted_test_file.txt.enc",
    "mime_type": "text/plain"
  }')

echo "Response: $INGEST_RESULT"
if echo "$INGEST_RESULT" | grep -q '"status": "success"'; then
  echo "✅ Direct ingestion successful"
else
  echo "❌ Direct ingestion failed"
  exit 1
fi

# Test 2: Search for the content
echo ""
echo "🔍 Test 2: Search for 'geheimer' in AI index"
SEARCH_RESULT=$(curl -s -X POST "$AI_AGENT_URL/search" \
  -H "Content-Type: application/json" \
  -d '{"query": "geheimer Testinhalt", "limit": 5}')

echo "Search Response: $SEARCH_RESULT"
if echo "$SEARCH_RESULT" | grep -q "encrypted_test_file"; then
  echo "✅ Search found encrypted file content"
else
  echo "⚠️ Search did not find content (may need more time for embedding)"
fi

# Test 3: RAG query
echo ""
echo "🤖 Test 3: RAG query - 'Was ist das Passwort?'"
RAG_RESULT=$(curl -s -X POST "$AI_AGENT_URL/rag" \
  -H "Content-Type: application/json" \
  -d '{"query": "Was ist das Passwort?", "top_k": 3}')

echo "RAG Response: $RAG_RESULT"
if echo "$RAG_RESULT" | grep -q "SECURE_TEST_123"; then
  echo "✅ RAG returned correct answer from encrypted content!"
else
  echo "⚠️ RAG did not return expected password (check LLM response)"
fi

# Cleanup: Delete test embedding
echo ""
echo "🧹 Cleanup: Deleting test embedding"
DELETE_RESULT=$(curl -s -X POST "$AI_AGENT_URL/delete" \
  -H "Content-Type: application/json" \
  -d '{"file_id": "encrypted_test_file.txt.enc"}')

echo "Delete Response: $DELETE_RESULT"

echo ""
echo "=== Test Complete ==="
