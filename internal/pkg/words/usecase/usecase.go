package usecase

type WordUsecase struct {
}

func NewAudioUsecase() *WordUsecase {
	return &WordUsecase{}
}

// SignUp handles the user registration process.
func (u *WordUsecase) GetWord() string {
	return ""
}
