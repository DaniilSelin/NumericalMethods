package main

import (
	"fmt"
	"os"
)
// Исходные данные: значения X и Y, по которым строится сплайн
var X = [...]float64{0.29,0.40,0.81,0.83,1.27,1.72,2.11}
var Y = [...]float64{1.336,1.494,2.247,2.293,3.560,5.584,8.248}

// Число шагов для вычислений в интервале [xi, xi+1]
var N int = 1000

// Имя файла для сохранения рассчитанных точек сплайна
var NameMethodData string = "LineSpline.dat"

// CalculateParamA возвращает значение параметра `a`
// yi_p - значение yi-1 для текущего интервала
func CalculateParamA(yi_p float64) float64 {
	return yi_p
}

// CalculateParamB вычисляет значение параметра `b`
// yi, yi_p - значения Y на концах интервала
// hi - расстояние между значениями X[i] и X[i-1]
func CalculateParamB(yi, yi_p, hi float64) (bi float64) {
	bi = (yi - yi_p) / hi
	return bi
}

// LineSpline вычисляет координаты точек линейного сплайна в интервале [xi, xi+1]
// ai, bi - параметры сплайна для интервала
// xi, xi_p - координаты начальной и конечной точек интервала
func LineSpline(ai, bi, xi, xi_p float64,) ([]float64, []float64) {
	Px := make([]float64, N) // Массив X координат точек сплайна
	Py := make([]float64, N) // Массив Y координат точек сплайна
	
	// Вычисляем шаг для интервала
	var h float64 = (xi - xi_p) / float64(N)
	for i := 0; i < N; i++ {
		x := xi_p + float64(i)*h // Позиция X для точки сплайна
		Px[i] = x
		Py[i] = ai + bi * (x - xi_p)  // Вычисляем Y на основе параметров
	}
	return Px, Py
}

// ParseToPoint преобразует массивы Px и Py в массив структур Point
// Px, Py - массивы координат точек, возвращает массив Point
func ParseToPoint(Px, Py []float64) []Point {
	points := make([]Point, N)

	for i, _ := range Px {
		points[i] = Point{
			X: Px[i],
			Y: Py[i],
		}
	}

	return points
}

// Основная функция программы, вычисляющая точки сплайна
func main() {
	var hi float64 // Переменная для шага между значениями X[i+1] и X[i]

	pointsSend := make(chan []Point)
	end := make(chan bool)

	// Удаляем файл с предыдущими результатами, если он существует
	err := os.Remove(NameMethodData)
	if err != nil {
		fmt.Println("Ошибка при очистке старых данных")
	}

	// Запускаем вторым независимым потоком программу для сохранения точек
	go ExportSplineDataToData(pointsSend, end, NameMethodData)

	// Обходим все интервалы [X[i], X[i+1]] для вычисления линейного сплайна
	for i := range Y[1:] {
		fmt.Println(fmt.Sprintf("ОБРАБОТКА ИНТЕРВАЛА [%f, %f]", X[i], X[i+1]))

		hi = X[i+1] - X[i] // Расстояние между X[i+1] и X[i]

		// Параметры a и b для текущего интервала
		ai := CalculateParamA(Y[i])
		bi := CalculateParamB(Y[i+1], Y[i], hi)
		fmt.Println("Значения параметров:")
		fmt.Println(fmt.Sprintf("%d) a = %f, b = %f", i, ai, bi))

		// Вычисляем точки линейного сплайна
		Px, Py := LineSpline(ai, bi, X[i+1], X[i])

		// Преобразуем точки в структуру Point и сохраняем в файл
		points := ParseToPoint(Px, Py)
		
		pointsSend <- points
	}

	// Создаем массив исходных точек для построения графика
	TablePoints := make([]Point, len(X))
	for i := range X {
		TablePoints[i] = Point{
			X: X[i],
			Y: Y[i],
		}
	}

	end <- true

	// Построение графика по данным сплайна и исходным точкам
	BuildGraph(NameMethodData, TablePoints)
}