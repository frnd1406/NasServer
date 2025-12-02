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

type embedQueryResponse struct {
	Embedding []float64 `json:"embedding"`
}

type searchResult struct {
	FilePath   string  `json:"file_path"`
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity"`
}

// SearchHandler handles semantic search queries using the AI knowledge agent.
func SearchHandler(db *database.DB, aiServiceURL string, httpClient *http.Client, logger *logrus.Logger) gin.HandlerFunc {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	baseURL := strings.TrimRight(aiServiceURL, "/")

	return func(c *gin.Context) {
		query := strings.TrimSpace(c.Query("q"))
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing query parameter 'q'"})
			return
		}

		embedding, err := fetchEmbedding(c.Request.Context(), client, baseURL+"/embed_query", query)
		if err != nil {
			logger.WithError(err).Error("Failed to fetch embedding from AI agent")
			c.JSON(http.StatusBadGateway, gin.H{"error": "embedding service unavailable"})
			return
		}

		// Convert embedding float64 slice to pgvector format string
		// pgvector requires format: '[0.1,0.2,0.3]'
		parts := make([]string, len(embedding))
		for i, v := range embedding {
			parts[i] = strconv.FormatFloat(v, 'f', -1, 64)
		}
		embeddingStr := "[" + strings.Join(parts, ",") + "]"

		rows, err := db.QueryContext(c.Request.Context(), `
			SELECT file_path, content, 1 - (embedding <=> $1::vector) as similarity
			FROM file_embeddings
			ORDER BY embedding <=> $1::vector
			LIMIT 10;
		`, embeddingStr)
		if err != nil {
			logger.WithError(err).Error("Failed to run similarity search query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database query failed"})
			return
		}
		defer rows.Close()

		results := make([]searchResult, 0, 10)
		for rows.Next() {
			var r searchResult
			if scanErr := rows.Scan(&r.FilePath, &r.Content, &r.Similarity); scanErr != nil {
				logger.WithError(scanErr).Error("Failed to scan search result row")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read results"})
				return
			}
			results = append(results, r)
		}
		if err := rows.Err(); err != nil {
			logger.WithError(err).Error("Row iteration error during search")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read results"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"query":   query,
			"results": results,
		})
	}
}

func fetchEmbedding(ctx context.Context, client *http.Client, url, query string) ([]float64, error) {
	payload := map[string]string{"text": query}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ai agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("ai agent status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed embedQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(parsed.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding")
	}

	return parsed.Embedding, nil
}
