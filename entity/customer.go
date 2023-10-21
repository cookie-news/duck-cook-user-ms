package entity

type Customer struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	User             string `json:"user"`
	Pass             string `json:"pass"`
	Name             string `json:"name"`
	ImageProfilePath string `json:"image_profile_path"`
}
