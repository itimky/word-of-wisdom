package shield

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	cfg          Config
	secret       string
	hashCalc     hashCalc
	quoteService quoteService
}

func NewService(
	cfg Config,
	hashCalc hashCalc,
	quoteService quoteService,
) *Service {
	return &Service{
		cfg:          cfg,
		hashCalc:     hashCalc,
		quoteService: quoteService,
	}
}

func (s *Service) Init() error {
	if err := s.updateSecret(); err != nil {
		return fmt.Errorf("update secret: %w", err)
	}

	go func() {
		ticker := time.NewTicker(s.cfg.SecretUpdateInterval)
		for range ticker.C {
			err := s.updateSecret()
			if err != nil {
				logrus.WithError(err).Error("periodic secret update")
			}
		}
	}()

	return nil
}

func (s *Service) updateSecret() error {
	secret, err := randomSecret(s.cfg.SecretLength)
	if err != nil {
		return fmt.Errorf("update secret: %w", err)
	}
	s.secret = secret
	return nil
}

func (s *Service) HandleInitial(clientIP string) InitialResult {
	initialHash := s.hashCalc.CalcInitialHash(clientIP, s.cfg.TourLength, s.secret)
	return InitialResult{InitialHash: initialHash, TourLength: s.cfg.TourLength}
}

func (s *Service) HandleTourComplete(clientIP string, request TourCompleteRequest) TourCompleteResult {
	var response TourCompleteResult
	if s.hashCalc.VerifyHash(request.InitialHash, request.LastHash, s.cfg.TourLength, clientIP, s.secret, s.cfg.GuideSecrets) {
		response.Granted = true
		response.Quote = s.quoteService.Get()
	} else {
		response.Granted = false
	}

	return response
}
