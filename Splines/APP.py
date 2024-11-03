from flask import Flask, render_template, request, redirect, url_for, jsonify
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from scipy.interpolate import CubicSpline
import subprocess
import os
from flask_cors import CORS

app = Flask(__name__)

CORS(app)  # Разрешение CORS для всех маршрутов

# Путь для сохранения данных
DATA_FILE = 'SetUpData.dat'


# Главная страница для ввода точек и выбора типа сплайна
@app.route('/')
def index():
    return render_template('index.html')


# Обработчик для сохранения точек и выбора сплайна
@app.route('/submit', methods=['POST'])
def submit():
    data = request.json  # Получаем JSON-данные
    points = data['points']
    splineInfo = data['splineInfo']

    try:
        spline_type = splineInfo.split(" ")[0]
        steps = int(splineInfo.split(" ")[1])
    except (IndexError, ValueError) as e:
        return jsonify(
            {'error': 'Ошибка в формате splineInfo. Убедитесь, что он содержит тип сплайна и количество шагов.'}), 400

    # Далее обрабатывайте точки и тип сплайна как вам нужно

    with open(DATA_FILE, 'w') as f:  # Открываем файл в режиме добавления
        f.write(f"{steps}\n")  # Сохраняем количество шагов

        # Парсим введенные точки
        points = [tuple(map(float, point.split(','))) for point in points.split(';')]
        df = pd.DataFrame(points, columns=['x', 'y'])

        # Сохраняем точки в CSV
        df.to_csv(f, index=False, header=False, sep=' ')

    # Построение графика
    if spline_type == 'line':
        # Путь к файлу с точками
        file_path = 'LineSpline.dat'
        subprocess.call(['./LineSpline'])
    elif spline_type == 'cubic':
        file_path = 'CubicSpline.dat'
        subprocess.call(['./CubicSpline'])
    else:
        file_path = 'ParabolicSpline.dat'
        subprocess.call(['./ParabolicSpline'])

    def load_data(file_path):
        x, y = [], []
        try:
            with open(file_path, 'r') as file:
                for line in file:
                    point = line.split()
                    if len(point) == 2:  # Убедитесь, что есть две координаты
                        x.append(float(point[0]))
                        y.append(float(point[1]))
            # Сортируем точки по значению x
            sorted_points = sorted(zip(x, y))
            x, y = zip(*sorted_points)
        except FileNotFoundError:
            print(f"Ошибка: файл {file_path} не найден.")
        except ValueError as e:
            print(f"Ошибка обработки данных: {e}")
        return x, y

    # Дополнительные точки для отображения
    x, y = load_data(file_path)

    additional_x = df['x'].tolist()  # Получение списка x
    additional_y = df['y'].tolist()  # Получение списка y

    # Создаем график
    plt.figure(figsize=(12, 6))
    plt.plot(x, y, linestyle='-', linewidth=2, color='blue', label=spline_type, alpha=0.8)  # Тонкая линия

    # Дополнительные точки
    plt.scatter(additional_x, additional_y, color='red', marker='o', s=100, edgecolor='black',
                label='УЗЛЫ')  # Дополнительные точки

    # Настройка графика
    plt.title('', fontsize=18, fontweight='bold')
    plt.xlabel('X', fontsize=14)
    plt.ylabel('Y', fontsize=14)
    plt.xlim(min(additional_x) - 0.5, max(additional_x) + 0.5)  # Установка границ по оси X
    plt.ylim(min(additional_y) - 0.5,
             max(additional_y) + 0.5)  # Установка границ по оси Y
    plt.axhline(0, color='black', linewidth=0.8, linestyle='--')  # Горизонтальная линия на нуле
    plt.axvline(0, color='black', linewidth=0.8, linestyle='--')  # Вертикальная линия на нуле
    plt.grid(True, linestyle='--', alpha=0.5)  # Сетка на фоне
    plt.legend()
    plt.tight_layout()  # Подгонка под размеры

    # Сохранение графика в файл
    plt.savefig('static/plot.png', dpi=300)  # Замените на нужное имя файла
    print("График сохранён как plot.png")

    # Запуск внешнего скрипта
    # subprocess.call(['python3', 'BuildGraph.py'])

    return jsonify({'message': 'Данные успешно обработаны'}), 200


@app.route('/load_data')
def load_data():
    data_file = 'SetUpData.dat'
    if os.path.exists(data_file):
        with open(data_file, 'r') as file:
            lines = file.readlines()
            points = [line.strip() for line in lines[1:]]  # Убираем первую строку
            return jsonify(points)  # Возвращаем только координаты
    return jsonify([])  # Если файл не найден, возвращаем пустой список


@app.route('/result')
def result():
    return render_template('result.html')


@app.route('/download_plot')
def download_plot():
    return redirect(url_for('static', filename='plot.png'))


if __name__ == '__main__':
    app.run(debug=True)
