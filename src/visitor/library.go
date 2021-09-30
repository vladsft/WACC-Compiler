package visitor

import (
	"strings"
	"sync"
	"wacc_32/ast"
)

type libManager struct {
	visited             map[string]struct{} //This is a set in Go
	aliases             map[string]string
	calledFunctions     chan string
	concurrentFunctions chan string //Functions called with the wacc keyword
	functions           chan *ast.Function
	userTypes           chan *ast.UserType
	mu                  sync.Mutex
}

func newLibManager() *libManager {
	return &libManager{
		visited:             make(map[string]struct{}),
		aliases:             make(map[string]string),
		calledFunctions:     make(chan string, 10), //10 is just an arbitrary buffer size
		concurrentFunctions: make(chan string, 10), //10 is just an arbitrary buffer size
		functions:           make(chan *ast.Function),
		userTypes:           make(chan *ast.UserType),
	}
}

func (l *libManager) derive() *libManager {
	return &libManager{
		visited:             l.visited,
		aliases:             make(map[string]string),
		calledFunctions:     l.calledFunctions,
		concurrentFunctions: l.concurrentFunctions,
		functions:           l.functions,
		userTypes:           l.userTypes,
	}
}

//Add a lock
func (l *libManager) addLib(filepath, alias string) (exists bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, exists = l.visited[filepath]
	l.visited[filepath] = struct{}{}
	l.aliases[alias] = filepath
	return exists
}

func (l *libManager) getLib(alias string) (string, bool) {
	lib, ok := l.aliases[alias]
	return lib, ok
}

//formatFilePath replaces all /'s from the filepath with $
func formatFilepath(filepath string) string {
	noSlashes := strings.Replace(filepath, "/", "$", -1)
	splitAtDots := strings.Split(noSlashes, ".")
	return splitAtDots[len(splitAtDots)-2] + "$"
}
