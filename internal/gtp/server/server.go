package server

import (
	"fmt"
	"time"

	"github.com/itimky/word-of-wisom/pkg/gtp"

	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg      Config
	secret   string
	hashCalc hashCalc
}

func NewServer(
	cfg Config,
	hashCalc hashCalc,
) *Server {
	return &Server{
		cfg:      cfg,
		hashCalc: hashCalc,
	}
}

func (s *Server) Init() error {
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

func (s *Server) updateSecret() error {
	secret, err := randomSecret(s.cfg.SecretLength)
	if err != nil {
		return fmt.Errorf("update secret: %w", err)
	}
	s.secret = secret
	return nil
}

func (s *Server) CheckPuzzle(clientIP string, solution *gtp.PuzzleSolution) PuzzleCheckResult {
	var result PuzzleCheckResult
	if solution == nil {
		result.Type = Restricted
		initialHash := s.hashCalc.CalcInitialHash(clientIP, s.cfg.TourLength, s.secret)
		result.Puzzle = &gtp.Puzzle{InitialHash: initialHash, TourLength: s.cfg.TourLength}
	} else {
		if s.hashCalc.VerifyHash(solution.InitialHash, solution.LastHash, s.cfg.TourLength, clientIP, s.secret, s.cfg.GuideSecrets) {
			result.Type = Ok
		} else {
			result.Type = WrongSolution
		}
	}

	return result
}
