package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dsbasko/url-shortener/internal/entity"
)

/*
getLastID
Вспомогательная функция получения последнего ID.
Так как по условиям задачи, требовалось сделать автоматически инкрементирующий
id, пришлось написать свою функцию его поиска. Для того что-бы не пробегать
по всему файлу с начала до конца, пришлось написать функцию реверсивного прохода
по файлу.
1. Перенос указателя в начало файла;
2. Создание сканера с помощью вспомогательной функции scanLinesReversed;
3. Анмаршалинг полученных данных;
4. Получение и инкрементация последнего id.
*/
func (s *Storage) getLastID() string {
	uuid := "1"

	_, err := s.file.Seek(0, 0)
	if err != nil {
		s.log.Debug(fmt.Errorf("failed to seek to the beginning of the file: %s", err))
		return uuid
	}

	var lastLine string
	scanner := bufio.NewScanner(s.file)
	scanner.Split(s.scanLinesReversed)
	if scanner.Scan() {
		lastLine = scanner.Text()
	}

	var dataJSON entity.URL
	err = json.Unmarshal([]byte(lastLine), &dataJSON)
	if err != nil {
		s.log.Debug(fmt.Errorf("failed to unmarshal JSON data: %s", err))
		return uuid
	}

	lastUUID, err := strconv.Atoi(dataJSON.ID)
	if err != nil {
		s.log.Debug(fmt.Errorf("failed to convert string to int: %s", err))
		return uuid
	}
	uuid = strconv.Itoa(lastUUID + 1)

	return uuid
}

/*
scanLinesReversed
Вспомогательная функция реверсивного поиска по файлу.
1. Запуск цикла в обратном порядке;
2. Поиск символа переноса строки;
3. Пропуск первого такого символа.
*/
func (s *Storage) scanLinesReversed(data []byte, atEOF bool) (advance int, token []byte, err error) {
	attempts := int8(1)

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != '\n' {
			continue
		}

		if attempts == 0 {
			return len(data) - i, data[i+1:], nil
		} else {
			attempts--
		}
	}

	return len(data), data, nil
}
