package delegate

import (
	"fmt"
	"github.com/luke513009828/crawlab-core/interfaces"
)

func GetEventName(d *ModelDelegate, method interfaces.ModelDelegateMethod) (eventName string) {
	return getEventName(d, method)
}

func getEventName(d *ModelDelegate, method interfaces.ModelDelegateMethod) (eventName string) {
	if method == interfaces.ModelDelegateMethodSave {
		hasChange := d.hasChange()
		if hasChange {
			method = interfaces.ModelDelegateMethodChange
		}
	}
	return fmt.Sprintf("model:%s:%s", d.colName, method)
}
