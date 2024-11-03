package main

import (
	"log"
	"fmt"
	"bufio"
	"os"
	"strconv"
)

type Point struct {
	X float64
	Y float64
}

// Функция для разделения строки
func split(s string, sep rune) []string {
    var result []string
    current := ""
    for _, char := range s {
        if char == sep {
            result = append(result, current)
            current = ""
        } else {
            current += string(char)
        }
    }
    result = append(result, current)
    return result
}

func LoadSetUpData() (N int, X, Y []float64) {
	// Открываем файл с данными
    file, err := os.Open("SetUpData.dat")
    if err != nil {
        log.Fatalf("Ошибка при открытии файла данных: %v", err)
    }
    defer file.Close()

	// Чтение данных из файла
    scanner := bufio.NewScanner(file)

    scanner.Scan()

    line := scanner.Text()
    N, err = strconv.Atoi(line)
    if err != nil {
        log.Printf("Ошибка при парсинге N: %v", err)
    }

    for scanner.Scan() {
        line = scanner.Text()
        parts := split(line, ' ')
        if len(parts) != 2 {
            continue
        }

        // Считываем значения X и Y
        x, err := strconv.ParseFloat(parts[0], 64)
        if err != nil {
            log.Printf("Ошибка при парсинге X: %v", err)
            continue
        }
        y, err := strconv.ParseFloat(parts[1], 64)
        if err != nil {
            log.Printf("Ошибка при парсинге Y: %v", err)
            continue
        }

        // Форматируем X для отображения и добавляем Y
        X = append(X, x)
        Y = append(Y, y)
    }
    return
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