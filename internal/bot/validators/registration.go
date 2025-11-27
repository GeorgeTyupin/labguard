package validators

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrNameEmpty    = errors.New("ФИО не может быть пустым")
	ErrNameTooLong  = errors.New("ФИО слишком длинное (макс. 100 символов)")
	ErrNameInvalid  = errors.New("ФИО должно содержать только буквы и пробелы")
	ErrGroupEmpty   = errors.New("группа не может быть пустой")
	ErrGroupInvalid = errors.New("неверный формат группы (пример: ИВТ-123)")
)

// ValidateName проверяет ФИО пользователя
func ValidateName(name string) error {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return ErrNameEmpty
	}

	if len(name) > 100 {
		return ErrNameTooLong
	}

	// Только буквы (кириллица, латиница), пробелы, дефисы
	matched, _ := regexp.MatchString(`^[а-яА-ЯёЁa-zA-Z\s\-]+$`, name)
	if !matched {
		return ErrNameInvalid
	}

	return nil
}

// ValidateGroup проверяет название группы
func ValidateGroup(group string) error {
	group = strings.TrimSpace(group)

	if len(group) == 0 {
		return ErrGroupEmpty
	}

	// Формат: буквы-цифры (например: ИВТ-123, ФИИТ-21)
	matched, _ := regexp.MatchString(`^[а-яА-ЯёЁa-zA-Z]+-\d+$`, group)
	if !matched {
		return ErrGroupInvalid
	}

	return nil
}
