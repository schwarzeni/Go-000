package appgroup

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

// App 应用需要实现 Start 和 Stop 接口，分别对应应用的启动和终止
type App interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// AppGroup 管理多个同时运行的应用
type AppGroup struct {
	apps        []App
	errgroup    errgroup.Group
	ctx         context.Context
	cancel      func()
	stopTimeout time.Duration
	quitSigs    []os.Signal
}

func NewAppGroup(options ...func(app *AppGroup)) *AppGroup {
	appGroup := &AppGroup{stopTimeout: time.Second * 2}
	appGroup.ctx, appGroup.cancel = context.WithCancel(context.Background())
	for _, option := range options {
		option(appGroup)
	}
	if len(appGroup.quitSigs) > 0 {
		appGroup.Add(&signalListener{
			ctx:     appGroup.ctx,
			cancel:  appGroup.cancel,
			signals: appGroup.quitSigs,
		})
	}
	return appGroup
}

// Add 添加一个应用
func (appGroup *AppGroup) Add(app App) { appGroup.apps = append(appGroup.apps, app) }

// Start 启动全部注册的应用
func (appGroup *AppGroup) Start() {
	for _, app := range appGroup.apps {
		currApp := app
		appGroup.errgroup.Go(func() error { return currApp.Start(appGroup.ctx) })
		appGroup.errgroup.Go(func() error {
			<-appGroup.ctx.Done()
			ctx, cancel := context.WithTimeout(context.Background(), appGroup.stopTimeout)
			defer cancel()
			return currApp.Stop(ctx)
		})
	}
}

// Wait 等待全部应用的退出
func (appGroup *AppGroup) Wait() error { return appGroup.errgroup.Wait() }

// WithStopTimeout 设置 App 停止的超时时间
func WithStopTimeout(d time.Duration) func(app *AppGroup) {
	return func(app *AppGroup) { app.stopTimeout = d }
}

// WithContext 设置 App 的上下文
func WithContext(ctx context.Context) func(app *AppGroup) {
	return func(app *AppGroup) { app.ctx, app.cancel = context.WithCancel(ctx) }
}

// WithGracefullyQuit 监听退出信号
func WithGracefullyQuit(signals ...os.Signal) func(app *AppGroup) {
	return func(app *AppGroup) { app.quitSigs = signals }
}

// signalListener 实现处理优雅退出中的系统信号监听
type signalListener struct {
	ctx     context.Context
	cancel  func()
	signals []os.Signal
}

func (s *signalListener) Start(ctx context.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, s.signals...)
	select {
	case <-ctx.Done():
	case <-sigs:
		log.Println("catch system term signal, quit all tasks in group")
		s.cancel()
	}
	return nil
}

func (s *signalListener) Stop(ctx context.Context) error { return nil }
