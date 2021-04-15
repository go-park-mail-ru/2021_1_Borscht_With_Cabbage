package models

import "github.com/microcosm-cc/bluemonday"

func (u *SuccessUserResponse) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Name = sanitizer.Sanitize(u.Name)
	u.Email = sanitizer.Sanitize(u.Email)
	u.Avatar = sanitizer.Sanitize(u.Avatar)
	u.Phone = sanitizer.Sanitize(u.Phone)
	u.MainAddress = sanitizer.Sanitize(u.MainAddress)
}

func (u *User) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Name = sanitizer.Sanitize(u.Name)
	u.Email = sanitizer.Sanitize(u.Email)
	u.Avatar = sanitizer.Sanitize(u.Avatar)
	u.Phone = sanitizer.Sanitize(u.Phone)
	u.MainAddress = sanitizer.Sanitize(u.MainAddress)
}

func (u *UserData) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Name = sanitizer.Sanitize(u.Name)
	u.Email = sanitizer.Sanitize(u.Email)
	u.Avatar = sanitizer.Sanitize(u.Avatar)
	u.Phone = sanitizer.Sanitize(u.Phone)
}

func (u *UserImageResponse) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Filename = sanitizer.Sanitize(u.Filename)
}

func (u *SuccessRestaurantResponse) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Title = sanitizer.Sanitize(u.Title)
	u.AdminEmail = sanitizer.Sanitize(u.AdminEmail)
	u.AdminPhone = sanitizer.Sanitize(u.AdminPhone)
	u.Avatar = sanitizer.Sanitize(u.Avatar)
	u.Description = sanitizer.Sanitize(u.Description)
}

func (u *Dish) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Name = sanitizer.Sanitize(u.Name)
	u.Description = sanitizer.Sanitize(u.Description)
	u.Image = sanitizer.Sanitize(u.Image)
}

func (u *DeleteSuccess) Sanitize() {
}

func (u *DishImageResponse) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Filename = sanitizer.Sanitize(u.Filename)
}

func (u *SectionWithDishes) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	for _, value := range u.Dishes {
		value.Sanitize()
	}
	u.SectionName = sanitizer.Sanitize(u.SectionName)
}

func (u *RestaurantImageResponse) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Filename = sanitizer.Sanitize(u.Filename)
}

func (u *Section) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Name = sanitizer.Sanitize(u.Name)
}

func (u *RestaurantWithDishes) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	for i := range u.Dishes {
		u.Dishes[i].Sanitize()
	}
	u.AdminEmail = sanitizer.Sanitize(u.AdminEmail)
	u.AdminPhone = sanitizer.Sanitize(u.AdminPhone)
	u.Avatar = sanitizer.Sanitize(u.Avatar)
	u.Description = sanitizer.Sanitize(u.Description)
	u.Title = sanitizer.Sanitize(u.Title)
}

func (u *RestaurantInfo) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.AdminEmail = sanitizer.Sanitize(u.AdminEmail)
	u.AdminPhone = sanitizer.Sanitize(u.AdminPhone)
	u.Avatar = sanitizer.Sanitize(u.Avatar)
	u.Description = sanitizer.Sanitize(u.Description)
	u.Title = sanitizer.Sanitize(u.Title)
}

func (u *DishInOrder) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Name = sanitizer.Sanitize(u.Name)
	u.Image = sanitizer.Sanitize(u.Image)
}

func (u *Order) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Address = sanitizer.Sanitize(u.Address)
	u.Restaurant = sanitizer.Sanitize(u.Restaurant)
	u.DeliveryTime = sanitizer.Sanitize(u.DeliveryTime)
	u.OrderTime = sanitizer.Sanitize(u.OrderTime)
	u.Status = sanitizer.Sanitize(u.Status)

	for i := range u.Foods {
		u.Foods[i].Sanitize()
	}
}

func (u *DishInBasket) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Image = sanitizer.Sanitize(u.Image)
	u.Name = sanitizer.Sanitize(u.Name)
}

func (u *BasketForUser) Sanitize() {
	sanitizer := bluemonday.UGCPolicy()
	u.Restaurant = sanitizer.Sanitize(u.Restaurant)
	u.RestaurantImage = sanitizer.Sanitize(u.RestaurantImage)

	for i := range u.Foods {
		u.Foods[i].Sanitize()
	}
}
