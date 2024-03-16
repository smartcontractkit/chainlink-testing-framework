package client

import (
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type AnvilPeriodicMiner struct {
	Client            *AnvilClient
	interval          time.Duration
	batchSendInterval time.Duration
	batchCapacity     int64
	stop              chan struct{}
}

func NewAnvilMiner(url string) *AnvilPeriodicMiner {
	return &AnvilPeriodicMiner{
		Client: NewAnvilClient(url),
		stop:   make(chan struct{}),
	}
}

func (m *AnvilPeriodicMiner) MinePeriodically(interval time.Duration) {
	m.interval = interval
	go func() {
		for {
			select {
			case <-m.stop:
				log.Info().Msg("anvil miner exiting")
				return
			default:
				if err := m.Client.Mine(nil); err != nil {
					log.Err(err).Send()
				}
			}
			time.Sleep(m.interval)
		}
	}()
}

func (m *AnvilPeriodicMiner) Stop() {
	m.stop <- struct{}{}
}

func (m *AnvilPeriodicMiner) MineBatch(capacity int64, checkInterval time.Duration, sendInterval time.Duration) {
	m.interval = checkInterval
	m.batchCapacity = capacity
	m.batchSendInterval = sendInterval
	ticker := time.NewTicker(m.batchSendInterval)
	go func() {
		for {
			resp, err := m.Client.TxPoolStatus(nil)
			if err != nil {
				log.Err(err).Send()
			}
			pendingTx, err := strconv.ParseInt(resp.Result.Pending[2:], 16, 64)
			if err != nil {
				log.Err(err).Msg("failed to convert pending tx from hex to dec")
			}
			log.Info().Int64("Pending", pendingTx).Msg("Batch has pending transactions")
			if pendingTx >= m.batchCapacity {
				if err := m.Client.Mine(nil); err != nil {
					log.Err(err).Send()
				}
				log.Info().Int64("Transactions", pendingTx).Msg("Block mined")
			}
			select {
			case <-m.stop:
				log.Info().Msg("anvil miner exiting")
				ticker.Stop()
				return
			case <-ticker.C:
				if err := m.Client.Mine(nil); err != nil {
					log.Err(err).Send()
				}
				log.Info().Int64("Transactions", pendingTx).Msg("Block mined by timeout")
			default:
			}
			time.Sleep(m.interval)
		}
	}()
}
