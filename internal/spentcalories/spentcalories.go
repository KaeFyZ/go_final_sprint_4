package spentcalories

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Функция принимает строку с данными формата "3456,Ходьба,3h00m",
// которая содержит количество шагов, вид активности и продолжительность активности.
func parseTraining(data string) (int, string, time.Duration, error) {
	//Разделить строку на слайс строк.
	var Str []string = strings.Split(data, ",")

	//Проверить, чтобы длина слайса была равна 3, так как в строке данных у нас количество шагов, вид активности и продолжительность.
	if len(Str) != 3 {
		return 0, "", 0, errors.New("некорректный формат входящей строки")
	}

	//Преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки.
	//При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	steps, err := strconv.Atoi(string(Str[0]))
	if err != nil {
		return 0, "", 0, err
	}
	activity := string(Str[1])

	//Преобразовать третий элемент слайса в time.Duration.
	//В пакете time есть метод для парсинга строки в time.Duration.
	//Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	duration, err := time.ParseDuration(string(Str[2]))
	if err != nil {
		return 0, "", 0, err
	}

	// проверки на корректность значений шагов и продолжительности
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше нуля")
	}

	//Если всё прошло без ошибок, верните количество шагов, вид активности, продолжительность и nil (для ошибки).
	return steps, activity, duration, nil

}

// Функция принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {
	//рассчитайте длину шага. Для этого умножьте высоту пользователя на коэффициент длины шага stepLengthCoefficient.
	//Соответствующая константа уже определена в пакете.
	stepLength := height * stepLengthCoefficient

	//умножьте пройденное количество шагов на длину шага.
	distanceMeters := float64(steps) * stepLength

	//разделите полученное значение на число метров в километре (mInKm, константа определена в пакете).
	return distanceMeters / mInKm
}

// Функция принимает количество шагов steps, рост пользователя height и продолжительность активности duration  и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	//Проверить, что продолжительность duration больше 0. Если это не так, вернуть 0.
	if duration <= 0 {
		return 0
	}

	//Вычислить дистанцию с помощью distance().
	dist := distance(steps, height)

	//Вычислить и вернуть среднюю скорость.
	//Для этого разделите дистанцию на продолжительность в часах.
	//Чтобы перевести продолжительность в часы, воспользуйтесь функцией из пакета time.
	return dist / (duration.Hours())

}

// Функция принимает:
// data string — строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов,
// вид активности и продолжительность активности.
// weight, height float64 — вес (кг.) и рост (м.) пользователя.
// И возвращает два значения:
// string — строка с информацией о тренировке в формате, приведенном ниже.
// error — ошибку, при ее возникновении внутри функции.
func TrainingInfo(data string, weight, height float64) (string, error) {
	//Получить значения из строки данных с помощью функции parseTraining(), обработать возможные ошибки и вывести их в лог с помощью log.Println(err).
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	//Проверить, какой вид тренировки был передан в строке, которую парсили (лучше использовать switch). Для каждого из видов тренировки рассчитать дистанцию, среднюю скорость и калории.
	var distanceKm float64
	var avgSpeed float64
	var calories float64
	switch activity {
	case "Бег":
		distanceKm = distance(steps, height)
		avgSpeed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		distanceKm = distance(steps, height)
		avgSpeed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default: //Если был передан неизвестный тип тренировки, вернуть ошибку с текстом неизвестный тип тренировки.
		return "", errors.New("неизвестный тип тренировки")
	}

	//Для каждого вида тренировки сформировать и вернуть строку, образец которой был представлен выше.
	result := "Тип тренировки: " + activity + "\n" +
		"Длительность: " + strconv.FormatFloat(duration.Hours(), 'f', 2, 64) + " ч.\n" +
		"Дистанция: " + strconv.FormatFloat(distanceKm, 'f', 2, 64) + " км.\n" +
		"Скорость: " + strconv.FormatFloat(avgSpeed, 'f', 2, 64) + " км/ч\n" +
		"Сожгли калорий: " + strconv.FormatFloat(calories, 'f', 2, 64) + "\n"

	return result, nil

}

// Функция принимает количество шагов steps, вес пользователя weight, рост пользователя height и продолжительность активности duration  и возвращает количество потраченных калорий при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//Проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
	switch {
	case weight <= 0:
		return 0, errors.New("вес должен быть больше нуля")
	case height <= 0:
		return 0, errors.New("рост должен быть больше нуля")
	case steps <= 0:
		return 0, errors.New("количество шагов должно быть больше нуля")
	case duration <= 0:
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	//Рассчитать среднюю скорость с помощью meanSpeed().
	avgspeed := meanSpeed(steps, height, duration)

	//Рассчитать и вернуть количество калорий. Для этого:
	var calories float64

	//Переведите продолжительность в минуты с помощью функции из пакета time.
	minutes := duration.Minutes()

	//Умножьте вес пользователя на среднюю скорость и продолжительность в минутах.
	caloriesFactor := weight * avgspeed * minutes

	//Разделите результат на число минут в часе для получения количества потраченных калорий.
	calories = caloriesFactor / minInH

	return calories, nil
}

// Функция принимает количество шагов steps, вес пользователя weight, рост пользователя height и продолжительность активности duration  и возвращает количество потраченных калорий при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	//Проверить входные параметры на корректность используя switch. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку.
	switch {
	case weight <= 0:
		return 0, errors.New("вес должен быть больше нуля")
	case height <= 0:
		return 0, errors.New("рост должен быть больше нуля")
	case steps <= 0:
		return 0, errors.New("количество шагов должно быть больше нуля")
	case duration <= 0:
		return 0, errors.New("продолжительность должна быть больше нуля")
	}

	//Рассчитать среднюю скорость с помощью meanSpeed().
	avgspeed := meanSpeed(steps, height, duration)

	//Рассчитать и вернуть количество калорий. Для этого:
	var calories float64

	//Переведите продолжительность в минуты с помощью функции из пакета time.
	minutes := duration.Minutes()

	//Умножьте вес пользователя на среднюю скорость и продолжительность в минутах.
	caloriesFactor := walkingCaloriesCoefficient * weight * avgspeed * minutes

	//Разделите результат на число минут в часе для получения количества потраченных калорий.
	calories = caloriesFactor / minInH

	return calories, nil
}
