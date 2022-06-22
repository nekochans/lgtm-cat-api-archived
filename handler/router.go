package handler

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	db "github.com/nekochans/lgtm-cat-api/db/sqlc"
	"github.com/nekochans/lgtm-cat-api/infrastructure"
	"github.com/nekochans/lgtm-cat-api/usecase/createltgmimage"
	"github.com/nekochans/lgtm-cat-api/usecase/fetchlgtmimages"
)

func NewRouter(uploader *manager.Uploader, q *db.Queries) *chi.Mux {
	uploadS3Bucket := os.Getenv("UPLOAD_S3_BUCKET_NAME")
	lgtmImagesCdnDomain := os.Getenv("LGTM_IMAGES_CDN_DOMAIN")

	s3Repository := infrastructure.NewS3Repository(uploader, uploadS3Bucket)
	createLgtmImageUseCase := createltgmimage.NewUseCase(s3Repository, lgtmImagesCdnDomain)
	createLgtmImageHandler := NewCreateLgtmImageHandler(createLgtmImageUseCase)

	lgtmImageRepository := infrastructure.NewLgtmImageRepository(q)
	extractRandomImagesUseCase := fetchlgtmimages.NewUseCase(lgtmImageRepository, lgtmImagesCdnDomain)
	extractRandomImagesHandler := NewFetchImagesHandler(extractRandomImagesUseCase)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:2222"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Post("/lgtm-images", createLgtmImageHandler.Create)
	r.Get("/lgtm-images", extractRandomImagesHandler.Extract)
	r.Get("/lgtm-images/recently-created", extractRandomImagesHandler.RetrieveRecentlyCreated)

	return r
}
