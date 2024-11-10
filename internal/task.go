package internal

type Task struct {
    Id          int `json:"id"` 
    Name        string `json:"name"`
    Description string `json:"description"`
    Complexity  string `json:"complexity"`
    Category      string `json:"category"`
}

