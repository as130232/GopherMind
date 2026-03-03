package handler

import (
	"GopherMind/internal/interface"
	"GopherMind/internal/service"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

func AnalyzeFoodHandler(ai _interface.AiProvider, foodSvc *service.FoodService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "無法讀取圖片", http.StatusBadRequest)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		// 將圖片讀取為位元組 (Bytes)
		imgBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "讀取圖片內容失敗", http.StatusInternalServerError)
			return
		}

		// 取得檔案的 MIME 類型 (例如 image/jpeg, image/png)
		mimeType := header.Header.Get("Content-Type")
		reply, err := foodSvc.AnalyzeFoodNutrition(r.Context(), mimeType, imgBytes)

		resp := ChatResponse{
			Reply:  reply,
			Source: "GopherMind-" + ai.Name(),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "讀取圖片內容失敗", http.StatusInternalServerError)
			return
		}
	}
}
