import matplotlib
matplotlib.use('Agg')  # Неинтерактивный бэкенд
import matplotlib.pyplot as plt

# Функция для загрузки данных из файла
def load_data(file_path):
    x, y = [], []
    with open(file_path, 'r') as file:
        for line in file:
            point = line.split()
            x.append(float(point[0]))
            y.append(float(point[1]))
    # Сортируем точки по значению x
    sorted_points = sorted(zip(x, y))
    x, y = zip(*sorted_points)
    return x, y

# Путь к файлу с точками
file_path_1 = 'LineSpline.dat'
file_path = 'CubicSpline.dat'
file_path_2 = 'ParabolicSpline.dat'

# Загружаем данные из файла
x, y = load_data(file_path)
x_1, y_1 = load_data(file_path_1)
x_2, y_2 = load_data(file_path_2)

# Дополнительные точки для отображения
additional_x = [0.29, 0.40, 0.81, 0.83, 1.27, 1.72, 2.11]
additional_y = [1.336, 1.494, 2.247, 2.293, 3.560, 5.584, 8.248]
#additional_x = [1, 2.5, 3.5, 5.5, 6]
#additional_y = [0.9108, 0.7237,-0.2004,-0.5184, -0.0848]

# Создаем график
plt.figure(figsize=(12, 6))
plt.plot(x, y, linestyle='-', linewidth=2, color='blue', label='КУБИЧЕСКИЙ СПЛАЙН', alpha=0.8)  # Тонкая линия
plt.plot(x_2, y_2, linestyle='-', linewidth=2, color='red', label='ПАРАБОЛИЧЕСКИЙ СПЛАЙН', alpha=0.8)  # Тонкая линия
plt.plot(x_1, y_1, linestyle='-', linewidth=2, color='black', label='ЛИНЕЙНЫЙ СПЛАЙН', alpha=1)  # Тонкая линия

# Дополнительные точки
plt.scatter(additional_x, additional_y, color='red', marker='o', s=100, edgecolor='black', label='УЗЛЫ')  # Дополнительные точки

# Настройка графика
plt.title('', fontsize=18, fontweight='bold')
plt.xlabel('X', fontsize=14)
plt.ylabel('Y', fontsize=14)
plt.xlim(0, 2.5)  # Установка границ по оси X
plt.ylim(min(min(y), min(additional_y)) - 0.5, max(max(y), max(additional_y)) + 0.5)  # Установка границ по оси Y
plt.axhline(0, color='black', linewidth=0.8, linestyle='--')  # Горизонтальная линия на нуле
plt.axvline(0, color='black', linewidth=0.8, linestyle='--')  # Вертикальная линия на нуле
plt.grid(True, linestyle='--', alpha=0.5)  # Сетка на фоне
plt.legend()
plt.tight_layout()  # Подгонка под размеры

# Сохранение графика в файл
plt.savefig('plot.png', dpi=300)  # Замените на нужное имя файла
print("График сохранён как plot.png")
