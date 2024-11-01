package models

type Project struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Annotation  string     `json:"annotation"`
	Time        string     `json:"time"`
	CardImgSrc  string     `json:"card_img_src"`
	BannerDesktop string   `json:"banner_desktop"`
	BannerMobile  string   `json:"banner_mobile"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
	Categories  []Category `json:"categories"`
}