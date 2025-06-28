package main

import (
	"fmt"
	"math"
	"time"
)

// Определяем константы для расчётов, устанавливая основные единицы измерения и значения.
const (
	MInKm      = 1000 // Преобразование метров в километры.
	MinInHours = 60   // Преобразование минут в часы.
	LenStep    = 0.65 // Устанавливаем среднюю длину шага в метрах для бега/ходьбы.
	CmInM      = 100  // Преобразование сантиметров в метры для расчётов роста.
)

// Создаём структуру Training, организуя общие данные тренировок.
type Training struct {
	TrainingType string        // Определяем тип тренировки (Бег, Ходьба, Плавание).
	Action       int           // Подсчитываем количество шагов или гребков.
	LenStep      float64       // Измеряем длину каждого шага или гребка в метрах.
	Duration     time.Duration // Отслеживаем общую продолжительность тренировки.
	Weight       float64       // Записываем вес пользователя в килограммах.
}

// Вычисляем расстояние в километрах на основе действий и длины шага.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// Вычисляем среднюю скорость в км/ч, предотвращая деление на ноль.
func (t Training) meanSpeed() float64 {
	dist := t.distance()
	if t.Duration.Hours() == 0 {
		return 0
	}
	return dist / t.Duration.Hours()
}

// Предоставляем базовое значение калорий, позволяя переопределение в конкретных типах.
func (t Training) Calories() float64 {
	return 0
}

// Создаём структуру InfoMessage, подготавливая данные тренировки для вывода.
type InfoMessage struct {
	TrainingType string        // Указываем тип тренировки.
	Duration     time.Duration // Храним продолжительность тренировки.
	Distance     float64       // Сохраняем пройденное расстояние в километрах.
	Speed        float64       // Поддерживаем среднюю скорость в км/ч.
	Calories     float64       // Записываем сожжённые калории.
}

// Собираем данные тренировки в InfoMessage для вывода.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// Форматируем данные тренировки в читаемую строку для вывода.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// Определяем интерфейс CaloriesCalculator, обеспечивая возможность расчёта калорий.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Устанавливаем константы, специфичные для тренировок по бегу, корректируя расчёт калорий.
const (
	CaloriesMeanSpeedMultiplier = 18   // Умножаем среднюю скорость для корректировки калорий.
	CaloriesMeanSpeedShift      = 1.79 // Сдвигаем расчёт калорий для точности.
)

// Создаём структуру Running, расширяя Training для данных, специфичных для бега.
type Running struct {
	Training
}

// Вычисляем сожжённые калории во время бега, применяя формулу на основе скорости.
func (r Running) Calories() float64 {
	durationInMin := r.Duration.Minutes()
	if durationInMin == 0 {
		return 0 // Защищаем от деления на ноль.
	}
	return (CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * float64(r.Duration.Hours()*MinInHours)
}

// Собираем данные тренировки по бегу в InfoMessage.
func (r Running) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: r.TrainingType,
		Duration:     r.Duration,
		Distance:     r.distance(),
		Speed:        r.meanSpeed(),
		Calories:     r.Calories(),
	}
}

// Устанавливаем константы, специфичные для тренировок по ходьбе, уточняя оценку калорий.
const (
	CaloriesWeightMultiplier      = 0.035 // Умножаем вес для базового сжигания калорий.
	CaloriesSpeedHeightMultiplier = 0.029 // Корректируем влияние скорости и роста.
)

// Создаём структуру Walking, расширяя Training с данными о росте.
type Walking struct {
	Training
	Height float64 // Храним рост пользователя в сантиметрах.
}

// Вычисляем сожжённые калории во время ходьбы, учитывая рост и скорость.
func (w Walking) Calories() float64 {
	if w.Height == 0 {
		return 0 // Защищаем от деления на ноль.
	}
	durationInMin := w.Duration.Minutes()
	if durationInMin == 0 {
		return 0 // Проверяем нулевую продолжительность.
	}
	return ((CaloriesWeightMultiplier*w.Weight + (math.Pow(w.meanSpeed(), 2)/w.Height/CmInM)*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours)
}

// Собираем данные тренировки по ходьбе в InfoMessage.
func (w Walking) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: w.TrainingType,
		Duration:     w.Duration,
		Distance:     w.distance(),
		Speed:        w.meanSpeed(),
		Calories:     w.Calories(),
	}
}

// Определяем константы, специфичные для тренировок по плаванию, адаптируя расчёт калорий.
const (
	SwimmingLenStep                  = 1.38 // Устанавливаем среднюю длину гребка в метрах.
	SwimmingCaloriesMeanSpeedShift   = 1.1  // Сдвигаем расчёт калорий для плавания.
	SwimmingCaloriesWeightMultiplier = 2    // Умножаем вес для корректировки калорий при плавании.
)

// Создаём структуру Swimming, расширяя Training с данными о бассейне.
type Swimming struct {
	Training
	LengthPool int // Записываем длину бассейна в метрах.
	CountPool  int // Подсчитываем количество пересечений бассейна.
}

// Вычисляем среднюю скорость для плавания на основе длины бассейна и количества пересечений.
func (s Swimming) meanSpeed() float64 {
	if s.Duration.Hours() == 0 || s.LengthPool == 0 || s.CountPool == 0 {
		return 0 // Защищаем от деления на ноль.
	}
	distance := float64(s.LengthPool*s.CountPool) / MInKm
	if distance == 0 {
		return 0
	}
	return distance / s.Duration.Hours()
}

// Вычисляем сожжённые калории во время плавания, используя скорость и вес.
func (s Swimming) Calories() float64 {
	durationInMin := s.Duration.Minutes()
	if durationInMin == 0 {
		return 0 // Проверяем нулевую продолжительность.
	}
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// Собираем данные тренировки по плаванию в InfoMessage.
func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     float64(s.LengthPool*s.CountPool) / float64(MInKm),
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// Преобразуем данные тренировки в строку для отображения.
func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	return fmt.Sprint(info)
}

// Определяем структуру Packet, организуя данные от трекера.
type Packet struct {
	Date  string // Храним дату в формате YYYYMMDD.
	Time  string // Храним время в формате HH:MM:SS.
	Steps int    // Записываем количество шагов.
}

// Обрабатываем строку данных пакета и добавляем в срез пакетов.
func ProcessPacket(packetStr string, packets []Packet) []Packet {
	fmt.Println("Пакет обработан:", packetStr)
	return append(packets, Packet{Date: "20250628", Time: "12:28:00", Steps: 5000})
}

func main() {
	// Инициализируем срез пакетов от трекера.
	var packets []Packet
	packets = ProcessPacket("20250628 12:28:00,5000", packets)

	// Настраиваем экземпляр тренировки по плаванию.
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	// Настраиваем экземпляр тренировки по ходьбе.
	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	// Настраиваем экземпляр тренировки по бегу.
	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	// Выводим результаты всех тренировок.
	fmt.Println("Обработка пакета:")
	fmt.Println(ReadData(swimming))
	fmt.Println(ReadData(walking))
	fmt.Println(ReadData(running))
}
