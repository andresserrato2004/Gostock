package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"BackEnd/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var bandera bool = false

type StockService struct {
	db     *gorm.DB
	Client *http.Client
}

func NewStockService(db *gorm.DB) *StockService {
	return &StockService{
		db:     db,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// IngestStocks descarga, limpia y almacena datos de acciones desde la API externa de KarenAI.
// Se ejecuta al inicio del servidor y maneja la paginación automática hasta traer todos los registros.
func (s *StockService) IngestStocks() error {

	/*var count int64
	if err := s.db.Model(&models.Stock{}).Count(&count).Error; err != nil {
		log.Println("Advertencia: No se pudo verificar el conteo de registros:", err)
	}
	if count > 0 {
		log.Printf("La base de datos ya contiene %d registros. Saltando ingestión inicial.\n", count)
		return nil
	}*/

	baseURL := os.Getenv("API_URL_Karen")
	if baseURL == "" {
		return fmt.Errorf("API_URL_Karen no esta en el archivo .env")
	}

	apiKey := os.Getenv("API_KEY_Karen")
	if apiKey == "" {
		return fmt.Errorf("API_KEY_Karen no esta en el archivo .env")
	}

	nextPage := ""

	for {

		FullUrl := baseURL

		if nextPage != "" {
			fmt.Println("entro aqui la primera vez")
			FullUrl = fmt.Sprintf("%s?next_page=%s", baseURL, nextPage)
		}

		// log.Println("Fetching URL:", FullUrl)

		var apiResp models.ApiResponse

		if !bandera {
			fmt.Println("entre por aui la primera vez")
			str, err := s.fetchWithRetries(FullUrl, apiKey, &apiResp, 3)
			if err != nil {
				return err
			}

			text := strings.Split(str, "\n")

			//fmt.Println(strings.Split(str, "\n"))

			//fmt.Println("next", string(rune(str[0]))+string(rune(str[1])))
			numeros := []string{}
			palabras := []string{}
			var nextPage2 string

			fmt.Println(text, "esto es lo que esta retornando la funcion")

			for index, palabra := range text {

				if index == 0 {
					//fmt.Println(palabra, " palabras")
					nextPage2 = palabra
				}

				if index == 1 {
					continue
				}

				if index >= 2 && index <= 11 {
					numeros = append(numeros, palabra)
				} else if index >= 12 && index <= 21 {
					palabras = append(palabras, palabra)
				}
			}

			/* 		numbers

			   		ini := 0
			*/

			var empresas []models.Empresas
			for k, nums := range numeros {
				arrnums := strings.Split(nums, " ")
				//		start := 0
				//var lentext int
				//fmt.Println((arrnums))
				intnums := []int{}

				for i, num := range arrnums {

					if i != 0 {
						val, _ := strconv.Atoi(num)
						intnums = append(intnums, val)
					}

				}
				result := []string{}
				//fmt.Println(intnums)
				start := 0
				i := 0
				for i < len(intnums) {

					cant := intnums[i]
					i++
					//fmt.Println(cant)
					var word string
					for j := 0; j < cant; j++ {

						longi := intnums[i]
						i++

						word += palabras[k][start : start+longi]
						start += longi
					}
					//fmt.Println(word)
					result = append(result, word)
				}

				empresa := models.Empresas{
					Ticker:      result[0],
					Target_from: result[1],
					Target_to:   result[2],
					Company:     result[3],
					Action:      result[4],
					Brokerage:   result[5],
					Rating_from: result[6],
					Rating_to:   result[7],
					Time:        result[8],
				}

				//	jsonData, _ := json.MarshalIndent(empresa, "", " ")
				empresas = append(empresas, empresa)

				//	fmt.Println(string(jsonData))
			}

			res := models.Creacion{
				Items:    empresas,
				NextPage: nextPage2,
			}

			jsonData, err := json.MarshalIndent(res, "", "    ")

			if err := json.Unmarshal(jsonData, &apiResp); err != nil {
				return err
			}

			//fmt.Println(string(jsonData))
			/* fmt.Println(numeros)
			fmt.Println(palabras)*/
			//fmt.Println(nextPage2)

			//fmt.Println(apiResp.Data)
			//fmt.Println(apiResp.NextPage)

			bandera = true
		} else {
			fmt.Println("entre por aui la segunda y etc vez")
			err := s.fetchWithRetries2(FullUrl, apiKey, &apiResp, 3)
			if err != nil {
				return err
			}
		}

		if len(apiResp.Data) > 0 {

			s.cleanData(apiResp.Data)
			// aqui es donde si ya hay datos entonces se actualiza solo una parte de la informacion
			// aqui es donde si ya hay datos entonces se actualiza solo una parte de la informacion
			// Upsert (Insertar o Actualizar): Si el Ticker ya existe, actualiza los campos clave.
			// Esto previene duplicados y mantiene la DB con la última info disponible.

			//fmt.Println(apiResp.NextPage, "sdkajhls")

			err := s.db.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "ticker"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"company", "brokerage", "action", "rating_from", "rating_to", "target_from_num", "target_to_num",
				}),
			}).Create(&apiResp.Data).Error

			if err != nil {
				log.Println("Error inserting stocks:", err)
			}
			//return json.NewDecoder(resp.Body).Decode(target)
		}

		nextPage = apiResp.NextPage
		fmt.Println(apiResp.NextPage, "hola123")
		/* fmt.Println(nextPage, "hola")
		 */
		if nextPage == "" {
			break
		}

		time.Sleep(500 * time.Millisecond) // Pausa para no sobrecargar el servidor
	}

	// log.Println("datos melos")

	return nil
}

// cleanData normaliza los datos "crudos" recibidos de la API.
// Convierte strings de precio ("$150.00") a float64 numérico para poder ordenar y filtrar.
func (s *StockService) cleanData(stocks []models.Stock) {
	for i := range stocks {
		stocks[i].TargetFromNum = parseCurrency(stocks[i].TargetFrom)
		stocks[i].TargetToNum = parseCurrency(stocks[i].TargetTo)
	}
}

func parseCurrency(value string) float64 {
	cleaned := strings.TrimSpace(strings.ReplaceAll(value, "$", ""))
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}
	return val
}

// fetchWithRetries realiza la petición HTTP con lógica de reintentos automática (Backoff).
func (s *StockService) fetchWithRetries(url, token string, target interface{}, maxRetries int) (string, error) {
	client := s.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	} /*
		fmt.Println("que paso aqui1")
		fmt.Println(url)
		fmt.Println("final de 1 ") */

	for i := 0; i < maxRetries; i++ {
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			str := string(body)
			fmt.Println(str, "esto es lo que llega en el bod")

			return str, err
			//return json.NewDecoder(resp.Body).Decode(target)
		} /*
			fmt.Println("que paso aqui")
			fmt.Println(url)
			fmt.Println("final de 1 ") */
		if err != nil {
			log.Printf("Error en la solicitud: %v", err)
		} else {
			log.Printf("Respuesta no exitosa: %s", resp.Status)
			resp.Body.Close()
		}

		time.Sleep(time.Duration(i+1) * time.Second)
	}
	a := " error"
	return a, fmt.Errorf("no se pudo obtener datos tras %d intentos", maxRetries)
}

func (s *StockService) fetchWithRetries2(url, token string, target interface{}, maxRetries int) error {
	client := s.Client
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	} /*
		fmt.Println("que paso aqui1")
		fmt.Println(url)
		fmt.Println("final de 1 ") */

	for i := 0; i < maxRetries; i++ {
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()

			/*body, _ := io.ReadAll(resp.Body)

			str := string(body)
			fmt.Println(str, "esto es lo que llega en el bod")
			*/
			return json.NewDecoder(resp.Body).Decode(target)
		} /*

			fmt.Println("que paso aqui")
			fmt.Println(url)
			fmt.Println("final de 1 ") */
		if err != nil {
			log.Printf("Error en la solicitud: %v", err)
		} else {
			log.Printf("Respuesta no exitosa: %s", resp.Status)
			resp.Body.Close()
		}

		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return fmt.Errorf("no se pudo obtener datos tras %d intentos", maxRetries)
}
