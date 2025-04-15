package main

import (
	"google.golang.org/grpc/resolver"
)

const (
	mySchme    = "ysh"
	myEndpoion = "resolver.ysh.com"
)

var addrs = []string{"127.0.0.1:8972", "127.0.0.1:8973", "127.0.0.1:8974"}

// 自定义name resolver，实现Resolver接口
type yshResolver struct {
	target      resolver.Target
	cc          resolver.ClientConn
	addresStore map[string][]string
}

func (r *yshResolver) ResolveNow(o resolver.ResolveNowOptions){
	addrsStrs := r.addresStore[r.target.Endpoint()]
	addrList := make([]resolver.Address, len(addrsStrs))
	for i, s := range addrsStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*yshResolver) Close() {}

// 需实现 Builder 接口
type yshResloverBuiler struct{}

func (*yshResloverBuiler) Build(target resolver.Target, cc resolver.ClientConn,  opts resolver.BuildOptions)  (resolver.Resolver, error) {
	r := &yshResolver{
		target: target,
		cc: cc,
		addresStore: map[string][]string{
			myEndpoion: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (*yshResloverBuiler) Scheme() string {return mySchme}

func init() {
	resolver.Register(&yshResloverBuiler{})
}