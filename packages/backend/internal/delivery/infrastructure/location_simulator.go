package infrastructure

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/tsucchinoko/food-delivery-tracker/internal/delivery/domain"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/events"
	"github.com/tsucchinoko/food-delivery-tracker/internal/shared/types"
)

// LocationSimulator は配達中のドライバー位置をシミュレーションする。
type LocationSimulator struct {
	publisher events.EventPublisher
	rng       *rand.Rand
}

// NewLocationSimulator は新しい LocationSimulator を生成する。
func NewLocationSimulator(publisher events.EventPublisher, rng *rand.Rand) *LocationSimulator {
	return &LocationSimulator{
		publisher: publisher,
		rng:       rng,
	}
}

// SimulateDeliveryRoute はレストランから配達先へ向かうランダムウォークをシミュレーションする。
// goroutine 内で呼び出し、1秒ごとに位置更新を配信する。
// コンテキストがキャンセルされるか配達完了まで動作する。
func (s *LocationSimulator) SimulateDeliveryRoute(
	ctx context.Context,
	driverID types.DriverID,
	from domain.Location,
	to domain.Location,
	onComplete func(),
) {
	const (
		updateInterval = 1 * time.Second
		totalSteps     = 30 // 約30秒で配達完了
		noise          = 0.0005 // ランダムノイズ（±約50m）
	)

	current := from
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	step := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			step++
			progress := float64(step) / float64(totalSteps)

			if progress >= 1.0 {
				// 配達完了: 最終位置を配信
				event := domain.NewDriverLocationUpdated(driverID.String(), to.Lat(), to.Lng())
				if err := s.publisher(ctx, event); err != nil {
					log.Printf("[location-sim] failed to publish final location: %v", err)
				}
				log.Printf("[location-sim] driver %s delivery complete", driverID)
				if onComplete != nil {
					onComplete()
				}
				return
			}

			// 目的地方向に進む + ランダムノイズ
			base := from.LerpTo(to, progress)
			noiseX := (s.rng.Float64() - 0.5) * 2 * noise
			noiseY := (s.rng.Float64() - 0.5) * 2 * noise
			current = domain.NewLocation(base.Lat()+noiseX, base.Lng()+noiseY)

			event := domain.NewDriverLocationUpdated(driverID.String(), current.Lat(), current.Lng())
			if err := s.publisher(ctx, event); err != nil {
				log.Printf("[location-sim] failed to publish location: %v", err)
				continue
			}
			log.Printf("[location-sim] driver %s step %d/%d (%.4f, %.4f)",
				driverID, step, totalSteps, current.Lat(), current.Lng())
		}
	}
}
