package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	webhookEndpoint := getEnv("WEBHOOK_ENDPOINT", "/github-webhook/")
	proxyPort := getEnv("PROXY_PORT", ":8080")
	jenkinsURL := getEnv("JENKINS_URL", "http://localhost:8081")
	jenkinsWebhookURL := fmt.Sprintf("%s%s", jenkinsURL, webhookEndpoint)

	http.HandleFunc(webhookEndpoint, func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, jenkinsWebhookURL)
	})

	log.Printf("Starting proxy server on port %s...\n", proxyPort)
	log.Fatal(http.ListenAndServe(proxyPort, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request, jenkinsWebhookURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	signature := r.Header.Get("X-Hub-Signature")
	if signature == "" {
		http.NotFound(w, r)
		return
	}

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("Missing GITHUB_WEBHOOK_SECRET environment variable")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if !validateSignature(secret, signature, body) {
		http.NotFound(w, r)
		return
	}

	proxyReq, err := http.NewRequest(http.MethodPost, jenkinsWebhookURL, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = r.Header.Clone()
	proxyReq.Header.Del("X-Forwarded-Proto")

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
}

func validateSignature(secret, signature string, payload []byte) bool {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := "sha1=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
