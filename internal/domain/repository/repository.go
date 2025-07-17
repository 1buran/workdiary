package repository

import (
	"github.com/1buran/workdiary/internal/domain/valueobject"
)

type WorkdiaryRepository interface {
	Add(d valueobject.Day)
	List() []valueobject.Day
	MaxDayHours() float32
	TotalHours() float32
	TotalAmount() float32
	Compact()
}
