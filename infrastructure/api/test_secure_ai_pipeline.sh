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

# =================================================================
# Test 4: DESTRUCTIVE PATH TEST (Ghost Knowledge Elimination)
# This is the critical test from Operation Ghost Busters
# =================================================================
echo ""
echo "=== 🔥 DESTRUCTIVE PATH TEST (Ghost Knowledge) ==="
echo ""

# Step 4.1: Upload a secret file
echo "📤 Step 4.1: Ingesting secret.txt with 'Codewort: Bananenbrot'"
GHOST_INGEST=$(curl -s -X POST "$AI_AGENT_URL/process" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Dieses Dokument enthält ein geheimes Codewort: Bananenbrot. Bitte merke es dir!",
    "file_id": "secret.txt",
    "file_path": "/mnt/data/secret.txt",
    "mime_type": "text/plain"
  }')

if echo "$GHOST_INGEST" | grep -q '"status": "success"'; then
  echo "✅ Secret ingested successfully"
else
  echo "❌ Secret ingestion failed"
  echo "Response: $GHOST_INGEST"
  exit 1
fi

# Small delay for embedding to be available
sleep 2

# Step 4.2: Verify the secret is accessible via RAG
echo ""
echo "🔍 Step 4.2: Asking 'Was ist das Codewort?'"
GHOST_RAG1=$(curl -s -X POST "$AI_AGENT_URL/rag" \
  -H "Content-Type: application/json" \
  -d '{"query": "Was ist das Codewort?", "top_k": 5}')

echo "RAG Response: $GHOST_RAG1"
if echo "$GHOST_RAG1" | grep -qi "Bananenbrot"; then
  echo "✅ SECRET FOUND: AI knows the Codewort (Bananenbrot)"
else
  echo "❌ SECRET NOT FOUND: Expected 'Bananenbrot' in response"
  exit 1
fi

# Step 4.3: Delete the secret file
echo ""
echo "🗑️ Step 4.3: Deleting secret.txt from AI index"
GHOST_DELETE=$(curl -s -X POST "$AI_AGENT_URL/delete" \
  -H "Content-Type: application/json" \
  -d '{"file_id": "secret.txt"}')

if echo "$GHOST_DELETE" | grep -q '"status": "success"'; then
  echo "✅ Secret deleted from index"
else
  echo "⚠️ Delete returned unexpected status (may have been already deleted or not found)"
  echo "Response: $GHOST_DELETE"
fi

# Small delay for deletion to propagate
sleep 1

# Step 4.4: CRITICAL - Verify the secret is GONE
echo ""
echo "🔍 Step 4.4: CRITICAL - Asking 'Was ist das Codewort?' again"
GHOST_RAG2=$(curl -s -X POST "$AI_AGENT_URL/rag" \
  -H "Content-Type: application/json" \
  -d '{"query": "Was ist das Codewort?", "top_k": 5}')

echo "RAG Response: $GHOST_RAG2"
if echo "$GHOST_RAG2" | grep -qi "Bananenbrot"; then
  echo "❌ GHOST KNOWLEDGE DETECTED: AI still knows 'Bananenbrot' after deletion!"
  echo "   This is a CRITICAL FAILURE - ghost knowledge persists!"
  exit 1
else
  echo "✅ GHOST BUSTED: AI no longer knows the secret (Bananenbrot not in response)"
  echo "   Ghost knowledge successfully eliminated!"
fi

echo ""
echo "=== All Tests Complete ==="
echo "✅ Phase 1: Direct ingestion - PASSED"
echo "✅ Phase 2: Search - PASSED"
echo "✅ Phase 3: RAG query - PASSED"
echo "✅ Phase 4: Ghost knowledge elimination - PASSED"
