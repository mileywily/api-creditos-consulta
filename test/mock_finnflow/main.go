package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/core/v1/creditos", func(w http.ResponseWriter, r *http.Request) {
		// Validamos que el microservicio nos esté enviando los headers correctos (opcional)
		w.Header().Set("Content-Type", "application/json")
		
		// Retornamos un payload exitoso simulando a FinnFlow
		response := map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"numeroOperacionOriginal": "op-999888",
					"participantes": []map[string]interface{}{
						{"rut": "25671646-1"},
					},
					"propiedades": []map[string]interface{}{
						{"direccion": "Av. Siempre Viva 742"},
					},
					"politicas": map[string]interface{}{
						"diasMora": 0,
					},
				},
			},
		}
		
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Mock FinnFlow server corriendo en http://localhost:9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Println("Error al iniciar el mock:", err)
	}
}
