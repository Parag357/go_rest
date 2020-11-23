package datastore

import (
	"github.com/jinzhu/gorm"
	"rest/model"
)

type ProductDataStore struct {
	db *gorm.DB
}

func NewProductDataStore(db *gorm.DB) ProductDataStore {
	return ProductDataStore{
		db: db,
	}
}

func (pd ProductDataStore) Create(model *model.Product) (err error) {
	return pd.db.Create(model).Error
}

func (pd ProductDataStore) Delete(model *model.Product, id string) {
	pd.db.Delete(model, id)
}

func (pd ProductDataStore) Save(model *model.Product) (err error) {
	return pd.db.Save(model).Error
}


func (pd ProductDataStore) GetCategorisedProducts(params map[string][]string) []model.Product{
	var prod []model.Product
	id, cat_ok := params["categoryId"]
	var db *gorm.DB
	if cat_ok{
		db = pd.db.Where("category_id = ?", id[0])
	}
	sort, sort_ok := params["sort"]
	if sort_ok{
		_, ord_ok := params["order"]
		if sort[0] == "price"{
			if ord_ok{
				if cat_ok{
					db = db.Order("price desc")
				}else{
					db = pd.db.Order("price desc")
				}
			}else{
				if cat_ok{
					db = db.Order("price")
				}else{
					db = pd.db.Order("price")
				}
			}
		}else{
			if ord_ok{
				if cat_ok{
					db = db.Order("expiry desc")
				}else{
					db = pd.db.Order("expiry desc")
				}
			}else{
				if cat_ok{
					db = db.Order("expiry")
				}else{
					db = pd.db.Order("expiry")
				}
			}
		}
	}
	if cat_ok || sort_ok{
		db = db.Find(&prod)
	}else{
		db=pd.db.Find(&prod)
	}
	//fmt.Println(prod)
	return prod
}

func (pd ProductDataStore) GetProductForUpdate(query string,id string,prod *model.Product)(err error){
	return pd.db.Where(query,id).Find(prod).Error
}