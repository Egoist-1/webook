package service

type MonitorLink interface {
	Analysis(str string) error
}
type monitorLink struct {
}

func (svc *monitorLink) Analysis(str string) error {

}
