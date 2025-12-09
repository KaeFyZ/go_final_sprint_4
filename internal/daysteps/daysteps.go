package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	. "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// Функция принимает строку с данными,
// которая содержит количество шагов и продолжительность прогулки в формате 3h50m(3 часа 50 минут).
func parsePackage(data string) (int, time.Duration, error) {
	//формат входящей строки: "678,0h50m"
	//Разделить строку на слайс строк.
	var str []string = strings.Split(data, ",")

	//Проверить, чтобы длина слайса была равна 2, так как в строке данных у нас количество шагов и продолжительность.
	if len(str) != 2 {
		return 0, 0, errors.New("")
	}

	//Преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки.
	//При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	steps, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, 0, err
	}

	//Проверить: количество шагов должно быть больше 0. Если это не так, вернуть нули и ошибку.
	if steps <= 0 {
		return 0, 0, errors.New("steps must be greater than zero")
	}

	//Преобразовать второй элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration.
	//Обработать возможные ошибки.
	//При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку.
	duration, err := time.ParseDuration(str[1])
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, errors.New("duration must be greater than zero")
	}

	//Если всё прошло без ошибок, верните количество шагов, продолжительность и nil (для ошибки).

	return steps, duration, nil

}

// Функция должна парсить строку с данными с помощью parsePackage(),
// вычислять дистанцию в километрах и количество потраченных калорий и возвращать строку в таком виде:
// Количество шагов: 792.
// Дистанция составила 0.51 км.
// Вы сожгли 221.33 ккал.
func DayActionInfo(data string, weight, height float64) string {
	//Получить данные о количестве шагов и продолжительности прогулки с помощью функции parsePackage().

	steps, duration, err := parsePackage(data)
	//В случае возникновения ошибки вывести её на экран и вернуть пустую строку.
	if err != nil {
		log.Println(err)
		return ""
	}

	//Проверить, чтобы количество шагов было больше 0. В противном случае вернуть пустую строку.
	if steps <= 0 {
		log.Println("steps must be greater than zero")
		return ""
	}

	//Вычислить дистанцию в метрах. Дистанция равна произведению количества шагов на длину шага.
	//Константа stepLength (длина шага) уже определена в коде.
	distance := float64(steps) * stepLength

	//Перевести дистанцию в километры, разделив её на число метров в километре (константа mInKm, определена в пакете).
	distanceKm := distance / mInKm

	//Вычислить количество калорий, потраченных на прогулке.
	//Функция для вычисления калорий func WalkingSpentCalories() будет определена в пакете spentcalories, которую вы тоже реализуете.
	calories, err := WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return err.Error()
	}

	//Сформировать строку, которую будете возвращать, пример которой был представлен выше.
	//Количество шагов: 792.
	//Дистанция составила 0.51 км.
	//Вы сожгли 221.33 ккал.
	//используя fmt
	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, calories)

	return result

}
