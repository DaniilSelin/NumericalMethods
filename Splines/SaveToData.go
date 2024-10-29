package main

import (
	"log"
	"fmt"
	"os"
)

type Point struct {
	X float64
	Y float64
}

func SaveFileData(points []Point, SaveFile string) {
	// Создание файла для сохранения данных
	dataFile, err := os.OpenFile(SaveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Ошибка при создании файла данных: %v", err)
	}
	defer dataFile.Close() // Отложенное закрытие файла

	// Запись данных в файл
	for _, point := range points {
		// Форматируем данные и записываем в файл
		_, err := fmt.Fprintf(dataFile, "%f %f\n", point.X, point.Y)
		if err != nil {
			log.Fatalf("ошибка при записи данных в файл: %v", err)
		}
	}
}

func ExportSplineDataToData(points chan []Point,end chan bool, SaveFile string) {
	for {
		select {
		case points_mus := <- points:
			SaveFileData(points_mus, SaveFile)
		case <- end:
			break
		}
	}
}