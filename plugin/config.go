package plugin

import (
	"github.com/luke513009828/crawlab-core/config"
	"path"
)

const DefaultPluginFsPathBase = "plugins"
const DefaultPluginDirName = "plugins"

var DefaultPluginDirPath = path.Join(config.DefaultConfigDirPath, DefaultPluginDirName)
