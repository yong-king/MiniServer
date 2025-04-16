package main

import "fmt"

type ServiceConfig struct {
	A string
	B string
	C int

	X struct{}
	Y Info
}

type Info struct {
	addr string
}

// NewServiceConfig 创建一个ServiceConfig的函数
func NewServiceConfig(a, b string, c int) *ServiceConfig {
	return &ServiceConfig{
		A: a,
		B: b,
		C: c,
	}
}

const defaultValueC = 1

// target: 想要A和B必须传入，C可以不传，不传就用默认值
func NewServiceConfig2(a, b string, c ...int) *ServiceConfig {
	valueC := defaultValueC
	if len(c) > 0 {
		valueC = c[0]
	}
	return &ServiceConfig{
		A: a,
		B: b,
		C: valueC,
	}
}

// Options模式
type FuncServiceConfigOption func(*ServiceConfig)

func NewServiceConfig3(a, b string, opst ...FuncServiceConfigOption) *ServiceConfig {
	sc := &ServiceConfig{
		A: a,
		B: b,
		C: defaultValueC,
	}
	// 针对可能传进来的FuncServiceConfigOption参数做处理
	for _, opt := range opst {
		opt(sc)
	}
	return sc
}

func WithC(c int) FuncServiceConfigOption {
	return func(sc *ServiceConfig) {
		sc.C = c
	}
}

func WithY(info Info) FuncServiceConfigOption {
	return func(sc *ServiceConfig) {
		sc.Y = info
	}
}

// 上面方法存在一些问题
// 我可以直接通过sc.C来修改，那我的WithC有什么用呢
// 不想被直接修改，只能通过我们提供的方法来进行修改

type config struct {
	name string
	age  int
}

const defaultName = "ysh"

func NewConfig(age int, opts ...ConfigOption) *config {
	cfg := &config{
		age:  age,
		name: defaultName,
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

type ConfigOption interface {
	apply(*config)
}

type funcConfigOption struct {
	f func(*config)
}

func (f funcConfigOption) apply(c *config) {
	f.f(c)
}

func NewfuncConfigOption(f func(*config)) funcConfigOption {
	return funcConfigOption{f: f}
}

func WithName(name string) ConfigOption {
	return NewfuncConfigOption(func(c *config) { c.name = name })
}

func main() {
	//info := Info{addr: "127.0.0.1"}
	//sc := NewServiceConfig3("ysh", "py", WithC(10), WithY(info))
	//fmt.Printf("sc:%#v\n", sc)

	cfg := NewConfig(18)
	fmt.Printf("cfg:%#v\n", cfg)
	cfg2 := NewConfig(18, WithName("张三"))
	fmt.Printf("cfg:%#v\n", cfg2)
}
