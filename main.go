package amocrmrepo

import (
	"fmt"
	"sync"
)

type MemoryStorage struct {
	mu sync.RWMutex
	//posts    map[string]*model.Post
	//comments map[string][]*model.Comment
}

// NewMemoryStorage создает новое in-memory хранилище.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		//posts:    make(map[string]*model.Post),
		//comments: make(map[string][]*model.Comment),
	}
}

type Accounts struct {
}

type Account_Integration struct {
}

func main() {
	fmt.Println("git")
}
