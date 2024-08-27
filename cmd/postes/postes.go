package postes

//post
type PostUser struct {
	Title string `json:"title"`
	Post  string `json:"post"`
}

//slice with postes
type Postes struct {
	Items []PostUser
}

// function
func New() *Postes { //init clear slice
	return &Postes{}
}

func (p *Postes) Add(post PostUser) {
	p.Items = append(p.Items, post)
}

func (p *Postes) GetAll() []PostUser {
	return p.Items
}
