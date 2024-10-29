package main

import (
	"fmt"
	"os"
)

// Заданные координаты X и Y для построения сплайна
var X = [...]float64{0.29, 0.40, 0.81, 0.83, 1.27, 1.72, 2.11}
var Y = [...]float64{1.336, 1.494, 2.247, 2.293, 3.560, 5.584, 8.248}

// Количество шагов для интерполяции на каждом интервале
var N int = 1000

// Имя файла для хранения данных, вычисленных точек параболического сплайна
var NameMethodData string = "ParabolicSpline.dat"

// CalculateParamA вычисляет параметр 'a' для параболического сплайна, 
// равный значению функции в предыдущей точке yi_p
func CalculateParamA(yi_p float64) float64 {
	return yi_p
}

// CalculateParamB вычисляет параметр 'b' для сплайна на основе текущего и 
// предыдущего значений y, шага интервала hi и параметра ci
func CalculateParamB(yi, yi_p, hi, ci float64) (bi float64) {
	bi = ((yi - yi_p) / hi) - hi*ci
	return bi
}

// CalculateParamCn вычисляет параметр 'c' для последнего интервала, 
// где требуется начальное значение ci для обратного расчета
func CalculateParamCn(yn, yn_p, xn, xn_p float64) (cn float64) {
	cn = (yn_p - yn) / ((xn - xn_p) * (xn - xn_p))
	return cn
}

// CalculateParamC вычисляет параметр 'c' для текущего интервала, 
// используя 'c' следующего интервала, значения y и x текущего и предыдущего интервалов
func CalculateParamC(ci_next, yi_n, yi, yi_p, xi_n, xi, xi_p float64) (ci float64) {
	hi_n := (xi_n - xi)
	hi := (xi - xi_p)
	gi := ((yi_n - yi) / hi_n) - ((yi - yi_p) / hi)
	ci = (gi - ci_next*hi_n) / hi
	return ci
}

// ParabolicSpline рассчитывает значения x и y для параболического сплайна 
// на интервале [xi_p, xi] с заданными параметрами a, b, c
func ParabolicSpline(ai, bi, ci, xi, xi_p float64) ([]float64, []float64) {
	Px := make([]float64, N)
	Py := make([]float64, N)
	// Шаг интерполяции
	var h float64 = (xi - xi_p) / float64(N)
	for i := 0; i < N; i++ {
		x := xi_p + float64(i)*h
		Px[i] = x
		Py[i] = ai + bi*(x - xi_p) + ci*(x - xi_p)*(x - xi_p)
	}
	return Px, Py
}

// ParseToPoint преобразует массивы x и y значений в срез точек типа Point
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

// main выполняет расчет параметров и построение параболического сплайна 
// для каждого интервала, используя обратный расчет значений параметра c
func main() {
	var hi float64 // Переменная для хранения шага интервала

	// Создаем каналы для передачи сообщений
	pointsSend := make(chan []Point)
	end := make(chan bool)

	// Удаляем предыдущие данные
	err := os.Remove(NameMethodData)
	if err != nil {
		fmt.Println("Ошибка при очистке старых данных")
	}
	// Кол-во узлов
	n := len(X)

	// Запускаем вторым независимым потоком программу для сохранения точек
	go ExportSplineDataToData(pointsSend, end, NameMethodData)

	// Вычисляем значение параметра c для последнего интервала
	ci := CalculateParamCn(Y[n-1], Y[n-2], X[n-1], X[n-2])

	// Цикл обратного вычисления параметров на интервалах [X[i], X[i+1]]
	for i := n - 2; i >= 1; i-- {
		fmt.Println(fmt.Sprintf("ОБРАБОТКА ИНТЕРВАЛА [%f, %f]", X[i+1], X[i]))
		hi = X[i+1] - X[i]

		// Вычисляем параметры a, b
		ai := CalculateParamA(Y[i])
		bi := CalculateParamB(Y[i+1], Y[i], hi, ci)
		fmt.Println("Значения параметров:")
		fmt.Println(fmt.Sprintf("%d) a = %f, b = %f, c = %f", i, ai, bi, ci))

		// Расчет значений на интервале с текущими параметрами
		Px, Py := ParabolicSpline(ai, bi, ci, X[i+1], X[i])
		points := ParseToPoint(Px, Py)

		// Отправляем точки на сохранение
		pointsSend <- points

		// Обновление параметра c для следующего интервала
		ci = CalculateParamC(ci, Y[i+1], Y[i], Y[i-1], X[i+1], X[i], X[i-1])
	}

	// Обработка первого интервала [X[0], X[1]]
	fmt.Println(fmt.Sprintf("ОБРАБОТКА ИНТЕРВАЛА [%f, %f]", X[1], X[0]))
	h1 := X[1] - X[0]
	a1 := CalculateParamA(Y[0])
	b1 := CalculateParamB(Y[1], Y[0], h1, ci)

	// Расчет значений на первом интервале
	Px, Py := ParabolicSpline(a1, b1, ci, X[1], X[0])
	points := ParseToPoint(Px, Py)

	fmt.Println("Значения параметров:")
	fmt.Println(fmt.Sprintf("0) a = %f, b = %f, c = %f", a1, b1, ci))

	pointsSend <- points

	// Сохраняем исходные точки X и Y для графика
	TablePoints := make([]Point, len(X))
	for i, _ := range X {
		TablePoints[i] = Point{
			X: X[i],
			Y: Y[i],
		}
	}
	// Сигнал о ззаврещении работы горутины
	end <- true

	// Построение графика (вызов закомментирован)
	//BuildGraph(NameMethodData, TablePoints)
}
