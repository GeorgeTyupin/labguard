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
	ErrGroupInvalid = errors.New("неверный формат группы (примеры: 111, ИВТ-123, М3О-111БВ-11)")
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

	// Формат: число, буквы/цифры-цифры, или буквы/цифры-цифры+буквы-цифры (например: 111, ИВТ-123, М3О-111БВ-11)
	matched, _ := regexp.MatchString(`^([0-9]+|[а-яА-ЯёЁa-zA-Z0-9]+-[0-9]+([а-яА-ЯёЁa-zA-Z]+-[0-9]+)?)$`, group)
	if !matched {
		return ErrGroupInvalid
	}

	return nil
}
