package bot

// HotUpdater 配置热更新接口
type HotUpdater interface {
	HotReload()
}

func HotUpdateModuleConfig(sign <-chan struct{}) {
	for {
		select {
		case <-sign:
			for _, module := range modules {
				if hotUpdater, ok := module.Instance.(HotUpdater); ok {
					hotUpdater.HotReload()
				}
			}
		default:

		}
	}
}
