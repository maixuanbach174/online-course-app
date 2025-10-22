module github.com/maixuanbach174/online-course-app/internal/education

go 1.25.2

require (
	github.com/maixuanbach174/online-course-app/internal/common v0.0.0-00010101-000000000000
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

replace github.com/maixuanbach174/online-course-app/internal/common => ../common/
