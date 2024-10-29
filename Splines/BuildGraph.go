package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/opts"
)

var (
    xMax float64 = 8
    xMin float64 = 0
    yMin float64 = -1
    yMax float64 = 1.2
)

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

//Функция для построения графика
func BuildGraph(DataFile string, TablePoints []Point) {
    // Открываем файл с данными
    file, err := os.Open(DataFile)
    if err != nil {
        log.Fatalf("Ошибка при открытии файла данных: %v", err)
    }
    defer file.Close()

    var xValues []string
    var yValues []opts.LineData

    // Чтение данных из файла
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
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
        xValues = append(xValues, fmt.Sprintf("%.6f", x))
        yValues = append(yValues, opts.LineData{Value: y})
    }

    // Создание нового графика
    lineChart := charts.NewLine()
    lineChart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Spline Points Plot"}))

    // Установка оси X и добавление основной серии данных
    lineChart.SetXAxis(xValues).AddSeries("Spline Data", yValues)

	// Установка границ и меток для осей
	lineChart.SetGlobalOptions(
	    charts.WithXAxisOpts(opts.XAxis{Name: "X", Min: xMin, Max: xMax + 1}),
	    charts.WithYAxisOpts(opts.YAxis{Name: "Y", Min: yMin, Max: yMax + 1}),
	)

    // Добавление дополнительных точек
    var additionalX []string
    var additionalY []opts.LineData
    for _, point := range TablePoints {
        additionalX = append(additionalX, fmt.Sprintf("%.6f", point.X)) // Форматируем X
        additionalY = append(additionalY, opts.LineData{Value: point.Y})
    }

    // Добавляем дополнительные точки
    lineChart.AddSeries("Additional Points", additionalY).SetXAxis(additionalX)

    // Сохранение графика в HTML файл
    f, err := os.Create("line_chart.html")
    if err != nil {
        log.Fatalf("Ошибка при создании файла для графика: %v", err)
    }
    defer f.Close()

    // Рендеринг графика в файл
    if err := lineChart.Render(f); err != nil {
        log.Fatalf("Ошибка при рендеринге графика: %v", err)
    }

    fmt.Println("График успешно построен и сохранен в ", DataFile)
}
