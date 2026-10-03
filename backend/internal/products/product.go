package products

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}
