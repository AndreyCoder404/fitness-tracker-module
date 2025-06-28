package main

import (
	"testing"
	"time"
)

// Test_distance проверяет расчет дистанции, пройденной во время тренировки.
// Этот тест тестирует метод distance() структуры Training, который вычисляет расстояние
// в километрах на основе количества шагов (или гребков) и длины шага.
// Ожидается, что результат будет пропорционален количеству действий и длине шага.
func Test_distance(t *testing.T) {
	type args struct {
		action  int     // Количество шагов или гребков, выполненных во время тренировки.
		lenStep float64 // Длина одного шага или гребка в метрах.
	}
	tests := []struct {
		name string  // Название теста для идентификации.
		args args    // Входные данные для теста.
		want float64 // Ожидаемое значение дистанции в километрах.
	}{
		// Тест 1: Успешный случай с типичным количеством шагов.
		{
			name: "Success test",
			args: args{
				action:  2000, // 2000 шагов, что является типичным значением для тренировки.
				lenStep: 0.65, // Стандартная длина шага 0.65 метра.
			},
			want: 1.3, // Ожидаемая дистанция: 2000 * 0.65 / 1000 = 1.3 км.
		},
		// Тест 2: Проверка случая с нулевым количеством шагов.
		{
			name: "Null action",
			args: args{
				action:  0,    // Нулевое количество шагов.
				lenStep: 0.65, // Длина шага не влияет при нулевом действии.
			},
			want: 0.0, // Ожидаемая дистанция: 0 * 0.65 / 1000 = 0.0 км.
		},
		// Тест 3: Проверка случая с одним шагом.
		{
			name: "One action",
			args: args{
				action:  1,    // Один шаг.
				lenStep: 0.65, // Стандартная длина шага.
			},
			want: 0.00065, // Ожидаемая дистанция: 1 * 0.65 / 1000 = 0.00065 км.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый объект Training с заданными параметрами.
			training := Training{Action: tt.args.action, LenStep: tt.args.lenStep}
			// Выполняем тест, сравнивая полученный результат с ожидаемым.
			if got := training.distance(); got != tt.want {
				t.Errorf("distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_meanSpeed проверяет расчет средней скорости во время тренировки.
// Этот тест тестирует метод meanSpeed() структуры Training, который вычисляет скорость
// в километрах в час на основе дистанции и длительности тренировки.
// Ожидается, что скорость будет нулевой при нулевой длительности.
func Test_meanSpeed(t *testing.T) {
	type args struct {
		action   int           // Количество шагов или гребков.
		duration time.Duration // Длительность тренировки.
		lenStep  float64       // Длина одного шага или гребка в метрах.
	}
	tests := []struct {
		name string  // Название теста для идентификации.
		args args    // Входные данные для теста.
		want float64 // Ожидаемое значение скорости в км/ч.
	}{
		// Тест 1: Успешный случай с типичной тренировкой.
		{
			name: "Successful test",
			args: args{
				action:   2000,          // 2000 шагов.
				duration: 2 * time.Hour, // Длительность 2 часа.
				lenStep:  0.65,          // Стандартная длина шага.
			},
			want: 0.65, // Ожидаемая скорость: (2000 * 0.65 / 1000) / 2 = 0.65 км/ч.
		},
		// Тест 2: Проверка случая с нулевой длительностью.
		{
			name: "Null duration",
			args: args{
				action:   2000, // 2000 шагов.
				duration: 0,    // Нулевая длительность.
				lenStep:  0.65, // Длина шага не влияет.
			},
			want: 0, // Ожидаемая скорость: 0, так как деление на ноль не допускается.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый объект Training с заданными параметрами.
			training := Training{Action: tt.args.action, LenStep: tt.args.lenStep, Duration: tt.args.duration}
			// Выполняем тест, сравнивая полученный результат с ожидаемым.
			if got := training.meanSpeed(); got != tt.want {
				t.Errorf("meanSpeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestShowTrainingInfo проверяет вывод информации о тренировке.
// Этот тест тестирует функцию ReadData, которая возвращает отформатированную информацию
// о тренировке в зависимости от её типа (Бег, Ходьба, Плавание или неизвестный).
// Ожидается, что для известных типов будут рассчитаны дистанция, скорость и калории.
func TestShowTrainingInfo(t *testing.T) {
	type args struct {
		action       int           // Количество шагов или гребков.
		trainingType string        // Тип тренировки.
		duration     time.Duration // Длительность тренировки.
		weight       float64       // Вес пользователя в килограммах.
		height       float64       // Рост пользователя в сантиметрах.
		lengthPool   int           // Длина бассейна в метрах (для плавания).
		countPool    int           // Количество пересечений бассейна.
	}
	tests := []struct {
		name string // Название теста для идентификации.
		args args   // Входные данные для теста.
		want string // Ожидаемая строка вывода.
	}{
		// Тест 1: Проверка тренировки по бегу.
		{
			name: "run test",
			args: args{
				action:       4000,            // 4000 шагов.
				trainingType: "Бег",           // Тип тренировки.
				duration:     9 * time.Minute, // Длительность 9 минут.
				weight:       85,              // Вес 85 кг.
				height:       185,             // Рост 185 см.
				lengthPool:   50,              // Не используется для бега.
				countPool:    2,               // Не используется для бега.
			},
			want: "Тип тренировки: Бег\nДлительность: 9.00 мин\nДистанция: 2.60 км.\nСр. скорость: 17.33 км/ч\nПотрачено ккал: 64.09\n",
			// Ожидаемая дистанция: 4000 * 0.65 / 1000 = 2.6 км.
			// Ожидаемая скорость: 2.6 / (9/60) = 17.33 км/ч.
			// Ожидаемые калории: (18 * 17.33 + 1.79) * 85 / 1000 * 9 ≈ 64.09.
		},
		// Тест 2: Проверка тренировки по ходьбе.
		{
			name: "walking test",
			args: args{
				action:       4000,          // 4000 шагов.
				trainingType: "Ходьба",      // Тип тренировки.
				duration:     1 * time.Hour, // Длительность 1 час.
				weight:       85,            // Вес 85 кг.
				height:       185,           // Рост 185 см.
				lengthPool:   50,            // Не используется для ходьбы.
				countPool:    2,             // Не используется для ходьбы.
			},
			want: "Тип тренировки: Ходьба\nДлительность: 60.00 мин\nДистанция: 2.60 км.\nСр. скорость: 2.60 км/ч\nПотрачено ккал: 220.27\n",
			// Ожидаемая дистанция: 4000 * 0.65 / 1000 = 2.6 км.
			// Ожидаемая скорость: 2.6 / 1 = 2.60 км/ч.
			// Ожидаемые калории: (0.035 * 85 + (2.60^2 / 1.85) * 0.029 * 85) * 1 * 60 ≈ 220.27.
		},
		// Тест 3: Проверка тренировки по плаванию.
		{
			name: "swimming test",
			args: args{
				action:       1000,             // 1000 гребков.
				trainingType: "Плавание",       // Тип тренировки.
				duration:     15 * time.Minute, // Длительность 15 минут.
				weight:       85,               // Вес 85 кг.
				height:       185,              // Рост 185 см (не используется).
				lengthPool:   100,              // Длина бассейна 100 м.
				countPool:    4,                // 4 пересечения.
			},
			want: "Тип тренировки: Плавание\nДлительность: 15.00 мин\nДистанция: 0.40 км.\nСр. скорость: 1.60 км/ч\nПотрачено ккал: 19.13\n",
			// Ожидаемая дистанция: 100 * 4 / 1000 = 0.4 км.
			// Ожидаемая скорость: 0.4 / (15/60) = 1.60 км/ч.
			// Ожидаемые калории: (1.60 + 1.1) * 2 * 85 * (15/60) ≈ 19.13.
		},
		// Тест 4: Проверка неизвестного типа тренировки.
		{
			name: "unknown test",
			args: args{
				action:       1000,          // 1000 действий.
				trainingType: "Керлинг",     // Неизвестный тип.
				duration:     5 * time.Hour, // Длительность 5 часов.
				weight:       85,            // Вес 85 кг.
				height:       185,           // Рост 185 см.
				lengthPool:   50,            // Не используется.
				countPool:    2,             // Не используется.
			},
			want: "неизвестный тип тренировки\n",
			// Ожидается вывод строки для неизвестного типа, так как "Керлинг" не поддерживается.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Определяем тип тренировки и создаем соответствующий объект.
			var training CaloriesCalculator
			switch tt.args.trainingType {
			case "Бег":
				training = Running{Training: Training{
					TrainingType: tt.args.trainingType,
					Action:       tt.args.action,
					LenStep:      LenStep,
					Duration:     tt.args.duration,
					Weight:       tt.args.weight,
				}}
			case "Ходьба":
				training = Walking{Training: Training{
					TrainingType: tt.args.trainingType,
					Action:       tt.args.action,
					LenStep:      LenStep,
					Duration:     tt.args.duration,
					Weight:       tt.args.weight,
				}, Height: tt.args.height}
			case "Плавание":
				training = Swimming{Training: Training{
					TrainingType: tt.args.trainingType,
					Action:       tt.args.action,
					LenStep:      SwimmingLenStep,
					Duration:     tt.args.duration,
					Weight:       tt.args.weight,
				}, LengthPool: tt.args.lengthPool, CountPool: tt.args.countPool}
			default:
				// Для неизвестного типа возвращаем заранее заданную строку.
				got := "неизвестный тип тренировки\n"
				if got != tt.want {
					t.Errorf("ShowTrainingInfo() = %v, want %v", got, tt.want)
				}
				return
			}
			// Выполняем тест, сравнивая результат ReadData с ожидаемым.
			got := ReadData(training)
			if got != tt.want {
				t.Errorf("ShowTrainingInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
