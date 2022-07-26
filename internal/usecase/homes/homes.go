package homes

type HomesUsecase struct {
	homesRepo homesRepo
}

func NewHomesUsecase(homeRepo homesRepo) *HomesUsecase {
	return &HomesUsecase{
		homesRepo: homeRepo,
	}
}
