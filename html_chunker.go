package kustomersdk

import (
	"fmt"
	"log"
	"strings"

	"github.com/rivo/uniseg"
)

// ChunkHTML safely chunks HTML while respecting emoji and HTML tag boundaries
func ChunkHTML(html string, maxSize int) []string {
	if len(html) == 0 {
		return []string{}
	}

	var chunks []string
	var currentChunk strings.Builder

	i := 0
	htmlBytes := []byte(html)

	for i < len(htmlBytes) {
		// Check if we're at the start of an HTML tag
		if htmlBytes[i] == '<' {
			tagEnd := findTagEnd(htmlBytes, i)
			if tagEnd == -1 {
				// Malformed tag, treat as regular text
				addTextSafely(&currentChunk, htmlBytes, &i, &chunks, maxSize)
				continue
			}

			tagLength := tagEnd - i + 1

			// If tag is too large for any chunk, log warning
			if tagLength > maxSize {
				log.Printf("Warning: HTML tag longer than %d characters", maxSize)
			}

			// If adding tag would exceed limit, start new chunk
			if currentChunk.Len()+tagLength > maxSize && currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
			}

			// Add the complete tag
			currentChunk.Write(htmlBytes[i : tagEnd+1])
			i = tagEnd + 1
		} else {
			// Handle regular text with Unicode safety
			addTextSafely(&currentChunk, htmlBytes, &i, &chunks, maxSize)
		}
	}

	// Add final chunk if it has content
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// addTextSafely adds text while respecting grapheme cluster boundaries
func addTextSafely(chunk *strings.Builder, data []byte, pos *int, chunks *[]string, maxSize int) {
	if *pos >= len(data) {
		return
	}

	// Get the next grapheme cluster
	cluster, _, _, _ := uniseg.FirstGraphemeCluster(data[*pos:], -1)
	if len(cluster) == 0 {
		return
	}

	// If adding this cluster would exceed maxSize, start new chunk
	if chunk.Len()+len(cluster) > maxSize {
		if chunk.Len() > 0 {
			*chunks = append(*chunks, chunk.String())
			chunk.Reset()
		}

		// If single cluster is larger than maxSize, log warning
		if len(cluster) > maxSize {
			log.Printf("Warning: Grapheme cluster longer than %d characters", maxSize)
		}
	}

	// Add the cluster to current chunk
	chunk.Write(cluster)
	*pos += len(cluster)
}

// findTagEnd finds the end position of an HTML tag
func findTagEnd(data []byte, start int) int {
	inQuotes := false
	quoteChar := byte(0)

	for i := start + 1; i < len(data); i++ {
		char := data[i]

		if !inQuotes {
			if char == '"' || char == '\'' {
				inQuotes = true
				quoteChar = char
			} else if char == '>' {
				return i
			}
		} else {
			if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			}
		}
	}

	return -1 // Malformed tag
}

// chunkHTMLBody is used to split an HTML body into chunks of 1024 characters that
// are then passed to Kustomer as custom variables in the conversation. It ensures
// that HTML tags and emojis are not split across chunk boundaries.
// Returns a map with keys like "htmlContent0Str", "htmlContent1Str", etc.
func chunkHTMLBody(input string) map[string]string {
	result := make(map[string]string)
	if len(input) == 0 {
		return result
	}

	chunks := ChunkHTML(input, 1024)
	for i, chunk := range chunks {
		key := fmt.Sprintf("htmlContent%dStr", i)
		result[key] = chunk
	}

	return result
}
