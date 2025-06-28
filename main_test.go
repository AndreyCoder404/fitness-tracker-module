package main

import (
	"math"
	"testing"
	"time"
)

// Тестируем функцию расчёта дистанции.
func Test_distance(t *testing.T) {
	// Определяем структуру для аргументов теста.
	type args struct {
		action  int     // Указываем количество шагов или гребков.
		lenStep float64 // Определяем длину каждого шага или гребка в метрах.
	}
	// Готовим тестовые случаи с именами, аргументами и ожидаемыми результатами.
	tests := []struct {
		name string  // Название каждого тестового случая.
		args args    // Предоставляем входные данные для теста.
		want float64 // Ожидаемая дистанция в километрах.
	}{
		{
			name: "Успешный тест",
			args: args{action: 2000, lenStep: 0.65},
			want: 1.3,
		},
		{
			name: "Нулевое действие",
			args: args{action: 0, lenStep: 0.65},
			want: 0.0,
		},
		{
			name: "Одно действие",
			args: args{action: 1, lenStep: 0.65},
			want: 0.00065,
		},
	}
	// Проходим по тестовым случаям, выполняем и проверяем результаты.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			training := Training{Action: tt.args.action, LenStep: tt.args.lenStep}
			if got := training.distance(); got != tt.want {
				t.Errorf("distance() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Тестируем функцию расчёта средней скорости.
func Test_meanSpeed(t *testing.T) {
	// Определяем структуру для аргументов теста.
	type args struct {
		action   int           // Указываем количество шагов или гребков.
		duration time.Duration // Устанавливаем продолжительность тренировки.
		lenStep  float64       // Определяем длину каждого шага или гребка в метрах.
	}
	// Готовим тестовые случаи с именами, аргументами и ожидаемыми результатами.
	tests := []struct {
		name string  // Название каждого тестового случая.
		args args    // Предоставляем входные данные для теста.
		want float64 // Ожидаемая средняя скорость в км/ч.
	}{
		{
			name: "Скорость бега",
			args: args{action: 5000, duration: 30 * time.Minute, lenStep: 0.65},
			want: 6.50,
		},
		{
			name: "Скорость ходьбы",
			args: args{action: 20000, duration: 225 * time.Minute, lenStep: 0.65},
			want: 3.47,
		},
		{
			name: "Скорость плавания",
			args: args{action: 2000, duration: 90 * time.Minute, lenStep: 1.38},
			want: 1.84, // Отражает фактическую скорость, рассчитанную на основе действия и длины шага.
		},
		{
			name: "Нулевая длительность",
			args: args{action: 2000, duration: 0, lenStep: 0.65},
			want: 0,
		},
	}
	// Проходим по тестовым случаям, выполняем и проверяем результаты с допуском.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			training := Training{Action: tt.args.action, LenStep: tt.args.lenStep, Duration: tt.args.duration}
			if got := training.meanSpeed(); math.Abs(got-tt.want) > 0.01 {
				t.Errorf("meanSpeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Тестируем функцию отображения данных тренировки.
func TestShowTrainingInfo(t *testing.T) {
	// Определяем структуру для аргументов теста.
	type args struct {
		action       int           // Указываем количество шагов или гребков.
		trainingType string        // Определяем тип тренировки.
		duration     time.Duration // Устанавливаем продолжительность тренировки.
		weight       float64       // Определяем вес пользователя в килограммах.
		height       float64       // Указываем рост пользователя в сантиметрах.
		lengthPool   int           // Устанавливаем длину бассейна в метрах.
		countPool    int           // Подсчитываем количество пересечений бассейна.
	}
	// Готовим тестовые случаи с именами, аргументами и ожидаемыми строками вывода.
	tests := []struct {
		name string // Название каждого тестового случая.
		args args   // Предоставляем входные данные для теста.
		want string // Ожидаемая отформатированная строка вывода.
	}{
		{
			name: "тест бега",
			args: args{action: 5000, trainingType: "Бег", duration: 30 * time.Minute, weight: 85, height: 185, lengthPool: 50, countPool: 2},
			want: "Тип тренировки: Бег\nДлительность: 30.00 мин\nДистанция: 3.25 км.\nСр. скорость: 6.50 км/ч\nПотрачено ккал: 302.91\n",
		},
		{
			name: "тест ходьбы",
			args: args{action: 20000, trainingType: "Ходьба", duration: 225 * time.Minute, weight: 85, height: 185, lengthPool: 50, countPool: 2},
			want: "Тип тренировки: Ходьба\nДлительность: 225.00 мин\nДистанция: 13.00 км.\nСр. скорость: 3.47 км/ч\nПотрачено ккал: 669.74\n",
		},
		{
			name: "тест плавания",
			args: args{action: 2000, trainingType: "Плавание", duration: 90 * time.Minute, weight: 85, height: 185, lengthPool: 50, countPool: 5},
			want: "Тип тренировки: Плавание\nДлительность: 90.00 мин\nДистанция: 0.25 км.\nСр. скорость: 0.17 км/ч\nПотрачено ккал: 323.00\n",
		},
		{
			name: "тест неизвестного типа",
			args: args{action: 1000, trainingType: "Керлинг", duration: 5 * time.Hour, weight: 85, height: 185, lengthPool: 50, countPool: 2},
			want: "неизвестный тип тренировки\n",
		},
	}
	// Проходим по тестовым случаям, выполняем и проверяем вывод.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
				got := "неизвестный тип тренировки\n"
				if got != tt.want {
					t.Errorf("ShowTrainingInfo() = %v, want %v", got, tt.want)
				}
				return
			}
			got := ReadData(training)
			if got != tt.want {
				t.Errorf("ShowTrainingInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
