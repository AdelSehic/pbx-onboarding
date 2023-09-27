package ari

import config "ari/configs"

type Ari struct {
	Cfg *config.Config
}

func New() *Ari { return &Ari{} }
