package db

import (
	"encoding/json"
	"io"
	"strconv"
)

type Product struct {
	Id          int32  `bson:"id,omitempty" json:"id,omitempty"`
	Brand       string `bson:"brand,omitempty" json:"brand,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Image       string `bson:"image,omitempty" json:"image,omitempty"`
	Price       int32  `bson:"price,omitempty" json:"price,omitempty"`
}

type Products []*Product

func (p *Product) ToString() string {
	s := "Product { id: " + strconv.Itoa(int(p.Id)) + ", brand: " + p.Brand + ", description: " + p.Description + ", image: " + p.Image + ", price: $" + strconv.Itoa(int(p.Price)) + " }"
	return s
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (px *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(px)
}

func (px *Products) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(px)
}
