package supabase

import (
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

func ConnectStorage() storage_go.Client {
	storage_client := storage_go.NewClient(os.Getenv("SUPABASE_PROJECTID"), os.Getenv("SUPABASE_SECRET_API_KEY"), nil)

	return *storage_client
}
