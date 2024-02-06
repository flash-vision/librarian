# librarian
 (Go) Manages object storage and retrieval across multipe services
// Assume Librarian struct and NewLibrarian constructor are defined as before

func main() {
	ctx := context.Background()

	// Initialize GCS storage
	gcsStorage, err := NewGCSStorage(ctx, "your-gcs-bucket-name")
	if err != nil {
		log.Fatalf("Failed to create GCS storage: %v", err)
	}

	// Setup librarian with GCS
	backends := []StorageBackend{
		{
			Type:   CloudStorageType,
			Storer: gcsStorage,
		},
		// Add other backends as needed
	}

	librarian, err := NewLibrarian(backends, "path/to/logfile.log")
	if err != nil {
		log.Fatalf("Failed to initialize Librarian: %v", err)
	}

	// Example usage
	key := "exampleKey"
	value := []byte("Hello, GCS!")

	// Store
	if err := librarian.Store(ctx, key, value); err != nil {
		log.Printf("Error storing data: %v", err)
	}

	// Retrieve
	retrieved, err := librarian.Retrieve(ctx, key)
	if err != nil {
		log.Printf("Error retrieving data: %v", err)
	} else {
		log.Printf("Retrieved data: %s", string(retrieved))
	}
}
