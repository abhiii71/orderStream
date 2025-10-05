package internal

import "context"

type AccountService interface {
	Register(ctx context.Context, name, email, password string) (string, error)
}

type accountService struct {
	accountRepo AccountRepository
}

func NewAccountService(r AccountRepository) AccountService {
	return &accountService{r}
}

func (s *accountService) Register() (string, error){
	_, err : s.accountRepo.GetAccountByEmail(ctx, )
	return 
}
