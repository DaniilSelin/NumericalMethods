package main

import (
	"fmt"
	"os"
)

// Массивы с исходными координатами (X, Y) для кубического сплайна
var X = [...]float64{0.29, 0.40, 0.81, 0.83, 1.27, 1.72, 2.11}
var Y = [...]float64{1.336, 1.494, 2.247, 2.293, 3.560, 5.584, 8.248}

// Количество шагов для интерполяции
var N int = 1000000

// Имя файла для сохранения результатов вычислений
var NameMethodData string = "CubicSpline.dat"

// Функция для вычисления параметра a
func CalculateParamA(yi_p float64) float64 {
	return yi_p
}

// Функция для вычисления параметра b
func CalculateParamB(yi, yi_p, hi, ci_next, ci float64) (bi float64) {
	bi = ((yi - yi_p) / hi) - (hi*(ci_next + 2*ci)) / 3 // Формула для b
	return bi
}

// Функция для вычисления параметра d
func CalculateParamD(ci_next, ci, hi float64) (di float64) {
	di = (ci_next - ci) / (3 * hi) // Формула для d
	return di
}

// Функция для вычисления параметра e
func CalculateParamE(hi, hi_p, ei_p float64) (ei float64) {
	ei = (-hi) / (hi_p*ei_p + 2*(hi_p + hi)) // Формула для e
	return ei
}

// Функция для вычисления параметра n
func CalculateParamN(hi, hi_p, ni_p, ei_p, yi, yi_p, yi_pp float64) (ni float64) {
	gi := 3 * (((yi - yi_p) / hi) - ((yi_p - yi_pp) / hi_p)) // Формула для g
	ni = (gi - hi_p*ni_p) / (hi_p*ei_p + 2*(hi_p + hi)) // Формула для n
	return ni
}

// Функция для вычисления параметра c
func CalculateParamC(ci_next, ei, ni float64) (ci float64) {
	ci = ei*ci_next + ni // Формула для c
	return ci
}

// Основная функция кубического сплайна
// N - количество шагов на интервале [xi, xi_p]
func CubicSpline(ai, bi, ci, di, xi, xi_p float64) ([]float64, []float64) {
	Px := make([]float64, N) // Массив для хранения x координат
	Py := make([]float64, N) // Массив для хранения y координат
	var h float64 = (xi - xi_p) / float64(N) // Вычисляем шаг интерполяции

	// Заполняем массивы Px и Py
	for i := 0; i < N; i++ {
		x := xi_p + float64(i)*h
		Px[i] = x // Запись x в массив
		// Вычисление y с использованием кубического полинома
		Py[i] = ai + bi*(x-xi_p) + ci*(x-xi_p)*(x-xi_p) + di*(x-xi_p)*(x-xi_p)*(x-xi_p)
	}
	return Px, Py // Возвращаем массивы координат
}

// Структура для представления точки
type Point struct {
	X float64
	Y float64
}

// Функция для преобразования массивов Px и Py в массив точек
func ParseToPoint(Px, Py []float64) []Point {
	points := make([]Point, N) // Массив точек

	// Заполняем массив точек
	for i, _ := range Px {
		points[i] = Point{
			X: Px[i],
			Y: Py[i],
		}
	}

	return points // Возвращаем массив точек
}

// Главная функция программы
func main() {
	var hi float64 // Разница между текущими x
	var hi_p float64 // Разница между предыдущими x
	var n int = len(X) // Количество точек

	// Создаем каналы для передачи сообщений
	pointsSend := make(chan []Point)
	end := make(chan bool)

	// Удаляем старые данные из файла, если он существует
	err := os.Remove(NameMethodData)
	if err != nil {
		fmt.Println("Ошибка при очистке старых данных")
	}

	// Запускаем вторым независимым потоком программу для сохранения точек
	go ExportSplineDataToData(pointsSend, end, NameMethodData)

	// Массивы для хранения прогоночных коэффициентов
	RunCoefN := make([]float64, len(X))
	RunCoefE := make([]float64, len(X))

	// Начинаем с 1 (первый индекс), и до предпоследнего (n-2), так как cn = 0
	for i := 2; i < n; i++ {
		hi = X[i] - X[i-1] // Вычисляем hi
		hi_p = X[i-1] - X[i-2] // Вычисляем hi_p
		RunCoefE[i-1] = CalculateParamE(hi, hi_p, RunCoefE[i-2]) // Вычисляем e
		RunCoefN[i-1] = CalculateParamN(hi, hi_p, RunCoefN[i-2], RunCoefE[i-2], Y[i], Y[i-1], Y[i-2]) // Вычисляем n
	}

	fmt.Println("Прогоночные коэфициенты: \n", RunCoefE, RunCoefN)

	// Задаем последний коэффициент c
	ci_next := RunCoefN[n-1]

	// Основной цикл для вычисления сплайнов по интервалам
	for i := n - 2; i >= 1; i-- {
		fmt.Println(fmt.Sprintf("ОБРАБОТКА ИНТЕРВАЛА [%f, %f]", X[i+1], X[i]))
		hi = X[i+1] - X[i] // Вычисляем текущий шаг

		// Вычисляем параметры сплайна
		ai := CalculateParamA(Y[i]) // Параметр a
		ci := CalculateParamC(ci_next, RunCoefE[i], RunCoefN[i]) // Параметр c
		bi := CalculateParamB(Y[i+1], Y[i], hi, ci_next, ci) // Параметр b
		di := CalculateParamD(ci_next, ci, hi) // Параметр d

		fmt.Println("Значения параметров:")
		fmt.Println(fmt.Sprintf("%d) a = %f, b = %f, c = %f, d = %f", i, ai, bi, ci, di))

		// Вычисляем точки кубического сплайна
		Px, Py := CubicSpline(ai, bi, ci, di, X[i+1], X[i])

		points := ParseToPoint(Px, Py) // Преобразуем в точки

		// Отправляем точки на сохранение
		pointsSend <- points

		// Обновляем ci_next для следующей итерации
		ci_next = ci
	}

	// Обработка первых значений параметров
	hi = X[1] - X[0]
	ai := CalculateParamA(Y[0])
	ci := CalculateParamC(ci_next, RunCoefE[0], RunCoefN[0])
	bi := CalculateParamB(Y[1], Y[0], hi, ci_next, ci)
	di := CalculateParamD(ci_next, ci, hi)

	fmt.Println("Значения параметров:")
	fmt.Println(fmt.Sprintf("%d) a = %f, b = %f, c = %f, d = %f", 0, ai, bi, ci, di))

	// Вычисляем точки для первых значений
	Px, Py := CubicSpline(ai, bi, ci, di, X[1], X[0])

	points := ParseToPoint(Px, Py) // Преобразуем в точки

	// Отправляем точки на сохранение
	pointsSend <- points

	// Создаем массив точек для исходных данных
	TablePoints := make([]Point, len(X))
	for i, _ := range X {
		TablePoints[i] = Point{
			X: X[i],
			Y: Y[i],
		}
	}

	// Сигнал о завершении работы
	end <- true

	// Вызов функции для построения графика (закомментировано)
	// BuildGraph(NameMethodData, TablePoints)
}
