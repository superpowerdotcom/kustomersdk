package kustomersdk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkHTMLBody(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: map[string]string{},
		},
		{
			name:  "Short string without HTML",
			input: "Hello World",
			expected: map[string]string{
				"htmlContent0Str": "Hello World",
			},
		},
		{
			name: "String with HTML tags across boundary",
			input: "<div><h1>Welcome to Our Website</h1><p>This is a long paragraph with lots of content that should be split across multiple chunks. " +
				strings.Repeat("We have lots of important information to share with you about our products and services. ", 10) +
				"</p><strong>Important Notice:</strong> Please read all terms and conditions carefully before proceeding with your purchase.</div>",
			expected: map[string]string{
				"htmlContent0Str": "<div><h1>Welcome to Our Website</h1><p>This is a long paragraph with lots of content that should be split across multiple chunks. " + 
					strings.Repeat("We have lots of important information to share with you about our products and services. ", 10) + 
					"</p>",
				"htmlContent1Str": "<strong>Important Notice:</strong> Please read all terms and conditions carefully before proceeding with your purchase.</div>",
			},
		},
		{
			name: "String with emoji across boundary",
			input: "<div><p>" + strings.Repeat("This is sample content with some text. ", 25) + "</p></div>👋<span>Hello there!</span>",
			expected: map[string]string{
				"htmlContent0Str": "<div><p>" + strings.Repeat("This is sample content with some text. ", 25) + "</p></div>👋<span>Hello there!</span>",
			},
		},
		{
			name:  "String with multiple HTML tags",
			input: "<div>Hello</div><p>This is a test with <strong>bold text</strong> and some <em>emphasis</em></p>",
			expected: map[string]string{
				"htmlContent0Str": "<div>Hello</div><p>This is a test with <strong>bold text</strong> and some <em>emphasis</em></p>",
			},
		},
		{
			name: "Long HTML email with emojis",
			input: "<html><body><div class='header'><h1>🎉 Special Offer!</h1></div>" +
				"<div class='content'><p>" + strings.Repeat("Take advantage of our amazing deals and discounts today! ", 15) + "</p>" +
				"<ul><li>Free shipping on orders over $50</li><li>30% off all electronics</li><li>Buy 2 get 1 free on books</li></ul></div>" +
				"👋<div class='footer'><p>" + strings.Repeat("Contact us at support@example.com for any questions. ", 12) + "</p></div>🎊</body></html>",
			expected: map[string]string{
				"htmlContent0Str": "<html><body><div class='header'><h1>🎉 Special Offer!</h1></div><div class='content'><p>" + 
					strings.Repeat("Take advantage of our amazing deals and discounts today! ", 15) + 
					"</p><ul><li>Free shipping on orders over $50</li><li>30% off all electronics",
				"htmlContent1Str": "</li><li>Buy 2 get 1 free on books</li></ul></div>👋<div class='footer'><p>" + 
					strings.Repeat("Contact us at support@example.com for any questions. ", 12) + 
					"</p></div>🎊</body></html>",
			},
		},
		{
			name: "Complex HTML with nested tags and emojis",
			input: "<div class='newsletter'><header><h2>📧 Monthly Newsletter</h2></header>" +
				"<article><section><h3>Featured Products</h3><p>" + strings.Repeat("Check out our latest and greatest products! ", 20) + "</p></section>" +
				"<section><h3>Customer Reviews</h3><blockquote>" + strings.Repeat("Amazing service and quality! Highly recommended. ", 15) + "</blockquote></section></article>" +
				"<footer>👨‍💼 From the team at Example Corp 🚀</footer></div>",
			expected: map[string]string{
				"htmlContent0Str": "<div class='newsletter'><header><h2>📧 Monthly Newsletter</h2></header><article><section><h3>Featured Products</h3><p>" + 
					strings.Repeat("Check out our latest and greatest products! ", 20) + 
					"</p></section><section>",
				"htmlContent1Str": "<h3>Customer Reviews</h3><blockquote>" + 
					strings.Repeat("Amazing service and quality! Highly recommended. ", 15) + 
					"</blockquote></section></article><footer>👨‍💼 From the team at Example Corp 🚀</footer></div>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chunkHTMLBody(tt.input)
			assert.Equal(t, tt.expected, result, "chunks should match expected output")
		})
	}
}
