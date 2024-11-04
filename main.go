package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
)

const (
	statsNamespace = "STATS"
)

// Request structure for the fizz-buzz parameters
type FizzBuzzRequest struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

// Response structure for the fizz-buzz result
type FizzBuzzResponse struct {
	Result []string `json:"result"`
}

// Response structure for the stats result
type StatsResponse struct {
	MostFrequentRequest FizzBuzzRequest `json:"most_frequent_request"`
	Hits                int    `json:"hits"`
}

func main() {
	r := chi.NewRouter()

	// Define a middleware for logging
	r.Use(loggingMiddleware)

	// FizzBuzz endpoint
	r.Post("/api/fizzbuzz", fizzBuzzHandler)

	// Stats endpoint
	r.Get("/api/stats", statsHandler)

	// Start the server with graceful shutdown logic
	workers.Serve(r)
}

// fizzBuzzHandler handles the fizz-buzz requests
func fizzBuzzHandler(w http.ResponseWriter, req *http.Request) {
	var fizzBuzzReq FizzBuzzRequest

	// Decode the request body into the FizzBuzzRequest struct
	if err := json.NewDecoder(req.Body).Decode(&fizzBuzzReq); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate input parameters
	if fizzBuzzReq.Int1 <= 0 || fizzBuzzReq.Int2 <= 0 || fizzBuzzReq.Limit <= 0 {
		http.Error(w, "int1, int2, and limit must be positive integers", http.StatusBadRequest)
		return
	}

	// Perform the fizz-buzz logic
	result := fizzBuzz(fizzBuzzReq.Int1, fizzBuzzReq.Int2, fizzBuzzReq.Limit, fizzBuzzReq.Str1, fizzBuzzReq.Str2)

	// Serialize request parameters as the KV key
	reqKey, _ := json.Marshal(fizzBuzzReq)

	// Initialize KV namespace
	kv, err := cloudflare.NewKVNamespace(statsNamespace)
	if err != nil {
		handleErr(w, "failed to initialize KV", err)
		return
	}

	// Increment or initialize count for this request
	incrementRequestCount(kv, string(reqKey))

	// Prepare the response
	fizzBuzzResp := FizzBuzzResponse{Result: result}

	// Set the response header and write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fizzBuzzResp)
}

// statsHandler returns the most frequent request parameters and hit count
func statsHandler(w http.ResponseWriter, req *http.Request) {
	// Initialize the KV namespace
	kv, err := cloudflare.NewKVNamespace(statsNamespace)
	if err != nil {
		handleErr(w, "failed to initialize KV", err)
		return
	}

	// List all keys in the namespace
	opts := &cloudflare.KVNamespaceListOptions{
		Limit: 1000, // Adjust this limit based on the expected number of entries
	}
	var mostFrequentKey string
	maxCount := 0

	for {
		result, err := kv.List(opts)
		if err != nil {
			handleErr(w, "failed to list keys in KV", err)
			return
		}

		// Check each key's count
		for _, item := range result.Keys {
			countStr, _ := kv.GetString(item.Name, nil)
			count, _ := strconv.Atoi(countStr)
			if count > maxCount {
				maxCount = count
				mostFrequentKey = item.Name
			}
		}

		// If there are no more keys to paginate through, break
		if result.ListComplete {
			break
		}

		// Update the cursor to continue with the next set of keys
		opts.Cursor = result.Cursor
	}

	// Prepare and send the response
	if mostFrequentKey == "" {
		http.Error(w, "No requests made yet.", http.StatusNotFound)
		return
	}
	
	// Decode the stats JSON string into a StatsRequest struct
	var mostFrequent FizzBuzzRequest
	if err := json.Unmarshal([]byte(mostFrequentKey), &mostFrequent); err != nil {
		http.Error(w, "failed to decode stats", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	statsResp := StatsResponse{MostFrequentRequest: mostFrequent, Hits: maxCount}

	// Set the response header and write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statsResp)
}

// FizzBuzz function that generates the desired output
func fizzBuzz(int1, int2, limit int, str1, str2 string) []string {
	var result []string
	for i := 1; i <= limit; i++ {
		if i%int1 == 0 && i%int2 == 0 {
			result = append(result, str1+str2)
		} else if i%int1 == 0 {
			result = append(result, str1)
		} else if i%int2 == 0 {
			result = append(result, str2)
		} else {
			result = append(result, fmt.Sprint(i))
		}
	}
	return result
}

// incrementRequestCount updates the request count in KV
func incrementRequestCount(kv *cloudflare.KVNamespace, key string) {
	countStr, _ := kv.GetString(key, nil)
	count, _ := strconv.Atoi(countStr)
	kv.PutString(key, strconv.Itoa(count+1), nil)
}

// loggingMiddleware logs incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Received request: %s %s", req.Method, req.URL.Path)
		next.ServeHTTP(w, req)
	})
}

func handleErr(w http.ResponseWriter, msg string, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}
