package logger

type tagSetting struct {
	defaultFlag bool
	tags        map[string]bool
}

func newTagSetting() *tagSetting {
	return &tagSetting{
		defaultFlag: true,
		tags:        map[string]bool{},
	}
}

func (this *tagSetting) SetAll(is bool) *tagSetting {
	this.defaultFlag = is
	for key, _ := range this.tags {
		this.tags[key] = is
	}
	return this
}

func (this *tagSetting) Set(tag string, is bool) *tagSetting {
	this.tags[tag] = is
	return this
}

func (this *tagSetting) Get(tag string) bool {
	if ret, exist := this.tags[tag]; exist {
		return ret
	}
	return this.defaultFlag
}
