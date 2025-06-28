package main

import (
	"fmt"
	"math"
	"time"
)

// Определяет основные единицы измерения для расчетов.
const (
	MInKm      = 1000 //  количество метров в одном километре.
	MinInHours = 60   // количество минут в одном часе.
	LenStep    = 0.65 //  стандартная длина шага в метрах (используется для ходьбы и бега).
	CmInM      = 100  //  количество сантиметров в одном метре (для перевода роста).
)

// Training — базовая структура, описывающая общие характеристики тренировки.
// Содержит тип тренировки, количество действий (шаги/гребки), длину шага,
// длительность и вес пользователя.
type Training struct {
	TrainingType string        // Тип тренировки (например, "Бег", "Ходьба", "Плавание").
	Action       int           // Количество выполненных шагов или гребков.
	LenStep      float64       // Длина одного шага или гребка в метрах.
	Duration     time.Duration // Общая продолжительность тренировки.
	Weight       float64       // Вес пользователя в килограммах.
}

// distance вычисляет дистанцию в километрах на основе количества действий и длины шага.
// Используется для всех типов тренировок, где применима концепция шагов/гребков.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed вычисляет среднюю скорость в километрах в час.
// Деление на ноль предотвращается проверкой Duration.Hours().
func (t Training) meanSpeed() float64 {
	dist := t.distance()
	if t.Duration.Hours() == 0 {
		return 0
	}
	return dist / t.Duration.Hours()
}

// Calories возвращает базовое значение калорий (0).
// Этот метод переопределяется в конкретных типах тренировок для точного расчета.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage — структура для хранения и вывода информации о тренировке.
type InfoMessage struct {
	TrainingType string        // Тип тренировки.
	Duration     time.Duration // Длительность тренировки.
	Distance     float64       // Пройденное расстояние в километрах.
	Speed        float64       // Средняя скорость в км/ч.
	Calories     float64       // Количество потраченных калорий.
}

// TrainingInfo собирает данные о тренировке в структуру InfoMessage.
// Использует методы distance, meanSpeed и Calories для расчета.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String форматирует данные тренировки в читаемый текстовый вид.
// Использует форматирование с двумя знаками после запятой для чисел.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator — интерфейс, определяющий обязательные методы для расчета калорий.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы, специфичные для тренировки по бегу.
// CaloriesMeanSpeedMultiplier и CaloriesMeanSpeedShift используются в формуле калорий.
const (
	CaloriesMeanSpeedMultiplier = 18   // Множитель для средней скорости.
	CaloriesMeanSpeedShift      = 1.79 // Сдвиг для корректировки скорости.
)

// Running — структура для тренировки по бегу, наследует Training.
type Running struct {
	Training
}

// Calories вычисляет потраченные калории для бега.
// Формула: (18 * скорость + 1.79) * вес * длительность в минутах, адаптирована для часов.
func (r Running) Calories() float64 {
	durationInMin := r.Duration.Minutes()
	if durationInMin == 0 {
		return 0 // Возвращаем 0, если длительность нулевая, чтобы избежать деления на ноль.
	}
	return (CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * float64(r.Duration.Hours()*MinInHours)
}

// TrainingInfo возвращает полную информацию о тренировке по бегу.
func (r Running) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: r.TrainingType,
		Duration:     r.Duration,
		Distance:     r.distance(),
		Speed:        r.meanSpeed(),
		Calories:     r.Calories(),
	}
}

// Константы, специфичные для тренировки по ходьбе.
// CaloriesWeightMultiplier и CaloriesSpeedHeightMultiplier — коэффициенты для расчета.
const (
	CaloriesWeightMultiplier      = 0.035 // Коэффициент, зависящий от веса.
	CaloriesSpeedHeightMultiplier = 0.029 // Коэффициент, зависящий от скорости и роста.
)

// Walking — структура для тренировки по ходьбе, наследует Training.
type Walking struct {
	Training
	Height float64 // Рост пользователя в сантиметрах, используется в расчете калорий.
}

// Calories вычисляет потраченные калории для ходьбы.
// Формула учитывает вес, скорость и рост, адаптирована для минутной длительности.
func (w Walking) Calories() float64 {
	if w.Height == 0 {
		return 0 // Возвращаем 0, если рост не указан, чтобы избежать деления на ноль.
	}
	durationInMin := w.Duration.Minutes()
	if durationInMin == 0 {
		return 0 // Проверка на нулевую длительность.
	}
	return ((CaloriesWeightMultiplier*w.Weight + (math.Pow(w.meanSpeed(), 2)/w.Height/CmInM)*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours)
}

// TrainingInfo возвращает полную информацию о тренировке по ходьбе.
func (w Walking) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: w.TrainingType,
		Duration:     w.Duration,
		Distance:     w.distance(),
		Speed:        w.meanSpeed(),
		Calories:     w.Calories(),
	}
}

// Константы, специфичные для тренировки по плаванию.
// SwimmingLenStep — длина гребка, остальные — коэффициенты для расчета.
const (
	SwimmingLenStep                  = 1.38 // Длина одного гребка в метрах.
	SwimmingCaloriesMeanSpeedShift   = 1.1  // Сдвиг для корректировки скорости.
	SwimmingCaloriesWeightMultiplier = 2    // Множитель для веса пользователя.
)

// Swimming — структура для тренировки по плаванию, наследует Training.
type Swimming struct {
	Training
	LengthPool int // Длина одного бассейна в метрах.
	CountPool  int // Количество пересечений бассейна.
}

// meanSpeed вычисляет среднюю скорость для плавания.
// Учитывает длину бассейна и количество пересечений, проверяет на нулевые значения.
func (s Swimming) meanSpeed() float64 {
	if s.Duration.Hours() == 0 || s.LengthPool == 0 || s.CountPool == 0 {
		return 0 // Возвращаем 0 при нулевых значениях, чтобы избежать ошибок.
	}
	return float64(s.LengthPool) * float64(s.CountPool) / float64(MInKm) / s.Duration.Hours()
}

// Calories вычисляет потраченные калории для плавания.
// Формула зависит от скорости, веса и длительности тренировки.
func (s Swimming) Calories() float64 {
	durationInMin := s.Duration.Minutes()
	if durationInMin == 0 {
		return 0 // Проверка на нулевую длительность.
	}
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// TrainingInfo возвращает полную информацию о тренировке по плаванию.
func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     float64(s.LengthPool*s.CountPool) / float64(MInKm),
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// ReadData преобразует данные тренировки в отформатированную строку.
// Принимает объект, реализующий интерфейс CaloriesCalculator.
func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	return fmt.Sprint(info)
}

// Packet — структура для хранения данных, полученных от фитнес-трекера.
type Packet struct {
	Date  string // Дата в формате YYYYMMDD.
	Time  string // Время в формате HH:MM:SS.
	Steps int    // Количество шагов.
}

// ProcessPacket обрабатывает строку данных от трекера и добавляет новый пакет.
// Выводит сообщение об обработке и возвращает обновленный срез пакетов.
func ProcessPacket(packetStr string, packets []Packet) []Packet {
	fmt.Println("Пакет обработан:", packetStr)
	return append(packets, Packet{Date: "20250628", Time: "12:28:00", Steps: 5000})
}

func main() {
	// Инициализация среза для хранения пакетов данных от трекера.
	var packets []Packet
	packets = ProcessPacket("20250628 12:28:00,5000", packets)

	// Создание объекта тренировки по плаванию с заданными параметрами.
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000, // Количество гребков.
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute, // Длительность 90 минут.
			Weight:       85,               // Вес пользователя 85 кг.
		},
		LengthPool: 50, // Длина бассейна 50 метров.
		CountPool:  5,  // Пять пересечений бассейна.
	}

	// Создание объекта тренировки по ходьбе с заданными параметрами.
	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000, // 20000 шагов.
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute, // Длительность 3 часа 45 минут.
			Weight:       85,                           // Вес пользователя 85 кг.
		},
		Height: 185, // Рост пользователя 185 см.
	}

	// Создание объекта тренировки по бегу с заданными параметрами.
	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000, // 5000 шагов.
			LenStep:      LenStep,
			Duration:     30 * time.Minute, // Длительность 30 минут.
			Weight:       85,               // Вес пользователя 85 кг.
		},
	}

	// Вывод результатов обработки пакета и данных тренировок.
	fmt.Println("Обработка пакета:")
	fmt.Println(ReadData(swimming))
	fmt.Println(ReadData(walking))
	fmt.Println(ReadData(running))
}
