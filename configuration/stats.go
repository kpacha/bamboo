package configuration

import (
	"log"
	"time"
)

type Stats struct {
	Enabled  bool
	InfluxDB *Influxdb
	StatsD   *StatsD
}

func (s *Stats) CreateClient() {
	if !s.Enabled {
		return
	}
	log.Println("Stats are enabled")
	if s.StatsD != nil {
		s.StatsD.CreateClient()
	}
	if s.InfluxDB != nil {
		s.InfluxDB.CreateClient()
	}
}

func (s *Stats) LogReloadError() {
	if s.StatsD != nil {
		s.StatsD.Increment(1.0, "haproxy.reload.error", 1)
	}
	if s.InfluxDB != nil {
		s.InfluxDB.Save("haproxy.reload.error", 1)
	}
}

func (s *Stats) LogReloadSuccess(duration time.Duration) {
	if s.StatsD != nil {
		s.StatsD.Timing(1.0, "haproxy.reload.marathon.duration", duration)
		s.StatsD.Increment(1.0, "haproxy.reload.marathon.reloaded", 1)
	}
	if s.InfluxDB != nil {
		s.InfluxDB.Save("haproxy.reload.marathon", duration.Nanoseconds())
	}
}

func (s *Stats) LogReloadSkipped() {
	if s.StatsD != nil {
		s.StatsD.Increment(1.0, "haproxy.reload.skipped", 1)
	}
	if s.InfluxDB != nil {
		s.InfluxDB.Save("haproxy.reload.skipped", 1)
	}
}

func (s *Stats) LogMarathonCallback() {
	if s.StatsD != nil {
		s.StatsD.Increment(1.0, "callback.marathon", 1)
	}
	if s.InfluxDB != nil {
		s.InfluxDB.Save("callback.marathon", 1)
	}
}

func (s *Stats) LogReloadDomain() {
	if s.StatsD != nil {
		s.StatsD.Increment(1.0, "reload.domain", 1)
	}
	if s.InfluxDB != nil {
		s.InfluxDB.Save("reload.domain", 1)
	}
}

func (s *Stats) LogRestart() {
	if s.StatsD != nil {
		s.StatsD.Increment(1.0, "restart", 1)
	}
	if s.InfluxDB != nil {
		s.InfluxDB.Save("restart", 1)
	}
}
