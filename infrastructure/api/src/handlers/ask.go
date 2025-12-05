package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/database"
	"github.com/sirupsen/logrus"
)

type ragResult struct {
	Answer  string   `json:"answer"`
	Sources []source `json:"sources"`
}

type source struct {
	FilePath string  `json:"file_path"`
	Score    float64 `json:"score"`
	Snippet  string  `json:"snippet"`
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// AskHandler implements RAG (Retrieval Augmented Generation):
// 1. Generate embedding for the question
// 2. Find top-k relevant documents
// 3. Send question + context to LLM
// 4. Return intelligent answer with sources
func AskHandler(db *database.DB, aiServiceURL, ollamaURL, llmModel string, httpClient *http.Client, logger *logrus.Logger) gin.HandlerFunc {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 60 * time.Second} // LLM needs more time
	}
	embedURL := strings.TrimRight(aiServiceURL, "/") + "/embed_query"
	generateURL := strings.TrimRight(ollamaURL, "/") + "/api/generate"

	return func(c *gin.Context) {
		question := strings.TrimSpace(c.Query("q"))
		if question == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing query parameter 'q'"})
			return
		}

		logger.Infof("RAG query: %s", question)

		// Step 1: Get embedding for the question
		embedding, err := fetchEmbedding(c.Request.Context(), client, embedURL, question)
		if err != nil {
			logger.WithError(err).Error("Failed to fetch embedding for RAG")
			c.JSON(http.StatusBadGateway, gin.H{"error": "embedding service unavailable"})
			return
		}

		// Convert to pgvector format
		parts := make([]string, len(embedding))
		for i, v := range embedding {
			parts[i] = strconv.FormatFloat(v, 'f', -1, 64)
		}
		embeddingStr := "[" + strings.Join(parts, ",") + "]"

		// Step 2: Find top 3 relevant documents
		rows, err := db.QueryContext(c.Request.Context(), `
			SELECT file_path, content, (1 - (embedding <=> $1::vector)) as similarity
			FROM file_embeddings
			WHERE file_path LIKE '/mnt/data/%'
			  AND file_path NOT LIKE '%/.trash/%'
			ORDER BY embedding <=> $1::vector
			LIMIT 3;
		`, embeddingStr)
		if err != nil {
			logger.WithError(err).Error("Failed to query documents for RAG")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database query failed"})
			return
		}
		defer rows.Close()

		var sources []source
		var contextParts []string
		for rows.Next() {
			var filePath, content string
			var similarity float64
			if err := rows.Scan(&filePath, &content, &similarity); err != nil {
				logger.WithError(err).Error("Failed to scan RAG result")
				continue
			}

			// Truncate content for snippet
			snippet := content
			if len(snippet) > 200 {
				snippet = snippet[:200] + "..."
			}

			sources = append(sources, source{
				FilePath: filePath,
				Score:    similarity,
				Snippet:  snippet,
			})

			// Add to context for LLM
			contextParts = append(contextParts, fmt.Sprintf("=== Datei: %s ===\n%s", filePath, content))
		}

		if len(sources) == 0 {
			c.JSON(http.StatusOK, ragResult{
				Answer:  "Keine relevanten Dokumente gefunden.",
				Sources: []source{},
			})
			return
		}

		// Step 3: Generate answer with LLM
		contextText := strings.Join(contextParts, "\n\n")
		prompt := fmt.Sprintf(`Du bist ein hilfreicher Assistent. Beantworte die Frage basierend auf den folgenden Dokumenten.
Antworte auf Deutsch, präzise und direkt. Wenn die Antwort nicht in den Dokumenten zu finden ist, sage das ehrlich.

DOKUMENTE:
%s

FRAGE: %s

ANTWORT:`, contextText, question)

		answer, err := generateWithOllama(c.Request.Context(), client, generateURL, llmModel, prompt)
		if err != nil {
			logger.WithError(err).Error("Failed to generate LLM response")
			c.JSON(http.StatusBadGateway, gin.H{"error": "LLM service unavailable"})
			return
		}

		c.JSON(http.StatusOK, ragResult{
			Answer:  strings.TrimSpace(answer),
			Sources: sources,
		})
	}
}

func generateWithOllama(ctx context.Context, client *http.Client, url, model, prompt string) (string, error) {
	reqBody := ollamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return "", fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return "", fmt.Errorf("ollama status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Response, nil
}
